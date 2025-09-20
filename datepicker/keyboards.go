package datepicker

import (
	"fmt"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (d *DatePicker) buildYearKeyboard(currentYear int) [][]tgbotapi.InlineKeyboardButton {
	var keyboard [][]tgbotapi.InlineKeyboardButton

	var topRow []tgbotapi.InlineKeyboardButton
	if d.from.IsZero() || d.from.Year() < currentYear-12 {
		topRow = append(topRow, tgbotapi.NewInlineKeyboardButtonData("\u2190 "+d.lang("Prev"), d.encodeState(state{cmd: cmdPrevYears, param: currentYear - 25})))
	}
	if d.to.IsZero() || d.to.Year() > currentYear+12 {
		topRow = append(topRow, tgbotapi.NewInlineKeyboardButtonData(d.lang("Next")+" \u2192", d.encodeState(state{cmd: cmdNextYears, param: currentYear + 25})))
	}
	if len(topRow) > 0 {
		keyboard = append(keyboard, topRow)
	}

	var row []tgbotapi.InlineKeyboardButton
	idx := 1
	for i := currentYear - 12; i <= currentYear+12; i++ {
		if idx > 5 {
			idx = 1
			keyboard = append(keyboard, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}

		yearText := strconv.Itoa(i)
		if !d.from.IsZero() && i < d.from.Year() {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData("-", d.encodeState(state{cmd: cmdNop})))
			idx++
			continue
		}
		if !d.to.IsZero() && i > d.to.Year() {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData("-", d.encodeState(state{cmd: cmdNop})))
			idx++
			continue
		}

		row = append(row, tgbotapi.NewInlineKeyboardButtonData(yearText, d.encodeState(state{cmd: cmdYearClick, param: i})))
		idx++
	}
	keyboard = append(keyboard, row)

	keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(d.lang("Back"), d.encodeState(state{cmd: cmdBack})),
	})

	return keyboard
}

func (d *DatePicker) buildMonthKeyboard() [][]tgbotapi.InlineKeyboardButton {
	var keyboard [][]tgbotapi.InlineKeyboardButton

	var row []tgbotapi.InlineKeyboardButton
	for i := 0; i < 12; i++ {
		if i > 0 && i%4 == 0 {
			keyboard = append(keyboard, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
		startOfMonth := time.Date(d.year, time.Month(i+1), 1, 0, 0, 0, 0, time.Local)
		if !d.from.IsZero() && startOfMonth.Before(d.from.AddDate(0, -1, 0)) {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData("-", d.encodeState(state{cmd: cmdNop})))
			continue
		}
		if !d.to.IsZero() && startOfMonth.After(d.to) {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData("-", d.encodeState(state{cmd: cmdNop})))
			continue
		}
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(d.lang(time.Month(i+1).String()), d.encodeState(state{cmd: cmdMonthClick, param: i + 1})))
	}
	keyboard = append(keyboard, row)
	keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(d.lang("Back"), d.encodeState(state{cmd: cmdBack}))})
	return keyboard
}

func (d *DatePicker) buildKeyboard() [][]tgbotapi.InlineKeyboardButton {
	var data [][]tgbotapi.InlineKeyboardButton

	top := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(d.lang(d.month.String()), d.encodeState(state{cmd: cmdSelectMonth})),
		tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(d.year), d.encodeState(state{cmd: cmdSelectYear})),
	}
	data = append(data, top)

	dow := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(d.lang("Monday"), d.encodeState(state{cmd: cmdNop})),
		tgbotapi.NewInlineKeyboardButtonData(d.lang("Tuesday"), d.encodeState(state{cmd: cmdNop})),
		tgbotapi.NewInlineKeyboardButtonData(d.lang("Wednesday"), d.encodeState(state{cmd: cmdNop})),
		tgbotapi.NewInlineKeyboardButtonData(d.lang("Thursday"), d.encodeState(state{cmd: cmdNop})),
		tgbotapi.NewInlineKeyboardButtonData(d.lang("Friday"), d.encodeState(state{cmd: cmdNop})),
		tgbotapi.NewInlineKeyboardButtonData(d.lang("Saturday"), d.encodeState(state{cmd: cmdNop})),
		tgbotapi.NewInlineKeyboardButtonData(d.lang("Sunday"), d.encodeState(state{cmd: cmdNop})),
	}
	if d.startFromSunday {
		dow = append([]tgbotapi.InlineKeyboardButton{dow[len(dow)-1]}, dow[:len(dow)-1]...)
	}
	data = append(data, dow)

	startOfMonth := time.Date(d.year, d.month, 1, 0, 0, 0, 0, time.Local)
	var row []tgbotapi.InlineKeyboardButton

	skipFirst := int(startOfMonth.Weekday())
	if !d.startFromSunday {
		if skipFirst == 0 {
			skipFirst = 6
		} else {
			skipFirst--
		}
	}

	idx := skipFirst + 1
	for i := 0; i < skipFirst; i++ {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(" ", d.encodeState(state{cmd: cmdNop})))
	}

	for i := 1; i <= startOfMonth.AddDate(0, 1, -1).Day(); i++ {
		if idx == 8 {
			idx = 1
			data = append(data, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}

		now := time.Date(d.year, d.month, i, 0, 0, 0, 0, time.Local)
		if !d.from.IsZero() && now.Before(d.from) {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData("-", d.encodeState(state{cmd: cmdNop})))
			idx++
			continue
		}
		if !d.to.IsZero() && now.After(d.to) {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData("-", d.encodeState(state{cmd: cmdNop})))
			idx++
			continue
		}

		var isInDates bool
		for _, dte := range d.dates {
			isInDates = dte.Day() == i && dte.Month() == d.month && dte.Year() == d.year
			if isInDates {
				break
			}
		}
		if len(d.dates) > 0 && ((isInDates && d.datesMode == DateModeExclude) || (!isInDates && d.datesMode == DateModeInclude)) {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData("-", d.encodeState(state{cmd: cmdNop})))
			idx++
			continue
		}

		row = append(row, tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(i), d.encodeState(state{cmd: cmdDayClick, param: i})))
		idx++
	}

	if idx < 8 {
		for i := 0; i < 8-idx; i++ {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(" ", d.encodeState(state{cmd: cmdNop})))
		}
	}
	data = append(data, row)

	var bottom []tgbotapi.InlineKeyboardButton
	if d.from.IsZero() || d.from.Before(startOfMonth) {
		prevMonth := startOfMonth.AddDate(0, -1, 0)
		bottom = append(bottom, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("\u2190 %s %d", d.lang(prevMonth.Month().String()), prevMonth.Year()), d.encodeState(state{cmd: cmdPrevMonth})))
	}
	bottom = append(bottom, tgbotapi.NewInlineKeyboardButtonData(d.lang("Cancel"), d.encodeState(state{cmd: cmdCancel})))
	if d.to.IsZero() || d.to.After(startOfMonth.AddDate(0, 1, -1)) {
		nextMonth := startOfMonth.AddDate(0, 1, 0)
		bottom = append(bottom, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s %d \u2192", d.lang(nextMonth.Month().String()), nextMonth.Year()), d.encodeState(state{cmd: cmdNextMonth})))
	}
	data = append(data, bottom)

	return data
}
