package main

import (
	"fmt"
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

	log.Printf("Creating msgraph client")
	client, adapter, err := auth.AzureGraphClient()
	if err != nil {
		log.Fatalf("unable to create msgraph client with error: %v", err)
	}

	allUsers, err := gopherusers.GetAllUsers(client, adapter)
	if err != nil {
		log.Fatalf("unable to grab all users with error: %v", err)
	}
	fmt.Printf("Found %v users", len(allUsers))

}
