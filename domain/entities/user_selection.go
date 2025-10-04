package entities

import "time"

// UserSelection represents a user's current selection state for creating reminders
type UserSelection struct {
	RecurrenceType  RecurrenceType `json:"recurrenceType" bson:"recurrenceType"`
	WeekOptions     []time.Weekday `json:"weekOptions" bson:"weekOptions"`
	MonthOptions    []int          `json:"monthOptions" bson:"monthOptions"`
	SelectedDate    time.Time      `json:"selectedDate" bson:"selectedDate"`
	SelectedTime    string         `json:"selectedTime" bson:"selectedTime"`
	IntervalDays    int            `json:"intervalDays" bson:"intervalDays"`
	ReminderMessage string         `json:"reminderMessage" bson:"reminderMessage"`
	CustomTime      bool           `json:"customTime" bson:"customTime"`
	CustomText      bool           `json:"customText" bson:"customText"`
	CustomInterval  bool           `json:"customInterval" bson:"customInterval"`
}

// NewUserSelection creates a new user selection with default values
func NewUserSelection() *UserSelection {
	return &UserSelection{
		WeekOptions:  []time.Weekday{},
		MonthOptions: []int{},
	}
}

// SetRecurrenceType sets the recurrence type and updates weekly flag
func (us *UserSelection) SetRecurrenceType(recurrenceType RecurrenceType) {
	us.RecurrenceType = recurrenceType
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

// Enables custom interval input
func (us *UserSelection) StartCustomInterval() {
	us.CustomInterval = true
}

// Enables custom interval input
func (us *UserSelection) SetCustomInterval(interval int) {
	us.CustomInterval = false
	us.IntervalDays = interval
}

// SetCustomText enables custom text input
func (us *UserSelection) SetSelectedDate(selectedDate time.Time) {
	us.SelectedDate = selectedDate
}

// Clear resets the user selection to default values
func (us *UserSelection) Clear() {
	*us = *NewUserSelection()
}
