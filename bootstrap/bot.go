package bootstrap

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/keyboards"
)

func NewBot(env *Env) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(env.Config.Bot.Token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	bot.Debug = env.Config.Bot.Debug

	wh, err := tgbotapi.NewWebhook(fmt.Sprintf("%s/telegram-webhook", env.Config.Bot.PublicURL))
	if err != nil {
		log.Fatalf("Failed to create webhook config: %v", err)
	}
	// Use bot.Request to send the WebhookConfig to Telegram
	_, err = bot.Request(wh)
	if err != nil {
		log.Fatalf("Failed to set webhook: %v", err)
	}

	// Get webhook info to confirm it's set and check for any errors from Telegram's side
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatalf("Failed to get webhook info: %v", err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram last webhook error: %s", info.LastErrorMessage)
	}
	log.Printf("Webhook set to: %s (pending: %t)", info.URL, info.PendingUpdateCount > 0)

	// Set up bot commands menu
	err = setupBotCommands(bot)
	if err != nil {
		log.Printf("Failed to setup bot commands: %v", err)
	}

	return bot
}

func setupBotCommands(bot *tgbotapi.BotAPI) error {
	// Set up commands for English (default)
	err := setupCommandsForLanguage(bot, keyboards.LangEN)
	if err != nil {
		return fmt.Errorf("failed to setup English commands: %w", err)
	}

	// Set up commands for Ukrainian
	err = setupCommandsForLanguage(bot, keyboards.LangUK)
	if err != nil {
		return fmt.Errorf("failed to setup Ukrainian commands: %w", err)
	}

	return nil
}

func setupCommandsForLanguage(bot *tgbotapi.BotAPI, langCode string) error {
	s := keyboards.T(langCode)

	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: s.CmdStartDesc},
		{Command: "list", Description: s.CmdListDesc},
		{Command: "setup", Description: s.CmdSetupDesc},
		{Command: "account", Description: s.CmdAccountDesc},
	}

	config := tgbotapi.SetMyCommandsConfig{
		Commands:     commands,
		LanguageCode: langCode,
	}

	_, err := bot.Request(config)
	return err
}
