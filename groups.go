package vk
import "fmt"

type Resource struct {
	APIClient
	RequestFactory
}

func (r Resource) Decode(request Request, v interface{}) error {
	res, err := r.Do(request)
	if err != nil {
		return err
	}
	return res.To(v)
}

type Groups struct {
	Resource
}

const (
	methodGroupsGetMembers = "groups.getMembers"
	methodGroupsGet        = "groups.get"
)

//go:generate stringer -type=GroupType
type GroupType int

type GroupDeactivatedStatus string

const (
	GroupOpen    GroupType = 0
	GroupClosed  GroupType = 1
	GroupPrivate GroupType = 2

	GroupDeactivated GroupDeactivatedStatus = "deleted"
	GroupBanned      GroupDeactivatedStatus = "banned"
	GroupActive      GroupDeactivatedStatus = ""
)

//go:generate stringer -type=GroupAdminLevel
type GroupAdminLevel int

const (
	GroupModerator     GroupAdminLevel = 1
	GroupRedactor      GroupAdminLevel = 2
	GroupAdministrator GroupAdminLevel = 3
)

type Group struct {
	ID           int                    `json:"id"`
	Name         string                 `json:"name`
	Slug         string                 `json:"screen_name"`
	Deactivated  GroupDeactivatedStatus `json:"deactivated"`
	IsClosed     GroupType              `json:"is_closed"`
	IsAdmin      Bool                   `json:"is_admin"`
	IsMember     Bool                   `json:"is_member"`
	AdminLevel   GroupAdminLevel        `json:"admin_level"`
	Type         string                 `json:"type"`
	Photo50      string                 `json:"photo_50"`
	Photo100     string                 `json:"photo_100"`
	Photo200     string                 `json:"photo_200"`
	Description  string                 `json:"description"`
	MembersCount int                    `json:"members_count"`
	Status       string                 `json:"status"`
}

func (g Group) GetStatus() string {
	if g.Deactivated == GroupActive {
		return "active"
	}
	return g.Status
}

func (g Group) String() string {
	return fmt.Sprintf("G:%s %s [count=%d,status=%s]", g.Slug, g.Name, g.MembersCount, g.GetStatus())
}

type GroupSearchFields struct {
	ID int `structs:"id"`
}

type GroupSearchResult struct {
	Count int    `json:"count"`
	Items []User `json:"items"`
}

type groupSearchResponse struct {
	Error    `json:"error"`
	Response GroupSearchResult `json:"response`
}

func (g Groups) GetMembers(q GroupSearchFields) (result GroupSearchResult, err error) {
	request := g.Request(methodGroupsGetMembers, q)
	return result, g.Decode(request, &result)
}

type GroupGetFields struct {
	Offset   int    `url:"offset,omitempty"`
	Count    int    `url:"count,omitempty"`
	UserID   int    `url:"user_id,omitempty"`
	Extended Bool   `url:"extended,omitempty"`
	Fields   string `url:"fields,omitempty"`
}

type GroupGetResult struct {
	Count int     `json:"count"`
	Items []Group `json:"items"`
}

func (g Groups) GetForUser(id int) ([]Group, error) {
	result := &GroupGetResult{}
	request := g.Request(methodGroupsGet, GroupGetFields{UserID: id,
		Count:    1000,
		Extended: true,
		Fields:   "description,members_count",
	})
	return result.Items, g.Decode(request, &result)
}

func (g Groups) Get(fields GroupGetFields) (result GroupGetResult, err error) {
	return result, g.Decode(g.Request(methodGroupsGet, fields), &result)
}
