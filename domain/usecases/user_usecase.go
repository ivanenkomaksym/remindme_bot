package usecases

import (
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/errors"
	"github.com/ivanenkomaksym/remindme_bot/domain/repositories"
)

// UserUseCase defines the interface for user business logic
type UserUseCase interface {
	GetUser(userID int64) (*entities.User, error)
	CreateOrUpdateUser(userID int64, userName, firstName, lastName, language string) (*entities.User, error)
	UpdateUserLanguage(userID int64, language string) error
	GetUserSelection(userID int64) (*entities.UserSelection, error)
	UpdateUserSelection(userID int64, selection *entities.UserSelection) error
	ClearUserSelection(userID int64) error
}

type userUseCase struct {
	userRepo repositories.UserRepository
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(userRepo repositories.UserRepository) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (u *userUseCase) GetUser(userID int64) (*entities.User, error) {
	user, err := u.userRepo.GetUser(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}
	return user, nil
}

func (u *userUseCase) CreateOrUpdateUser(userID int64, userName, firstName, lastName, language string) (*entities.User, error) {
	if userID <= 0 {
		return nil, errors.NewDomainError("INVALID_USER_ID", "User ID must be positive", nil)
	}
	if userName == "" && firstName == "" && lastName == "" {
		return nil, errors.NewDomainError("INVALID_USER_DATA", "At least one name field must be provided", nil)
	}

	user, err := u.userRepo.CreateOrUpdateUser(userID, userName, firstName, lastName, language)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUseCase) UpdateUserLanguage(userID int64, language string) error {
	if userID <= 0 {
		return errors.NewDomainError("INVALID_USER_ID", "User ID must be positive", nil)
	}
	if language == "" {
		return errors.NewDomainError("INVALID_LANGUAGE", "Language cannot be empty", nil)
	}

	return u.userRepo.UpdateUserLanguage(userID, language)
}

func (u *userUseCase) GetUserSelection(userID int64) (*entities.UserSelection, error) {
	if userID <= 0 {
		return nil, errors.NewDomainError("INVALID_USER_ID", "User ID must be positive", nil)
	}

	selection, err := u.userRepo.GetUserSelection(userID)
	if err != nil {
		return nil, err
	}
	if selection == nil {
		// Return empty selection if none exists
		return entities.NewUserSelection(), nil
	}
	return selection, nil
}

func (u *userUseCase) UpdateUserSelection(userID int64, selection *entities.UserSelection) error {
	if userID <= 0 {
		return errors.NewDomainError("INVALID_USER_ID", "User ID must be positive", nil)
	}
	if selection == nil {
		return errors.NewDomainError("INVALID_SELECTION", "User selection cannot be nil", nil)
	}

	return u.userRepo.UpdateUserSelection(userID, selection)
}

func (u *userUseCase) ClearUserSelection(userID int64) error {
	if userID <= 0 {
		return errors.NewDomainError("INVALID_USER_ID", "User ID must be positive", nil)
	}

	return u.userRepo.ClearUserSelection(userID)
}
