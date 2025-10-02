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

type reminderUseCaseMock struct {
	usecases.ReminderUseCase
	createReminderFn     func(userID int64, selection *entities.UserSelection) (*entities.Reminder, error)
	getUserRemindersFn   func(userID int64) ([]entities.Reminder, error)
	getAllRemindersFn    func() ([]entities.Reminder, error)
	deleteReminderFn     func(reminderID, userID int64) error
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
func (m *reminderUseCaseMock) GetActiveReminders() ([]entities.Reminder, error) {
	if m.getActiveRemindersFn != nil {
		return m.getActiveRemindersFn()
	}
	return []entities.Reminder{}, nil
}

func TestReminderController_MethodGuards(t *testing.T) {
	c := NewReminderController(&reminderUseCaseMock{})

	// Wrong method for GetUserReminders
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/reminders/1", nil)
	req.SetPathValue("user_id", "123")
	c.ProcessUserReminders(rw, req)
	if rw.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rw.Code)
	}

	// Wrong method for GetAllReminders
	rw = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/reminders/all", nil)
	c.GetAllReminders(rw, req)
	if rw.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rw.Code)
	}

	// Wrong method for DeleteReminder
	rw = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/reminders/1/1", nil)
	c.DeleteReminder(rw, req)
	if rw.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rw.Code)
	}

	// Wrong method for GetActiveReminders
	rw = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/reminders/active", nil)
	c.GetActiveReminders(rw, req)
	if rw.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rw.Code)
	}
}

func TestReminderController_ProcessUserReminders_Validation(t *testing.T) {
	c := NewReminderController(&reminderUseCaseMock{})

	// missing user_id
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/reminders", nil)
	c.ProcessUserReminders(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}

	// bad user_id
	rw = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/reminders/abc", nil)
	req.SetPathValue("user_id", "abc")
	c.ProcessUserReminders(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}
}

func TestReminderController_DeleteReminder_Validation(t *testing.T) {
	c := NewReminderController(&reminderUseCaseMock{})

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

	c := NewReminderController(mock)

	// Test successful retrieval
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/reminders/123", nil)
	req.SetPathValue("user_id", "123")

	c.ProcessUserReminders(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
}

func TestReminderController_CreateUserReminder_Success(t *testing.T) {
	// Setup mock reminder with expected data
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

	c := NewReminderController(mock)

	// Create request body with user selection data
	body := `{
		"recurrenceType": "Daily",
		"selectedTime": "10:00",
		"reminderMessage": "Test reminder message"
	}`

	// Test successful creation
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/reminders/123", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("user_id", "123")

	c.ProcessUserReminders(rw, req)

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

	c := NewReminderController(mock)

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

	c := NewReminderController(mock)

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/reminders/active", nil)

	c.GetActiveReminders(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
}

func TestReminderController_DeleteReminder_Success(t *testing.T) {
	deleteCount := 0
	mock := &reminderUseCaseMock{
		deleteReminderFn: func(reminderID, userID int64) error {
			if reminderID == 1 && userID == 123 {
				deleteCount++
				return nil
			}
			return nil
		},
	}

	c := NewReminderController(mock)

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/reminders/123/1", nil)
	req.SetPathValue("user_id", "123")
	req.SetPathValue("reminder_id", "1")

	c.DeleteReminder(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}

	if deleteCount != 1 {
		t.Fatalf("expected delete to be called once, got %d", deleteCount)
	}
}
