package bot

import (
	"log"

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

	authorizer.TdlibParameters

	tdlibClient, err := client.NewClient(authorizer)
	if err != nil {
		log.Fatalf("Error initializing Telegram client: %v", err)
	}

	return &TdLib{
		TdLibClient: tdlibClient,
	}
}
