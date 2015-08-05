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

func (c Country) Is(id CountryID) bool {
	return c.ID == id
}

func (country Country) String() string {
	if country.ID == CountryUnknown {
		return "unknown"
	}
	return country.Title
}

type User struct {
	ID        int     `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Sex       Sex     `json:"sex"`
	Country   Country `json:"country"`
}
