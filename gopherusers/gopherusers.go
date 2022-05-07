package gopherusers

import (
	"fmt"
	"math/rand"
	"time"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	models "github.com/microsoftgraph/msgraph-sdk-go/models"
	msgraph_errors "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

/*
- Get All Azure Users
- Get User by Email
- Create User
- Delete User
- Disable User
*/

type GopherUser struct {
	AccountEnabled                bool
	FirstName                     string
	ForceChangePasswordNextSignIn bool
	LastName                      string
	DisplayName                   string
	UserPrincipalName             string
	MailNickname                  string
}

func GetUserByID(c *msgraphsdk.GraphServiceClient, uid string) (models.Userable, error) {
	user, err := c.UsersById(uid).Get()
	if err != nil {
		return nil, fmt.Errorf("error finding user via objectid=%v: %v", uid, err)
	}
	return user, nil
}

// func GetAllUsers(c *msgraphsdkgo.GraphServiceClient) (models.Userable, error) {
// 	user, err := c.Users().Get(nil)
// 	if err != nil {
// 		fmt.Printf("Error getting users: %v\n", err)
// 		return err
// 	}
// 	user.GetNextLink()
// }

func (user GopherUser) NewUser(c *msgraphsdk.GraphServiceClient) (models.Userable, error) {

	//password := NewRandomPassword(18)
	password := "pass1234"
	requestBody := models.NewUser()
	passProfile := models.NewPasswordProfile()
	passProfile.SetForceChangePasswordNextSignIn(&user.ForceChangePasswordNextSignIn)
	passProfile.SetPassword(&password)
	requestBody.SetPasswordProfile(passProfile)
	requestBody.SetAccountEnabled(&user.AccountEnabled)
	requestBody.SetDisplayName(&user.DisplayName)
	requestBody.SetUserPrincipalName(&user.UserPrincipalName)
	requestBody.SetMailNickname(&user.MailNickname)

	results, err := c.Users().Post(requestBody)
	if err != nil {
		oderr := err.(*msgraph_errors.ODataError)
		return nil, fmt.Errorf("error creating new user\nCode=%v\nmessage=%v", *oderr.GetError().GetCode(), *oderr.GetError().GetMessage())
	}

	fmt.Println("Created new user:", results)
	return results, nil
}

//NewRandomPassword is used to generate a temporary password
func NewRandomPassword(length int) string {
	var password string

	for i := 0; i != length; i++ {
		password += newRandomASCII()
	}
	return password
}

func newRandomASCII() string {
	rand.Seed(time.Now().UTC().UnixNano())
	i := rand.Intn(126-33) + 33
	return string(i)
}
