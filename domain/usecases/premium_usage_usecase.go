package usecases

import (
	"fmt"
	"log"

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
	// NLP-specific methods
	ValidateCanMakeRequest(userID int64) error
	ConsumeRequest(userID int64) error
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

// ValidateCanMakeRequest validates if user can make an NLP request
func (p *premiumUsageUseCase) ValidateCanMakeRequest(userID int64) error {
	// Get or create user usage record
	usage, err := p.premiumUsageRepo.GetOrCreateUserUsage(userID)
	if err != nil {
		log.Printf("Failed to get NLP usage for user %d: %v", userID, err)
		return fmt.Errorf("failed to check usage limits")
	}

	// Check if user can make a request
	if usage.CanMakeRequest() {
		return nil
	}

	log.Printf("User %d exceeded NLP rate limit: %d/%d requests used", userID, usage.RequestsUsed, usage.RequestsLimit)

	remainingDays := p.getDaysUntilReset(usage)
	var errorMsg string

	switch usage.PremiumStatus {
	case entities.PremiumStatusFree:
		errorMsg = fmt.Sprintf("You've reached your monthly limit of %d AI text reminders. Upgrade to Premium for more requests or try again in %d days.",
			usage.RequestsLimit, remainingDays)
	case entities.PremiumStatusBasic:
		errorMsg = fmt.Sprintf("You've reached your monthly limit of %d AI text reminders. Upgrade to Pro for unlimited requests or try again in %d days.",
			usage.RequestsLimit, remainingDays)
	default:
		errorMsg = "Rate limit exceeded. Please try again later."
	}

	return fmt.Errorf("%s", errorMsg)
}

// ConsumeRequest consumes one request from the user's quota
func (p *premiumUsageUseCase) ConsumeRequest(userID int64) error {
	// Get or create user usage record
	usage, err := p.premiumUsageRepo.GetOrCreateUserUsage(userID)
	if err != nil {
		log.Printf("Failed to get NLP usage for user %d: %v", userID, err)
		return fmt.Errorf("failed to update usage")
	}

	// Increment usage count
	usage.RequestsUsed++

	// Update in repository
	err = p.premiumUsageRepo.UpdateUserUsage(usage)
	if err != nil {
		log.Printf("Failed to update NLP usage for user %d: %v", userID, err)
		return fmt.Errorf("failed to update usage")
	}

	return nil
}

// getDaysUntilReset calculates days until the monthly reset
func (p *premiumUsageUseCase) getDaysUntilReset(usage *entities.PremiumUsage) int {
	now := usage.LastReset
	nextMonth := now.AddDate(0, 1, 0)

	// Get the first day of next month
	firstOfNextMonth := nextMonth.AddDate(0, 0, -nextMonth.Day()+1)

	days := int(firstOfNextMonth.Sub(now).Hours() / 24)
	if days < 1 {
		days = 1
	}
	return days
}
