package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ivanenkomaksym/remindme_bot/domain/usecases"
	"github.com/ivanenkomaksym/remindme_bot/keyboards"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type botUseCaseMock struct {
	usecases.BotUseCase
	processKeyboardSelectionFn func(cb *tgbotapi.CallbackQuery) (*keyboards.SelectionResult, error)
	processUserInputFn         func(msg *tgbotapi.Message) (*keyboards.SelectionResult, error)
}

func (m *botUseCaseMock) ProcessKeyboardSelection(cb *tgbotapi.CallbackQuery) (*keyboards.SelectionResult, error) {
	if m.processKeyboardSelectionFn != nil {
		return m.processKeyboardSelectionFn(cb)
	}
	return nil, nil
}
func (m *botUseCaseMock) ProcessUserInput(msg *tgbotapi.Message) (*keyboards.SelectionResult, error) {
	if m.processUserInputFn != nil {
		return m.processUserInputFn(msg)
	}
	return nil, nil
}

type botAPINop struct{ t *testing.T }

func (b *botAPINop) Send(ch interface{}) (tgbotapi.Message, error) { return tgbotapi.Message{}, nil }

func TestHandleWebhook_MethodGuard(t *testing.T) {
	ctrl := NewBotController(&botUseCaseMock{}, &tgbotapi.BotAPI{})
	req := httptest.NewRequest(http.MethodGet, "/webhook", nil)
	rw := httptest.NewRecorder()

	ctrl.HandleWebhook(rw, req)

	if rw.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rw.Code)
	}
}

func TestHandleWebhook_BadJSON(t *testing.T) {
	ctrl := NewBotController(&botUseCaseMock{}, &tgbotapi.BotAPI{})
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBufferString("{"))
	rw := httptest.NewRecorder()

	ctrl.HandleWebhook(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}
}

func TestHandleWebhook_ValidUpdate(t *testing.T) {
	// Minimal update with text message
	update := tgbotapi.Update{UpdateID: 1, Message: &tgbotapi.Message{MessageID: 2, From: &tgbotapi.User{ID: 10}, Chat: &tgbotapi.Chat{ID: 10}, Text: "hi"}}
	body, _ := json.Marshal(update)

	// Return nil selection to avoid calling Send on real bot
	mockUC := &botUseCaseMock{processUserInputFn: func(msg *tgbotapi.Message) (*keyboards.SelectionResult, error) { return nil, nil }}
	ctrl := NewBotController(mockUC, &tgbotapi.BotAPI{})
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	rw := httptest.NewRecorder()

	ctrl.HandleWebhook(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
}
