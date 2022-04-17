package users

import (
	"log"

	"github.com/amarbut24/gopherland/auth"
)

func ExampleGetUserByID() {
	client, err := auth.AzureGraphClient()
	u, err := users.GetUserByID(client, "8e086b41-7bc0-4b4a-9c3e-0a7ff59d710b")
	if err != nil {
		log.Fatalf("unable to locate user: %v", err)
	}
	// Output: olleh
}
