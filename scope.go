package vk

import "strings"
import "sort"

// Permission of application
type Permission string

const (
	PermOffline Permission = "offline"
	PermFriends Permission = "friends"
	PermPhotos  Permission = "photos"
	PermGroups  Permission = "groups"
)

func (p Permission) String() string {
	return string(p)
}

type Scope map[Permission]bool

func (s Scope) Has(p Permission) bool {
	if s == nil {
		return false
	}
	return s[p]
}

func (s Scope) Add(permissions ...Permission) {
	for _, v := range permissions {
		s[v] = true
	}
}

func (s Scope) Del(permissions ...Permission) {
	for _, v := range permissions {
		delete(s, v)
	}
}

func (s Scope) String() string {
	var permissions []string
	for k := range s {
		permissions = append(permissions, k.String())
	}
	sort.Strings(permissions)
	return strings.Join(permissions, ",")
}

func NewScope(permissions ...Permission) Scope {
	s := Scope{}
	s.Add(permissions...)
	return s
}
