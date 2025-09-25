package entities

import "time"

// User represents a user in the system
type User struct {
	ID        int64     `json:"id" bson:"id"`
	UserName  string    `json:"userName" bson:"userName"`
	FirstName string    `json:"firstName" bson:"firstName"`
	LastName  string    `json:"lastName" bson:"lastName"`
	Language  string    `json:"language" bson:"language"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
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
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdateLanguage updates the user's language
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
