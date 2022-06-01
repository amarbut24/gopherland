package gopherusers

import (
	"testing"

	"github.com/amarbut24/gopherland/auth"
	models "github.com/microsoftgraph/msgraph-sdk-go/models"
)

//Uncomment below block if you need to run test locally
/*
func init() {
	envvars := make(map[string]string)
	envvars["AZURE_TENANT_ID"] = "02e9f3a0-53a5-4898-bb6e-e97008b17be7"
	envvars["AZURE_CLIENT_ID"] = "98b51714-780b-41ab-b0a9-aaa8833b6be2"
	envvars["AZURE_CLIENT_CERTIFICATE_PATH"] = "/home/anthony/selfsigned.crt"
	auth.SetAzureEnv(envvars)
}
*/

//testUser will be populated vai TestNewUser
//and will be used throughout the remaining tests
var testUser models.Userable

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
		testUser = newUser
	}
}

func TestGetUserByID(t *testing.T) {
	client, _, err := auth.AzureGraphClient()
	if err != nil {
		t.Errorf("unable to authenticate to azure ad %v", err)
	}

	byID, err := GetUserByID(client, *testUser.GetId())
	if err != nil {
		t.Errorf("unable to locate user: %v", err)
	}
	t.Logf("found user based on Id = %v", *byID.GetId())
}

func TestGetUserByUPN(t *testing.T) {
	client, _, err := auth.AzureGraphClient()
	if err != nil {
		t.Errorf("unable to authenticate to azure ad %v", err)
	}

	byUPN, err := GetUserByID(client, *testUser.GetUserPrincipalName())
	if err != nil {
		t.Errorf("unable to locate user: %v", err)
	}
	t.Logf("found user based on Id = %v", *byUPN.GetId())
}

func TestDeleteUserByID(t *testing.T) {
	client, _, err := auth.AzureGraphClient()
	if err != nil {
		t.Errorf("unable to authenticate to azure ad %v", err)
	}

	err = DeleteUserByID(client, *testUser.GetId())
	if err != nil {
		t.Errorf("unable to locate user: %v", err)
	}
	t.Logf("user %v was deleted\n", *testUser.GetDisplayName())
}
