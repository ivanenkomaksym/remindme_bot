package keyboards

import (
	"testing"
	"time"
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
	var weekdays []time.Weekday
	m := GetWeekRangeMarkup(weekdays, LangEN)
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
	var weekdays []time.Weekday

	// Toggle Monday (index 0)
	res := HandleWeekSelection(CallbackWeekDay+"Mon", &weekdays, LangEN)
	if res == nil || weekdays[0] != time.Monday {
		t.Fatalf("monday should be toggled on")
	}
	if res.Markup == nil || len(res.Markup.InlineKeyboard) != 9 {
		t.Fatalf("expected 9 rows after toggle")
	}

	// Confirm selection
	res = HandleWeekSelection(CallbackWeekSelect, &weekdays, LangEN)
	if res == nil || res.Markup == nil {
		t.Fatalf("expected markup on select")
	}
}
