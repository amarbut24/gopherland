package gopherusers

import (
	"testing"

	"github.com/amarbut24/gopherland/auth"
)

var envvars = make(map[string]string)

// used for local testing. Go tests ran via Github Workflows will use
// <insert workflow block example>
func init() {
	envvars["AZURE_TENANT_ID"] = "02e9f3a0-53a5-4898-bb6e-e97008b17be7"
	envvars["AZURE_CLIENT_ID"] = "98b51714-780b-41ab-b0a9-aaa8833b6be2"
	envvars["AZURE_CLIENT_CERTIFICATE_PATH"] = "/home/anthony/selfsigned.crt"
	auth.SetAzureEnv(envvars)
}

func TestClient(t *testing.T) {
	_, _, err := auth.AzureGraphClient()
	if err != nil {
		t.Errorf("unable to authenticate to azure ad %v", err)
	}
}

func TestNewUser(t *testing.T) {
	u1 := GopherUser{
		FirstName:                     "Test",
		LastName:                      "User",
		MailNickname:                  "tuser",
		UserPrincipalName:             "tuser@gopherland.onmicrosoft.com",
		ForceChangePasswordNextSignIn: true,
		DisplayName:                   "Test User",
		AccountEnabled:                false,
	}

	client, _, err := auth.AzureGraphClient()
	if err != nil {
		t.Errorf("unable to authenticate to azure ad %v", err)
	}

	newUser, err := u1.NewUser(client)
	if err != nil {
		t.Errorf("got err %v, wanted no errors", err)
	} else if newUser == nil {
		t.Logf("existing user %v was found when attempting to create new user", u1.UserPrincipalName)
	} else {
		t.Logf("created user %s\n", *newUser.GetUserPrincipalName())
	}
}

// func ExampleGetUserByID() {
// 	client, err := auth.AzureGraphClient()
// 	u, err := GetUserByID(client, "8e086b41-7bc0-4b4a-9c3e-0a7ff59d710b")
// 	if err != nil {
// 		log.Fatalf("unable to locate user: %v", err)
// 	}
// 	// Output: olleh
// }
