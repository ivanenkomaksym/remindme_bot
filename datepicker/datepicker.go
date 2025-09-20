package datepicker

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DatesMode int

const (
	DateModeExclude DatesMode = iota
	DateModeInclude
)

type OnSelectHandler func(bot *tgbotapi.BotAPI, mes *tgbotapi.Message, date time.Time)
type OnCancelHandler func(bot *tgbotapi.BotAPI, mes *tgbotapi.Message)
type OnErrorHandler func(err error)

type DatePicker struct {
	// configurable params
	startFromSunday bool
	language        string
	langs           LangsData
	deleteOnSelect  bool
	deleteOnCancel  bool
	from            time.Time
	to              time.Time
	dates           []time.Time
	datesMode       DatesMode
	onSelect        OnSelectHandler
	onCancel        OnCancelHandler
	onError         OnErrorHandler

	// current date
	month time.Month
	year  int

	// internal
	prefix string
}

func New(onSelect OnSelectHandler, opts ...Option) *DatePicker {
	year, month, _ := time.Now().Date()

	dp := &DatePicker{
		language:       "en",
		langs:          loadLangs(),
		deleteOnSelect: true,
		deleteOnCancel: true,
		onSelect:       onSelect,
		onCancel:       defaultOnCancel,
		onError:        defaultOnError,
		month:          month,
		year:           year,
		prefix:         randomString(16),
	}

	for _, opt := range opts {
		opt(dp)
	}

	return dp
}

// Prefix returns the callback data prefix of the widget
func (d *DatePicker) Prefix() string {
	return d.prefix
}

// Keyboard returns the inline keyboard markup for the current state
func (d *DatePicker) Keyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.InlineKeyboardMarkup{InlineKeyboard: d.buildKeyboard()}
}

// MarshalJSON enables embedding directly into messages in this repo's helpers
func (d *DatePicker) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Keyboard())
}

// HandleCallback processes a callback update. It returns true if the update belonged to this widget.
func (d *DatePicker) HandleCallback(bot *tgbotapi.BotAPI, upd *tgbotapi.Update) bool {
	if upd == nil || upd.CallbackQuery == nil || upd.CallbackQuery.Data == "" {
		return false
	}
	if !strings.HasPrefix(upd.CallbackQuery.Data, d.prefix) {
		return false
	}

	d.callback(bot, upd)
	return true
}

func defaultOnError(err error) {
	log.Printf("[TG-UI-DATEPICKER-TBA] [ERROR] %s", err)
}

func defaultOnCancel(_ *tgbotapi.BotAPI, _ *tgbotapi.Message) {}
