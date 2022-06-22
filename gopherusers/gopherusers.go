package gopherusers

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/amarbut24/gopherland/gophererrors"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	models "github.com/microsoftgraph/msgraph-sdk-go/models"
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

type ConcurrentResult struct {
	Success    bool
	ObjectName string
	Error      error
}

// GetUserByID can used to return an Azure AD user via ObjectID
func GetUserByID(c *msgraphsdk.GraphServiceClient, uid string) (models.Userable, error) {
	user, err := c.UsersById(uid).Get()
	if err != nil {
		odataerr := gophererrors.HandleODataErr(err, "error finding user via objectid")
		return nil, odataerr
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
		odataerr := gophererrors.HandleODataErr(err, "error finding user via UserPrincipalName")
		return nil, odataerr
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
		odataerr := gophererrors.HandleODataErr(err, "error deleting user via objectid")
		return odataerr
	}
	return nil
}

// GetAllUsers returns all Azure AD users
func GetAllUsers(c *msgraphsdk.GraphServiceClient, adapter *msgraphsdk.GraphRequestAdapter) ([]models.Userable, error) {
	users, err := c.Users().Get()
	if err != nil {
		odataerr := gophererrors.HandleODataErr(err, "error retrieving all users")
		return nil, odataerr
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
		log.Printf("found user %v, skipping creation\n", user.UserPrincipalName)
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
		_, m := gophererrors.GetODataDetails(err)
		// this seems to be a bug with the SDK
		// created an issue here https://github.com/microsoftgraph/msgraph-sdk-go/issues/203
		if m == "Unable to read JSON request payload. Please ensure Content-Type header is set and payload is of valid JSON format." {
			return user.NewUser(c)
		} else {
			odataerr := gophererrors.HandleODataErr(err, "error creating new user")
			return nil, odataerr
		}
	}
	log.Printf("created new user %v\n", user.UserPrincipalName)
	return newUser, nil
}

// CNewUser adds channels to the NewUser
// which can be used to create many users at once via CNewUsers
func CNewUser(user GopherUser, c *msgraphsdk.GraphServiceClient, ch chan ConcurrentResult) {
	_, err := user.NewUser(c)
	if err != nil {
		ch <- ConcurrentResult{false, user.UserPrincipalName, err}
	}
	ch <- ConcurrentResult{true, user.UserPrincipalName, nil}
}

// CNewUsers build on CNewUser
// It takes a slice of users and creates seperate goroutines for each user
func CNewUsers(ch chan ConcurrentResult, users []GopherUser, client *msgraphsdk.GraphServiceClient) {
	for _, u := range users {
		go CNewUser(u, client, ch)
		//log.Printf("number of go routines %v", runtime.NumGoroutine())
	}

	f := []ConcurrentResult{}
	for i := 0; i < len(users); i++ {
		r := <-ch
		if !r.Success {
			log.Printf("failed to concurrently create %v", r.ObjectName)
			f = append(f, r)
		}
	}
	if len(f) > 0 {
		fmt.Println("outputting failed user creations")
	}
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
