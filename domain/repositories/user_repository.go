package repositories

import "github.com/ivanenkomaksym/remindme_bot/domain/entities"

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// User management
	GetUsers() ([]*entities.User, error)
	GetUser(userID int64) (*entities.User, error)
	CreateOrUpdateUser(userID int64, userName, firstName, lastName, language string) (*entities.User, error)
	UpdateUserLanguage(userID int64, language string) error
	UpdateUserInfo(userID int64, userName, firstName, lastName string) error

	// User selection management
	GetUserSelection(userID int64) (*entities.UserSelection, error)
	SetUserSelection(userID int64, selection *entities.UserSelection) error
	UpdateUserSelection(userID int64, selection *entities.UserSelection) error
	ClearUserSelection(userID int64) error

	// Combined operations
	GetUserWithSelection(userID int64) (*entities.User, *entities.UserSelection, error)
	CreateOrUpdateUserWithSelection(userID int64, userName, firstName, lastName, language string) (*entities.User, *entities.UserSelection, error)
}
