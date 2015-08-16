package vk

import "fmt"
import "bytes"
import "text/template"

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
	GroupID  int    `url:"group_id,omitempty"`
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

// batch get
func (g Groups) GetBatch(getFields GroupGetFields) ([]User, int, error) {
	js := `var group_id = {{.GroupID}};
	var count = 1000;
	var offset = {{.Offset}};
	var calls = 1;
	var response = API.groups.getMembers({"count": count, "offset": offset, "group_id": group_id, fields: "{{ .Fields }}"});
	var members_count = response.count;
	var members = response.items;
	offset = offset + count;
	while ((offset < members_count) && (calls < 25)) {
		response = API.groups.getMembers({"count": count, "offset": offset, "group_id": group_id, fields: "{{ .Fields }}"});
		members = members + response.items;
		members_count = response.count;
		offset = offset + count;
		calls = calls + 1;
	}
	return {count: members_count, members: members};`

	// rendering js from template
	code := new(bytes.Buffer)
	t := template.Must(template.New("get20k.js").Parse(js))
	if err := t.Execute(code, getFields); err != nil {
		return nil, 0, err
	}

	// preparing request fields
	fields := struct {
		Code string `url:"code"`
	}{Code: code.String()}
	req := g.Request(methodExecute, fields)
	result := struct {
		Count   int    `json:"members_count"`
		Members []User `json:"members"`
	}{}
	return result.Members, result.Count, g.Decode(req, &result)
}
