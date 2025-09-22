package entities

import "time"

// UserSelection represents a user's current selection state for creating reminders
type UserSelection struct {
	RecurrenceType  RecurrenceType `json:"recurrenceType"`
	WeekOptions     [7]bool        `json:"weekOptions"`
	MonthOptions    [28]bool       `json:"monthOptions"`
	SelectedDate    time.Time      `json:"selectedDate"`
	SelectedTime    string         `json:"selectedTime"`
	ReminderMessage string         `json:"reminderMessage"`
	CustomTime      bool           `json:"customTime"`
	CustomText      bool           `json:"customText"`
}

// NewUserSelection creates a new user selection with default values
func NewUserSelection() *UserSelection {
	return &UserSelection{
		WeekOptions: [7]bool{false, false, false, false, false, false, false},
		MonthOptions: [28]bool{
			false, false, false, false, false, false, false,
			false, false, false, false, false, false, false,
			false, false, false, false, false, false, false,
			false, false, false, false, false, false,
		},
	}
}

// SetRecurrenceType sets the recurrence type and updates weekly flag
func (us *UserSelection) SetRecurrenceType(recurrenceType RecurrenceType) {
	us.RecurrenceType = recurrenceType
}

// SetWeekOption sets a specific day of the week
func (us *UserSelection) SetWeekOption(day int, selected bool) {
	if day >= 0 && day < 7 {
		us.WeekOptions[day] = selected
	}
}

// SetSelectedTime sets the selected time
func (us *UserSelection) SetSelectedTime(time string) {
	us.SelectedTime = time
	us.CustomTime = false
}

// SetCustomTime enables custom time input
func (us *UserSelection) SetCustomTime() {
	us.CustomTime = true
	us.SelectedTime = ""
}

// SetReminderMessage sets the reminder message
func (us *UserSelection) SetReminderMessage(message string) {
	us.ReminderMessage = message
	us.CustomText = false
}

// SetCustomText enables custom text input
func (us *UserSelection) SetCustomText() {
	us.CustomText = true
	us.ReminderMessage = ""
}

// SetCustomText enables custom text input
func (us *UserSelection) SetSelectedDate(selectedDate time.Time) {
	us.SelectedDate = selectedDate
}

// Clear resets the user selection to default values
func (us *UserSelection) Clear() {
	*us = *NewUserSelection()
}
