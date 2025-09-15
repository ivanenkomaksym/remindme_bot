package keyboards

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
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
	s := T(LangEN)
	// Default messages (6) + custom row + back row => 8 rows
	if len(m.InlineKeyboard) != len(s.DefaultMessages)+2 {
		t.Fatalf("unexpected rows: %d", len(m.InlineKeyboard))
	}
}

func TestHandleMessageSelection_DefaultAndCustom(t *testing.T) {
	user := &entities.User{}
	userSelection := &entities.UserSelection{}
	s := T(LangEN)

	// Custom path
	res, done := HandleMessageSelection(CallbackMessageCustom, user, userSelection)
	if res == nil || done {
		t.Fatalf("custom should return result, false")
	}
	if !userSelection.CustomText {
		t.Fatalf("state.CustomText should be true")
	}

	// Pick default message index 1
	res, done = HandleMessageSelection(CallbackPrefixMessage+string(rune(1)), user, userSelection)
	if res != nil || !done {
		t.Fatalf("default select should return nil, true")
	}
	if userSelection.ReminderMessage != s.DefaultMessages[1] {
		t.Fatalf("unexpected message: %s", userSelection.ReminderMessage)
	}
}

func TestHadleCustomText(t *testing.T) {
	var msg tgbotapi.MessageConfig
	user := &entities.User{}
	userSelection := &entities.UserSelection{}
	_, done := HandleCustomText("hello", &msg, user, userSelection)
	if !done || userSelection.ReminderMessage != "hello" {
		t.Fatalf("custom text failed")
	}
}
