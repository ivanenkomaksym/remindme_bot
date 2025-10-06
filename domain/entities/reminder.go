package entities

import (
	"time"
)

// Reminder represents a reminder in the system
type Reminder struct {
	ID          int64       `json:"id,string" bson:"id"`
	UserID      int64       `json:"userId,string" bson:"userId"`
	Message     string      `json:"message" bson:"message"`
	CreatedAt   time.Time   `json:"createdAt" bson:"createdAt"`
	NextTrigger *time.Time  `json:"nextTrigger" bson:"nextTrigger"`
	Recurrence  *Recurrence `json:"recurrence" bson:"recurrence"`
	IsActive    bool        `json:"isActive" bson:"isActive"`
}

// NewReminder creates a new reminder entity
func NewReminder(id, userID int64, message string, recurrence *Recurrence, nextTrigger *time.Time) *Reminder {
	return &Reminder{
		ID:          id,
		UserID:      userID,
		Message:     message,
		CreatedAt:   time.Now(),
		NextTrigger: nextTrigger,
		Recurrence:  recurrence,
		IsActive:    true,
	}
}

// Deactivate marks the reminder as inactive
func (r *Reminder) Deactivate() {
	r.IsActive = false
}

// UpdateNextTrigger updates the next trigger time
func (r *Reminder) UpdateNextTrigger(nextTrigger *time.Time) {
	r.NextTrigger = nextTrigger
}
