package repositories

import (
	"github.com/ivanenkomaksym/remindme_bot/models"
)

type UserRepository interface {
	// User management
	GetUser(userID int64) *models.User
	CreateOrUpdateUser(userID int64, userName, firstName, lastName, language string) *models.User
	UpdateUserLanguage(userID int64, language string) bool

	// User state management
	GetUserSelection(userID int64) *models.UserSelection
	SetUserSelection(userID int64, state *models.UserSelection)
	UpdateUserSelection(userID int64, state *models.UserSelection)
	ClearUserSelection(userID int64)

	// Combined operations
	GetUserWithSelection(userID int64) (*models.User, *models.UserSelection)
	CreateOrUpdateUserWithSelection(userID int64, userName, firstName, lastName, language string) (*models.User, *models.UserSelection)
}
