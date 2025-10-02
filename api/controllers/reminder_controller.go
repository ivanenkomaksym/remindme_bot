package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/usecases"
)

// ReminderController handles reminder-related HTTP requests
type ReminderController struct {
	reminderUseCase usecases.ReminderUseCase
}

// NewReminderController creates a new reminder controller
func NewReminderController(reminderUseCase usecases.ReminderUseCase) *ReminderController {
	return &ReminderController{
		reminderUseCase: reminderUseCase,
	}
}

// ProcessUserReminders returns all reminders for a specific user or creates a new reminder
func (c *ReminderController) ProcessUserReminders(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from path
	userIDStr := r.PathValue("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id parameter is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPost {
		var userSelection entities.UserSelection
		if err := json.NewDecoder(r.Body).Decode(&userSelection); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		reminder, err := c.reminderUseCase.CreateReminder(userID, &userSelection)
		if err != nil {
			log.Printf("Failed to create reminder: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reminder)
		return
	}

	if r.Method == http.MethodGet {
		reminders, err := c.reminderUseCase.GetUserReminders(userID)
		if err != nil {
			log.Printf("Failed to get user reminders: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reminders)
		return
	}

	http.Error(w, "Only GET and POST requests are allowed", http.StatusMethodNotAllowed)
	return
}

// GetAllReminders returns all reminders (admin endpoint)
func (c *ReminderController) GetAllReminders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	reminders, err := c.reminderUseCase.GetAllReminders()
	if err != nil {
		log.Printf("Failed to get all reminders: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reminders)
}

// DeleteReminder deletes a specific reminder
func (c *ReminderController) DeleteReminder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract reminder ID and user ID from path
	reminderIDStr := r.PathValue("reminder_id")
	userIDStr := r.PathValue("user_id")

	if reminderIDStr == "" || userIDStr == "" {
		http.Error(w, "reminder_id and user_id parameters are required", http.StatusBadRequest)
		return
	}

	reminderID, err := strconv.ParseInt(reminderIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid reminder_id", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	err = c.reminderUseCase.DeleteReminder(reminderID, userID)
	if err != nil {
		log.Printf("Failed to delete reminder: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reminder deleted successfully"})
}

// GetActiveReminders returns all active reminders
func (c *ReminderController) GetActiveReminders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	reminders, err := c.reminderUseCase.GetActiveReminders()
	if err != nil {
		log.Printf("Failed to get active reminders: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reminders)
}
