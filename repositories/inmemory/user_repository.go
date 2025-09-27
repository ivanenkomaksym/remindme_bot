package inmemory

import (
	"sync"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/repositories"
)

type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[int64]*entities.User
}

func NewInMemoryUserRepository() repositories.UserRepository {
	return &InMemoryUserRepository{
		users: make(map[int64]*entities.User),
	}
}

func (r *InMemoryUserRepository) GetUsers() ([]*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var users = []*entities.User{}

	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
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

func (r *InMemoryUserRepository) GetOrCreateUser(userID int64, userName, firstName, lastName, language string) (*entities.User, error) {
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

func (r *InMemoryUserRepository) UpdateLocation(userID int64, timezone string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[userID]
	if !exists {
		return nil
	}

	location, err := time.LoadLocation(timezone)
	if err != nil {
		return nil
	}
	user.SetLocation(location)
	user.UpdatedAt = time.Now()
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
