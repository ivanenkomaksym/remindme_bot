package keyboards

import (
	"github.com/ivanenkomaksym/remindme_bot/internal/models"
)

type KeyboardType int64

const (
	Reccurence KeyboardType = iota
	Time
	Week
)

func (k KeyboardType) String() string {
	switch k {
	case Reccurence:
		return "reccurence"
	case Time:
		return "time"
	case Week:
		return "week"
	default:
		return "unknown"
	}
}

func GetKeyboardType(callbackData string) KeyboardType {
	if IsTimeSelectionCallback(callbackData) {
		return Time
	}
	if IsWeekSelectionCallback(callbackData) {
		return Week
	}
	_, err := models.ToRecurrenceType(callbackData)
	if err == nil {
		return Reccurence
	}
	return -1
}
