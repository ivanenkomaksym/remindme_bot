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

	if usage.RequestsLimit != 20 {
		t.Errorf("Expected RequestsLimit 20, got %d", usage.RequestsLimit)
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
		name      string
		lastReset time.Time
		expected  bool
	}{
		{
			name:      "Same month",
			lastReset: now.AddDate(0, 0, -10), // 10 days ago
			expected:  false,
		},
		{
			name:      "Previous month",
			lastReset: now.AddDate(0, -1, 0), // 1 month ago
			expected:  true,
		},
		{
			name:      "Previous year",
			lastReset: now.AddDate(-1, 0, 0), // 1 year ago
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usage := &PremiumUsage{
				LastReset: tt.lastReset,
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
	if usage.RequestsLimit != 500 {
		t.Errorf("Expected RequestsLimit 500, got %d", usage.RequestsLimit)
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
