package datepicker

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	cmdPrevMonth = iota
	cmdNextMonth
	cmdPrevYears
	cmdNextYears
	cmdCancel
	cmdBack
	cmdMonthClick
	cmdYearClick
	cmdNop

	cmdDayClick
	cmdSelectMonth
	cmdSelectYear
)

func (d *DatePicker) callback(bot *tgbotapi.BotAPI, upd *tgbotapi.Update) {
	st := d.decodeState(upd.CallbackQuery.Data)

	switch st.cmd {
	case cmdYearClick:
		d.year = st.param
		d.showMain(bot, upd)
	case cmdMonthClick:
		d.month = time.Month(st.param)
		d.showMain(bot, upd)
	case cmdDayClick:
		// Prevent selecting dates before dp.from (including past when ForbidPastDates is used)
		selected := time.Date(d.year, d.month, st.param, 0, 0, 0, 0, time.Local)
		if !d.from.IsZero() && selected.Before(d.from) {
			// ignore selection, just refresh main view to indicate no-op
			d.showMain(bot, upd)
			break
		}
		if !d.to.IsZero() && selected.After(d.to) {
			d.showMain(bot, upd)
			break
		}
		d.onSelect(bot, upd.CallbackQuery.Message, selected)
	case cmdCancel:
		/*if d.deleteOnCancel {
			deleteCfg := tgbotapi.NewDeleteMessage(upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID)
			if _, err := bot.Request(deleteCfg); err != nil {
				d.onError(fmt.Errorf("failed to delete message onCancel: %w", err))
			}
		}*/
		d.onCancel(bot, upd.CallbackQuery.Message)
	case cmdPrevYears:
		d.showSelectYear(bot, upd, st.param)
	case cmdNextYears:
		d.showSelectYear(bot, upd, st.param)
	case cmdPrevMonth:
		d.month--
		if d.month == 0 {
			d.month = 12
			d.year--
		}
		d.showMain(bot, upd)
	case cmdNextMonth:
		d.month++
		if d.month > 12 {
			d.month = 1
			d.year++
		}
		d.showMain(bot, upd)
	case cmdBack:
		d.showMain(bot, upd)
	case cmdSelectMonth:
		d.showSelectMonth(bot, upd)
	case cmdSelectYear:
		d.showSelectYear(bot, upd, d.year)
	case cmdNop:
		// do nothing
	default:
		d.onError(fmt.Errorf("unknown command: %d", st.cmd))
	}

	// answer callback
	_ = d.answerCallback(bot, upd.CallbackQuery)
}

func (d *DatePicker) showSelectMonth(bot *tgbotapi.BotAPI, upd *tgbotapi.Update) {
	edit := tgbotapi.NewEditMessageReplyMarkup(upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID, tgbotapi.InlineKeyboardMarkup{InlineKeyboard: d.buildMonthKeyboard()})
	if _, err := bot.Request(edit); err != nil {
		d.onError(fmt.Errorf("error edit message in showSelectMonth, %w", err))
	}
}

func (d *DatePicker) showSelectYear(bot *tgbotapi.BotAPI, upd *tgbotapi.Update, currentYear int) {
	edit := tgbotapi.NewEditMessageReplyMarkup(upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID, tgbotapi.InlineKeyboardMarkup{InlineKeyboard: d.buildYearKeyboard(currentYear)})
	if _, err := bot.Request(edit); err != nil {
		d.onError(fmt.Errorf("error edit message in showSelectYear, %w", err))
	}
}

func (d *DatePicker) showMain(bot *tgbotapi.BotAPI, upd *tgbotapi.Update) {
	edit := tgbotapi.NewEditMessageReplyMarkup(upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID, d.Keyboard())
	if _, err := bot.Request(edit); err != nil {
		d.onError(fmt.Errorf("error edit message in showMain, %w", err))
	}
}

func (d *DatePicker) answerCallback(bot *tgbotapi.BotAPI, cq *tgbotapi.CallbackQuery) error {
	cfg := tgbotapi.NewCallback(cq.ID, "")
	_, err := bot.Request(cfg)
	return err
}
