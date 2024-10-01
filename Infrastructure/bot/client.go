package bot

import (
	"encoding/json"
	config "kesbekes/Config"
	"log"
	"path/filepath"

	"github.com/zelenin/go-tdlib/client"
)

type TdLib struct {
	TdLibClient *client.Client
}

func NewTdLib() *TdLib {
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
