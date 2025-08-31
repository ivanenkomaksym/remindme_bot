// bootstrap telegram bot
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/ivanenkomaksym/offerforyou_bot/internal/models"
	"github.com/joho/godotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	bot *tgbotapi.BotAPI

	welcomeMessage = "Welcome to the Reminder Bot!"
)

func buildMainMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(models.Daily.String(), models.Daily.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(models.Weekly.String(), models.Weekly.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(models.Monthly.String(), models.Monthly.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(models.Interval.String(), models.Interval.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(models.Custom.String(), models.Custom.String()),
		),
	)
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the bot
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN is not set in .env file")
	}
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	bot.Debug = true

	// Get the WEBHOOK_URL from environment variables
	// This will be the URL of your deployed Cloud Run service + the webhook path
	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		log.Fatal("WEBHOOK_URL environment variable not set. This should be your Cloud Run service URL including the path (e.g., https://<service-url>/telegram-webhook).")
	}

	wh, err := tgbotapi.NewWebhook(webhookURL)
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

	// Get the port from environment variables, default to 8080 for Cloud Run
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	// Create a new HTTP multiplexer
	mux := http.NewServeMux()

	// Define the webhook endpoint that Telegram will send updates to
	// It's good practice to use a non-root path for webhooks.
	mux.HandleFunc("/telegram-webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		var update tgbotapi.Update
		// Decode the JSON request body into a tgbotapi.Update struct
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			log.Printf("ERROR: Could not decode update: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Process the update in a non-blocking way if possible,
		// but for simple bots, direct handling is fine.
		handleUpdate(update)

		// Respond with 200 OK to Telegram immediately
		// This acknowledges receipt of the update and prevents Telegram from retrying.
		w.WriteHeader(http.StatusOK)
	})

	// Add a health check endpoint for Cloud Run
	// Cloud Run sends requests to the root path by default for health checks.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Starting HTTP server on %s", addr)
	// Start the HTTP server. This will block indefinitely, serving requests.
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}

func handleUpdate(update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		user := update.CallbackQuery.From
		text := update.CallbackQuery.Data

		log.Printf("'[%s] %s %s' selected '%s'", user.UserName, user.FirstName, user.LastName, text)
	}

	if update.Message != nil && update.Message.IsCommand() {
		user := update.Message.From
		log.Printf("'[%s] %s %s' started chat", user.UserName, user.FirstName, user.LastName)

		if update.Message.Command() == "start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, welcomeMessage)
			msg.ParseMode = "HTML"
			msg.ReplyMarkup = buildMainMenu()
			bot.Send(msg)
		}
	}
}
