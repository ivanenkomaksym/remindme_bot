package usecases

import (
	"log"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/datepicker"
	"github.com/ivanenkomaksym/remindme_bot/keyboards"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

type DateUseCase interface {
	CreateDatepicker(message *tgbotapi.Message, user *entities.User, userSelection *entities.UserSelection) *keyboards.SelectionResult
	HandleDatepickerCallback(callbackQuery *tgbotapi.CallbackQuery) bool
}

type dateUseCase struct {
	userUseCase UserUseCase
	datepickers map[int64]*datepicker.DatePicker // userID -> datepicker instance
	bot         *tgbotapi.BotAPI
}

// NewDateUseCase creates a new date use case
func NewDateUseCase(userUseCase UserUseCase, bot *tgbotapi.BotAPI) DateUseCase {
	return &dateUseCase{
		userUseCase: userUseCase,
		bot:         bot,
		datepickers: make(map[int64]*datepicker.DatePicker),
	}
}

func (d *dateUseCase) CreateDatepicker(message *tgbotapi.Message,
	user *entities.User,
	userSelection *entities.UserSelection) *keyboards.SelectionResult {

	onSelect := func(bot *tgbotapi.BotAPI, m *tgbotapi.Message, date time.Time) {
		userSelection.SetSelectedDate(date)
		err := d.userUseCase.UpdateUserSelection(user.ID, userSelection)
		if err != nil {
			log.Printf("Failed to update user selection: %v", err)
		}

		s := keyboards.T(user.Language)
		text := s.MsgSelectTime
		markup := keyboards.GetHourRangeMarkup(user.Language)

		msg := tgbotapi.NewEditMessageText(
			user.ID,
			message.MessageID,
			"", // Text will be set later
		)

		msg.Text = text
		msg.ReplyMarkup = markup
		msg.ParseMode = "HTML"
		d.bot.Send(msg)

		delete(d.datepickers, user.ID)
	}

	onCancel := func(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
		// Remove datepicker from active datepickers on cancel
		delete(d.datepickers, user.ID)

		// Navigate back to setup menu
		s := keyboards.T(user.Language)
		text := "⚙️ " + s.NavSetup + ":"
		markup := keyboards.GetSetupMenuMarkup(user.Language)

		msg := tgbotapi.NewEditMessageText(
			user.ID,
			message.MessageID,
			text,
		)
		msg.ReplyMarkup = markup
		msg.ParseMode = "HTML"
		d.bot.Send(msg)
	}

	dp := datepicker.New(onSelect,
		datepicker.Language(user.Language),
		datepicker.OnCancel(onCancel),
		datepicker.WithCallbackPrefix("datepicker_selection"),
		datepicker.ForbidPastDates())

	// Store datepicker instance for this user
	d.datepickers[user.ID] = dp

	s := keyboards.T(user.Language)
	text := s.MsgSelectDate
	markup := dp.Keyboard()

	return &keyboards.SelectionResult{Text: text, Markup: &markup}
}

// HandleDatepickerCallback handles datepicker-specific callbacks
func (d *dateUseCase) HandleDatepickerCallback(callbackQuery *tgbotapi.CallbackQuery) bool {
	userID := callbackQuery.From.ID

	// Check if user has an active datepicker
	datepicker, exists := d.datepickers[userID]
	if !exists {
		return false
	}

	// Create update from callback query
	update := &tgbotapi.Update{
		CallbackQuery: callbackQuery,
	}

	// Let the datepicker handle the callback
	return datepicker.HandleCallback(d.bot, update)
}
