package bootstrap

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewBot(env *Env) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(env.Config.Bot.Token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	bot.Debug = env.Config.Bot.Debug

	wh, err := tgbotapi.NewWebhook(env.Config.Bot.WebhookURL)
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

	return bot
}
