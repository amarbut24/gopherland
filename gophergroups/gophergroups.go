package gophergroups

import (
	"fmt"

	"github.com/amarbut24/gopherland/gophererrors"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	groups "github.com/microsoftgraph/msgraph-sdk-go/groups"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// GopherGroup struct
type GopherGroup struct {
	DisplayName     string
	Description     string
	MailEnabled     bool
	MailNickname    string
	ObjectID        string
	SecurityEnabled bool
	GroupTypes      []string
}

// NewGroup allows you to create a new Azure AD group
func (group GopherGroup) NewGroup(c *msgraphsdk.GraphServiceClient) (models.Groupable, error) {

	foundGroup, _ := GetGroupByDisplayName(c, group.DisplayName)
	if foundGroup.DisplayName != "" {
		fmt.Printf("found user %v, skipping creation\n", group.DisplayName)
		return nil, nil
	}

	requestBody := models.NewGroup()
	requestBody.SetDescription(&group.Description)
	requestBody.SetDisplayName(&group.DisplayName)
	requestBody.SetGroupTypes(group.GroupTypes)
	requestBody.SetMailEnabled(&group.MailEnabled)
	requestBody.SetMailNickname(&group.MailNickname)
	requestBody.SetSecurityEnabled(&group.SecurityEnabled)

	newGroup, err := c.Groups().Post(requestBody)
	if err != nil {
		odataerr := gophererrors.HandleODataErr(err, "error creating new group")
		return nil, odataerr
	}
	return newGroup, nil
}

// GetGroupByID can used to return an Azure AD group via ObjectID
func GetGroupByID(c *msgraphsdk.GraphServiceClient, uid string) (models.Groupable, error) {
	group, err := c.GroupsById(uid).Get()
	if err != nil {
		odataerr := gophererrors.HandleODataErr(err, "error finding user via objectid")
		return nil, odataerr
	}
	return group, nil
}

// GetGroupByDisplayName can used to return an Azure AD group via DisplayName
func GetGroupByDisplayName(c *msgraphsdk.GraphServiceClient, displayname string) (GopherGroup, error) {
	filter := fmt.Sprintf("displayName eq '%s'", displayname)
	requestParameters := &groups.GroupsRequestBuilderGetQueryParameters{
		Filter: &filter,
	}

	options := &groups.GroupsRequestBuilderGetRequestConfiguration{
		QueryParameters: requestParameters,
	}

	group, err := c.Groups().GetWithRequestConfigurationAndResponseHandler(options, nil)
	if err != nil {
		odataerr := gophererrors.HandleODataErr(err, "error finding group via displayName")
		return GopherGroup{}, odataerr
	}
	if len(group.GetValue()) > 0 {
		if len(group.GetValue()) > 1 {
			return GopherGroup{}, fmt.Errorf("more than one value was returned when matching displayName %v, cleanup duplicate groups", displayname)
		}
		return ConvertToGopherGroup(group.GetValue()[0]), nil
	}
	return GopherGroup{}, nil

}

// DeleteGroupByID can used to delete an Azure AD group using Object ID
func DeleteUserByID(c *msgraphsdk.GraphServiceClient, uid string) error {
	err := c.GroupsById(uid).Delete()
	if err != nil {
		odataerr := gophererrors.HandleODataErr(err, "error deleting group via objectid")
		return odataerr
	}
	return nil
}

// GetAllGroups returns all Azure AD groups
func GetAllGroups(c *msgraphsdk.GraphServiceClient, adapter *msgraphsdk.GraphRequestAdapter) ([]models.Groupable, error) {
	groups, err := c.Groups().Get()
	if err != nil {
		odataerr := gophererrors.HandleODataErr(err, "error retrieving all groups")
		return nil, odataerr
	}

	pageIterator, err := msgraphcore.NewPageIterator(groups, adapter, models.CreateUserCollectionResponseFromDiscriminatorValue)
	if err != nil {
		return nil, fmt.Errorf("unable to create new pageIterator: %v", err)
	}

	var allGroups []models.Groupable
	err = pageIterator.Iterate(func(pageItem interface{}) bool {
		allGroups = append(allGroups, pageItem.(models.Groupable))
		// Return true to continue the iteration
		return true
	})
	if err != nil {
		return nil, fmt.Errorf("error occured when iterating over pages: %v", err)
	}

	return allGroups, nil
}

// ConvertToGoperGroup converts the models.Groupable object
// to a GopherGroup struct
func ConvertToGopherGroup(g models.Groupable) GopherGroup {
	return GopherGroup{
		DisplayName:     *g.GetDisplayName(),
		Description:     *g.GetDescription(),
		MailEnabled:     *g.GetMailEnabled(),
		MailNickname:    *g.GetMailNickname(),
		ObjectID:        *g.GetId(),
		SecurityEnabled: *g.GetSecurityEnabled(),
		GroupTypes:      g.GetGroupTypes(),
	}
}

//func (group GopherGroup) AddMembers(c *msgraphsdk.GraphServiceClient, memberids []string) {

// Per https://github.com/microsoftgraph/msgraph-sdk-go/issues/155#issuecomment-1142048347
// below is the correct way to implement a group add, however I've been unable to get it working
// namely the 'Post' method is not available
// waiting for a reply here https://github.com/microsoftgraph/msgraph-sdk-go/issues/197#issuecomment-1151696495

// reference := membersAdd.NewRef()
// for _, id := range memberids {
// 	reference.GetAdditionalData()["@odata.id"] = "https://graph.microsoft.com/v1.0/directoryObjects/" + id
// }

// graphClient := c
// r, err := graphClient.GroupsById(group.ObjectID).Members().Post(reference)
// fmt.Println(r, err)

// for now below is workaround
// url := fmt.Sprintf("https://graph.microsoft.com/v1.0/groups/%v/members/$ref", group.ObjectID)
// requestBuilder := groups.NewGroupsRequestBuilder(url, a)
// objectToAdd := models.NewGroup()

// for id := range memberids {
// 	objectToAdd.GetAdditionalData()["@odata.id"] = "https://graph.microsoft.com/v1.0/directoryObjects/" + string(id)
// }
// res, err := requestBuilder.Post()
//}
