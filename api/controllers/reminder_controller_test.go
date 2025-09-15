package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/usecases"
)

type reminderUseCaseMock struct {
	usecases.ReminderUseCase
	getUserRemindersFn   func(userID int64) ([]entities.Reminder, error)
	getAllRemindersFn    func() ([]entities.Reminder, error)
	deleteReminderFn     func(reminderID, userID int64) error
	getActiveRemindersFn func() ([]entities.Reminder, error)
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
	req := httptest.NewRequest(http.MethodPost, "/reminders?user_id=1", nil)
	c.GetUserReminders(rw, req)
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
	req = httptest.NewRequest(http.MethodGet, "/reminders/delete?reminder_id=1&user_id=1", nil)
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

func TestReminderController_GetUserReminders_Validation(t *testing.T) {
	c := NewReminderController(&reminderUseCaseMock{})

	// missing user_id
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/reminders", nil)
	c.GetUserReminders(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}

	// bad user_id
	rw = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/reminders?user_id=abc", nil)
	c.GetUserReminders(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}
}

func TestReminderController_DeleteReminder_Validation(t *testing.T) {
	c := NewReminderController(&reminderUseCaseMock{})

	// missing params
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/reminders/delete", nil)
	c.DeleteReminder(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}

	// bad reminder_id
	rw = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, "/reminders/delete?reminder_id=x&user_id=1", nil)
	c.DeleteReminder(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}

	// bad user_id
	rw = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, "/reminders/delete?reminder_id=1&user_id=x", nil)
	c.DeleteReminder(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}
}
