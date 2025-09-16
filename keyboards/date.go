package keyboards

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

// Callback prefixes for date selection
const (
	CallbackDateStart = "date_start"
	CallbackDateMonth = "date_month:" // YYYY-MM
	CallbackDateDay   = "date_day:"   // YYYY-MM-DD
	CallbackDatePrev  = "date_prev"
	CallbackDateNext  = "date_next"
)

func IsDateSelectionCallback(callbackData string) bool {
	if strings.HasPrefix(callbackData, CallbackDateMonth) {
		return true
	}
	if strings.HasPrefix(callbackData, CallbackDateDay) {
		return true
	}
	return callbackData == CallbackDateStart || callbackData == CallbackDatePrev || callbackData == CallbackDateNext
}

// GetDateSelectionStart shows calendar for the current month
func GetDateSelectionStart(lang string) *SelectionResult {
	now := time.Now()
	return buildCalendar(now.Year(), int(now.Month()), lang)
}

// HandleDateSelection processes date-related callbacks
func HandleDateSelection(callbackData string,
	user *entities.User,
	userSelection *entities.UserSelection) *SelectionResult {
	s := T(user.Language)

	// Determine current displayed month
	now := time.Now()
	year, month := now.Year(), int(now.Month())

	if strings.HasPrefix(callbackData, CallbackDateMonth) {
		// Explicit month provided
		var y, m int
		_, _ = fmt.Sscanf(callbackData[len(CallbackDateMonth):], "%d-%d", &y, &m)
		if y > 0 && m >= 1 && m <= 12 {
			year, month = y, m
		}
		return buildCalendar(year, month, user.Language)
	}

	if callbackData == CallbackDatePrev || callbackData == CallbackDateNext {
		if callbackData == CallbackDatePrev {
			// Do not allow navigating into the past months beyond current month
			prev := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -1, 0)
			if prev.Year() < now.Year() || (prev.Year() == now.Year() && prev.Month() < now.Month()) {
				return buildCalendar(now.Year(), int(now.Month()), user.Language)
			}
			year, month = prev.Year(), int(prev.Month())
		} else {
			next := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0)
			year, month = next.Year(), int(next.Month())
		}
		return buildCalendar(year, month, user.Language)
	}

	if strings.HasPrefix(callbackData, CallbackDateDay) {
		// Set the selected date and move to time picker
		dateStr := callbackData[len(CallbackDateDay):] // YYYY-MM-DD
		userSelection.SelectedDate = dateStr
		return &SelectionResult{Text: s.MsgSelectTime, Markup: GetHourRangeMarkup(user.Language)}
	}

	// Start
	return buildCalendar(year, month, user.Language)
}

func buildCalendar(year int, month int, lang string) *SelectionResult {
	s := T(lang)
	// Header row: Month Year and navigation
	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	title := firstOfMonth.Format("January 2006")

	var rows [][]tgbotapi.InlineKeyboardButton
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("«", CallbackDatePrev),
		tgbotapi.NewInlineKeyboardButtonData(title, fmt.Sprintf("%s%04d-%02d", CallbackDateMonth, year, month)),
		tgbotapi.NewInlineKeyboardButtonData("»", CallbackDateNext),
	))

	// Weekday headers (Mon..Sun) using localized names first letters
	weekdayNames := s.WeekdayNames
	if len(weekdayNames) == 7 {
		var headerButtons []tgbotapi.InlineKeyboardButton
		for _, name := range weekdayNames {
			label := name
			if len(name) > 2 {
				label = name[:2]
			}
			headerButtons = append(headerButtons, tgbotapi.NewInlineKeyboardButtonData(label, CallbackDateMonth))
		}
		rows = append(rows, headerButtons)
	}

	// Calendar grid
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	firstWeekday := int(firstOfMonth.Weekday())
	if firstWeekday == 0 {
		firstWeekday = 7 // Make Sunday=7 to start from Monday=1
	}
	dim := daysInLocal(time.Month(month), year)

	// Fill leading blanks
	day := 1
	if firstWeekday > 1 {
		var blanks []tgbotapi.InlineKeyboardButton
		for i := 1; i < firstWeekday; i++ {
			blanks = append(blanks, tgbotapi.NewInlineKeyboardButtonData(" ", CallbackDateMonth))
		}
		rows = append(rows, blanks)
	}

	// Fill day buttons
	for day <= dim {
		var weekRow []tgbotapi.InlineKeyboardButton
		for len(weekRow) < 7 && day <= dim {
			current := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
			label := fmt.Sprintf("%d", day)
			callback := fmt.Sprintf("%s%04d-%02d-%02d", CallbackDateDay, year, month, day)
			if current.Before(today) {
				// Mark past days as x / disabled
				label = "x"
				callback = CallbackDateMonth
			}
			weekRow = append(weekRow, tgbotapi.NewInlineKeyboardButtonData(label, callback))
			day++
		}
		rows = append(rows, weekRow)
	}

	// Footer with back button
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(s.BtnBack, MainMenu),
	))

	menu := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &SelectionResult{Text: s.MsgSelectDate, Markup: &menu}
}

// daysInLocal returns number of days in the given month/year in local calendar terms
func daysInLocal(month time.Month, year int) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
