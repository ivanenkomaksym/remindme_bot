package entities

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID           int64          `json:"id" bson:"id"`
	UserName     string         `json:"userName" bson:"userName"`
	FirstName    string         `json:"firstName" bson:"firstName"`
	LastName     string         `json:"lastName" bson:"lastName"`
	Language     string         `json:"language" bson:"language"`
	LocationName string         `json:"location" bson:"location"`
	Location     *time.Location `bson:"-"` // Ignore by mongo
	CreatedAt    time.Time      `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt" bson:"updatedAt"`
}

// NewUser creates a new user entity
func NewUser(id int64, userName, firstName, lastName, language string) *User {
	now := time.Now()
	return &User{
		ID:        id,
		UserName:  userName,
		FirstName: firstName,
		LastName:  lastName,
		Language:  language,
		Location:  time.Now().Location(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (u *User) UpdateLanguage(language string) {
	u.Language = language
	u.UpdatedAt = time.Now()
}

// UpdateInfo updates the user's basic information
func (u *User) UpdateInfo(userName, firstName, lastName string) {
	u.UserName = userName
	u.FirstName = firstName
	u.LastName = lastName
	u.UpdatedAt = time.Now()
}

func (u *User) GetLocation() *time.Location {
	// If the private field is nil, try to load it from the stored string.
	if u.Location == nil && u.LocationName != "" {
		loc, err := time.LoadLocation(u.LocationName)
		if err == nil {
			u.Location = loc
		}
	}
	return u.Location
}

func (u *User) SetLocation(loc *time.Location) {
	u.Location = loc
	if loc != nil {
		u.LocationName = loc.String()
	} else {
		u.LocationName = ""
	}
}
