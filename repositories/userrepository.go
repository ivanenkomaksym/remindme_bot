package repositories

import "github.com/ivanenkomaksym/remindme_bot/domain/entities"

type UserRepository interface {
	// User management
	GetUser(userID int64) *entities.User
	CreateOrUpdateUser(userID int64, userName, firstName, lastName, language string) *entities.User
	UpdateUserLanguage(userID int64, language string) bool

	// User state management
	GetUserSelection(userID int64) *entities.UserSelection
	SetUserSelection(userID int64, state *entities.UserSelection)
	UpdateUserSelection(userID int64, state *entities.UserSelection)
	ClearUserSelection(userID int64)

	// Combined operations
	GetUserWithSelection(userID int64) (*entities.User, *entities.UserSelection)
	CreateOrUpdateUserWithSelection(userID int64, userName, firstName, lastName, language string) (*entities.User, *entities.UserSelection)
}
