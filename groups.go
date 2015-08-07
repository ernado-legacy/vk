package vk

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

type Group struct {
	ID       int    `json:"id"`
	Name     string `json:"name`
	IsClosed Bool   `json:"is_closed"`
	IsAdmin  Bool   `json:"is_admin"`
	IsMember Bool   `json:"is_member"`
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
	Offset   int  `url:"offset,omitempty"`
	Count    int  `url:"count,omitempty"`
	UserID   int  `url:"user_id,omitempty"`
	Extended Bool `url:"extended,omitempty"`
}

type GroupGetResult struct {
	Count int     `json:"count"`
	Items []Group `json:"items"`
}

func (g Groups) Get(fields GroupGetFields) (result GroupGetResult, err error) {
	return result, g.Decode(g.Request(methodGroupsGet, fields), &result)
}
