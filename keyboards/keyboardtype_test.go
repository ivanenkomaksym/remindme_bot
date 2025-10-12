package keyboards

import (
	"testing"
)

func TestKeyboardTypeString(t *testing.T) {
	tests := []struct {
		k    KeyboardType
		want string
	}{
		{Setup, "setup"},
		{Reccurence, "reccurence"},
		{Time, "time"},
		{Week, "week"},
		{Month, "month"},
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

	if got := GetKeyboardType(SetupMenu); got != Setup {
		t.Fatalf("GetKeyboardType(SetupMenu) = %v, want %v", got, Setup)
	}
	if got := GetKeyboardType(CallbackTimeStart); got != Time {
		t.Fatalf("GetKeyboardType(time) = %v, want %v", got, Time)
	}
	if got := GetKeyboardType(CallbackWeekDay + s.WeekdayNames[0]); got != Week {
		t.Fatalf("GetKeyboardType(week day) = %v, want %v", got, Week)
	}
	if got := GetKeyboardType(CallbackMonthDay + "1"); got != Month {
		t.Fatalf("GetKeyboardType(month day) = %v, want %v", got, Month)
	}
	if got := GetKeyboardType(CallbackPrefixMessage + "0"); got != Message {
		t.Fatalf("GetKeyboardType(message) = %v, want %v", got, Message)
	}
}
