package gopherusers

import (
	"fmt"
	"math/rand"
	"time"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	models "github.com/microsoftgraph/msgraph-sdk-go/models"
	msgraph_errors "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	users "github.com/microsoftgraph/msgraph-sdk-go/users"
)

// GopherUser can be used to intilzae a new user
type GopherUser struct {
	AccountEnabled                bool
	FirstName                     string
	ForceChangePasswordNextSignIn bool
	LastName                      string
	DisplayName                   string
	UserPrincipalName             string
	MailNickname                  string
}

// GetUserByID can used to return an Azure AD user via ObjectID
func GetUserByID(c *msgraphsdk.GraphServiceClient, uid string) (models.Userable, error) {
	user, err := c.UsersById(uid).Get()
	if err != nil {
		oderr := err.(*msgraph_errors.ODataError).GetError()
		c := *oderr.GetCode()
		m := *oderr.GetMessage()
		return nil, fmt.Errorf("error finding user via objectid\nCode=%v\nmessage=%v", c, m)
	}
	return user, nil
}

// GetUserByUPN can used to return an Azure AD user via UPN
func GetUserByUPN(c *msgraphsdk.GraphServiceClient, upn string) (models.Userable, error) {
	filter := fmt.Sprintf("userPrincipalName eq '%s'", upn)
	requestParameters := &users.UsersRequestBuilderGetQueryParameters{
		Filter: &filter,
	}

	options := &users.UsersRequestBuilderGetRequestConfiguration{
		QueryParameters: requestParameters,
	}

	user, err := c.Users().GetWithRequestConfigurationAndResponseHandler(options, nil)
	if err != nil {
		oderr := err.(*msgraph_errors.ODataError).GetError()
		c := *oderr.GetCode()
		m := *oderr.GetMessage()
		return nil, fmt.Errorf("error finding user via UserPrincipalName=%v\nCode=%v\nmessage=%v", upn, c, m)
	}
	if len(user.GetValue()) > 0 {
		if len(user.GetValue()) > 1 {
			return nil, fmt.Errorf("more than one value was returned when matching userPrincipalName %v, this should not happen", upn)
		}

		return user.GetValue()[0], nil
	}
	return nil, nil

}

// DeleteUserByID can used to delete an Azure AD user using Object ID
func DeleteUserByID(c *msgraphsdk.GraphServiceClient, uid string) error {
	err := c.UsersById(uid).Delete()
	if err != nil {
		oderr := err.(*msgraph_errors.ODataError).GetError()
		c := *oderr.GetCode()
		m := *oderr.GetMessage()
		return fmt.Errorf("error finding user via objectid\nCode=%v\nmessage=%v", c, m)
	}
	return nil
}

// GetAllUsers returns all Azure AD users
func GetAllUsers(c *msgraphsdk.GraphServiceClient, adapter *msgraphsdk.GraphRequestAdapter) ([]models.Userable, error) {
	users, err := c.Users().Get()
	if err != nil {
		oderr := err.(*msgraph_errors.ODataError).GetError()
		c := *oderr.GetCode()
		m := *oderr.GetMessage()
		return nil, fmt.Errorf("error creating new user\nCode=%v\nmessage=%v", c, m)
	}

	pageIterator, err := msgraphcore.NewPageIterator(users, adapter, models.CreateUserCollectionResponseFromDiscriminatorValue)
	if err != nil {
		return nil, fmt.Errorf("unable to create new pageIterator: %v", err)
	}

	var allUsers []models.Userable
	err = pageIterator.Iterate(func(pageItem interface{}) bool {
		allUsers = append(allUsers, pageItem.(models.Userable))
		// Return true to continue the iteration
		return true
	})
	if err != nil {
		return nil, fmt.Errorf("error occured when iterating over pages: %v", err)
	}

	return allUsers, nil
}

// NewUser allows you to create a new Azure AD user
func (user GopherUser) NewUser(c *msgraphsdk.GraphServiceClient) (models.Userable, error) {

	foundUser, _ := GetUserByUPN(c, user.UserPrincipalName)
	if foundUser != nil {
		fmt.Printf("found user %v, skipping creation\n", user.UserPrincipalName)
		return nil, nil
	}

	password := newRandomPassword(18)
	requestBody := models.NewUser()
	passProfile := models.NewPasswordProfile()
	passProfile.SetForceChangePasswordNextSignIn(&user.ForceChangePasswordNextSignIn)
	passProfile.SetPassword(&password)
	requestBody.SetPasswordProfile(passProfile)
	requestBody.SetAccountEnabled(&user.AccountEnabled)
	requestBody.SetDisplayName(&user.DisplayName)
	requestBody.SetUserPrincipalName(&user.UserPrincipalName)
	requestBody.SetMailNickname(&user.MailNickname)

	newUser, err := c.Users().Post(requestBody)
	if err != nil {
		oderr := err.(*msgraph_errors.ODataError).GetError()
		c := *oderr.GetCode()
		m := *oderr.GetMessage()
		return nil, fmt.Errorf("error creating new user\nCode=%v\nmessage=%v", c, m)
	}
	return newUser, nil
}

func newRandomPassword(length int) string {
	var password string

	for i := 0; i != length; i++ {
		password += newRandomASCII()
	}
	return password
}

func newRandomASCII() string {
	rand.Seed(time.Now().UTC().UnixNano())
	i := 0
	for {
		i = rand.Intn(126-33) + 33
		// we can't use '<' , '>', '"' as characters in a password

		if i != 62 && i != 60 && i != 34 {
			break
		}
	}
	return fmt.Sprintf("%c", i)
}
