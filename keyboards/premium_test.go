package keyboards

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

func TestIsPremiumCallback(t *testing.T) {
	tests := []struct {
		name         string
		callbackData string
		expected     bool
	}{
		{"Premium view callback", CallbackAccountViewPremium, true},
		{"Premium upgrade callback", CallbackPremiumUpgrade, true},
		{"Non-premium callback", CallbackAccountChangeLanguage, false},
		{"Random callback", "random_callback", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPremiumCallback(tt.callbackData)
			if result != tt.expected {
				t.Errorf("IsPremiumCallback(%s) = %v, want %v", tt.callbackData, result, tt.expected)
			}
		})
	}
}

func TestHandlePremiumSelection_ViewPremium(t *testing.T) {
	user := &tgbotapi.User{ID: 123}
	userEntity := &entities.User{Language: LangEN}
	userUsage := entities.NewPremiumUsage(123)

	result, err := HandlePremiumSelection(user, CallbackAccountViewPremium, userEntity, userUsage)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Text == "" {
		t.Error("Expected non-empty text")
	}

	if result.Markup == nil {
		t.Error("Expected markup, got nil")
	}
}

func TestHandlePremiumSelection_Upgrade(t *testing.T) {
	user := &tgbotapi.User{ID: 123}
	userEntity := &entities.User{Language: LangEN}
	userUsage := entities.NewPremiumUsage(123)

	result, err := HandlePremiumSelection(user, CallbackPremiumUpgrade, userEntity, userUsage)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	s := T(LangEN)
	if result.Text != s.PremiumUpgradeComingSoon {
		t.Errorf("Expected upgrade coming soon text, got: %s", result.Text)
	}
}
