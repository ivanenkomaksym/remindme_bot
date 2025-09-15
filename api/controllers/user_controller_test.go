package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/usecases"
)

type userUseCaseMock struct {
	usecases.UserUseCase
	getUserFn            func(userID int64) (*entities.User, error)
	updateUserLanguageFn func(userID int64, language string) error
	getUserSelectionFn   func(userID int64) (*entities.UserSelection, error)
	clearUserSelectionFn func(userID int64) error
}

func (m *userUseCaseMock) GetUser(userID int64) (*entities.User, error) {
	if m.getUserFn != nil {
		return m.getUserFn(userID)
	}
	return &entities.User{ID: userID}, nil
}
func (m *userUseCaseMock) UpdateUserLanguage(userID int64, language string) error {
	if m.updateUserLanguageFn != nil {
		return m.updateUserLanguageFn(userID, language)
	}
	return nil
}
func (m *userUseCaseMock) GetUserSelection(userID int64) (*entities.UserSelection, error) {
	if m.getUserSelectionFn != nil {
		return m.getUserSelectionFn(userID)
	}
	return entities.NewUserSelection(), nil
}
func (m *userUseCaseMock) ClearUserSelection(userID int64) error {
	if m.clearUserSelectionFn != nil {
		return m.clearUserSelectionFn(userID)
	}
	return nil
}

func TestUserController_MethodGuards(t *testing.T) {
	c := NewUserController(&userUseCaseMock{})

	// GetUser wrong method
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/user?user_id=1", nil)
	c.GetUser(rw, req)
	if rw.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rw.Code)
	}

	// UpdateUserLanguage wrong method
	rw = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/user/language?user_id=1", nil)
	c.UpdateUserLanguage(rw, req)
	if rw.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rw.Code)
	}

	// GetUserSelection wrong method
	rw = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/user/selection?user_id=1", nil)
	c.GetUserSelection(rw, req)
	if rw.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rw.Code)
	}

	// ClearUserSelection wrong method
	rw = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/user/selection?user_id=1", nil)
	c.ClearUserSelection(rw, req)
	if rw.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rw.Code)
	}
}

func TestUserController_GetUser_Validation(t *testing.T) {
	c := NewUserController(&userUseCaseMock{})

	// missing user_id
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	c.GetUser(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}

	// bad user_id
	rw = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/user?user_id=abc", nil)
	c.GetUser(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}
}

func TestUserController_UpdateUserLanguage_Validation(t *testing.T) {
	c := NewUserController(&userUseCaseMock{})

	// missing user_id
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/user/language", bytes.NewBufferString(`{"language":"en"}`))
	c.UpdateUserLanguage(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}

	// invalid json
	rw = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPut, "/user/language?user_id=1", bytes.NewBufferString("{"))
	c.UpdateUserLanguage(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}

	// missing field
	rw = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPut, "/user/language?user_id=1", bytes.NewBufferString(`{}`))
	c.UpdateUserLanguage(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}
}

func TestUserController_GetUserSelection_Validation(t *testing.T) {
	c := NewUserController(&userUseCaseMock{})

	// missing user_id
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/user/selection", nil)
	c.GetUserSelection(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}

	// bad user_id
	rw = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/user/selection?user_id=abc", nil)
	c.GetUserSelection(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}
}

func TestUserController_ClearUserSelection_Validation(t *testing.T) {
	c := NewUserController(&userUseCaseMock{})

	// missing user_id
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/user/selection", nil)
	c.ClearUserSelection(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}

	// bad user_id
	rw = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, "/user/selection?user_id=abc", nil)
	c.ClearUserSelection(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}
}
