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
	UserID           int64         `json:"userId" bson:"userId"`
	RequestsUsed     int           `json:"requestsUsed" bson:"requestsUsed"`
	RequestsLimit    int           `json:"requestsLimit" bson:"requestsLimit"`
	LastReset        time.Time     `json:"lastReset" bson:"lastReset"`
	PremiumStatus    PremiumStatus `json:"premiumStatus" bson:"premiumStatus"`
	PremiumUpgradeAt *time.Time    `json:"premiumUpgradeAt,omitempty" bson:"premiumUpgradeAt,omitempty"`
	PremiumExpiresAt *time.Time    `json:"premiumExpiresAt,omitempty" bson:"premiumExpiresAt,omitempty"`
	CreatedAt        time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt        time.Time     `json:"updatedAt" bson:"updatedAt"`
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

// ShouldReset checks if the usage should be reset based on subscription timing
func (pu *PremiumUsage) ShouldReset() bool {
	now := time.Now()

	switch pu.PremiumStatus {
	case PremiumStatusFree:
		// Free users reset on calendar month boundary
		lastResetMonth := pu.LastReset.Month()
		lastResetYear := pu.LastReset.Year()
		currentMonth := now.Month()
		currentYear := now.Year()
		return currentYear > lastResetYear || (currentYear == lastResetYear && currentMonth > lastResetMonth)

	case PremiumStatusBasic, PremiumStatusPro:
		// Premium users reset based on subscription cycle (30 days from upgrade)
		if pu.PremiumUpgradeAt == nil {
			// Fallback to calendar month if no upgrade date
			lastResetMonth := pu.LastReset.Month()
			lastResetYear := pu.LastReset.Year()
			currentMonth := now.Month()
			currentYear := now.Year()
			return currentYear > lastResetYear || (currentYear == lastResetYear && currentMonth > lastResetMonth)
		}

		// Calculate next reset date (30 days from upgrade or last reset, whichever is more recent)
		var cycleStart time.Time
		if pu.LastReset.After(*pu.PremiumUpgradeAt) {
			cycleStart = pu.LastReset
		} else {
			cycleStart = *pu.PremiumUpgradeAt
		}

		nextResetDate := cycleStart.AddDate(0, 0, 30) // 30 days cycle
		return now.After(nextResetDate)

	default:
		// Unknown status, fallback to calendar month
		lastResetMonth := pu.LastReset.Month()
		lastResetYear := pu.LastReset.Year()
		currentMonth := now.Month()
		currentYear := now.Year()
		return currentYear > lastResetYear || (currentYear == lastResetYear && currentMonth > lastResetMonth)
	}
} // ResetUsage resets the monthly usage counter
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
	now := time.Now()

	// If upgrading from free to premium, record upgrade timestamp
	if pu.PremiumStatus == PremiumStatusFree && (status == PremiumStatusBasic || status == PremiumStatusPro) {
		pu.PremiumUpgradeAt = &now

		// Calculate expiration (30 days from upgrade)
		expiresAt := now.AddDate(0, 0, 30)
		pu.PremiumExpiresAt = &expiresAt

		// Reset usage immediately when upgrading to premium
		pu.RequestsUsed = 0
		pu.LastReset = now
	}

	// If changing premium levels, preserve upgrade date but update expiration
	if (pu.PremiumStatus == PremiumStatusBasic || pu.PremiumStatus == PremiumStatusPro) &&
		(status == PremiumStatusBasic || status == PremiumStatusPro) {
		// Keep existing upgrade date, update expiration from now
		expiresAt := now.AddDate(0, 0, 30)
		pu.PremiumExpiresAt = &expiresAt
	}

	// If downgrading to free, clear premium dates
	if status == PremiumStatusFree {
		pu.PremiumUpgradeAt = nil
		pu.PremiumExpiresAt = nil
	}

	pu.PremiumStatus = status
	pu.RequestsLimit = GetDefaultRequestLimit(status)
	pu.UpdatedAt = now
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

// IsPremiumExpired checks if the premium subscription has expired
func (pu *PremiumUsage) IsPremiumExpired() bool {
	if pu.PremiumStatus == PremiumStatusFree {
		return false // Free users don't have expiration
	}

	if pu.PremiumExpiresAt == nil {
		return false // No expiration date set
	}

	return time.Now().After(*pu.PremiumExpiresAt)
}

// GetDaysUntilExpiration returns days until premium expires (-1 if not premium or no expiration)
func (pu *PremiumUsage) GetDaysUntilExpiration() int {
	if pu.PremiumStatus == PremiumStatusFree || pu.PremiumExpiresAt == nil {
		return -1
	}

	if pu.IsPremiumExpired() {
		return 0
	}

	duration := time.Until(*pu.PremiumExpiresAt)
	days := int(duration.Hours() / 24)
	if days < 0 {
		return 0
	}
	return days
}

// GetDaysUntilReset returns days until next usage reset
func (pu *PremiumUsage) GetDaysUntilReset() int {
	now := time.Now()

	switch pu.PremiumStatus {
	case PremiumStatusFree:
		// Free users reset on first day of next month
		nextMonth := now.AddDate(0, 1, 0)
		firstOfNextMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, now.Location())
		duration := time.Until(firstOfNextMonth)
		days := int(duration.Hours() / 24)
		if days < 1 {
			return 1
		}
		return days

	case PremiumStatusBasic, PremiumStatusPro:
		// Premium users reset 30 days from upgrade or last reset
		if pu.PremiumUpgradeAt == nil {
			// Fallback to calendar month
			nextMonth := now.AddDate(0, 1, 0)
			firstOfNextMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, now.Location())
			duration := time.Until(firstOfNextMonth)
			days := int(duration.Hours() / 24)
			if days < 1 {
				return 1
			}
			return days
		}

		var cycleStart time.Time
		if pu.LastReset.After(*pu.PremiumUpgradeAt) {
			cycleStart = pu.LastReset
		} else {
			cycleStart = *pu.PremiumUpgradeAt
		}

		nextResetDate := cycleStart.AddDate(0, 0, 30)
		duration := time.Until(nextResetDate)
		days := int(duration.Hours() / 24)
		if days < 1 {
			return 1
		}
		return days

	default:
		return 1
	}
}

// HandleExpiredPremium downgrades expired premium users to free
func (pu *PremiumUsage) HandleExpiredPremium() bool {
	if pu.IsPremiumExpired() {
		pu.SetPremiumStatus(PremiumStatusFree)
		return true
	}
	return false
}
