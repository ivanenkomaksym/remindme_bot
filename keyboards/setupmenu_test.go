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

func TestIsNlpTextInputCallback(t *testing.T) {
	if !IsNlpTextInputCallback(CallbackNlpTextInput) {
		t.Fatalf("CallbackNlpTextInput should be recognized")
	}
	if IsNlpTextInputCallback("not_nlp") {
		t.Fatalf("Unexpected NLP callback recognition")
	}
}

func TestGetSetupMenuMarkup(t *testing.T) {
	expectedRows := 5
	expectedButtonsPerRow := []int{1, 2, 2, 2, 2} // NLP (1), Once+Daily (2), Weekly+Monthly (2), Interval+Spaced (2), MyReminders+Back (2)

	m := GetSetupMenuMarkup(LangEN)
	if len(m.InlineKeyboard) != expectedRows {
		t.Fatalf("expected %d rows, got %d", expectedRows, len(m.InlineKeyboard))
	}

	// Check button count for each row
	for i := 0; i < expectedRows; i++ {
		expected := expectedButtonsPerRow[i]
		actual := len(m.InlineKeyboard[i])
		if actual != expected {
			t.Fatalf("row %d expected %d buttons, got %d", i, expected, actual)
		}
	}
}
