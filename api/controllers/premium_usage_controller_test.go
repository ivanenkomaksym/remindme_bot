package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/errors"
)

// Mock repositories for testing
type mockPremiumUsageRepository struct {
	usages map[int64]*entities.PremiumUsage
}

func newMockPremiumUsageRepository() *mockPremiumUsageRepository {
	return &mockPremiumUsageRepository{
		usages: make(map[int64]*entities.PremiumUsage),
	}
}

func (m *mockPremiumUsageRepository) GetUserUsage(userID int64) (*entities.PremiumUsage, error) {
	if usage, exists := m.usages[userID]; exists {
		return usage, nil
	}
	return nil, errors.ErrUserNotFound
}

func (m *mockPremiumUsageRepository) CreateUserUsage(usage *entities.PremiumUsage) error {
	m.usages[usage.UserID] = usage
	return nil
}

func (m *mockPremiumUsageRepository) UpdateUserUsage(usage *entities.PremiumUsage) error {
	m.usages[usage.UserID] = usage
	return nil
}

func (m *mockPremiumUsageRepository) GetOrCreateUserUsage(userID int64) (*entities.PremiumUsage, error) {
	if usage, exists := m.usages[userID]; exists {
		return usage, nil
	}
	usage := entities.NewPremiumUsage(userID)
	m.usages[userID] = usage
	return usage, nil
}

func (m *mockPremiumUsageRepository) GetAllUsage() ([]entities.PremiumUsage, error) {
	var usages []entities.PremiumUsage
	for _, usage := range m.usages {
		usages = append(usages, *usage)
	}
	return usages, nil
}

func (m *mockPremiumUsageRepository) DeleteUserUsage(userID int64) error {
	delete(m.usages, userID)
	return nil
}

func (m *mockPremiumUsageRepository) GetUsageByPremiumStatus(status entities.PremiumStatus) ([]entities.PremiumUsage, error) {
	var filtered []entities.PremiumUsage
	for _, usage := range m.usages {
		if usage.PremiumStatus == status {
			filtered = append(filtered, *usage)
		}
	}
	return filtered, nil
}

type mockUserRepository struct {
	users map[int64]*entities.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[int64]*entities.User),
	}
}

func (m *mockUserRepository) GetUser(userID int64) (*entities.User, error) {
	if user, exists := m.users[userID]; exists {
		return user, nil
	}
	return nil, errors.ErrUserNotFound
}

func (m *mockUserRepository) GetUsers() ([]*entities.User, error) {
	var users []*entities.User
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *mockUserRepository) GetOrCreateUser(userID int64, userName, firstName, lastName, language string) (*entities.User, error) {
	if user, exists := m.users[userID]; exists {
		return user, nil
	}
	user := &entities.User{
		ID:        userID,
		FirstName: firstName,
		LastName:  lastName,
		UserName:  userName,
		Language:  language,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.users[userID] = user
	return user, nil
}

func (m *mockUserRepository) UpdateUserLanguage(userID int64, language string) error {
	if user, exists := m.users[userID]; exists {
		user.Language = language
		user.UpdatedAt = time.Now()
		return nil
	}
	return errors.ErrUserNotFound
}

func (m *mockUserRepository) CreateUser(userID int64, userName, firstName, lastName, language string) (*entities.User, error) {
	if _, exists := m.users[userID]; exists {
		return nil, errors.ErrUserExists
	}
	user := &entities.User{
		ID:        userID,
		FirstName: firstName,
		LastName:  lastName,
		UserName:  userName,
		Language:  language,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.users[userID] = user
	return user, nil
}

func (m *mockUserRepository) UpdateLocation(userID int64, location string) error {
	if user, exists := m.users[userID]; exists {
		user.LocationName = location
		user.UpdatedAt = time.Now()
		return nil
	}
	return errors.ErrUserNotFound
}

func (m *mockUserRepository) UpdateUserInfo(userID int64, userName, firstName, lastName string) error {
	if user, exists := m.users[userID]; exists {
		user.UserName = userName
		user.FirstName = firstName
		user.LastName = lastName
		user.UpdatedAt = time.Now()
		return nil
	}
	return errors.ErrUserNotFound
}

func (m *mockUserRepository) DeleteUser(userID int64) error {
	delete(m.users, userID)
	return nil
}

func TestPremiumUsageController_GetUserPremiumUsage(t *testing.T) {
	premiumRepo := newMockPremiumUsageRepository()
	userRepo := newMockUserRepository()
	controller := NewPremiumUsageController(premiumRepo, userRepo)

	// Create a test user
	userRepo.GetOrCreateUser(123, "testuser", "Test", "User", "en")

	t.Run("Get user premium usage - success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/premium/123", nil)
		req.SetPathValue("user_id", "123")
		w := httptest.NewRecorder()

		controller.GetUserPremiumUsage(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		json.NewDecoder(w.Body).Decode(&response)

		if !response["success"].(bool) {
			t.Error("Expected success to be true")
		}
	})

	t.Run("Get user premium usage - user not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/premium/999", nil)
		req.SetPathValue("user_id", "999")
		w := httptest.NewRecorder()

		controller.GetUserPremiumUsage(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("Get user premium usage - invalid user ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/premium/invalid", nil)
		req.SetPathValue("user_id", "invalid")
		w := httptest.NewRecorder()

		controller.GetUserPremiumUsage(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestPremiumUsageController_UpgradeUserPremium(t *testing.T) {
	premiumRepo := newMockPremiumUsageRepository()
	userRepo := newMockUserRepository()
	controller := NewPremiumUsageController(premiumRepo, userRepo)

	// Create a test user
	userRepo.GetOrCreateUser(123, "testuser", "Test", "User", "en")

	t.Run("Upgrade user premium - success", func(t *testing.T) {
		requestBody := UpgradePremiumRequest{
			PremiumStatus: entities.PremiumStatusBasic,
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/api/premium/123/upgrade", bytes.NewReader(body))
		req.SetPathValue("user_id", "123")
		w := httptest.NewRecorder()

		controller.UpgradeUserPremium(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		// Verify the upgrade worked
		usage, _ := premiumRepo.GetUserUsage(123)
		if usage.PremiumStatus != entities.PremiumStatusBasic {
			t.Errorf("Expected premium status to be basic, got %s", usage.PremiumStatus)
		}
	})

	t.Run("Upgrade user premium - invalid status", func(t *testing.T) {
		requestBody := UpgradePremiumRequest{
			PremiumStatus: "invalid",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/api/premium/123/upgrade", bytes.NewReader(body))
		req.SetPathValue("user_id", "123")
		w := httptest.NewRecorder()

		controller.UpgradeUserPremium(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestPremiumUsageController_ResetUserUsage(t *testing.T) {
	premiumRepo := newMockPremiumUsageRepository()
	userRepo := newMockUserRepository()
	controller := NewPremiumUsageController(premiumRepo, userRepo)

	// Create a test user and usage
	userRepo.GetOrCreateUser(123, "testuser", "Test", "User", "en")
	usage := entities.NewPremiumUsage(123)
	usage.RequestsUsed = 10
	premiumRepo.CreateUserUsage(usage)

	t.Run("Reset user usage - success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/premium/123/reset", nil)
		req.SetPathValue("user_id", "123")
		w := httptest.NewRecorder()

		controller.ResetUserUsage(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		// Verify the reset worked
		updatedUsage, _ := premiumRepo.GetUserUsage(123)
		if updatedUsage.RequestsUsed != 0 {
			t.Errorf("Expected requests used to be 0, got %d", updatedUsage.RequestsUsed)
		}
	})
}

func TestPremiumUsageController_GetPremiumUsageByStatus(t *testing.T) {
	premiumRepo := newMockPremiumUsageRepository()
	userRepo := newMockUserRepository()
	controller := NewPremiumUsageController(premiumRepo, userRepo)

	// Create test data
	usage1 := entities.NewPremiumUsage(123)
	usage1.SetPremiumStatus(entities.PremiumStatusFree)
	premiumRepo.CreateUserUsage(usage1)

	usage2 := entities.NewPremiumUsage(456)
	usage2.SetPremiumStatus(entities.PremiumStatusBasic)
	premiumRepo.CreateUserUsage(usage2)

	t.Run("Get premium usage by status - success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/premium/status/free", nil)
		req.SetPathValue("status", "free")
		w := httptest.NewRecorder()

		controller.GetPremiumUsageByStatus(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		json.NewDecoder(w.Body).Decode(&response)

		data := response["data"].([]interface{})
		if len(data) != 1 {
			t.Errorf("Expected 1 usage record, got %d", len(data))
		}
	})

	t.Run("Get premium usage by status - invalid status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/premium/status/invalid", nil)
		req.SetPathValue("status", "invalid")
		w := httptest.NewRecorder()

		controller.GetPremiumUsageByStatus(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}
