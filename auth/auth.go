package auth

import (
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	a "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

//AzureGraphClient uses default azure credentials
//and returns a msgraph client
func AzureGraphClient() (*msgraphsdk.GraphServiceClient, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("authentication failure: %+v", err)
	}

	auth, err := a.NewAzureIdentityAuthenticationProviderWithScopes(cred, nil)
	if err != nil {
		return nil, fmt.Errorf("error authentication provider: %v", err)
	}

	adapter, err := msgraphsdk.NewGraphRequestAdapter(auth)
	if err != nil {
		return nil, fmt.Errorf("error creating adapter: %v", err)
	}
	client := msgraphsdk.NewGraphServiceClient(adapter)
	return client, nil
}

//SetAzureEnv sets env vars required to authenticate to Azure AD
func SetAzureEnv(env map[string]string) error {
	for k, v := range env {
		log.Printf("setting env var key=%v +value=%v", k, v)
		err := os.Setenv(k, v)
		if err != nil {
			return fmt.Errorf("ran into an error when trying to set env for key:%v, value:%v", k, v)
		}

		ev := os.Getenv(k)
		if ev != v {
			return fmt.Errorf("was unable to locate env var for key=%v", k)
		} else {
			log.Printf("env var for key=%v was succesfully set to value=%v", k, ev)
		}
	}
	return nil
}
