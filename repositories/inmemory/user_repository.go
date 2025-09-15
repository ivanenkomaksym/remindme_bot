package inmemory

import (
	"sync"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/repositories"
)

type InMemoryUserRepository struct {
	mu             sync.RWMutex
	users          map[int64]*entities.User
	userSelections map[int64]*entities.UserSelection
}

func NewInMemoryUserRepository() repositories.UserRepository {
	return &InMemoryUserRepository{
		users:          make(map[int64]*entities.User),
		userSelections: make(map[int64]*entities.UserSelection),
	}
}

// User management methods
func (r *InMemoryUserRepository) GetUser(userID int64) (*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[userID]
	if !exists {
		return nil, nil
	}

	// Return a copy to prevent external modifications
	userCopy := *user
	return &userCopy, nil
}

func (r *InMemoryUserRepository) CreateOrUpdateUser(userID int64, userName, firstName, lastName, language string) (*entities.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user := entities.NewUser(userID, userName, firstName, lastName, language)
	r.users[userID] = user

	// Return a copy to prevent external modifications
	userCopy := *user
	return &userCopy, nil
}

func (r *InMemoryUserRepository) UpdateUserLanguage(userID int64, language string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[userID]
	if !exists {
		return nil // User doesn't exist, nothing to update
	}

	user.UpdateLanguage(language)
	return nil
}

func (r *InMemoryUserRepository) UpdateUserInfo(userID int64, userName, firstName, lastName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[userID]
	if !exists {
		return nil // User doesn't exist, nothing to update
	}

	user.UpdateInfo(userName, firstName, lastName)
	return nil
}

// User selection management methods
func (r *InMemoryUserRepository) GetUserSelection(userID int64) (*entities.UserSelection, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	selection, exists := r.userSelections[userID]
	if !exists {
		return nil, nil
	}

	// Return a copy to prevent external modifications
	selectionCopy := *selection
	return &selectionCopy, nil
}

func (r *InMemoryUserRepository) SetUserSelection(userID int64, selection *entities.UserSelection) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Store a copy to prevent external modifications
	selectionCopy := *selection
	r.userSelections[userID] = &selectionCopy
	return nil
}

func (r *InMemoryUserRepository) UpdateUserSelection(userID int64, selection *entities.UserSelection) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Store a copy to prevent external modifications
	selectionCopy := *selection
	r.userSelections[userID] = &selectionCopy
	return nil
}

func (r *InMemoryUserRepository) ClearUserSelection(userID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if selection, exists := r.userSelections[userID]; exists {
		selection.Clear()
	}
	return nil
}

// Combined operations
func (r *InMemoryUserRepository) GetUserWithSelection(userID int64) (*entities.User, *entities.UserSelection, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, userExists := r.users[userID]
	selection, selectionExists := r.userSelections[userID]

	var userCopy *entities.User
	var selectionCopy *entities.UserSelection

	if userExists {
		userCopyVal := *user
		userCopy = &userCopyVal
	}

	if selectionExists {
		selectionCopyVal := *selection
		selectionCopy = &selectionCopyVal
	}

	return userCopy, selectionCopy, nil
}

func (r *InMemoryUserRepository) CreateOrUpdateUserWithSelection(userID int64, userName, firstName, lastName, language string) (*entities.User, *entities.UserSelection, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Create or update user
	user := entities.NewUser(userID, userName, firstName, lastName, language)
	r.users[userID] = user

	// Get or create user selection
	selection, exists := r.userSelections[userID]
	if !exists {
		selection = entities.NewUserSelection()
		r.userSelections[userID] = selection
	}

	// Return copies to prevent external modifications
	userCopy := *user
	selectionCopy := *selection
	return &userCopy, &selectionCopy, nil
}
