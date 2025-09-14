package repositories

import (
	"sync"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

type InMemoryUserRepository struct {
	mu             sync.RWMutex
	users          map[int64]*entities.User
	userSelections map[int64]*entities.UserSelection
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:          make(map[int64]*entities.User),
		userSelections: make(map[int64]*entities.UserSelection),
	}
}

// User management methods
func (r *InMemoryUserRepository) GetUser(userID int64) *entities.User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[userID]
	if !exists {
		return nil
	}

	// Return a copy to prevent external modifications
	userCopy := *user
	return &userCopy
}

func (r *InMemoryUserRepository) CreateOrUpdateUser(userID int64, userName, firstName, lastName, language string) *entities.User {
	r.mu.Lock()
	defer r.mu.Unlock()

	user := &entities.User{
		ID:        userID,
		UserName:  userName,
		FirstName: firstName,
		LastName:  lastName,
		Language:  language,
	}

	r.users[userID] = user

	// Return a copy to prevent external modifications
	userCopy := *user
	return &userCopy
}

func (r *InMemoryUserRepository) UpdateUserLanguage(userID int64, language string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[userID]
	if !exists {
		return false
	}

	user.Language = language
	return true
}

// User state management methods
func (r *InMemoryUserRepository) GetUserSelection(userID int64) *entities.UserSelection {
	r.mu.RLock()
	defer r.mu.RUnlock()

	state, exists := r.userSelections[userID]
	if !exists {
		return nil
	}

	// Return a copy to prevent external modifications
	stateCopy := *state
	return &stateCopy
}

func (r *InMemoryUserRepository) SetUserSelection(userID int64, state *entities.UserSelection) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Store a copy to prevent external modifications
	stateCopy := *state
	r.userSelections[userID] = &stateCopy
}

func (r *InMemoryUserRepository) UpdateUserSelection(userID int64, state *entities.UserSelection) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Store a copy to prevent external modifications
	stateCopy := *state
	r.userSelections[userID] = &stateCopy
}

func (r *InMemoryUserRepository) ClearUserSelection(userID int64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if state, exists := r.userSelections[userID]; exists {
		// Reset all fields to default values
		*state = entities.UserSelection{}
		state.WeekOptions = [7]bool{false, false, false, false, false, false, false}
	}
}

// Combined operations
func (r *InMemoryUserRepository) GetUserWithSelection(userID int64) (*entities.User, *entities.UserSelection) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, userExists := r.users[userID]
	state, stateExists := r.userSelections[userID]

	var userCopy *entities.User
	var stateCopy *entities.UserSelection

	if userExists {
		userCopyVal := *user
		userCopy = &userCopyVal
	}

	if stateExists {
		stateCopyVal := *state
		stateCopy = &stateCopyVal
	}

	return userCopy, stateCopy
}

func (r *InMemoryUserRepository) CreateOrUpdateUserWithSelection(userID int64, userName, firstName, lastName, language string) (*entities.User, *entities.UserSelection) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Create or update user
	user := &entities.User{
		ID:        userID,
		UserName:  userName,
		FirstName: firstName,
		LastName:  lastName,
		Language:  language,
	}
	r.users[userID] = user

	// Get or create user state
	state, exists := r.userSelections[userID]
	if !exists {
		state = &entities.UserSelection{
			WeekOptions: [7]bool{false, false, false, false, false, false, false},
		}
		r.userSelections[userID] = state
	}

	// Return copies to prevent external modifications
	userCopy := *user
	stateCopy := *state
	return &userCopy, &stateCopy
}
