package usecases

import (
	"testing"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/mocks"
	"github.com/ivanenkomaksym/remindme_bot/repositories/inmemory"
	"github.com/stretchr/testify/assert"
)

func TestUserUseCase_GetUser(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository()
	selRepo := inmemory.NewInMemoryUserSelectionRepository()
	useCase := NewUserUseCase(mockRepo, selRepo)

	// Test case: User exists
	user := entities.NewUser(1, "testuser", "Test", "User", "en")
	mockRepo.Users[1] = user

	result, err := useCase.GetUser(1)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "testuser", result.UserName)

	// Test case: User doesn't exist
	result, err = useCase.GetUser(999)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUserUseCase_CreateOrUpdateUser(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository()
	selRepo := inmemory.NewInMemoryUserSelectionRepository()
	useCase := NewUserUseCase(mockRepo, selRepo)

	// Test case: Valid user creation
	result, err := useCase.CreateOrUpdateUser(1, "testuser", "Test", "User", "en")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "testuser", result.UserName)
	assert.Equal(t, "Test", result.FirstName)
	assert.Equal(t, "User", result.LastName)
	assert.Equal(t, "en", result.Language)

	// Test case: Invalid user ID
	result, err = useCase.CreateOrUpdateUser(0, "testuser", "Test", "User", "en")
	assert.Error(t, err)
	assert.Nil(t, result)

	// Test case: Empty names
	result, err = useCase.CreateOrUpdateUser(1, "", "", "", "en")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUserUseCase_UpdateUserLanguage(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository()
	selRepo := inmemory.NewInMemoryUserSelectionRepository()
	useCase := NewUserUseCase(mockRepo, selRepo)

	// Create a user first
	user := entities.NewUser(1, "testuser", "Test", "User", "en")
	mockRepo.Users[1] = user

	// Test case: Valid language update
	err := useCase.UpdateUserLanguage(1, "es")
	assert.NoError(t, err)
	assert.Equal(t, "es", user.Language)

	// Test case: Invalid user ID
	err = useCase.UpdateUserLanguage(0, "es")
	assert.Error(t, err)

	// Test case: Empty language
	err = useCase.UpdateUserLanguage(1, "")
	assert.Error(t, err)
}

func TestUserUseCase_GetUserSelection(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository()
	selRepo := inmemory.NewInMemoryUserSelectionRepository()
	useCase := NewUserUseCase(mockRepo, selRepo)

	// Test case: User selection exists
	selection := entities.NewUserSelection()
	selection.SetRecurrenceType(entities.Daily)
	_ = selRepo.SetUserSelection(1, selection)

	result, err := useCase.GetUserSelection(1)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, entities.Daily, result.RecurrenceType)

	// Test case: User selection doesn't exist (should return empty selection)
	result, err = useCase.GetUserSelection(999)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, entities.RecurrenceType(0), result.RecurrenceType)

	// Test case: Invalid user ID
	result, err = useCase.GetUserSelection(0)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUserUseCase_UpdateUserSelection(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository()
	selRepo := inmemory.NewInMemoryUserSelectionRepository()
	useCase := NewUserUseCase(mockRepo, selRepo)

	selection := entities.NewUserSelection()
	selection.SetRecurrenceType(entities.Weekly)

	// Test case: Valid update
	err := useCase.UpdateUserSelection(1, selection)
	assert.NoError(t, err)

	// Test case: Invalid user ID
	err = useCase.UpdateUserSelection(0, selection)
	assert.Error(t, err)

	// Test case: Nil selection
	err = useCase.UpdateUserSelection(1, nil)
	assert.Error(t, err)
}

func TestUserUseCase_ClearUserSelection(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository()
	selRepo := inmemory.NewInMemoryUserSelectionRepository()
	useCase := NewUserUseCase(mockRepo, selRepo)

	// Create a user selection first
	selection := entities.NewUserSelection()
	selection.SetRecurrenceType(entities.Daily)
	_ = selRepo.SetUserSelection(1, selection)

	// Test case: Valid clear
	err := useCase.ClearUserSelection(1)
	assert.NoError(t, err)

	// Test case: Invalid user ID
	err = useCase.ClearUserSelection(0)
	assert.Error(t, err)
}
