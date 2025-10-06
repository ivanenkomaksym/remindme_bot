package mocks

import (
	"time"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	Users                   map[int64]*entities.User
	UserSelections          map[int64]*entities.UserSelection
	GetUsersFunc            func() ([]*entities.User, error)
	GetUserFunc             func(userID int64) (*entities.User, error)
	CreateUserFunc          func(userID int64, userName, firstName, lastName, language string) (*entities.User, error)
	GetOrCreateUserFunc     func(userID int64, userName, firstName, lastName, language string) (*entities.User, error)
	UpdateUserLanguageFunc  func(userID int64, language string) error
	UpdateLocationFunc      func(userID int64, location string) error
	UpdateUserInfoFunc      func(userID int64, userName, firstName, lastName string) error
	DeleteUserFunc          func(userID int64) error
	GetUserSelectionFunc    func(userID int64) (*entities.UserSelection, error)
	UpdateUserSelectionFunc func(userID int64, selection *entities.UserSelection) error
	ClearUserSelectionFunc  func(userID int64) error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Users:          make(map[int64]*entities.User),
		UserSelections: make(map[int64]*entities.UserSelection),
	}
}

func (m *MockUserRepository) GetUsers() ([]*entities.User, error) {
	var users = []*entities.User{}

	for _, user := range m.Users {
		users = append(users, user)
	}

	return users, nil
}

func (m *MockUserRepository) GetUser(userID int64) (*entities.User, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(userID)
	}
	user, exists := m.Users[userID]
	if !exists {
		return nil, nil
	}
	userCopy := *user
	return &userCopy, nil
}

func (m *MockUserRepository) CreateUser(userID int64, userName, firstName, lastName, language string) (*entities.User, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(userID, userName, firstName, lastName, language)
	}
	// Check if user already exists
	if _, exists := m.Users[userID]; exists {
		return nil, nil // User already exists
	}
	user := entities.NewUser(userID, userName, firstName, lastName, language)
	m.Users[userID] = user
	userCopy := *user
	return &userCopy, nil
}

func (m *MockUserRepository) GetOrCreateUser(userID int64, userName, firstName, lastName, language string) (*entities.User, error) {
	if m.GetOrCreateUserFunc != nil {
		return m.GetOrCreateUserFunc(userID, userName, firstName, lastName, language)
	}
	user := entities.NewUser(userID, userName, firstName, lastName, language)
	m.Users[userID] = user
	userCopy := *user
	return &userCopy, nil
}

func (m *MockUserRepository) UpdateUserLanguage(userID int64, language string) error {
	if m.UpdateUserLanguageFunc != nil {
		return m.UpdateUserLanguageFunc(userID, language)
	}
	if user, exists := m.Users[userID]; exists {
		user.UpdateLanguage(language)
	}
	return nil
}

func (m *MockUserRepository) UpdateLocation(userID int64, location string) error {
	if m.UpdateLocationFunc != nil {
		return m.UpdateLocationFunc(userID, location)
	}
	loc, _ := time.LoadLocation(location)
	if user, exists := m.Users[userID]; exists {
		user.SetLocation(loc)
	}
	return nil
}

func (m *MockUserRepository) UpdateUserInfo(userID int64, userName, firstName, lastName string) error {
	if m.UpdateUserInfoFunc != nil {
		return m.UpdateUserInfoFunc(userID, userName, firstName, lastName)
	}
	if user, exists := m.Users[userID]; exists {
		user.UpdateInfo(userName, firstName, lastName)
	}
	return nil
}

func (m *MockUserRepository) DeleteUser(userID int64) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(userID)
	}
	delete(m.Users, userID)
	return nil
}

func (m *MockUserRepository) GetUserSelection(userID int64) (*entities.UserSelection, error) {
	if m.GetUserSelectionFunc != nil {
		return m.GetUserSelectionFunc(userID)
	}
	selection, exists := m.UserSelections[userID]
	if !exists {
		return nil, nil
	}
	selectionCopy := *selection
	return &selectionCopy, nil
}

func (m *MockUserRepository) SetUserSelection(userID int64, selection *entities.UserSelection) error {
	selectionCopy := *selection
	m.UserSelections[userID] = &selectionCopy
	return nil
}

func (m *MockUserRepository) UpdateUserSelection(userID int64, selection *entities.UserSelection) error {
	if m.UpdateUserSelectionFunc != nil {
		return m.UpdateUserSelectionFunc(userID, selection)
	}
	selectionCopy := *selection
	m.UserSelections[userID] = &selectionCopy
	return nil
}

func (m *MockUserRepository) ClearUserSelection(userID int64) error {
	if m.ClearUserSelectionFunc != nil {
		return m.ClearUserSelectionFunc(userID)
	}
	if selection, exists := m.UserSelections[userID]; exists {
		selection.Clear()
	}
	return nil
}

func (m *MockUserRepository) GetUserWithSelection(userID int64) (*entities.User, *entities.UserSelection, error) {
	user, err := m.GetUser(userID)
	if err != nil {
		return nil, nil, err
	}
	selection, err := m.GetUserSelection(userID)
	if err != nil {
		return nil, nil, err
	}
	return user, selection, nil
}

func (m *MockUserRepository) GetOrCreateUserWithSelection(userID int64, userName, firstName, lastName, language string) (*entities.User, *entities.UserSelection, error) {
	user, err := m.GetOrCreateUser(userID, userName, firstName, lastName, language)
	if err != nil {
		return nil, nil, err
	}
	selection := entities.NewUserSelection()
	m.UserSelections[userID] = selection
	selectionCopy := *selection
	return user, &selectionCopy, nil
}
