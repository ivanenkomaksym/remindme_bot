package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ivanenkomaksym/remindme_bot/domain/errors"
	"github.com/ivanenkomaksym/remindme_bot/domain/usecases"
)

// UserController handles user-related HTTP requests
type UserController struct {
	userUseCase usecases.UserUseCase
}

// NewUserController creates a new user controller
func NewUserController(userUseCase usecases.UserUseCase) *UserController {
	return &UserController{
		userUseCase: userUseCase,
	}
}

// GetUser returns user information
func (c *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := c.userUseCase.GetUsers()
	if err != nil {
		log.Printf("Failed to get users: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from query parameters
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

	user, err := c.userUseCase.GetUser(userID)
	if err == errors.ErrUserNotFound {
		w.WriteHeader(404)
		return
	} else if err != nil {
		log.Printf("Failed to get user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateUserLanguage updates user's language preference
func (c *UserController) UpdateUserLanguage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Only PUT requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from query parameters
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id parameter is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	// Parse request body
	var request struct {
		Language string `json:"language"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if request.Language == "" {
		http.Error(w, "language field is required", http.StatusBadRequest)
		return
	}

	err = c.userUseCase.UpdateUserLanguage(userID, request.Language)
	if err != nil {
		log.Printf("Failed to update user language: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Language updated successfully"})
}

// GetUserSelection returns user's current selection state
func (c *UserController) GetUserSelection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from query parameters
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id parameter is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	selection, err := c.userUseCase.GetUserSelection(userID)
	if err != nil {
		log.Printf("Failed to get user selection: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(selection)
}

// ClearUserSelection clears user's selection state
func (c *UserController) ClearUserSelection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from query parameters
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id parameter is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	err = c.userUseCase.ClearUserSelection(userID)
	if err != nil {
		log.Printf("Failed to clear user selection: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User selection cleared successfully"})
}
