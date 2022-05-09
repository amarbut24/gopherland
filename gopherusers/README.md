# users

# Create new user
```
u1 := gopherusers.GopherUser{
		FirstName:                     "Anthony",
		LastName:                      "Marbut",
		MailNickname:                  "amarbut",
		UserPrincipalName:             "amarbut@gopherland.onmicrosoft.com",
		ForceChangePasswordNextSignIn: true,
		DisplayName:                   "Anthony Marbut",
		AccountEnabled:                false,
}

u, err := u1.NewUser(client)
if err != nil {
    log.Print(err)
} else {
    fmt.Printf("Created user %s", *u.GetUserPrincipalName())
}

```
# Get user by Object ID

```
u, err := gopherusers.GetUserByID(client, "8e086b41-7bc0-4b4a-9c3e-0a7ff59d710b")
if err != nil {
    log.Fatalf("unable to locate user: %v", err)
}
fmt.Println(*u.GetUserPrincipalName())

```
