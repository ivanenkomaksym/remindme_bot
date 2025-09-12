package keyboards

import "github.com/ivanenkomaksym/remindme_bot/internal/models"

type Strings struct {
	Welcome               string
	RecurrenceTypes       map[models.RecurrenceType]string
	BtnBack               string
	BtnCustomTime         string
	MsgSelectTime         string
	MsgSelectHour         string
	MsgSelectWithinHour   string
	MsgSelectMessage      string
	MsgEnterCustomTime    string
	MsgEnterCustomMessage string
	MsgInvalidTimeFormat  string
	BtnMyReminders        string
	NoReminders           string
	YourReminders         string
	BtnDelete             string
	DefaultMessages       []string
	ReminderSet           string
	Frequency             string
	Days                  string
	NoneSelected          string
	Time                  string
	Message               string
	ReminderScheduled     string
	At                    string
}

var stringsByLang = map[string]Strings{
	LangEN: {
		Welcome: "Welcome to the Reminder Bot!",
		RecurrenceTypes: map[models.RecurrenceType]string{
			models.Daily:    "Daily",
			models.Weekly:   "Weekly",
			models.Monthly:  "Monthly",
			models.Interval: "Interval",
			models.Custom:   "Custom",
		},
		BtnBack:               "← Back",
		BtnCustomTime:         "Custom",
		MsgSelectTime:         "Select time for daily reminders:",
		MsgSelectHour:         "Select time range:",
		MsgSelectWithinHour:   "Select time within %02d:00-%02d:00:",
		MsgSelectMessage:      "Select your reminder message:",
		MsgEnterCustomTime:    "Please type your custom time in HH:MM format (e.g., 14:30):",
		MsgEnterCustomMessage: "Please type your custom reminder message:",
		MsgInvalidTimeFormat:  "Invalid time format.",
		BtnMyReminders:        "My reminders",
		NoReminders:           "You have no reminders yet.",
		YourReminders:         "Your reminders:\n\n",
		BtnDelete:             "Delete",
		DefaultMessages: []string{"Time to take a break!",
			"Don't forget your medication",
			"Check your email",
			"Drink some water",
			"Stand up and stretch",
			"Review your tasks"},
		ReminderSet:       "Reminder Set",
		Frequency:         "Frequency",
		Days:              "Days",
		NoneSelected:      "None selected",
		Time:              "Time",
		Message:           "Message",
		ReminderScheduled: "Your reminder has been scheduled!",
		At:                "at",
	},
	LangUK: {
		Welcome: "Ласкаво просимо до бота-нагадувача!",
		RecurrenceTypes: map[models.RecurrenceType]string{
			models.Daily:    "Щодня",
			models.Weekly:   "Щотижня",
			models.Monthly:  "Щомісяця",
			models.Interval: "Інтервал",
			models.Custom:   "Власний",
		},
		BtnBack:               "← Назад",
		BtnCustomTime:         "Свій час",
		MsgSelectTime:         "Оберіть час для щоденних нагадувань:",
		MsgSelectHour:         "Оберіть діапазон часу:",
		MsgSelectWithinHour:   "Оберіть час між %02d:00-%02d:00:",
		MsgSelectMessage:      "Оберіть текст нагадування:",
		MsgEnterCustomTime:    "Введіть час у форматі HH:MM (напр., 14:30):",
		MsgEnterCustomMessage: "Введіть власний текст нагадування:",
		MsgInvalidTimeFormat:  "Неправильний формат часу.",
		BtnMyReminders:        "Мої нагадування",
		NoReminders:           "У вас ще немає нагадувань.",
		YourReminders:         "Ваші нагадування:\n\n",
		BtnDelete:             "Видалити",
		DefaultMessages: []string{"Час зробити перерву!",
			"Не забудьте прийняти ліки",
			"Перевірте свою електронну пошту",
			"Випийте трохи води",
			"Встаньте і розімніться",
			"Перегляньте свої завдання"},
		ReminderSet:       "Нагадування встановлено",
		Frequency:         "Частота",
		Days:              "Дні",
		NoneSelected:      "Нічого не вибрано",
		Time:              "Час",
		Message:           "Повідомлення",
		ReminderScheduled: "Ваше нагадування заплановано!",
		At:                "в",
	},
}

func T(lang string) Strings {
	if s, ok := stringsByLang[lang]; ok {
		return s
	}
	return stringsByLang[LangEN]
}

// RecurrenceTypeLabel returns a localized string for a given recurrence type.
// Falls back to the enum String() if translation is missing.
func RecurrenceTypeLabel(lang string, rt models.RecurrenceType) string {
	s := T(lang)
	if s.RecurrenceTypes != nil {
		if v, ok := s.RecurrenceTypes[rt]; ok {
			return v
		}
	}
	return rt.String()
}
