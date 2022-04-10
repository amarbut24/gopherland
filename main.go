package main

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/amarbut24/gopherland/auth"
	a "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

var envvars = make(map[string]string)

func main() {
	// set env vars
	envvars["AZURE_TENANT_ID"] = "02e9f3a0-53a5-4898-bb6e-e97008b17be7"
	envvars["AZURE_CLIENT_ID"] = "98b51714-780b-41ab-b0a9-aaa8833b6be2"
	envvars["AZURE_CLIENT_CERTIFICATE_PATH"] = "/home/anthony/selfsigned.crt"
	auth.SetAzureEnv(envvars)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("Authentication failure: %+v", err)
	}

	if err != nil {
		fmt.Printf("Error creating credentials: %v\n", err)
	}

	auth, err := a.NewAzureIdentityAuthenticationProviderWithScopes(cred, []string{"Files.Read"})
	if err != nil {
		fmt.Printf("Error authentication provider: %v\n", err)
		return
	}

	adapter, err := msgraphsdk.NewGraphRequestAdapter(auth)
	if err != nil {
		fmt.Printf("Error creating adapter: %v\n", err)
		return
	}
	client := msgraphsdk.NewGraphServiceClient(adapter)

	result, err := client.Me().Drive().Get(nil)
	if err != nil {
		fmt.Printf("Error getting the drive: %v\n", err)
	}
	fmt.Printf("Found Drive : %v\n", *result.GetId())
}
