package keyboards

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestIsWeekSelectionCallback(t *testing.T) {
	if !IsWeekSelectionCallback(CallbackWeekDay + "0") {
		t.Fatalf("day callback should be recognized")
	}
	if !IsWeekSelectionCallback(CallbackWeekSelect) {
		t.Fatalf("select should be recognized")
	}
	if IsWeekSelectionCallback("nope") {
		t.Fatalf("unexpected recognition")
	}
}

func TestGetWeekRangeMarkup(t *testing.T) {
	var opts [7]bool
	m := GetWeekRangeMarkup(opts, LangEN)
	// 7 day rows + select + back
	if len(m.InlineKeyboard) != 9 {
		t.Fatalf("expected 9 rows, got %d", len(m.InlineKeyboard))
	}
	for i := 0; i < 7; i++ {
		if len(m.InlineKeyboard[i]) != 1 {
			t.Fatalf("day row %d should have 1 button", i)
		}
	}
}

func TestHandleWeekSelection_ToggleAndSelect(t *testing.T) {
	var msg tgbotapi.EditMessageTextConfig
	opts := [7]bool{}

	// Toggle Monday (index 0)
	mk := HandleWeekSelection(CallbackWeekDay+"0", &msg, &opts, LangEN)
	if mk == nil || !opts[0] {
		t.Fatalf("monday should be toggled on")
	}
	if len(mk.InlineKeyboard) != 9 {
		t.Fatalf("expected 9 rows after toggle")
	}

	// Confirm selection
	mk = HandleWeekSelection(CallbackWeekSelect, &msg, &opts, LangEN)
	if mk == nil {
		t.Fatalf("expected markup on select")
	}
}
