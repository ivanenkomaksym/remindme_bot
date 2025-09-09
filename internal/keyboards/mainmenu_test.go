package keyboards

import "testing"

func TestIsMainMenuSelection(t *testing.T) {
	if !IsMainMenuSelection(MainMenu) {
		t.Fatalf("MainMenu should be recognized")
	}
	if IsMainMenuSelection("not_main") {
		t.Fatalf("Unexpected main menu recognition")
	}
}

func TestGetMainMenuMarkup(t *testing.T) {
	m := GetMainMenuMarkup()
	if len(m.InlineKeyboard) != 6 {
		t.Fatalf("expected 6 rows, got %d", len(m.InlineKeyboard))
	}
	// Each row should contain exactly 1 button
	for i := 0; i < 6; i++ {
		if len(m.InlineKeyboard[i]) != 1 {
			t.Fatalf("row %d expected 1 button, got %d", i, len(m.InlineKeyboard[i]))
		}
	}
}
