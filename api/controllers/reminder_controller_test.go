package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/usecases"
)

func TestReminderController_ValidationChecks(t *testing.T) {
	c := NewReminderController(&reminderUseCaseMock{}, &mockNLPService{}, &mockUserUseCase{})

	tests := []struct {
		name       string
		method     string
		path       string
		userID     string
		reminderID string
		handler    func(http.ResponseWriter, *http.Request)
	}{
		{
			name:    "GetUserReminders - Missing user_id",
			method:  http.MethodGet,
			path:    "/reminders",
			handler: c.GetUserReminders,
		},
		{
			name:    "GetUserReminders - Invalid user_id",
			method:  http.MethodGet,
			path:    "/reminders/abc",
			userID:  "abc",
			handler: c.GetUserReminders,
		},
		{
			name:    "DeleteReminder - Missing params",
			method:  http.MethodDelete,
			path:    "/reminders",
			handler: c.DeleteReminder,
		},
		{
			name:       "DeleteReminder - Invalid reminder_id",
			method:     http.MethodDelete,
			path:       "/reminders/1/x",
			userID:     "1",
			reminderID: "x",
			handler:    c.DeleteReminder,
		},
		{
			name:       "DeleteReminder - Invalid user_id",
			method:     http.MethodDelete,
			path:       "/reminders/x/1",
			userID:     "x",
			reminderID: "1",
			handler:    c.DeleteReminder,
		},
		{
			name:    "GetReminder - Missing params",
			method:  http.MethodGet,
			path:    "/reminders",
			handler: c.GetReminder,
		},
		{
			name:       "UpdateReminder - Invalid reminder_id",
			method:     http.MethodPut,
			path:       "/reminders/1/x",
			userID:     "1",
			reminderID: "x",
			handler:    c.UpdateReminder,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, nil)
			if tt.userID != "" {
				req.SetPathValue("user_id", tt.userID)
			}
			if tt.reminderID != "" {
				req.SetPathValue("reminder_id", tt.reminderID)
			}
			tt.handler(rw, req)
			if rw.Code != http.StatusBadRequest {
				t.Fatalf("%s: expected 400, got %d", tt.name, rw.Code)
			}
		})
	}
}

type reminderUseCaseMock struct {
	usecases.ReminderUseCase
	createReminderFn     func(userID int64, selection *entities.UserSelection) (*entities.Reminder, error)
	getUserRemindersFn   func(userID int64) ([]entities.Reminder, error)
	getReminderFn        func(userID, reminderID int64) (*entities.Reminder, error)
	getAllRemindersFn    func() ([]entities.Reminder, error)
	deleteReminderFn     func(reminderID, userID int64) error
	updateReminderFn     func(userID, reminderID int64, reminder *entities.Reminder) (*entities.Reminder, error)
	getActiveRemindersFn func() ([]entities.Reminder, error)
}

func (m *reminderUseCaseMock) CreateReminder(userID int64, selection *entities.UserSelection) (*entities.Reminder, error) {
	if m.createReminderFn != nil {
		return m.createReminderFn(userID, selection)
	}
	return &entities.Reminder{}, nil
}

func (m *reminderUseCaseMock) GetUserReminders(userID int64) ([]entities.Reminder, error) {
	if m.getUserRemindersFn != nil {
		return m.getUserRemindersFn(userID)
	}
	return []entities.Reminder{}, nil
}
func (m *reminderUseCaseMock) GetAllReminders() ([]entities.Reminder, error) {
	if m.getAllRemindersFn != nil {
		return m.getAllRemindersFn()
	}
	return []entities.Reminder{}, nil
}
func (m *reminderUseCaseMock) DeleteReminder(reminderID, userID int64) error {
	if m.deleteReminderFn != nil {
		return m.deleteReminderFn(reminderID, userID)
	}
	return nil
}
func (m *reminderUseCaseMock) GetReminder(userID, reminderID int64) (*entities.Reminder, error) {
	if m.getReminderFn != nil {
		return m.getReminderFn(userID, reminderID)
	}
	return &entities.Reminder{}, nil
}

func (m *reminderUseCaseMock) UpdateReminder(userID, reminderID int64, reminder *entities.Reminder) (*entities.Reminder, error) {
	if m.updateReminderFn != nil {
		return m.updateReminderFn(userID, reminderID, reminder)
	}
	return reminder, nil
}

func (m *reminderUseCaseMock) GetActiveReminders() ([]entities.Reminder, error) {
	if m.getActiveRemindersFn != nil {
		return m.getActiveRemindersFn()
	}
	return []entities.Reminder{}, nil
}

func TestReminderController_MethodGuards(t *testing.T) {
	c := NewReminderController(&reminderUseCaseMock{}, &mockNLPService{}, &mockUserUseCase{})

	tests := []struct {
		name    string
		method  string
		path    string
		handler func(http.ResponseWriter, *http.Request)
	}{
		{
			name:    "Wrong method for GetUserReminders",
			method:  http.MethodPatch,
			path:    "/reminders/1",
			handler: c.GetUserReminders,
		},
		{
			name:    "Wrong method for GetAllReminders",
			method:  http.MethodPost,
			path:    "/reminders/all",
			handler: c.GetAllReminders,
		},
		{
			name:    "Wrong method for DeleteReminder",
			method:  http.MethodGet,
			path:    "/reminders/1/1",
			handler: c.DeleteReminder,
		},
		{
			name:    "Wrong method for GetActiveReminders",
			method:  http.MethodPost,
			path:    "/reminders/active",
			handler: c.GetActiveReminders,
		},
		{
			name:    "Wrong method for GetReminder",
			method:  http.MethodPut,
			path:    "/reminders/1/1",
			handler: c.GetReminder,
		},
		{
			name:    "Wrong method for UpdateReminder",
			method:  http.MethodGet,
			path:    "/reminders/1/1",
			handler: c.UpdateReminder,
		},
		{
			name:    "Wrong method for CreateReminder",
			method:  http.MethodPut,
			path:    "/reminders/1",
			handler: c.CreateReminder,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.SetPathValue("user_id", "1")
			if strings.Contains(tt.path, "/1/1") {
				req.SetPathValue("reminder_id", "1")
			}
			tt.handler(rw, req)
			if rw.Code != http.StatusMethodNotAllowed {
				t.Fatalf("expected 405, got %d", rw.Code)
			}
		})
	}
}

func TestReminderController_DeleteReminder_Validation(t *testing.T) {
	c := NewReminderController(&reminderUseCaseMock{}, &mockNLPService{}, &mockUserUseCase{})

	// missing params
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/reminders", nil)
	c.DeleteReminder(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}

	// bad reminder_id
	rw = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, "/reminders/1/x", nil)
	req.SetPathValue("user_id", "1")
	req.SetPathValue("reminder_id", "x")
	c.DeleteReminder(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}

	// bad user_id
	rw = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, "/reminders/x/1", nil)
	req.SetPathValue("user_id", "x")
	req.SetPathValue("reminder_id", "1")
	c.DeleteReminder(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}
}

func TestReminderController_GetUserReminders_Success(t *testing.T) {
	// Setup mock reminder with expected data
	expectedReminders := []entities.Reminder{
		{
			ID:      1,
			UserID:  123,
			Message: "Test reminder 1",
		},
		{
			ID:      2,
			UserID:  123,
			Message: "Test reminder 2",
		},
	}

	mock := &reminderUseCaseMock{
		getUserRemindersFn: func(userID int64) ([]entities.Reminder, error) {
			if userID == 123 {
				return expectedReminders, nil
			}
			return nil, nil
		},
	}

	c := NewReminderController(mock, &mockNLPService{}, &mockUserUseCase{})

	// Test successful retrieval
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/reminders/123", nil)
	req.SetPathValue("user_id", "123")

	c.GetUserReminders(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
}

func TestReminderController_GetReminder_Success(t *testing.T) {
	expectedReminder := &entities.Reminder{
		ID:      1,
		UserID:  123,
		Message: "Test reminder 1",
	}

	mock := &reminderUseCaseMock{
		getReminderFn: func(userID, reminderID int64) (*entities.Reminder, error) {
			if userID == 123 && reminderID == 1 {
				return expectedReminder, nil
			}
			return nil, nil
		},
	}

	c := NewReminderController(mock, &mockNLPService{}, &mockUserUseCase{})

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/reminders/123/1", nil)
	req.SetPathValue("user_id", "123")
	req.SetPathValue("reminder_id", "1")

	c.GetReminder(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
}

func TestReminderController_CreateReminder_Success(t *testing.T) {
	expectedReminder := &entities.Reminder{
		ID:        1,
		UserID:    123,
		Message:   "Test reminder message",
		CreatedAt: time.Now(),
	}

	mock := &reminderUseCaseMock{
		createReminderFn: func(userID int64, selection *entities.UserSelection) (*entities.Reminder, error) {
			if userID == 123 && selection != nil {
				// Verify the selection data matches what we sent
				if selection.RecurrenceType != entities.Daily ||
					selection.ReminderMessage != "Test reminder message" ||
					selection.SelectedTime != "10:00" {
					t.Fatalf("received unexpected selection data")
				}
				return expectedReminder, nil
			}
			return nil, nil
		},
	}

	c := NewReminderController(mock, &mockNLPService{}, &mockUserUseCase{})

	body := `{
		"recurrenceType": "Daily",
		"selectedTime": "10:00",
		"reminderMessage": "Test reminder message"
	}`

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/reminders/123", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("user_id", "123")

	c.CreateReminder(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
}

func TestReminderController_GetAllReminders_Success(t *testing.T) {
	expectedReminders := []entities.Reminder{
		{
			ID:      1,
			UserID:  123,
			Message: "Test reminder 1",
		},
		{
			ID:      2,
			UserID:  456,
			Message: "Test reminder 2",
		},
	}

	mock := &reminderUseCaseMock{
		getAllRemindersFn: func() ([]entities.Reminder, error) {
			return expectedReminders, nil
		},
	}

	c := NewReminderController(mock, &mockNLPService{}, &mockUserUseCase{})

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/reminders/all", nil)

	c.GetAllReminders(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
}

func TestReminderController_GetActiveReminders_Success(t *testing.T) {
	expectedReminders := []entities.Reminder{
		{
			ID:      1,
			UserID:  123,
			Message: "Active reminder 1",
		},
		{
			ID:      2,
			UserID:  456,
			Message: "Active reminder 2",
		},
	}

	mock := &reminderUseCaseMock{
		getActiveRemindersFn: func() ([]entities.Reminder, error) {
			return expectedReminders, nil
		},
	}

	c := NewReminderController(mock, &mockNLPService{}, &mockUserUseCase{})

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/reminders/active", nil)

	c.GetActiveReminders(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
}

func TestReminderController_DeleteReminder_Success(t *testing.T) {
	mock := &reminderUseCaseMock{
		deleteReminderFn: func(reminderID, userID int64) error {
			if reminderID == 1 && userID == 123 {
				return nil
			}
			return nil
		},
	}

	c := NewReminderController(mock, &mockNLPService{}, &mockUserUseCase{})

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/reminders/123/1", nil)
	req.SetPathValue("user_id", "123")
	req.SetPathValue("reminder_id", "1")

	c.DeleteReminder(rw, req)

	if rw.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rw.Code)
	}
}

// Mock services for testing
type mockNLPService struct{}

func (m *mockNLPService) ParseReminderText(userID int64, text string, userTimezone string, userLanguage string) (*entities.UserSelection, error) {
	return nil, nil
}

type mockUserUseCase struct{}

func (m *mockUserUseCase) GetUsers() ([]*entities.User, error) {
	return nil, nil
}

func (m *mockUserUseCase) GetUser(userID int64) (*entities.User, error) {
	return nil, nil
}

func (m *mockUserUseCase) CreateUser(userID int64, userName, firstName, lastName, language string) (*entities.User, error) {
	return nil, nil
}

func (m *mockUserUseCase) GetOrCreateUser(userID int64, userName, firstName, lastName, language string) (*entities.User, error) {
	return nil, nil
}

func (m *mockUserUseCase) UpdateUserLanguage(userID int64, language string) error {
	return nil
}

func (m *mockUserUseCase) UpdateLocation(userID int64, location string) error {
	return nil
}

func (m *mockUserUseCase) GetUserSelection(userID int64) (*entities.UserSelection, error) {
	return nil, nil
}

func (m *mockUserUseCase) UpdateUserSelection(userID int64, selection *entities.UserSelection) error {
	return nil
}

func (m *mockUserUseCase) ClearUserSelection(userID int64) error {
	return nil
}

func (m *mockUserUseCase) DeleteUser(userID int64) error {
	return nil
}
