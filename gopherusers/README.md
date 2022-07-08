# Azure AD Users Package

# Create new user
```
u1 := gopherusers.GopherUser{
		FirstName:                     "John",
		LastName:                      "Doe",
		MailNickname:                  "jdoe",
		UserPrincipalName:             "jdoe@contoso.com",
		ForceChangePasswordNextSignIn: true,
		DisplayName:                   "John Doe",
		AccountEnabled:                false,
}


newUser, err := u1.NewUser(client)
if err != nil {
	log.Print(err)
} else {
	log.Printf("created user %s\n", newUser.UserPrincipalName)
}

```
# Get user by Object ID

```
byID, err := gopherusers.GetUserByID(client, newUser.ObjectID)
if err != nil {
	log.Fatalf("unable to locate user: %v", err)
}
log.Println("found user based on Id =", byID.ObjectID)

```

# Get user by UserPrincipalName

```go
byUPN, err := gopherusers.GetUserByUPN(client, newUser.UserPrincipalName)
if err != nil {
	log.Fatalf("unable to locate user: %v", err)
}
log.Println("found user based on UPN =", byUPN.UserPrincipalName)

```

# Get all users will return []models.Userable

```go
allUsers, err := gopherusers.GetAllUsers(client, adapter)
if err != nil {
	log.Fatalf("unable to grab all users with error: %v", err)
}
fmt.Printf("Found %v users", len(allUsers))

```

# Delete user by object id

```go

err = gopherusers.DeleteUserByID(client, byUPN.ObjectID)
if err != nil {
	log.Fatalf("unable to locate user: %v", err)
} else {
	log.Printf("user %v was deleted\n", byUPN.DisplayName)
}

```

# Example - Create, query, enable, and delet a new user

```go
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
		FirstName:                     "John",
		LastName:                      "Doe",
		MailNickname:                  "jdoe",
		UserPrincipalName:             "jdoe@contoso.com",
		ForceChangePasswordNextSignIn: true,
		DisplayName:                   "John Doe",
		AccountEnabled:                false,
    }

	// Create new user
	newUser, err := u1.NewUser(client)
	if err != nil {
		log.Print(err)
	} else {
		log.Printf("created user %s\n", newUser.UserPrincipalName)
	}

	// Find the user by objectID
	byId, err := gopherusers.GetUserByID(client, newUser.ObjectID)
	if err != nil {
		log.Fatalf("unable to locate user: %v", err)
	}
	log.Println("found user based on Id =", byId.ObjectID)

	// Find the user by UserPrincipalName
	byUPN, err := gopherusers.GetUserByUPN(client, newUser.UserPrincipalName)
	if err != nil {
		log.Fatalf("unable to locate user: %v", err)
	}
	log.Println("found user based on UPN =", byUPN.UserPrincipalName)

	// Now lets delete the account
	err = gopherusers.DeleteUserByID(client, byUPN.ObjectID)
	if err != nil {
		log.Fatalf("unable to locate user: %v", err)
	} else {
		log.Printf("user %v was deleted\n", byUPN.DisplayName)
	}
}

```