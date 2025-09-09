package keyboards

type Strings struct {
	Welcome               string
	BtnDaily              string
	BtnWeekly             string
	BtnMonthly            string
	BtnInterval           string
	BtnCustom             string
	BtnBack               string
	BtnCustomTime         string
	MsgSelectTime         string
	MsgSelectHour         string
	MsgSelectWithinHour   string
	MsgSelectMessage      string
	MsgEnterCustomTime    string
	MsgEnterCustomMessage string
	BtnMyReminders        string
	NoReminders           string
	YourReminders         string
	BtnDelete             string
}

var stringsByLang = map[string]Strings{
	LangEN: {
		Welcome:               "Welcome to the Reminder Bot!",
		BtnDaily:              "daily",
		BtnWeekly:             "weekly",
		BtnMonthly:            "monthly",
		BtnInterval:           "interval",
		BtnCustom:             "custom",
		BtnBack:               "← Back",
		BtnCustomTime:         "Custom",
		MsgSelectTime:         "Select time for daily reminders:",
		MsgSelectHour:         "Select time range:",
		MsgSelectWithinHour:   "Select time within %02d:00-%02d:00:",
		MsgSelectMessage:      "Select your reminder message:",
		MsgEnterCustomTime:    "Please type your custom time in HH:MM format (e.g., 14:30):",
		MsgEnterCustomMessage: "Please type your custom reminder message:",
		BtnMyReminders:        "My reminders",
		NoReminders:           "You have no reminders yet.",
		YourReminders:         "Your reminders:\n\n",
		BtnDelete:             "Delete",
	},
	LangUK: {
		Welcome:               "Ласкаво просимо до бота-нагадувача!",
		BtnDaily:              "Щодня",
		BtnWeekly:             "Щотижня",
		BtnMonthly:            "Щомісяця",
		BtnInterval:           "Інтервал",
		BtnCustom:             "власний",
		BtnBack:               "← Назад",
		BtnCustomTime:         "Свій час",
		MsgSelectTime:         "Оберіть час для щоденних нагадувань:",
		MsgSelectHour:         "Оберіть діапазон часу:",
		MsgSelectWithinHour:   "Оберіть час між %02d:00-%02d:00:",
		MsgSelectMessage:      "Оберіть текст нагадування:",
		MsgEnterCustomTime:    "Введіть час у форматі HH:MM (напр., 14:30):",
		MsgEnterCustomMessage: "Введіть власний текст нагадування:",
		BtnMyReminders:        "Мої нагадування",
		NoReminders:           "У вас ще немає нагадувань.",
		YourReminders:         "Ваші нагадування:\n\n",
		BtnDelete:             "Видалити",
	},
}

func T(lang string) Strings {
	if s, ok := stringsByLang[lang]; ok {
		return s
	}
	return stringsByLang[LangEN]
}
