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

	// API endpoints - Users
	mux.HandleFunc("GET /api/users", app.Container.UserController.GetUsers)
	mux.HandleFunc("POST /api/users", app.Container.UserController.CreateUser)
	mux.HandleFunc("GET /api/users/{user_id}", app.Container.UserController.GetUser)
	mux.HandleFunc("DELETE /api/users/{user_id}", app.Container.UserController.DeleteUser)
	mux.HandleFunc("PUT /api/users/{user_id}/language", app.Container.UserController.UpdateUserLanguage)
	mux.HandleFunc("PUT /api/users/{user_id}/location", app.Container.UserController.UpdateUserLocation)
	mux.HandleFunc("GET /api/users/{user_id}/selection", app.Container.UserController.GetUserSelection)
	mux.HandleFunc("DELETE /api/users/{user_id}/selection", app.Container.UserController.ClearUserSelection)

	// API endpoints - Reminders
	mux.HandleFunc("GET /api/reminders", app.Container.ReminderController.GetAllReminders)
	mux.HandleFunc("GET /api/reminders/{user_id}", app.Container.ReminderController.GetUserReminders)
	mux.HandleFunc("POST /api/reminders/{user_id}", app.Container.ReminderController.CreateReminder)
	mux.HandleFunc("POST /api/reminders/{user_id}/from-text", app.Container.ReminderController.CreateReminderFromText)
	mux.HandleFunc("GET /api/reminders/{user_id}/{reminder_id}", app.Container.ReminderController.GetReminder)
	mux.HandleFunc("PUT /api/reminders/{user_id}/{reminder_id}", app.Container.ReminderController.UpdateReminder)
	mux.HandleFunc("DELETE /api/reminders/{user_id}/{reminder_id}", app.Container.ReminderController.DeleteReminder)
	mux.HandleFunc("GET /api/reminders/{user_id}/active", app.Container.ReminderController.GetActiveReminders)

	// Add a health check endpoint for Cloud Run
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	if app.Env.Config.Bot.Enabled {
		go notifier.StartReminderNotifier(app.Container.ReminderRepo, app.Env.Config.App, app.Bot)
	}

	log.Printf("Starting HTTP server on %s", addr)
	// Start the HTTP server. This will block indefinitely, serving requests.
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
