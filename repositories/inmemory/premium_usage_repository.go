package inmemory

import (
	"fmt"
	"sync"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/repositories"
)

// InMemoryPremiumUsageRepository provides in-memory storage for premium usage
type InMemoryPremiumUsageRepository struct {
	usage map[int64]*entities.PremiumUsage
	mutex sync.RWMutex
}

// NewInMemoryPremiumUsageRepository creates a new in-memory premium usage repository
func NewInMemoryPremiumUsageRepository() repositories.PremiumUsageRepository {
	return &InMemoryPremiumUsageRepository{
		usage: make(map[int64]*entities.PremiumUsage),
	}
}

// GetUserUsage retrieves premium usage for a specific user
func (r *InMemoryPremiumUsageRepository) GetUserUsage(userID int64) (*entities.PremiumUsage, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	usage, exists := r.usage[userID]
	if !exists {
		return nil, fmt.Errorf("premium usage not found for user %d", userID)
	}

	// Return a copy to prevent external modification
	usageCopy := *usage
	return &usageCopy, nil
}

// CreateUserUsage creates a new premium usage record for a user
func (r *InMemoryPremiumUsageRepository) CreateUserUsage(usage *entities.PremiumUsage) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.usage[usage.UserID]; exists {
		return fmt.Errorf("premium usage already exists for user %d", usage.UserID)
	}

	// Store a copy to prevent external modification
	usageCopy := *usage
	r.usage[usage.UserID] = &usageCopy
	return nil
}

// UpdateUserUsage updates an existing premium usage record
func (r *InMemoryPremiumUsageRepository) UpdateUserUsage(usage *entities.PremiumUsage) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.usage[usage.UserID]; !exists {
		return fmt.Errorf("premium usage not found for user %d", usage.UserID)
	}

	// Store a copy to prevent external modification
	usageCopy := *usage
	r.usage[usage.UserID] = &usageCopy
	return nil
}

// GetOrCreateUserUsage gets existing usage or creates a new one
func (r *InMemoryPremiumUsageRepository) GetOrCreateUserUsage(userID int64) (*entities.PremiumUsage, error) {
	// First try to get existing usage
	if usage, err := r.GetUserUsage(userID); err == nil {
		return usage, nil
	}

	// Create new usage if not found
	newUsage := entities.NewPremiumUsage(userID)
	if err := r.CreateUserUsage(newUsage); err != nil {
		return nil, fmt.Errorf("failed to create premium usage for user %d: %w", userID, err)
	}

	return r.GetUserUsage(userID)
}

// GetAllUsage retrieves all premium usage records
func (r *InMemoryPremiumUsageRepository) GetAllUsage() ([]entities.PremiumUsage, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	usages := make([]entities.PremiumUsage, 0, len(r.usage))
	for _, usage := range r.usage {
		usages = append(usages, *usage)
	}

	return usages, nil
}

// DeleteUserUsage deletes premium usage record for a user
func (r *InMemoryPremiumUsageRepository) DeleteUserUsage(userID int64) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.usage[userID]; !exists {
		return fmt.Errorf("premium usage not found for user %d", userID)
	}

	delete(r.usage, userID)
	return nil
}

// GetUsageByPremiumStatus gets usage records filtered by premium status
func (r *InMemoryPremiumUsageRepository) GetUsageByPremiumStatus(status entities.PremiumStatus) ([]entities.PremiumUsage, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var filtered []entities.PremiumUsage
	for _, usage := range r.usage {
		if usage.PremiumStatus == status {
			filtered = append(filtered, *usage)
		}
	}

	return filtered, nil
}
