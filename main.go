package main

import (
	"log"

	"github.com/amarbut24/gopherland/auth"
	"github.com/amarbut24/gopherland/gopherusers"
)

var envvars = make(map[string]string)

func main() {
	//Set env vars for NewDefaultAzureCredential
	envvars["AZURE_TENANT_ID"] = "02e9f3a0-53a5-4898-bb6e-e97008b17be7"
	envvars["AZURE_CLIENT_ID"] = "98b51714-780b-41ab-b0a9-aaa8833b6be2"
	envvars["AZURE_CLIENT_CERTIFICATE_PATH"] = "/home/anthony/selfsigned.crt"
	auth.SetAzureEnv(envvars)

	log.Printf("creating msgraph client")
	client, _, err := auth.AzureGraphClient()
	if err != nil {
		log.Fatalf("unable to create msgraph client with error: %v", err)
	}

	// Create new user struct
	u1 := gopherusers.GopherUser{
		FirstName:                     "Anthony",
		LastName:                      "Marbut",
		MailNickname:                  "amarbut",
		UserPrincipalName:             "amarbut@gopherland.onmicrosoft.com",
		ForceChangePasswordNextSignIn: true,
		DisplayName:                   "Anthony Marbut",
		AccountEnabled:                false,
	}

	// Create new user
	newUser, err := u1.NewUser(client)
	if err != nil {
		log.Print(err)
	} else if newUser == nil {
		log.Printf("existing user %v was found when attempting to create new user", u1.UserPrincipalName)
	} else {
		log.Printf("created user %s\n", *newUser.GetUserPrincipalName())
	}

	// Find the user by objectID
	byId, err := gopherusers.GetUserByID(client, *newUser.GetId())
	if err != nil {
		log.Fatalf("unable to locate user: %v", err)
	}
	log.Println("found user based on Id =", *byId.GetId())

	// Find the user by UserPrincipalName
	byUPN, err := gopherusers.GetUserByUPN(client, *newUser.GetUserPrincipalName())
	if err != nil {
		log.Fatalf("unable to locate user: %v", err)
	}
	log.Println("found user based on UPN =", *byUPN.GetUserPrincipalName())

	// Enable account
	d := *byUPN.GetDisplayName()
	log.Printf("our new user %v is disabled by default, lets enable\n", d)
	b := true
	byUPN.SetAccountEnabled(&b)
	log.Printf("account enabled status = %v\n", *byUPN.GetAccountEnabled())

	// Now lets delete the account
	err = gopherusers.DeleteUserByID(client, *byUPN.GetId())
	if err != nil {
		log.Fatalf("unable to locate user: %v", err)
	} else {
		log.Printf("user %v was deleted\n", *byUPN.GetDisplayName())
	}
}
