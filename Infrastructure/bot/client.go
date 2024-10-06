package bot

import (
	"encoding/json"
	"fmt"
	config "kesbekes/Config"
	"log"
	"path/filepath"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zelenin/go-tdlib/client"
)

type TdLib struct {
	TdLibClient *client.Client
	Bot         *tgbotapi.BotAPI
}

func NewTdLib(bot *tgbotapi.BotAPI) *TdLib {
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

	return &TdLib{
		TdLibClient: tdlibClient,
		Bot:         bot,
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
	listener := t.TdLibClient.GetListener()
	defer listener.Close()

	// Channel for passing new messages to workers
	messageChannel := make(chan *client.UpdateNewMessage, 100) // Buffer size as needed

	// Start a fixed number of worker goroutines
	workerCount := 5 // Adjust this value based on the expected load
	for i := 0; i < workerCount; i++ {
		go func() {
			for update := range messageChannel {
				processMessage(update, chatIDs, t.Bot, userID)
			}
		}()
	}

	// Process updates and send new messages to the channel
	fmt.Println("New message listener activated")
	for update := range listener.Updates {
		switch u := update.(type) {
		case *client.UpdateNewMessage:
			// Send the new message to the messageChannel for processing by a worker
			messageChannel <- u
		default:
			fmt.Printf("Unknown update type received: %v\n", u)
		}
	}

	// Close the message channel when listener is done
	close(messageChannel)
}

// processMessage processes each new message and checks if it belongs to a chat in chatIDs
func processMessage(update *client.UpdateNewMessage, chatIDs []int64, bot *tgbotapi.BotAPI, userID int64) {
	newMessage := update.Message
	// Check if the message has text content
	if messageText, ok := newMessage.Content.(*client.MessageText); ok {
		fmt.Printf("New message text received: %s\n", messageText.Text.Text)
		checkIfChatIDExists(chatIDs, newMessage.ChatId, bot, userID)
	} else {
		fmt.Println("Received new message but it does not contain text.")
	}
}

// checkIfChatIDExists checks if the chatID exists in the chatIDs slice and sends a message if it does
func checkIfChatIDExists(chatIDs []int64, chatID int64, bot *tgbotapi.BotAPI, userID int64) {
	for _, id := range chatIDs {
		if id == chatID {
			bot.Send(tgbotapi.NewMessage(userID, "Interested chat found"))
			return
		}
	}
}
