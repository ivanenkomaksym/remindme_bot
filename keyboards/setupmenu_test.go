package keyboards

import "testing"

func TestIsSetupMenuSelection(t *testing.T) {
	if !IsSetupMenuSelection(SetupMenu) {
		t.Fatalf("SetupMenu should be recognized")
	}
	if IsSetupMenuSelection("not_setup") {
		t.Fatalf("Unexpected setup menu recognition")
	}
}

func TestGetSetupMenuMarkup(t *testing.T) {
	expectedRows := 7
	m := GetSetupMenuMarkup(LangEN)
	if len(m.InlineKeyboard) != expectedRows {
		t.Fatalf("expected %d rows, got %d", expectedRows, len(m.InlineKeyboard))
	}
	// Each row should contain exactly 1 button
	for i := 0; i < expectedRows; i++ {
		if len(m.InlineKeyboard[i]) != 1 {
			t.Fatalf("row %d expected 1 button, got %d", i, len(m.InlineKeyboard[i]))
		}
	}
}
