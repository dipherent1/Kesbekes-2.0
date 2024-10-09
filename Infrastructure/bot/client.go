package bot

import (
	"context"
	"encoding/json"
	"fmt"
	config "kesbekes/Config"
	ai "kesbekes/Infrastructure/AI"
	repositories "kesbekes/Repositories"
	"log"
	"path/filepath"
	"time"

	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zelenin/go-tdlib/client"
)

type TdLib struct {
	TdLibClient *client.Client
	Bot         *tgbotapi.BotAPI
	RedisClient *redis.Client
	BotRepo     *repositories.TelegramRepository
	ai          *ai.AI
}

func NewTdLib(bot *tgbotapi.BotAPI, botRepo *repositories.TelegramRepository, ai *ai.AI) *TdLib {
	// Initialize authorizer

	// Configure TDLib
	authorizer := client.ClientAuthorizer()
	go client.CliInteractor(authorizer)

	authorizer.TdlibParameters <- &client.SetTdlibParametersRequest{
		UseTestDc:           false,
		DatabaseDirectory:   filepath.Join(".tdlib", "database"),
		FilesDirectory:      filepath.Join(".tdlib", "files"),
		UseFileDatabase:     true,
		UseChatInfoDatabase: true,
		UseMessageDatabase:  true,
		UseSecretChats:      false,
		ApiId:               config.APIID,
		ApiHash:             config.APIHash,
		SystemLanguageCode:  "en",
		DeviceModel:         "Server",
		SystemVersion:       "1.0.0",
		ApplicationVersion:  "1.0.0",
		// EnableStorageOptimizer: true,
		// IgnoreFileNames:        false,
	}

	_, err := client.SetLogVerbosityLevel(&client.SetLogVerbosityLevelRequest{
		NewVerbosityLevel: 1,
	})
	if err != nil {
		log.Fatalf("SetLogVerbosityLevel error: %s", err)
	}

	tdlibClient, err := client.NewClient(authorizer)
	if err != nil {
		log.Fatalf("NewClient error: %s", err)
	}

	optionValue, err := client.GetOption(&client.GetOptionRequest{
		Name: "version",
	})
	if err != nil {
		log.Fatalf("GetOption error: %s", err)
	}

	log.Printf("TDLib version: %s", optionValue.(*client.OptionValueString).Value)

	me, err := tdlibClient.GetMe()
	if err != nil {
		log.Fatalf("GetMe error: %s", err)
	}

	log.Printf("Me: %s %s ", me.FirstName, me.LastName)

	redis := config.CreateRedisClient()

	return &TdLib{
		TdLibClient: tdlibClient,
		Bot:         bot,
		RedisClient: redis,
		BotRepo:     botRepo,
		ai:          ai,
	}
}

func (t *TdLib) Get10Updates() {
	updates, err := t.TdLibClient.GetChatHistory(&client.GetChatHistoryRequest{
		ChatId: -1001132580599,
		Limit:  20,
	})

	if err != nil {
		log.Fatalf("GetUpdates error: %s", err)
	}

	log.Printf("Total messages retrieved: %d", len(updates.Messages))

	for _, update := range updates.Messages {
		content, err := json.MarshalIndent(update.Content, "", "  ")
		if err != nil {
			log.Printf("Error marshalling update content: %s", err)
			continue
		}
		log.Printf("Update: %s", content)
	}
}

func (t *TdLib) Listen(chatIDs []int64, userID int64) {
	// Create a listener
	prefernces, err := t.BotRepo.GetUserPreferences(userID)
	if err != nil {
		log.Fatalf("Error getting user preferences: %s", err)
	}

	listener := t.TdLibClient.GetListener()
	defer listener.Close()

	// Channel for passing new messages to workers
	messageChannel := make(chan *client.UpdateNewMessage, 100) // Buffer size as needed

	// Start a fixed number of worker goroutines for checking preferences
	workerCount := 5 // Adjust this value based on expected load
	for i := 0; i < workerCount; i++ {
		go t.CheckIsPreference(userID, prefernces) // Worker continuously checks Redis for messages
	}

	// Start goroutine to process updates and enqueue messages to Redis
	fmt.Println("New message listener activated")
	for update := range listener.Updates {
		switch u := update.(type) {
		case *client.UpdateNewMessage:
			// Call processMessage here to filter and enqueue relevant messages
			go t.processMessage(u, chatIDs)
		default:
			fmt.Printf("Unknown update type received: %v\n", u)
		}
	}

	// Close the message channel when listener is done
	close(messageChannel)
}

// processMessage processes each new message and checks if it belongs to a chat in chatIDs
func (t *TdLib) processMessage(update *client.UpdateNewMessage, chatIDs []int64) {
	newMessage := update.Message
	// Check if the message has text content
	if messageText, ok := newMessage.Content.(*client.MessageText); ok {
		fmt.Printf("New message text received: %s\n", messageText.Text.Text)

		if ChatIDExists(chatIDs, newMessage.ChatId) {
			// Serialize and enqueue the message for processing
			EnqueueMessage(newMessage, t.RedisClient)
		}
	} else {
		fmt.Println("Received new message but it does not contain text.")
	}
}

// checkIfChatIDExists checks if the chatID exists in the chatIDs slice and sends a message if it does
func ChatIDExists(chatIDs []int64, chatID int64) bool {
	for _, id := range chatIDs {
		if id == chatID {
			return true
		}
	}
	return false
}

func (t *TdLib) ProcessTxt(txt string, prefernces []string) (bool, error) {
	return t.ai.IsPreferred(txt, prefernces)
}

func EnqueueMessage(newMessage *client.Message, redisClient *redis.Client) {
	ctx := context.Background()
	// Serialize the message into JSON
	messageBytes, err := json.Marshal(newMessage)
	if err != nil {
		log.Printf("Error serializing message: %v", err)
		return
	}
	// Push to Redis list
	_, err = redisClient.LPush(ctx, "messageQueue", messageBytes).Result()
	if err != nil {
		log.Printf("Error enqueuing message: %s", err)
	}
}

func (t *TdLib) CheckIsPreference(userID int64, prefernces []string) {
	// get preference
	for {
		// Pop a message from the Redis queue
		messageJSON, err := t.RedisClient.BRPop(context.Background(), 0*time.Second, "messageQueue").Result()
		if err != nil {
			log.Printf("Error retrieving message from queue: %s", err)
			continue
		}

		// Deserialize the message JSON
		var newMessage client.Message
		err = json.Unmarshal([]byte(messageJSON[1]), &newMessage)
		if err != nil {
			log.Printf("Error deserializing message: %s", err)
			continue
		}

		// Process the message text to check if it matches preference
		if messageText, ok := newMessage.Content.(*client.MessageText); ok {
			isPreferred, err := t.ProcessTxt(messageText.Text.Text, prefernces)
			if err != nil {
				log.Printf("Error processing text: %s", err)
				continue
			}
			if isPreferred {
				msg := fmt.Sprintf("Interested topic in chat ID: %d", newMessage.ChatId)
				t.Bot.Send(tgbotapi.NewMessage(userID, msg))
			}
		}
	}
}
