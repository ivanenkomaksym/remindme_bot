package keyboards

import (
	"strings"
)

// Callback prefixes for Once flow
const (
	CallbackDatepickerSelection = "datepicker_selection"
)

func IsDateCallback(callbackData string) bool {
	return strings.HasPrefix(callbackData, CallbackDatepickerSelection)
}
