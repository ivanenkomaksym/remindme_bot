package keyboards

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/internal/types"
)

func TestIsMessageSelectionCallback(t *testing.T) {
	if !IsMessageSelectionCallback(CallbackPrefixMessage + "0") {
		t.Fatalf("prefix message should be recognized")
	}
	if IsMessageSelectionCallback("x") {
		t.Fatalf("unexpected recognition")
	}
}

func TestGetMessageSelectionMarkup(t *testing.T) {
	m := GetMessageSelectionMarkup(LangEN)
	// Default messages (6) + custom row + back row => 8 rows
	if len(m.InlineKeyboard) != len(DefaultMessages)+2 {
		t.Fatalf("unexpected rows: %d", len(m.InlineKeyboard))
	}
}

func TestHandleMessageSelection_DefaultAndCustom(t *testing.T) {
	var msg tgbotapi.EditMessageTextConfig
	state := &types.UserSelectionState{}

	// Custom path
	mk, done := HandleMessageSelection(CallbackMessageCustom, &msg, state)
	if mk != nil || done {
		t.Fatalf("custom should return nil, false")
	}
	if !state.CustomText {
		t.Fatalf("state.CustomText should be true")
	}

	// Pick default message index 1
	mk, done = HandleMessageSelection(CallbackPrefixMessage+string(rune(1)), &msg, state)
	if mk != nil || !done {
		t.Fatalf("default select should return nil, true")
	}
	if state.ReminderMessage != DefaultMessages[1] {
		t.Fatalf("unexpected message: %s", state.ReminderMessage)
	}
}

func TestHadleCustomText(t *testing.T) {
	var msg tgbotapi.MessageConfig
	state := &types.UserSelectionState{}
	_, done := HadleCustomText("hello", &msg, state)
	if !done || state.ReminderMessage != "hello" {
		t.Fatalf("custom text failed")
	}
}
