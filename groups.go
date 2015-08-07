package vk

import (
	"errors"
	"io"
)

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

// Bool is special format for vk bool values
// that are represented as integers - 1,0
type Bool bool

const (
	byteOne  = 49
	byteZero = 48
)

func (v Bool) MarshalJSON() ([]byte, error) {
	if v {
		return []byte{byteOne}, nil
	}
	return []byte{byteZero}, nil
}

func (v *Bool) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		return nil
	}
	if len(data) != 1 {
		return io.ErrUnexpectedEOF
	}
	if data[0] == byteOne {
		*v = true
	} else if data[0] == byteZero {
		*v = false
	} else {
		return errors.New("bool value overflow")
	}
	return nil
}

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

type OffsetCount struct {
	Offset int `structs:"offset"`
	Count  int `structs:"count"`
}

type GroupGetFields struct {
	Offset   int  `structs:"offset"`
	Count    int  `structs:"count"`
	UserID   int  `structs:"user_id,omitempty"`
	Extended Bool `structs:"extended,omitempty"`
}

type GroupGetResult struct {
	Count int `json:"count"`
	Items []Group `json:"items"`
}

func (g Groups) Get(fields GroupGetFields) (result GroupGetResult, err error) {
	return result, g.Decode(g.Request(methodGroupsGet, fields), &result)
}
