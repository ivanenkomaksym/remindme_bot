package keyboards

import "testing"

func TestIsMonthSelectionCallback(t *testing.T) {
	if !IsMonthSelectionCallback(CallbackMonthDay + "1") {
		t.Fatalf("day callback should be recognized")
	}
	if !IsMonthSelectionCallback(CallbackMonthSelect) {
		t.Fatalf("select should be recognized")
	}
	if IsMonthSelectionCallback("nope") {
		t.Fatalf("unexpected recognition")
	}
}

func TestGetMonthRangeMarkup(t *testing.T) {
	var opts [28]bool
	m := GetMonthRangeMarkup(opts, LangEN)
	// 4 day rows + select + back
	if len(m.InlineKeyboard) != 6 {
		t.Fatalf("expected 6 rows, got %d", len(m.InlineKeyboard))
	}
	for i := 0; i < 4; i++ {
		if len(m.InlineKeyboard[i]) != 7 {
			t.Fatalf("day row %d should have 7 buttons", i)
		}
	}
}

func TestHandleMonthSelection_ToggleAndSelect(t *testing.T) {
	opts := [28]bool{}

	// Toggle day 1
	res := HandleMonthSelection(CallbackMonthDay+"1", &opts, LangEN)
	if res == nil || !opts[0] {
		t.Fatalf("day 1 should be toggled on")
	}
	if res.Markup == nil || len(res.Markup.InlineKeyboard) != 6 {
		t.Fatalf("expected 6 rows after toggle")
	}

	// Confirm selection
	res = HandleMonthSelection(CallbackMonthSelect, &opts, LangEN)
	if res == nil || res.Markup == nil {
		t.Fatalf("expected markup on select")
	}
}
