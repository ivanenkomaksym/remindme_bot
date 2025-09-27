package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ivanenkomaksym/remindme_bot/config"
	"github.com/ivanenkomaksym/remindme_bot/domain/usecases"
	"github.com/ivanenkomaksym/remindme_bot/keyboards"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BotController handles bot-related HTTP requests
type BotController struct {
	botUseCase  usecases.BotUseCase
	userUsecase usecases.UserUseCase
	dateUseCase usecases.DateUseCase
	config      config.Config
	bot         *tgbotapi.BotAPI
}

// NewBotController creates a new bot controller
func NewBotController(botUseCase usecases.BotUseCase, userUsecase usecases.UserUseCase, dateUseCase usecases.DateUseCase, config config.Config, bot *tgbotapi.BotAPI) *BotController {
	return &BotController{
		botUseCase:  botUseCase,
		userUsecase: userUsecase,
		dateUseCase: dateUseCase,
		config:      config,
		bot:         bot,
	}
}

// HandleWebhook handles incoming webhook requests from Telegram
func (c *BotController) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("ERROR: Could not decode update: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Process the update
	if err := c.processUpdate(update); err != nil {
		log.Printf("ERROR: Failed to process update: %v", err)
	}

	// Respond with 200 OK to Telegram immediately
	w.WriteHeader(http.StatusOK)
}

// processUpdate processes a Telegram update
func (c *BotController) processUpdate(update tgbotapi.Update) error {
	if update.CallbackQuery != nil {
		return c.processCallbackQuery(update.CallbackQuery)
	}

	if update.Message != nil {
		return c.processMessage(update.Message)
	}

	return nil
}

// processCallbackQuery processes callback queries (button presses)
func (c *BotController) processCallbackQuery(callbackQuery *tgbotapi.CallbackQuery) error {
	msg := tgbotapi.NewEditMessageText(
		callbackQuery.Message.Chat.ID,
		callbackQuery.Message.MessageID,
		"", // Text will be set later
	)

	response, err := c.botUseCase.ProcessKeyboardSelection(callbackQuery)
	if err != nil {
		log.Printf("Failed to process callback query: %v", err)
		return err
	}

	if response != nil {
		log.Printf("Callback response: %+v", response)
		msg.Text = response.Text
		msg.ReplyMarkup = response.Markup
		msg.ParseMode = "HTML"
		c.bot.Send(msg)
	}

	return nil
}

// processMessage processes text messages
func (c *BotController) processMessage(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(
		message.Chat.ID,
		"", // Text will be set later
	)

	user := message.From

	// Create or get user
	userEntity, err := c.userUsecase.GetOrCreateUser(
		user.ID,
		user.UserName,
		user.FirstName,
		user.LastName,
		"",
	)
	if err != nil {
		return err
	}

	var response *keyboards.SelectionResult

	// Ask to set timezone if not set yet.
	if userEntity.GetLocation() == nil {
		url := fmt.Sprintf("%s/set-timezone?user_id=%d", c.config.Bot.PublicURL, user.ID)
		response, err = c.botUseCase.ProcessTimezone(userEntity, url)
		if err != nil {
			log.Printf("Failed to process message: %v", err)
			return err
		}

		return c.processResponse(msg, response)
	}

	response, err = c.botUseCase.ProcessUserInput(message)
	if err != nil {
		log.Printf("Failed to process message: %v", err)
		return err
	}

	return c.processResponse(msg, response)
}

func (c *BotController) processResponse(msg tgbotapi.MessageConfig, response *keyboards.SelectionResult) error {
	log.Printf("Message response: %+v", response)
	msg.Text = response.Text
	msg.ReplyMarkup = response.Markup
	msg.ParseMode = "HTML"
	c.bot.Send(msg)

	return nil
}
