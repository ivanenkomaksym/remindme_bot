package repositories

import (
	"sync"

	"github.com/ivanenkomaksym/remindme_bot/internal/models"
	"github.com/ivanenkomaksym/remindme_bot/internal/types"
)

type InMemoryUserRepository struct {
	mu         sync.RWMutex
	users      map[int64]*models.User
	userStates map[int64]*types.UserSelectionState
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:      make(map[int64]*models.User),
		userStates: make(map[int64]*types.UserSelectionState),
	}
}

// User management methods
func (r *InMemoryUserRepository) GetUser(userID int64) *models.User {
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

func (r *InMemoryUserRepository) CreateOrUpdateUser(userID int64, userName, firstName, lastName, language string) *models.User {
	r.mu.Lock()
	defer r.mu.Unlock()

	user := &models.User{
		Id:        userID,
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
func (r *InMemoryUserRepository) GetUserState(userID int64) *types.UserSelectionState {
	r.mu.RLock()
	defer r.mu.RUnlock()

	state, exists := r.userStates[userID]
	if !exists {
		return nil
	}

	// Return a copy to prevent external modifications
	stateCopy := *state
	return &stateCopy
}

func (r *InMemoryUserRepository) SetUserState(userID int64, state *types.UserSelectionState) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Store a copy to prevent external modifications
	stateCopy := *state
	r.userStates[userID] = &stateCopy
}

func (r *InMemoryUserRepository) UpdateUserState(userID int64, state *types.UserSelectionState) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Store a copy to prevent external modifications
	stateCopy := *state
	r.userStates[userID] = &stateCopy
}

func (r *InMemoryUserRepository) ClearUserState(userID int64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if state, exists := r.userStates[userID]; exists {
		// Reset all fields to default values
		*state = types.UserSelectionState{}
		state.WeekOptions = [7]bool{false, false, false, false, false, false, false}
	}
}

// Combined operations
func (r *InMemoryUserRepository) GetUserWithState(userID int64) (*models.User, *types.UserSelectionState) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, userExists := r.users[userID]
	state, stateExists := r.userStates[userID]

	var userCopy *models.User
	var stateCopy *types.UserSelectionState

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

func (r *InMemoryUserRepository) CreateOrUpdateUserWithState(userID int64, userName, firstName, lastName, language string) (*models.User, *types.UserSelectionState) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Create or update user
	user := &models.User{
		Id:        userID,
		UserName:  userName,
		FirstName: firstName,
		LastName:  lastName,
		Language:  language,
	}
	r.users[userID] = user

	// Get or create user state
	state, exists := r.userStates[userID]
	if !exists {
		state = &types.UserSelectionState{
			WeekOptions: [7]bool{false, false, false, false, false, false, false},
		}
		r.userStates[userID] = state
	}

	// Return copies to prevent external modifications
	userCopy := *user
	stateCopy := *state
	return &userCopy, &stateCopy
}
