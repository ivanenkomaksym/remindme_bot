package keyboards

import (
	"testing"
)

func TestKeyboardTypeString(t *testing.T) {
	tests := []struct {
		k    KeyboardType
		want string
	}{
		{Main, "main"},
		{Reccurence, "reccurence"},
		{Time, "time"},
		{Week, "week"},
		{Message, "message"},
		{KeyboardType(-1), "unknown"},
	}
	for _, tt := range tests {
		if got := tt.k.String(); got != tt.want {
			t.Fatalf("String() = %s, want %s", got, tt.want)
		}
	}
}

func TestGetKeyboardType(t *testing.T) {
	s := T(LangEN)

	if got := GetKeyboardType(MainMenu); got != Main {
		t.Fatalf("GetKeyboardType(MainMenu) = %v, want %v", got, Main)
	}
	if got := GetKeyboardType(CallbackTimeStart); got != Time {
		t.Fatalf("GetKeyboardType(time) = %v, want %v", got, Time)
	}
	if got := GetKeyboardType(CallbackWeekDay + s.WeekdayNames[0]); got != Week {
		t.Fatalf("GetKeyboardType(week day) = %v, want %v", got, Week)
	}
	if got := GetKeyboardType(CallbackPrefixMessage + "0"); got != Message {
		t.Fatalf("GetKeyboardType(message) = %v, want %v", got, Message)
	}
}
