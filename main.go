package main

import (
	"github.com/amarbut24/gopherland/auth"
)

var envvars = make(map[string]string)

func main() {
	// set env vars
	envvars["AZURE_TENANT_ID"] = "02e9f3a0-53a5-4898-bb6e-e97008b17be7"
	envvars["AZURE_CLIENT_ID"] = "98b51714-780b-41ab-b0a9-aaa8833b6be2"
	envvars["AZURE_CLIENT_CERTIFICATE_PATH"] = "/etc/ssl/private/selfsigned.key"
	auth.SetAzureEnv(envvars)
}
