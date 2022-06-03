# Azure AD Auth Package

# Create a Graph Client
```go
//Below block is using ClientCertificateCredential to auth to Azure AD using a Service Principal  
envvars := make(map[string]string)
envvars["AZURE_TENANT_ID"] = "02e9f3a0-53a5-4898-bb6e-e97008b17be7"
envvars["AZURE_CLIENT_ID"] = "98b51714-780b-41ab-b0a9-aaa8833b6be2"
envvars["AZURE_CLIENT_CERTIFICATE_PATH"] = "/home/anthony/selfsigned.crt"
auth.SetAzureEnv(envvars)

//Once you have the credentials configured below block will return a client
client, _, err := auth.AzureGraphClient()
if err != nil {
    t.Errorf("unable to authenticate to azure ad %v", err)
}

//Optionally you can also have an adapter returned
//an adapter is required for particular functions
client, adapter, err := auth.AzureGraphClient()
if err != nil {
    t.Errorf("unable to authenticate to azure ad %v", err)
}

allUsers, err := GetAllUsers(client, adapter)
if err != nil {
    t.Errorf("unable to grab all users with error: %v", err)
}
t.Logf("Found %v users", len(allUsers))


```
