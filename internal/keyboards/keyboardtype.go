package keyboards

import (
	"github.com/ivanenkomaksym/remindme_bot/internal/models"
)

type KeyboardType int64

const (
	Main KeyboardType = iota
	Reccurence
	Time
	Week
	Message
)

func (k KeyboardType) String() string {
	switch k {
	case Main:
		return "main"
	case Reccurence:
		return "reccurence"
	case Time:
		return "time"
	case Week:
		return "week"
	case Message:
		return "message"
	default:
		return "unknown"
	}
}

func GetKeyboardType(callbackData string) KeyboardType {
	if IsMainMenuSelection(callbackData) {
		return Main
	}

	if IsTimeSelectionCallback(callbackData) {
		return Time
	}
	if IsWeekSelectionCallback(callbackData) {
		return Week
	}
	if IsMessageSelectionCallback(callbackData) {
		return Message
	}
	_, err := models.ToRecurrenceType(callbackData)
	if err == nil {
		return Reccurence
	}
	return -1
}
