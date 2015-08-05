package vk

type Groups struct {
	factory RequestFactory
	api     APIClient
}

const (
	methodGroupsGetMembers = "groups.getMembers"
)

type GroupSearchFields struct {
	ID int
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
	request := g.factory.Request(methodGroupsGetMembers, q)
	response := groupSearchResponse{}
	return response.Response, g.api.Do(request, &response)
}
