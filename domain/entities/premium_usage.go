package entities

import "time"

// PremiumStatus represents different user subscription levels
type PremiumStatus string

const (
	PremiumStatusFree  PremiumStatus = "free"
	PremiumStatusBasic PremiumStatus = "basic"
	PremiumStatusPro   PremiumStatus = "pro"
)

// String returns the string representation of PremiumStatus
func (ps PremiumStatus) String() string {
	return string(ps)
}

// PremiumUsage represents a user's premium features usage tracking
type PremiumUsage struct {
	UserID        int64         `json:"userId" bson:"userId"`
	RequestsUsed  int           `json:"requestsUsed" bson:"requestsUsed"`
	RequestsLimit int           `json:"requestsLimit" bson:"requestsLimit"`
	LastReset     time.Time     `json:"lastReset" bson:"lastReset"`
	PremiumStatus PremiumStatus `json:"premiumStatus" bson:"premiumStatus"`
	CreatedAt     time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt" bson:"updatedAt"`
}

// NewPremiumUsage creates a new premium usage record for a user
func NewPremiumUsage(userID int64) *PremiumUsage {
	now := time.Now()
	return &PremiumUsage{
		UserID:        userID,
		RequestsUsed:  0,
		RequestsLimit: GetDefaultRequestLimit(PremiumStatusFree),
		LastReset:     now,
		PremiumStatus: PremiumStatusFree,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// GetDefaultRequestLimit returns the default monthly request limit for a premium status
func GetDefaultRequestLimit(status PremiumStatus) int {
	switch status {
	case PremiumStatusFree:
		return 20 // 20 free requests per month
	case PremiumStatusBasic:
		return 500 // 500 requests per month for basic premium
	case PremiumStatusPro:
		return -1 // Unlimited for pro users
	default:
		return 20
	}
}

// CanMakeRequest checks if the user can make another premium request
func (pu *PremiumUsage) CanMakeRequest() bool {
	// Check if we need to reset (new month)
	if pu.ShouldReset() {
		pu.ResetUsage()
	}

	// Pro users have unlimited requests
	if pu.PremiumStatus == PremiumStatusPro {
		return true
	}

	// Check if under limit
	return pu.RequestsUsed < pu.RequestsLimit
}

// ShouldReset checks if the usage should be reset (monthly reset)
func (pu *PremiumUsage) ShouldReset() bool {
	now := time.Now()
	lastResetMonth := pu.LastReset.Month()
	lastResetYear := pu.LastReset.Year()
	currentMonth := now.Month()
	currentYear := now.Year()

	return currentYear > lastResetYear || (currentYear == lastResetYear && currentMonth > lastResetMonth)
}

// ResetUsage resets the monthly usage counter
func (pu *PremiumUsage) ResetUsage() {
	pu.RequestsUsed = 0
	pu.LastReset = time.Now()
	pu.UpdatedAt = time.Now()
}

// IncrementUsage increments the usage counter
func (pu *PremiumUsage) IncrementUsage() {
	pu.RequestsUsed++
	pu.UpdatedAt = time.Now()
}

// SetPremiumStatus updates the premium status and adjusts limits
func (pu *PremiumUsage) SetPremiumStatus(status PremiumStatus) {
	pu.PremiumStatus = status
	pu.RequestsLimit = GetDefaultRequestLimit(status)
	pu.UpdatedAt = time.Now()
}

// GetRemainingRequests returns the number of remaining requests
func (pu *PremiumUsage) GetRemainingRequests() int {
	if pu.PremiumStatus == PremiumStatusPro {
		return -1 // Unlimited
	}

	remaining := pu.RequestsLimit - pu.RequestsUsed
	if remaining < 0 {
		return 0
	}
	return remaining
}

// IsOverLimit checks if the user has exceeded their limit
func (pu *PremiumUsage) IsOverLimit() bool {
	if pu.PremiumStatus == PremiumStatusPro {
		return false
	}
	return pu.RequestsUsed >= pu.RequestsLimit
}
