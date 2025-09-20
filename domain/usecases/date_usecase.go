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
	HandleDatepickerSelection(user *entities.User, userSelection *entities.UserSelection) *keyboards.SelectionResult
	HandleDatepickerCallback(callbackQuery *tgbotapi.CallbackQuery) bool
}

type dateUseCase struct {
	userUseCase UserUseCase
	datepickers map[int64]*datepicker.DatePicker // userID -> datepicker instance
	bot         *tgbotapi.BotAPI
}

// NewBotUseCase creates a new bot use case
func NewDateUseCase(userUseCase UserUseCase, bot *tgbotapi.BotAPI) DateUseCase {
	return &dateUseCase{
		userUseCase: userUseCase,
		bot:         bot,
	}
}

func (d *dateUseCase) HandleDatepickerSelection(user *entities.User,
	userSelection *entities.UserSelection) *keyboards.SelectionResult {

	onSelect := func(bot *tgbotapi.BotAPI, m *tgbotapi.Message, date time.Time) {
		userSelection.SelectedDate = date.Format("2006-01-02")
		err := d.userUseCase.UpdateUserSelection(user.ID, userSelection)
		if err != nil {
			log.Printf("Failed to update user selection: %v", err)
		}

		msg := tgbotapi.NewMessage(m.Chat.ID, "Selected date: "+date.Format("2006-01-02"))
		bot.Send(msg)

		delete(d.datepickers, user.ID)
	}

	onCancel := func(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
		// Remove datepicker from active datepickers on cancel
		delete(d.datepickers, user.ID)
	}

	dp := datepicker.New(onSelect,
		datepicker.Language(user.Language),
		datepicker.OnCancel(onCancel))

	// Store datepicker instance for this user
	d.datepickers[user.ID] = dp

	return &keyboards.SelectionResult{}
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
