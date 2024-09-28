package bot

import (
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

	// optionValue, err := client.GetOption(&client.GetOptionRequest{
	// 	Name: "version",
	// })
	// if err != nil {
	// 	log.Fatalf("GetOption error: %s", err)
	// }

	// log.Printf("TDLib version: %s", optionValue.(*client.OptionValueString).Value)

	me, err := tdlibClient.GetMe()
	if err != nil {
		log.Fatalf("GetMe error: %s", err)
	}

	log.Printf("Me: %s %s", me.FirstName, me.LastName)

	return &TdLib{
		TdLibClient: tdlibClient,
	}

	// Wait for the next authorization state change

}
