package users

import (
	"fmt"

	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func GetUserByID(c *msgraphsdkgo.GraphServiceClient, uid string) (models.Userable, error) {
	user, err := c.UsersById(uid).Get(nil)
	if err != nil {
		return nil, fmt.Errorf("error finding user via objectid=%v: %v", uid, err)
	}
	return user, nil
}

func GetAllUsers(c *msgraphsdkgo.GraphServiceClient) (models.Userable, error) {

}
