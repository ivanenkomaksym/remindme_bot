package repositories

import "github.com/ivanenkomaksym/remindme_bot/domain/entities"

// PremiumUsageRepository defines the interface for premium usage persistence
type PremiumUsageRepository interface {
	// GetUserUsage retrieves premium usage for a specific user
	GetUserUsage(userID int64) (*entities.PremiumUsage, error)

	// CreateUserUsage creates a new premium usage record for a user
	CreateUserUsage(usage *entities.PremiumUsage) error

	// UpdateUserUsage updates an existing premium usage record
	UpdateUserUsage(usage *entities.PremiumUsage) error

	// GetOrCreateUserUsage gets existing usage or creates a new one
	GetOrCreateUserUsage(userID int64) (*entities.PremiumUsage, error)

	// GetAllUsage retrieves all premium usage records (for admin purposes)
	GetAllUsage() ([]entities.PremiumUsage, error)

	// DeleteUserUsage deletes premium usage record for a user
	DeleteUserUsage(userID int64) error

	// GetUsageByPremiumStatus gets usage records filtered by premium status
	GetUsageByPremiumStatus(status entities.PremiumStatus) ([]entities.PremiumUsage, error)
}
