package keyboards

import (
	"github.com/ivanenkomaksym/remindme_bot/internal/models"
)

type KeyboardType int64

const (
	Reccurence KeyboardType = iota
	Time
	DaysOfWeek
)

func (k KeyboardType) String() string {
	switch k {
	case Reccurence:
		return "reccurence"
	case Time:
		return "time"
	case DaysOfWeek:
		return "days_of_week"
	default:
		return "unknown"
	}
}

func GetKeyboardType(callbackData string) KeyboardType {
	if IsTimeSelectionCallback(callbackData) {
		return Time
	}
	_, err := models.ToRecurrenceType(callbackData)
	if err == nil {
		return Reccurence
	}
	return -1
}
