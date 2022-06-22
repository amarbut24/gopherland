package main

import (
	"fmt"
	"log"
	"strconv"

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

	users := []gopherusers.GopherUser{}
	for i := 0; i < 10; i++ {
		un := strconv.Itoa(i)
		u1 := gopherusers.GopherUser{
			FirstName:                     "Test",
			LastName:                      "User" + un,
			MailNickname:                  "testuser" + un,
			UserPrincipalName:             fmt.Sprintf("testuser%v@gopherland.onmicrosoft.com", un),
			ForceChangePasswordNextSignIn: true,
			DisplayName:                   "Test User" + un,
			AccountEnabled:                false,
		}
		users = append(users, u1)
	}

	ch := make(chan gopherusers.ConcurrentResult)
	gopherusers.CNewUsers(ch, users, client)

	// ag, _ := gophergroups.GetAllGroups(client, adapter)
	// fmt.Println("Found all groups", ag)

	// g, _ := gophergroups.GetGroupByDisplayName(client, "Test Group")
	// g.AddMembers(client, []string{"73e66338-4510-47a7-bfe7-25a2ae9d6024"})
}
