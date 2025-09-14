package keyboards

import (
	"testing"
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
	opts := [7]bool{}

	// Toggle Monday (index 0)
	res := HandleWeekSelection(CallbackWeekDay+"0", &opts, LangEN)
	if res == nil || !opts[0] {
		t.Fatalf("monday should be toggled on")
	}
	if res.Markup == nil || len(res.Markup.InlineKeyboard) != 9 {
		t.Fatalf("expected 9 rows after toggle")
	}

	// Confirm selection
	res = HandleWeekSelection(CallbackWeekSelect, &opts, LangEN)
	if res == nil || res.Markup == nil {
		t.Fatalf("expected markup on select")
	}
}
