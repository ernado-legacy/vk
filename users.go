package vk

// Sex of a user
type Sex int

const (
	SexUnknown Sex = 0
	Female     Sex = 1
	Male       Sex = 2
)

func (sex Sex) String() string {
	if sex == Male {
		return "male"
	}
	if sex == Female {
		return "female"
	}
	return "unknown"
}

type CountryID int

const (
	CountryUnknown CountryID = 0
	Russia         CountryID = 1
)

type Country struct {
	ID    CountryID `json:"id"`
	Title string    `json:"title"`
}

type City struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func (c Country) Is(id CountryID) bool {
	return c.ID == id
}

func (country Country) String() string {
	if country.ID == CountryUnknown {
		return "unknown"
	}
	return country.Title
}

//go:generate stringer -type=Relation
type Relation int

const (
	RelationUnknown      Relation = 0
	RelationSingle       Relation = 1
	RelationHasFriend    Relation = 2
	RelationEngaged      Relation = 3
	RelationMarried      Relation = 4
	RelationComplicated  Relation = 5
	RelationActiveSearch Relation = 6
	RelationInLove       Relation = 7
)

type User struct {
	ID        int     `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Sex       Sex     `json:"sex"`
	Country   Country `json:"country"`
	City      City    `json:"city"`
	Hidden    Bool    `json:"hidden"`
	Birthday  string  `json:"bdate"`
	PhotoMax  string  `json:"photo_max"`
	Status    string  `json:"status"`
	LastSeen  struct {
		Time     int64 `json:"time"`
		Platform int   `json:"platform"`
	} `json:"last_seen"`
	Books string `json:"books"`
	About string `json:"about"`
}

// UserFields all fields that are in User struct
const UserFields = "id,first_name,last_name,sex,country,city,photo_max,last_seen"
