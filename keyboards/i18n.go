package keyboards

import (
	"time"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

var WeekdayNameToKeyMap = map[string]time.Weekday{
	"Mon": time.Monday,
	"Tue": time.Tuesday,
	"Wed": time.Wednesday,
	"Thu": time.Thursday,
	"Fri": time.Friday,
	"Sat": time.Saturday,
	"Sun": time.Sunday,
}

type Strings struct {
	Welcome                  string
	RecurrenceTypes          map[entities.RecurrenceType]string
	BtnBack                  string
	BtnCustomTime            string
	MsgSelectTime            string
	MsgSelectHour            string
	MsgSelectWithinHour      string
	MsgSelectMessage         string
	MsgEnterCustomTime       string
	MsgEnterCustomMessage    string
	MsgInvalidTimeFormat     string
	MsgInvalidIntervalFormat string
	BtnMyReminders           string
	NoReminders              string
	YourReminders            string
	BtnDelete                string
	DefaultMessages          []string
	ReminderSet              string
	Frequency                string
	Days                     string
	NoneSelected             string
	Date                     string
	Time                     string
	Message                  string
	ReminderScheduled        string
	At                       string
	// Week-related i18n
	WeekdayNames        map[time.Weekday]string
	WeekdayNamesShort   map[time.Weekday]string
	MsgSelectWeekdays   string
	MsgSelectTimeWeekly string
	BtnSelect           string
	// Date-related i18n
	MsgSelectDate string
	// Interval-related i18n
	MsgIntervalPrompt          string // e.g., "Every N days"
	MsgEveryNDays              string // e.g., "Every %d days"
	MsgEveryNDaysSpaced        string // e.g., "Every %s days"
	MsgParsingFailed           string
	MsgTimezoneAutoDetect      string
	MsgTimezoneAutoDetectDescr string
	MsgTimezoneSet             string
	// Navigation-related i18n
	NavList         string
	NavSetup        string
	NavAccount      string
	NavChooseOption string
	// Bot command descriptions
	CmdStartDesc   string
	CmdListDesc    string
	CmdSetupDesc   string
	CmdAccountDesc string
	// Account management i18n
	AccTitle          string
	AccUsername       string
	AccLanguage       string
	AccTimezone       string
	AccCreatedAt      string
	AccNoUsername     string
	AccNoTimezone     string
	AccChangeLanguage string
	AccChangeTimezone string
	// Timezone selection i18n
	TzManualSelect string
	TzSelectPrompt string
}

var stringsByLang = map[string]Strings{
	LangEN: {
		Welcome: "Welcome to the Reminder Bot!",
		RecurrenceTypes: map[entities.RecurrenceType]string{
			entities.Once:                  "Once",
			entities.Daily:                 "Daily",
			entities.Weekly:                "Weekly",
			entities.Monthly:               "Monthly",
			entities.Interval:              "Interval",
			entities.SpacedBasedRepetition: "Spaced Repetition",
		},
		BtnBack:                  "← Back",
		BtnCustomTime:            "Custom",
		MsgSelectTime:            "Select time for daily reminders:",
		MsgSelectHour:            "Select time range:",
		MsgSelectWithinHour:      "Select time within %02d:00-%02d:00:",
		MsgSelectMessage:         "Select your reminder message:",
		MsgEnterCustomTime:       "Please type your custom time in HH:MM format (e.g., 14:30):",
		MsgEnterCustomMessage:    "Please type your custom reminder message:",
		MsgInvalidTimeFormat:     "Invalid time format.",
		MsgInvalidIntervalFormat: "Invalid interval format. Expected 1-7",
		BtnMyReminders:           "My reminders",
		NoReminders:              "You have no reminders yet.",
		YourReminders:            "Your reminders:\n\n",
		BtnDelete:                "Delete",
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
		Date:              "Date",
		Time:              "Time",
		Message:           "Message",
		ReminderScheduled: "Your reminder has been scheduled!",
		At:                "at",
		WeekdayNames: map[time.Weekday]string{
			time.Monday:    "Monday",
			time.Tuesday:   "Tuesday",
			time.Wednesday: "Wednesday",
			time.Thursday:  "Thursday",
			time.Friday:    "Friday",
			time.Saturday:  "Saturday",
			time.Sunday:    "Sunday",
		},
		WeekdayNamesShort: map[time.Weekday]string{
			time.Monday:    "Mon",
			time.Tuesday:   "Tue",
			time.Wednesday: "Wed",
			time.Thursday:  "Thu",
			time.Friday:    "Fri",
			time.Saturday:  "Sat",
			time.Sunday:    "Sun",
		},
		MsgSelectWeekdays:          "Select weekdays:",
		MsgSelectTimeWeekly:        "Select time for weekly reminders:",
		BtnSelect:                  "Select",
		MsgSelectDate:              "Select a date:",
		MsgIntervalPrompt:          "Every N days",
		MsgEveryNDays:              "Every %d days",
		MsgEveryNDaysSpaced:        "Every %s days",
		MsgParsingFailed:           "I didn't understand that. Please use the menu buttons.",
		MsgTimezoneAutoDetect:      "🌍 Set Timezone Automatically",
		MsgTimezoneAutoDetectDescr: "Click the button to detect your timezone.",
		MsgTimezoneSet:             "✅ Your timezone is set to",
		NavList:                    "Show reminders",
		NavSetup:                   "Setup reminder",
		NavAccount:                 "Account",
		NavChooseOption:            "Choose an option:",
		CmdStartDesc:               "Start the bot and show main menu",
		CmdListDesc:                "Show or remove reminders",
		CmdSetupDesc:               "Set up time, recurrence, and reminder settings",
		CmdAccountDesc:             "Manage account settings",
		AccTitle:                   "👤 Account Information",
		AccUsername:                "Username",
		AccLanguage:                "Language",
		AccTimezone:                "Timezone",
		AccCreatedAt:               "Created",
		AccNoUsername:              "Not set",
		AccNoTimezone:              "Not set",
		AccChangeLanguage:          "🌐 Change Language",
		AccChangeTimezone:          "🌍 Change Timezone",
		TzManualSelect:             "📍 Select Manually",
		TzSelectPrompt:             "Select your timezone:",
	},
	LangUK: {
		Welcome: "Ласкаво просимо до бота-нагадувача!",
		RecurrenceTypes: map[entities.RecurrenceType]string{
			entities.Once:                  "Один раз",
			entities.Daily:                 "Щодня",
			entities.Weekly:                "Щотижня",
			entities.Monthly:               "Щомісяця",
			entities.Interval:              "Інтервал",
			entities.SpacedBasedRepetition: "Інтервал з повторенням",
		},
		BtnBack:                  "↩️ Повернутися до меню",
		BtnCustomTime:            "Свій час",
		MsgSelectTime:            "Оберіть час для нагадувань:",
		MsgSelectHour:            "Оберіть діапазон часу:",
		MsgSelectWithinHour:      "Оберіть час між %02d:00-%02d:00:",
		MsgSelectMessage:         "Оберіть текст нагадування:",
		MsgEnterCustomTime:       "Введіть час у форматі HH:MM (напр., 14:30):",
		MsgEnterCustomMessage:    "Введіть власний текст нагадування:",
		MsgInvalidTimeFormat:     "Неправильний формат часу.",
		MsgInvalidIntervalFormat: "Неправильний формат інтервалу. Очікується 1-7",
		BtnMyReminders:           "Мої нагадування",
		NoReminders:              "У вас ще немає нагадувань.",
		YourReminders:            "Ваші нагадування:\n\n",
		BtnDelete:                "Видалити",
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
		Date:              "Дата",
		Time:              "Час",
		Message:           "Повідомлення",
		ReminderScheduled: "Ваше нагадування заплановано!",
		At:                "в",
		WeekdayNames: map[time.Weekday]string{
			time.Monday:    "Понеділок",
			time.Tuesday:   "Вівторок",
			time.Wednesday: "Середа",
			time.Thursday:  "Четвер",
			time.Friday:    "П'ятниця",
			time.Saturday:  "Субота",
			time.Sunday:    "Неділя",
		},
		WeekdayNamesShort: map[time.Weekday]string{
			time.Monday:    "Пн",
			time.Tuesday:   "Вт",
			time.Wednesday: "Ср",
			time.Thursday:  "Чт",
			time.Friday:    "Пт",
			time.Saturday:  "Сб",
			time.Sunday:    "Нд",
		},
		MsgSelectWeekdays:          "Оберіть дні тижня:",
		MsgSelectTimeWeekly:        "Оберіть час для щотижневих нагадувань:",
		BtnSelect:                  "Обрати",
		MsgSelectDate:              "Оберіть дату:",
		MsgIntervalPrompt:          "Кожні N днів",
		MsgEveryNDays:              "Кожні %d днів",
		MsgEveryNDaysSpaced:        "Кожні %s днів",
		MsgParsingFailed:           "Я не зрозумів. Будь ласка, скористайтеся кнопками меню.",
		MsgTimezoneAutoDetect:      "🌍 Автоматично встановити часовий пояс",
		MsgTimezoneAutoDetectDescr: "Натисніть кнопку, щоб визначити свій часовий пояс.",
		MsgTimezoneSet:             "✅ Часовий пояс встановлено на",
		NavList:                    "Показати нагадування",
		NavSetup:                   "Встановити нагадування",
		NavAccount:                 "Акаунт",
		NavChooseOption:            "Оберіть опцію:",
		CmdStartDesc:               "Запустити бота та показати головне меню",
		CmdListDesc:                "Показати або видалити нагадування",
		CmdSetupDesc:               "Налаштувати час, повторення та параметри нагадувань",
		CmdAccountDesc:             "Управління налаштуваннями акаунту",
		AccTitle:                   "👤 Інформація про рахунок",
		AccUsername:                "Ім'я користувача",
		AccLanguage:                "Мова",
		AccTimezone:                "Часовий пояс",
		AccCreatedAt:               "Створено",
		AccNoUsername:              "Не встановлено",
		AccNoTimezone:              "Не встановлено",
		AccChangeLanguage:          "🌐 Змінити мову",
		AccChangeTimezone:          "🌍 Змінити часовий пояс",
		TzManualSelect:             "📍 Обрати вручну",
		TzSelectPrompt:             "Оберіть свій часовий пояс:",
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
func RecurrenceTypeLabel(lang string, rt entities.RecurrenceType) string {
	s := T(lang)
	if s.RecurrenceTypes != nil {
		if v, ok := s.RecurrenceTypes[rt]; ok {
			return v
		}
	}
	return rt.String()
}
