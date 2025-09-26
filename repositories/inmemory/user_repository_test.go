package inmemory

import (
	"testing"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryUserRepository_GetOrCreateUser(t *testing.T) {
	repo := NewInMemoryUserRepository()

	// Test case: Create new user
	user, err := repo.GetOrCreateUser(1, "testuser", "Test", "User", "en")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "testuser", user.UserName)
	assert.Equal(t, "Test", user.FirstName)
	assert.Equal(t, "User", user.LastName)
	assert.Equal(t, "en", user.Language)

	// Test case: Update existing user
	user, err = repo.GetOrCreateUser(1, "updateduser", "Updated", "User", "es")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "updateduser", user.UserName)
	assert.Equal(t, "Updated", user.FirstName)
	assert.Equal(t, "es", user.Language)
}

func TestInMemoryUserRepository_GetUser(t *testing.T) {
	repo := NewInMemoryUserRepository()

	// Test case: User doesn't exist
	user, err := repo.GetUser(1)
	assert.NoError(t, err)
	assert.Nil(t, user)

	// Create a user
	_, err = repo.GetOrCreateUser(1, "testuser", "Test", "User", "en")
	assert.NoError(t, err)

	// Test case: User exists
	user, err = repo.GetUser(1)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "testuser", user.UserName)
}

func TestInMemoryUserRepository_UpdateUserLanguage(t *testing.T) {
	repo := NewInMemoryUserRepository()

	// Create a user
	_, err := repo.GetOrCreateUser(1, "testuser", "Test", "User", "en")
	assert.NoError(t, err)

	// Test case: Update language
	err = repo.UpdateUserLanguage(1, "es")
	assert.NoError(t, err)

	// Verify the update
	user, err := repo.GetUser(1)
	assert.NoError(t, err)
	assert.Equal(t, "es", user.Language)
}

func TestInMemoryUserSelectionRepository_SelectionLifecycle(t *testing.T) {
	selRepo := NewInMemoryUserSelectionRepository()

	// Test case: Get non-existent selection
	selection, err := selRepo.GetUserSelection(1)
	assert.NoError(t, err)
	assert.Nil(t, selection)

	// Test case: Set user selection
	newSelection := entities.NewUserSelection()
	newSelection.SetRecurrenceType(entities.Daily)
	newSelection.SetSelectedTime("09:00")
	newSelection.SetReminderMessage("Test reminder")

	err = selRepo.SetUserSelection(1, newSelection)
	assert.NoError(t, err)

	// Test case: Get user selection
	selection, err = selRepo.GetUserSelection(1)
	assert.NoError(t, err)
	assert.NotNil(t, selection)
	assert.Equal(t, entities.Daily, selection.RecurrenceType)
	assert.Equal(t, "09:00", selection.SelectedTime)
	assert.Equal(t, "Test reminder", selection.ReminderMessage)

	// Test case: Update user selection
	selection.SetRecurrenceType(entities.Weekly)
	err = selRepo.UpdateUserSelection(1, selection)
	assert.NoError(t, err)

	// Verify the update
	updatedSelection, err := selRepo.GetUserSelection(1)
	assert.NoError(t, err)
	assert.Equal(t, entities.Weekly, updatedSelection.RecurrenceType)

	// Test case: Clear user selection
	err = selRepo.ClearUserSelection(1)
	assert.NoError(t, err)

	// Verify the clear
	clearedSelection, err := selRepo.GetUserSelection(1)
	assert.NoError(t, err)
	assert.Nil(t, clearedSelection)
}
