package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/services"
	"github.com/ivanenkomaksym/remindme_bot/domain/usecases"
)

// ReminderController handles reminder-related HTTP requests
type ReminderController struct {
	reminderUseCase usecases.ReminderUseCase
	nlpService      services.NLPService
	userUseCase     usecases.UserUseCase
}

// NewReminderController creates a new reminder controller
func NewReminderController(reminderUseCase usecases.ReminderUseCase, nlpService services.NLPService, userUseCase usecases.UserUseCase) *ReminderController {
	return &ReminderController{
		reminderUseCase: reminderUseCase,
		nlpService:      nlpService,
		userUseCase:     userUseCase,
	}
}

// GetUserReminders returns all reminders for a specific user
func (c *ReminderController) GetUserReminders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

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

	reminders, err := c.reminderUseCase.GetUserReminders(userID)
	if err != nil {
		log.Printf("Failed to get user reminders: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reminders)
}

// CreateReminder creates a new reminder for a user
func (c *ReminderController) CreateReminder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

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
}

// CreateReminderFromTextRequest represents the request body for natural language reminder creation
type CreateReminderFromTextRequest struct {
	Text     string `json:"text"`
	Timezone string `json:"timezone,omitempty"`
	Language string `json:"language,omitempty"`
}

// CreateReminderFromText creates a new reminder from natural language text using OpenAI
func (c *ReminderController) CreateReminderFromText(w http.ResponseWriter, r *http.Request) {
	result, user, req := validateCreateReminderFromTextRequest(r, w, c)
	if !result {
		return
	}

	// Use provided timezone/language or fall back to user's settings
	timezone := req.Timezone
	if timezone == "" {
		timezone = user.LocationName
		if timezone == "" {
			timezone = "UTC" // fallback to UTC if no timezone is set
		}
	}
	language := req.Language
	if language == "" {
		language = user.Language
		if language == "" {
			language = "en" // fallback to English
		}
	}

	// Parse the natural language text using NLP service
	userSelection, err := c.nlpService.ParseReminderText(req.Text, timezone, language)
	if err != nil {
		log.Printf("NLP parsing failed for text '%s': %v", req.Text, err)
		http.Error(w, "Failed to parse reminder text: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Create the reminder using the parsed selection
	reminder, err := c.reminderUseCase.CreateReminder(user.ID, userSelection)
	if err != nil {
		log.Printf("Failed to create reminder from NLP: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create reminder from text: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reminder)
}

func validateCreateReminderFromTextRequest(r *http.Request, w http.ResponseWriter, c *ReminderController) (bool, *entities.User, *CreateReminderFromTextRequest) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return false, nil, nil
	}

	userIDStr := r.PathValue("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id parameter is required", http.StatusBadRequest)
		return false, nil, nil
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return false, nil, nil
	}

	var req CreateReminderFromTextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return false, nil, nil
	}

	if req.Text == "" {
		http.Error(w, "text field is required", http.StatusBadRequest)
		return false, nil, nil
	}

	user, err := c.userUseCase.GetUser(userID)
	if err != nil {
		log.Printf("Failed to get user %d: %v", userID, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return false, nil, nil
	}

	return true, user, &req
}

// GetReminder returns a specific reminder by ID
func (c *ReminderController) GetReminder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.PathValue("user_id")
	reminderIDStr := r.PathValue("reminder_id")

	if userIDStr == "" || reminderIDStr == "" {
		http.Error(w, "user_id and reminder_id parameters are required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	reminderID, err := strconv.ParseInt(reminderIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid reminder_id", http.StatusBadRequest)
		return
	}

	reminder, err := c.reminderUseCase.GetReminder(userID, reminderID)
	if err != nil {
		log.Printf("Failed to get reminder: %v", err)
		http.Error(w, "Internal Server Error", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reminder)
}

// UpdateReminder updates a specific reminder
func (c *ReminderController) UpdateReminder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Only PUT requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.PathValue("user_id")
	reminderIDStr := r.PathValue("reminder_id")

	if userIDStr == "" || reminderIDStr == "" {
		http.Error(w, "user_id and reminder_id parameters are required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	reminderID, err := strconv.ParseInt(reminderIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid reminder_id", http.StatusBadRequest)
		return
	}

	var reminder entities.Reminder
	if err := json.NewDecoder(r.Body).Decode(&reminder); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedReminder, err := c.reminderUseCase.UpdateReminder(userID, reminderID, &reminder)
	if err != nil {
		log.Printf("Failed to update reminder: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedReminder)
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
		log.Printf("Failed to delete reminder: %v. Not found", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
