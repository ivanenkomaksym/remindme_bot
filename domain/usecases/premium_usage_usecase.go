package usecases

import (
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/repositories"
)

// PremiumUsageUseCase defines the interface for premium usage business logic
type PremiumUsageUseCase interface {
	GetUserUsage(userID int64) (*entities.PremiumUsage, error)
	GetOrCreateUserUsage(userID int64) (*entities.PremiumUsage, error)
	UpdateUserUsage(usage *entities.PremiumUsage) error
	UpgradeUser(userID int64, status entities.PremiumStatus) (*entities.PremiumUsage, error)
	ResetUserUsage(userID int64) (*entities.PremiumUsage, error)
}

type premiumUsageUseCase struct {
	premiumUsageRepo repositories.PremiumUsageRepository
}

// NewPremiumUsageUseCase creates a new premium usage use case
func NewPremiumUsageUseCase(premiumUsageRepo repositories.PremiumUsageRepository) PremiumUsageUseCase {
	return &premiumUsageUseCase{
		premiumUsageRepo: premiumUsageRepo,
	}
}

func (p *premiumUsageUseCase) GetUserUsage(userID int64) (*entities.PremiumUsage, error) {
	return p.premiumUsageRepo.GetUserUsage(userID)
}

func (p *premiumUsageUseCase) GetOrCreateUserUsage(userID int64) (*entities.PremiumUsage, error) {
	return p.premiumUsageRepo.GetOrCreateUserUsage(userID)
}

func (p *premiumUsageUseCase) UpdateUserUsage(usage *entities.PremiumUsage) error {
	return p.premiumUsageRepo.UpdateUserUsage(usage)
}

func (p *premiumUsageUseCase) UpgradeUser(userID int64, status entities.PremiumStatus) (*entities.PremiumUsage, error) {
	usage, err := p.premiumUsageRepo.GetOrCreateUserUsage(userID)
	if err != nil {
		return nil, err
	}

	usage.SetPremiumStatus(status)

	err = p.premiumUsageRepo.UpdateUserUsage(usage)
	if err != nil {
		return nil, err
	}

	return usage, nil
}

func (p *premiumUsageUseCase) ResetUserUsage(userID int64) (*entities.PremiumUsage, error) {
	usage, err := p.premiumUsageRepo.GetOrCreateUserUsage(userID)
	if err != nil {
		return nil, err
	}

	usage.RequestsUsed = 0
	usage.LastReset = usage.LastReset.AddDate(0, 1, 0) // Add one month

	err = p.premiumUsageRepo.UpdateUserUsage(usage)
	if err != nil {
		return nil, err
	}

	return usage, nil
}
