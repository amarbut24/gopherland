package gopherusers

import (
	"fmt"

	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

/*
- Get All Azure Users
- Get User by Email
- Create User
- Delete User
- Disable User
*/

func GetUserByID(c *msgraphsdkgo.GraphServiceClient, uid string) (models.Userable, error) {
	user, err := c.UsersById(uid).Get(nil)
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
