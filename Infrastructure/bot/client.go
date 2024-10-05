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
	// Create a listener and defer closing it to ensure resources are freed
	listener := t.TdLibClient.GetListener()
	defer listener.Close()

	// Notify when a new message listener is active
	fmt.Println("New message listener activated")
	fmt.Printf("this is the listener %v\n", listener)

	// Start listening for updates in a goroutine so it doesn't block the main thread
	fmt.Printf("listener.Updates: %v\n", listener.Updates) // Check what kind of updates you're receiving

	for update := range listener.Updates {
		switch u := update.(type) {
		case *client.UpdateNewMessage:
			newMessage := u.Message
			// Check if the message has text content
			if messageText, ok := newMessage.Content.(*client.MessageText); ok {
				fmt.Printf("New message text received: %s\n", messageText.Text.Text)
				go checkIfChatIDExists(chatIDs, newMessage.ChatId, t.Bot, userID)

			} else {
				fmt.Println("Received new message but it does not contain text.")
			}
		default:
			fmt.Printf("Unknown update type received: %v\n", u)
		}
	}

}

// checkIfChatIDExists checks if the chatID exists in the chatIDs slice if it does send using bot
func checkIfChatIDExists(chatIDs []int64, chatID int64, bot *tgbotapi.BotAPI, userID int64) {
	for _, id := range chatIDs {
		if id == chatID {
			bot.Send(tgbotapi.NewMessage(userID, "interested chat found"))
			return
		}
	}
}
