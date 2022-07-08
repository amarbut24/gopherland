package gopherusers

import (
	"testing"
	"time"

	"github.com/amarbut24/gopherland/auth"
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
var testUser GopherUser

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
		t.Fatalf("got err %v, wanted no errors", err)
	} else {
		t.Logf("created user %s\n", newUser.DisplayName)
		testUser = newUser
	}

	// throw in a sleep as remaining test depend on the new user
	time.Sleep(time.Second * 30)
}

func TestGetUserByID(t *testing.T) {
	client, _, err := auth.AzureGraphClient()
	if err != nil {
		t.Errorf("unable to authenticate to azure ad %v", err)
	}

	byID, err := GetUserByID(client, testUser.ObjectID)
	if err != nil {
		t.Errorf("unable to locate user: %v", err)
	}
	t.Logf("found user based on Id = %v", byID.ObjectID)
}

func TestGetUserByUPN(t *testing.T) {
	client, _, err := auth.AzureGraphClient()
	if err != nil {
		t.Errorf("unable to authenticate to azure ad %v", err)
	}

	byUPN, err := GetUserByID(client, testUser.UserPrincipalName)
	if err != nil {
		t.Errorf("unable to locate user: %v", err)
	}
	t.Logf("found user based on UPN = %v", byUPN.UserPrincipalName)
}

func TestDeleteUserByID(t *testing.T) {
	client, _, err := auth.AzureGraphClient()
	if err != nil {
		t.Errorf("unable to authenticate to azure ad %v", err)
	}

	err = DeleteUserByID(client, testUser.ObjectID)
	if err != nil {
		t.Errorf("unable to locate user: %v", err)
	}
	t.Logf("user %v was deleted\n", testUser.DisplayName)
}

func TestGetAllUsers(t *testing.T) {
	client, adapter, err := auth.AzureGraphClient()
	if err != nil {
		t.Errorf("unable to authenticate to azure ad %v", err)
	}

	allUsers, err := GetAllUsers(client, adapter)
	if err != nil {
		t.Errorf("unable to grab all users with error: %v", err)
	}
	t.Logf("Found %v users", len(allUsers))
}
