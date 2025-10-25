package entities

import (
	"testing"
	"time"
)

func TestNewPremiumUsage(t *testing.T) {
	userID := int64(123)
	usage := NewPremiumUsage(userID)

	if usage.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, usage.UserID)
	}

	if usage.RequestsUsed != 0 {
		t.Errorf("Expected RequestsUsed 0, got %d", usage.RequestsUsed)
	}

	if usage.RequestsLimit != RequestLimitFree {
		t.Errorf("Expected RequestsLimit %d, got %d", RequestLimitFree, usage.RequestsLimit)
	}

	if usage.PremiumStatus != PremiumStatusFree {
		t.Errorf("Expected PremiumStatus %s, got %s", PremiumStatusFree, usage.PremiumStatus)
	}
}

func TestPremiumUsage_CanMakeRequest(t *testing.T) {
	tests := []struct {
		name          string
		requestsUsed  int
		requestsLimit int
		premiumStatus PremiumStatus
		expected      bool
	}{
		{
			name:          "Free user under limit",
			requestsUsed:  10,
			requestsLimit: 20,
			premiumStatus: PremiumStatusFree,
			expected:      true,
		},
		{
			name:          "Free user at limit",
			requestsUsed:  20,
			requestsLimit: 20,
			premiumStatus: PremiumStatusFree,
			expected:      false,
		},
		{
			name:          "Free user over limit",
			requestsUsed:  25,
			requestsLimit: 20,
			premiumStatus: PremiumStatusFree,
			expected:      false,
		},
		{
			name:          "Pro user over limit",
			requestsUsed:  1000,
			requestsLimit: 500,
			premiumStatus: PremiumStatusPro,
			expected:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usage := &PremiumUsage{
				RequestsUsed:  tt.requestsUsed,
				RequestsLimit: tt.requestsLimit,
				PremiumStatus: tt.premiumStatus,
				LastReset:     time.Now(),
			}

			result := usage.CanMakeRequest()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestPremiumUsage_ShouldReset(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name           string
		lastReset      time.Time
		premiumStatus  PremiumStatus
		premiumUpgrade *time.Time
		expected       bool
	}{
		{
			name:          "Free user - same month",
			lastReset:     now.AddDate(0, 0, -10), // 10 days ago
			premiumStatus: PremiumStatusFree,
			expected:      false,
		},
		{
			name:          "Free user - previous month",
			lastReset:     now.AddDate(0, -1, 0), // 1 month ago
			premiumStatus: PremiumStatusFree,
			expected:      true,
		},
		{
			name:           "Premium user - within 30 days of upgrade",
			lastReset:      now.AddDate(0, 0, -10), // 10 days ago
			premiumStatus:  PremiumStatusBasic,
			premiumUpgrade: &[]time.Time{now.AddDate(0, 0, -10)}[0], // upgraded 10 days ago
			expected:       false,
		},
		{
			name:           "Premium user - after 30 days from upgrade",
			lastReset:      now.AddDate(0, 0, -35), // 35 days ago
			premiumStatus:  PremiumStatusBasic,
			premiumUpgrade: &[]time.Time{now.AddDate(0, 0, -35)}[0], // upgraded 35 days ago
			expected:       true,
		},
		{
			name:           "Premium user - no upgrade date (fallback)",
			lastReset:      now.AddDate(0, -1, 0), // 1 month ago
			premiumStatus:  PremiumStatusBasic,
			premiumUpgrade: nil,
			expected:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usage := &PremiumUsage{
				LastReset:        tt.lastReset,
				PremiumStatus:    tt.premiumStatus,
				PremiumUpgradeAt: tt.premiumUpgrade,
			}

			result := usage.ShouldReset()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestPremiumUsage_SetPremiumStatus(t *testing.T) {
	usage := NewPremiumUsage(123)

	// Test upgrading to basic
	usage.SetPremiumStatus(PremiumStatusBasic)
	if usage.PremiumStatus != PremiumStatusBasic {
		t.Errorf("Expected PremiumStatus %s, got %s", PremiumStatusBasic, usage.PremiumStatus)
	}
	if usage.RequestsLimit != RequestLimitBasic {
		t.Errorf("Expected RequestsLimit %d, got %d", RequestLimitBasic, usage.RequestsLimit)
	}

	// Test upgrading to pro
	usage.SetPremiumStatus(PremiumStatusPro)
	if usage.PremiumStatus != PremiumStatusPro {
		t.Errorf("Expected PremiumStatus %s, got %s", PremiumStatusPro, usage.PremiumStatus)
	}
	if usage.RequestsLimit != -1 {
		t.Errorf("Expected RequestsLimit -1 (unlimited), got %d", usage.RequestsLimit)
	}
}

func TestPremiumUsage_GetRemainingRequests(t *testing.T) {
	tests := []struct {
		name          string
		requestsUsed  int
		requestsLimit int
		premiumStatus PremiumStatus
		expected      int
	}{
		{
			name:          "Free user with remaining requests",
			requestsUsed:  10,
			requestsLimit: 20,
			premiumStatus: PremiumStatusFree,
			expected:      10,
		},
		{
			name:          "Free user at limit",
			requestsUsed:  20,
			requestsLimit: 20,
			premiumStatus: PremiumStatusFree,
			expected:      0,
		},
		{
			name:          "Pro user unlimited",
			requestsUsed:  100,
			requestsLimit: -1,
			premiumStatus: PremiumStatusPro,
			expected:      -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usage := &PremiumUsage{
				RequestsUsed:  tt.requestsUsed,
				RequestsLimit: tt.requestsLimit,
				PremiumStatus: tt.premiumStatus,
			}

			result := usage.GetRemainingRequests()
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestPremiumUsage_PremiumSubscriptionFeatures(t *testing.T) {
	t.Run("Upgrade from free to premium", func(t *testing.T) {
		usage := NewPremiumUsage(123)

		// Initially free
		if usage.PremiumStatus != PremiumStatusFree {
			t.Errorf("Expected PremiumStatusFree, got %s", usage.PremiumStatus)
		}

		// Upgrade to basic
		usage.SetPremiumStatus(PremiumStatusBasic)

		if usage.PremiumStatus != PremiumStatusBasic {
			t.Errorf("Expected PremiumStatusBasic, got %s", usage.PremiumStatus)
		}

		if usage.PremiumUpgradeAt == nil {
			t.Error("Expected PremiumUpgradeAt to be set")
		}

		if usage.PremiumExpiresAt == nil {
			t.Error("Expected PremiumExpiresAt to be set")
		}

		if usage.RequestsUsed != 0 {
			t.Errorf("Expected RequestsUsed to be reset to 0, got %d", usage.RequestsUsed)
		}
	})

	t.Run("Premium expiration check", func(t *testing.T) {
		usage := NewPremiumUsage(123)

		// Set as premium with past expiration
		pastTime := time.Now().AddDate(0, 0, -1) // 1 day ago
		usage.PremiumStatus = PremiumStatusBasic
		usage.PremiumExpiresAt = &pastTime

		if !usage.IsPremiumExpired() {
			t.Error("Expected premium to be expired")
		}

		// Handle expiration
		if !usage.HandleExpiredPremium() {
			t.Error("Expected HandleExpiredPremium to return true")
		}

		if usage.PremiumStatus != PremiumStatusFree {
			t.Errorf("Expected status to be downgraded to free, got %s", usage.PremiumStatus)
		}
	})

	t.Run("Days until expiration", func(t *testing.T) {
		usage := NewPremiumUsage(123)
		usage.SetPremiumStatus(PremiumStatusBasic)

		days := usage.GetDaysUntilExpiration()
		if days < 29 || days > 30 {
			t.Errorf("Expected approximately 30 days until expiration, got %d", days)
		}

		// Free user should return -1
		usage.SetPremiumStatus(PremiumStatusFree)
		days = usage.GetDaysUntilExpiration()
		if days != -1 {
			t.Errorf("Expected -1 for free user, got %d", days)
		}
	})
}
