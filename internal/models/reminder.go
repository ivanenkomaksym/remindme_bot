package models

import "time"

type Reminder struct {
	ID          int64       // Unique identifier
	User        User        // Telegram user
	Message     string      // Reminder message
	CreatedAt   time.Time   // When the reminder was created
	NextTrigger time.Time   // Next time the reminder should fire
	Recurrence  *Recurrence // Recurrence pattern, nil if one-time
	IsActive    bool        // Whether the reminder is active
}
