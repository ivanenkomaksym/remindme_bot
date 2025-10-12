package keyboards

import "github.com/ivanenkomaksym/remindme_bot/domain/entities"

type KeyboardType int64

const (
	Main KeyboardType = iota
	Setup
	Reccurence
	Date
	Time
	Week
	Month
	Message
	Reminders
)

func (kt KeyboardType) String() string {
	switch kt {
	case Main:
		return "main"
	case Setup:
		return "setup"
	case Reccurence:
		return "reccurence"
	case Date:
		return "date"
	case Time:
		return "time"
	case Week:
		return "week"
	case Month:
		return "month"
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

	if IsSetupMenuSelection(callbackData) {
		return Setup
	}

	if IsDateCallback(callbackData) {
		return Date
	}

	if IsTimeSelectionCallback(callbackData) {
		return Time
	}
	if IsWeekSelectionCallback(callbackData) {
		return Week
	}
	if IsMonthSelectionCallback(callbackData) {
		return Month
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
