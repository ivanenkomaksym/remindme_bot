package repositories

import (
	"github.com/ivanenkomaksym/remindme_bot/internal/models"
	"github.com/ivanenkomaksym/remindme_bot/internal/types"
)

type UserRepository interface {
	// User management
	GetUser(userID int64) *models.User
	CreateOrUpdateUser(userID int64, userName, firstName, lastName, language string) *models.User
	UpdateUserLanguage(userID int64, language string) bool

	// User state management
	GetUserState(userID int64) *types.UserSelectionState
	SetUserState(userID int64, state *types.UserSelectionState)
	UpdateUserState(userID int64, state *types.UserSelectionState)
	ClearUserState(userID int64)

	// Combined operations
	GetUserWithState(userID int64) (*models.User, *types.UserSelectionState)
	CreateOrUpdateUserWithState(userID int64, userName, firstName, lastName, language string) (*models.User, *types.UserSelectionState)
}
