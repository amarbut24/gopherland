package main

import (
	"fmt"
	"log"

	"github.com/amarbut24/gopherland/auth"
	"github.com/amarbut24/gopherland/gophergroups"
)

var envvars = make(map[string]string)

func main() {
	//Set env vars for NewDefaultAzureCredential
	envvars["AZURE_TENANT_ID"] = "02e9f3a0-53a5-4898-bb6e-e97008b17be7"
	envvars["AZURE_CLIENT_ID"] = "98b51714-780b-41ab-b0a9-aaa8833b6be2"
	envvars["AZURE_CLIENT_CERTIFICATE_PATH"] = "/home/anthony/selfsigned.crt"
	auth.SetAzureEnv(envvars)

	log.Printf("creating msgraph client")
	client, adapter, err := auth.AzureGraphClient()
	if err != nil {
		log.Fatalf("unable to create msgraph client with error: %v", err)
	}

	// g, _ := gophergroups.GetGroupByID(client, "46307819-3f04-46ef-bdfe-cd4891498454")
	// fmt.Println(g)

	// Create new user struct
	// g1 := gophergroups.GopherGroup{
	// 	DisplayName:     "Test Group",
	// 	Description:     "Group created using Go",
	// 	SecurityEnabled: true,
	// 	MailNickname:    "testgroup",
	// }

	// // Create new group
	// newGroup, err := g1.NewGroup(client)
	// if err != nil {
	// 	log.Print(err)
	// } else {
	// 	log.Printf("created group %s\n", *newGroup.GetDisplayName())
	// }
	g, _ := gophergroups.GetGroupByDisplayName(client, "Test Group")
	fmt.Println("Return group via DisplayName", *g.GetId())

	ag, _ := gophergroups.GetAllGroups(client, adapter)
	fmt.Println("Found all groups", ag)
}
