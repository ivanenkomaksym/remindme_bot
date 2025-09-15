package keyboards

import "github.com/ivanenkomaksym/remindme_bot/domain/entities"

type KeyboardType int64

const (
	Main KeyboardType = iota
	Reccurence
	Time
	Week
	Message
	Reminders
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
	case Reminders:
		return "reminders"
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
	if IsRemindersCallback(callbackData) {
		return Reminders
	}
	_, err := entities.ToRecurrenceType(callbackData)
	if err == nil {
		return Reccurence
	}
	return -1
}
