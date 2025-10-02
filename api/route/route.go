package route

import (
	"log"
	"net/http"

	"github.com/ivanenkomaksym/remindme_bot/api/middleware"
	"github.com/ivanenkomaksym/remindme_bot/bootstrap"
	"github.com/ivanenkomaksym/remindme_bot/notifier"
)

func Setup(app *bootstrap.Application) {
	addr := app.Env.Config.GetServerAddress()

	// Create a new HTTP multiplexer
	mux := http.NewServeMux()

	// Apply middleware
	handler := middleware.LoggingMiddleware(mux)
	handler = middleware.CORSMiddleware(handler)
	handler = middleware.APIKeyMiddleware(app, handler)
	handler = middleware.RecoveryMiddleware(handler)

	// Define the webhook endpoint that Telegram will send updates to
	mux.HandleFunc("/telegram-webhook", app.Container.BotController.HandleWebhook)

	// Public endpoints for timezone capture
	mux.HandleFunc("/set-timezone", app.Container.TimezoneController.ServePage)
	mux.HandleFunc("/set-timezone/callback", app.Container.TimezoneController.Callback)

	// API endpoints
	mux.HandleFunc("/api/users", app.Container.UserController.GetUsers)
	mux.HandleFunc("/api/users/{user_id}", app.Container.UserController.GetUser)
	mux.HandleFunc("/api/users/{user_id}/language", app.Container.UserController.UpdateUserLanguage)
	mux.HandleFunc("/api/users/{user_id}/selection", app.Container.UserController.GetUserSelection)
	mux.HandleFunc("/api/users/{user_id}/selection/clear", app.Container.UserController.ClearUserSelection)

	mux.HandleFunc("/api/reminders", app.Container.ReminderController.GetAllReminders)
	mux.HandleFunc("/api/reminders/{user_id}", app.Container.ReminderController.ProcessUserReminders)
	mux.HandleFunc("/api/reminders/{user_id}/active", app.Container.ReminderController.GetActiveReminders)
	mux.HandleFunc("/api/reminders/{user_id}/delete/{reminder_id}", app.Container.ReminderController.DeleteReminder)

	// Add a health check endpoint for Cloud Run
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	go notifier.StartReminderNotifier(app.Container.ReminderRepo, app.Env.Config.App, app.Bot)

	log.Printf("Starting HTTP server on %s", addr)
	// Start the HTTP server. This will block indefinitely, serving requests.
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
