package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/repositories"
	"github.com/ivanenkomaksym/remindme_bot/domain/response"
)

// PremiumUsageController handles HTTP requests for premium usage
type PremiumUsageController struct {
	premiumUsageRepo repositories.PremiumUsageRepository
	userRepo         repositories.UserRepository
}

// NewPremiumUsageController creates a new premium usage controller
func NewPremiumUsageController(premiumUsageRepo repositories.PremiumUsageRepository, userRepo repositories.UserRepository) *PremiumUsageController {
	return &PremiumUsageController{
		premiumUsageRepo: premiumUsageRepo,
		userRepo:         userRepo,
	}
}

// PremiumUsageResponse represents the API response for premium usage
type PremiumUsageResponse struct {
	UserID            int64                  `json:"userId"`
	RequestsUsed      int                    `json:"requestsUsed"`
	RequestsLimit     int                    `json:"requestsLimit"`
	RemainingRequests int                    `json:"remainingRequests"`
	LastReset         time.Time              `json:"lastReset"`
	PremiumStatus     entities.PremiumStatus `json:"premiumStatus"`
	PremiumUpgradeAt  *time.Time             `json:"premiumUpgradeAt,omitempty"`
	PremiumExpiresAt  *time.Time             `json:"premiumExpiresAt,omitempty"`
	DaysUntilReset    int                    `json:"daysUntilReset"`
	DaysUntilExpiry   int                    `json:"daysUntilExpiry"`
	IsExpired         bool                   `json:"isExpired"`
	CreatedAt         time.Time              `json:"createdAt"`
	UpdatedAt         time.Time              `json:"updatedAt"`
}

// UpgradePremiumRequest represents the request to upgrade premium status
type UpgradePremiumRequest struct {
	PremiumStatus entities.PremiumStatus `json:"premiumStatus"`
}

// toPremiumUsageResponse converts PremiumUsage entity to API response
func toPremiumUsageResponse(usage *entities.PremiumUsage) *PremiumUsageResponse {
	return &PremiumUsageResponse{
		UserID:            usage.UserID,
		RequestsUsed:      usage.RequestsUsed,
		RequestsLimit:     usage.RequestsLimit,
		RemainingRequests: usage.GetRemainingRequests(),
		LastReset:         usage.LastReset,
		PremiumStatus:     usage.PremiumStatus,
		PremiumUpgradeAt:  usage.PremiumUpgradeAt,
		PremiumExpiresAt:  usage.PremiumExpiresAt,
		DaysUntilReset:    usage.GetDaysUntilReset(),
		DaysUntilExpiry:   usage.GetDaysUntilExpiration(),
		IsExpired:         usage.IsPremiumExpired(),
		CreatedAt:         usage.CreatedAt,
		UpdatedAt:         usage.UpdatedAt,
	}
}

// GetUserPremiumUsage handles GET /api/premium/{user_id}
func (c *PremiumUsageController) GetUserPremiumUsage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	userIDStr := r.PathValue("user_id")
	if userIDStr == "" {
		response.WriteBadRequest(w, "Missing user_id parameter")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.WriteBadRequest(w, "Invalid user_id parameter")
		return
	}

	// Check if user exists
	_, err = c.userRepo.GetUser(userID)
	if err != nil {
		response.WriteNotFound(w, "User not found")
		return
	}

	// Get or create premium usage
	usage, err := c.premiumUsageRepo.GetOrCreateUserUsage(userID)
	if err != nil {
		response.WriteInternalError(w, "Failed to get premium usage", err)
		return
	}

	// Handle expired premium automatically
	if usage.HandleExpiredPremium() {
		if err := c.premiumUsageRepo.UpdateUserUsage(usage); err != nil {
			// Log error but continue with response
			fmt.Printf("Failed to update expired premium status for user %d: %v\n", userID, err)
		}
	}

	response.WriteSuccess(w, "Premium usage retrieved successfully", toPremiumUsageResponse(usage))
}

// GetAllPremiumUsage handles GET /api/premium
func (c *PremiumUsageController) GetAllPremiumUsage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	usages, err := c.premiumUsageRepo.GetAllUsage()
	if err != nil {
		response.WriteInternalError(w, "Failed to get premium usage", err)
		return
	}

	var responses []*PremiumUsageResponse
	for _, usage := range usages {
		// Handle expired premium automatically
		if usage.HandleExpiredPremium() {
			if err := c.premiumUsageRepo.UpdateUserUsage(&usage); err != nil {
				// Log error but continue
				fmt.Printf("Failed to update expired premium status for user %d: %v\n", usage.UserID, err)
			}
		}
		responses = append(responses, toPremiumUsageResponse(&usage))
	}

	response.WriteSuccess(w, "Premium usage list retrieved successfully", responses)
}

// UpgradeUserPremium handles PUT /api/premium/{user_id}/upgrade
func (c *PremiumUsageController) UpgradeUserPremium(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		response.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	userIDStr := r.PathValue("user_id")
	if userIDStr == "" {
		response.WriteBadRequest(w, "Missing user_id parameter")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.WriteBadRequest(w, "Invalid user_id parameter")
		return
	}

	// Check if user exists
	_, err = c.userRepo.GetUser(userID)
	if err != nil {
		response.WriteNotFound(w, "User not found")
		return
	}

	// Parse request body
	var req UpgradePremiumRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteBadRequest(w, "Invalid request body")
		return
	}

	// Validate premium status
	if req.PremiumStatus != entities.PremiumStatusFree &&
		req.PremiumStatus != entities.PremiumStatusBasic &&
		req.PremiumStatus != entities.PremiumStatusPro {
		response.WriteBadRequest(w, "Invalid premium status")
		return
	}

	// Get or create premium usage
	usage, err := c.premiumUsageRepo.GetOrCreateUserUsage(userID)
	if err != nil {
		response.WriteInternalError(w, "Failed to get premium usage", err)
		return
	}

	// Update premium status
	usage.SetPremiumStatus(req.PremiumStatus)

	// Save changes
	if err := c.premiumUsageRepo.UpdateUserUsage(usage); err != nil {
		response.WriteInternalError(w, "Failed to update premium status", err)
		return
	}

	response.WriteSuccess(w, "Premium status updated successfully", toPremiumUsageResponse(usage))
}

// ResetUserUsage handles POST /api/premium/{user_id}/reset
func (c *PremiumUsageController) ResetUserUsage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	userIDStr := r.PathValue("user_id")
	if userIDStr == "" {
		response.WriteBadRequest(w, "Missing user_id parameter")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.WriteBadRequest(w, "Invalid user_id parameter")
		return
	}

	// Check if user exists
	_, err = c.userRepo.GetUser(userID)
	if err != nil {
		response.WriteNotFound(w, "User not found")
		return
	}

	// Get premium usage
	usage, err := c.premiumUsageRepo.GetUserUsage(userID)
	if err != nil {
		response.WriteNotFound(w, "Premium usage not found")
		return
	}

	// Reset usage
	usage.ResetUsage()

	// Save changes
	if err := c.premiumUsageRepo.UpdateUserUsage(usage); err != nil {
		response.WriteInternalError(w, "Failed to reset usage", err)
		return
	}

	response.WriteSuccess(w, "Usage reset successfully", toPremiumUsageResponse(usage))
}

// GetPremiumUsageByStatus handles GET /api/premium/status/{status}
func (c *PremiumUsageController) GetPremiumUsageByStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	statusStr := r.PathValue("status")
	if statusStr == "" {
		response.WriteBadRequest(w, "Missing status parameter")
		return
	}

	status := entities.PremiumStatus(statusStr)
	if status != entities.PremiumStatusFree &&
		status != entities.PremiumStatusBasic &&
		status != entities.PremiumStatusPro {
		response.WriteBadRequest(w, "Invalid premium status")
		return
	}

	usages, err := c.premiumUsageRepo.GetUsageByPremiumStatus(status)
	if err != nil {
		response.WriteInternalError(w, "Failed to get premium usage by status", err)
		return
	}

	var responses []*PremiumUsageResponse
	for _, usage := range usages {
		// Handle expired premium automatically
		if usage.HandleExpiredPremium() {
			if err := c.premiumUsageRepo.UpdateUserUsage(&usage); err != nil {
				// Log error but continue
				fmt.Printf("Failed to update expired premium status for user %d: %v\n", usage.UserID, err)
			}
		}
		responses = append(responses, toPremiumUsageResponse(&usage))
	}

	response.WriteSuccess(w, "Premium usage by status retrieved successfully", responses)
}

// DeleteUserPremiumUsage handles DELETE /api/premium/{user_id}
func (c *PremiumUsageController) DeleteUserPremiumUsage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		response.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	userIDStr := r.PathValue("user_id")
	if userIDStr == "" {
		response.WriteBadRequest(w, "Missing user_id parameter")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.WriteBadRequest(w, "Invalid user_id parameter")
		return
	}

	// Delete premium usage
	if err := c.premiumUsageRepo.DeleteUserUsage(userID); err != nil {
		response.WriteInternalError(w, "Failed to delete premium usage", err)
		return
	}

	response.WriteSuccess(w, "Premium usage deleted successfully", nil)
}
