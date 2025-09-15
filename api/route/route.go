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
	handler = middleware.RecoveryMiddleware(handler)

	// Define the webhook endpoint that Telegram will send updates to
	mux.HandleFunc("/telegram-webhook", app.Container.BotController.HandleWebhook)

	// API endpoints
	mux.HandleFunc("/api/users", app.Container.UserController.GetUser)
	mux.HandleFunc("/api/users/language", app.Container.UserController.UpdateUserLanguage)
	mux.HandleFunc("/api/users/selection", app.Container.UserController.GetUserSelection)
	mux.HandleFunc("/api/users/selection/clear", app.Container.UserController.ClearUserSelection)

	mux.HandleFunc("/api/reminders", app.Container.ReminderController.GetUserReminders)
	mux.HandleFunc("/api/reminders/all", app.Container.ReminderController.GetAllReminders)
	mux.HandleFunc("/api/reminders/active", app.Container.ReminderController.GetActiveReminders)
	mux.HandleFunc("/api/reminders/delete", app.Container.ReminderController.DeleteReminder)

	// Add a health check endpoint for Cloud Run
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	go notifier.StartReminderNotifier(app.Container.ReminderRepo, app.Bot)

	log.Printf("Starting HTTP server on %s", addr)
	// Start the HTTP server. This will block indefinitely, serving requests.
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
