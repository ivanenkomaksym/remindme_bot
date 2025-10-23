package usecases

import (
	"testing"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/repositories/inmemory"
)

func TestPremiumUsageUseCase_GetOrCreateUserUsage(t *testing.T) {
	repo := inmemory.NewInMemoryPremiumUsageRepository()
	useCase := NewPremiumUsageUseCase(repo)

	userID := int64(123)

	// First call should create new usage
	usage1, err := useCase.GetOrCreateUserUsage(userID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if usage1.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, usage1.UserID)
	}

	if usage1.PremiumStatus != entities.PremiumStatusFree {
		t.Errorf("Expected free status, got %s", usage1.PremiumStatus)
	}

	// Second call should return existing usage
	usage2, err := useCase.GetOrCreateUserUsage(userID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if usage2.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, usage2.UserID)
	}
}

func TestPremiumUsageUseCase_UpgradeUser(t *testing.T) {
	repo := inmemory.NewInMemoryPremiumUsageRepository()
	useCase := NewPremiumUsageUseCase(repo)

	userID := int64(123)

	// Upgrade user to basic premium
	usage, err := useCase.UpgradeUser(userID, entities.PremiumStatusBasic)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if usage.PremiumStatus != entities.PremiumStatusBasic {
		t.Errorf("Expected basic premium status, got %s", usage.PremiumStatus)
	}

	if usage.RequestsLimit != 500 {
		t.Errorf("Expected 500 requests limit, got %d", usage.RequestsLimit)
	}
}

func TestPremiumUsageUseCase_ResetUserUsage(t *testing.T) {
	repo := inmemory.NewInMemoryPremiumUsageRepository()
	useCase := NewPremiumUsageUseCase(repo)

	userID := int64(123)

	// Create usage and use some requests
	usage, _ := useCase.GetOrCreateUserUsage(userID)
	usage.RequestsUsed = 10
	useCase.UpdateUserUsage(usage)

	// Reset usage
	resetUsage, err := useCase.ResetUserUsage(userID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resetUsage.RequestsUsed != 0 {
		t.Errorf("Expected 0 requests used after reset, got %d", resetUsage.RequestsUsed)
	}
}
