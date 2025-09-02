package repositories

import (
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/internal/models"
)

type InMemoryReminderRepository struct {
	mu        sync.Mutex
	nextID    int64
	reminders []models.Reminder
}

func NewInMemoryReminderRepository() *InMemoryReminderRepository {
	return &InMemoryReminderRepository{
		nextID:    1,
		reminders: make([]models.Reminder, 0),
	}
}

func (r *InMemoryReminderRepository) CreateDailyReminder(timeStr string, user models.User, message string) *models.Reminder {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	recurrence := models.DailyAt(timeStr)
	next := nextDailyTrigger(now, timeStr)

	reminder := models.Reminder{
		ID:          r.nextID,
		User:        user,
		Message:     message,
		CreatedAt:   now,
		NextTrigger: next,
		Recurrence:  recurrence,
		IsActive:    true,
	}

	r.nextID++
	r.reminders = append(r.reminders, reminder)
	return &r.reminders[len(r.reminders)-1]
}

func (r *InMemoryReminderRepository) CreateWeeklyReminder(daysOfWeek []time.Weekday, timeStr string, user models.User, message string) *models.Reminder {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	recurrence := models.CustomWeekly(daysOfWeek, timeStr)
	next := nextWeeklyTrigger(now, daysOfWeek, timeStr)

	reminder := models.Reminder{
		ID:          r.nextID,
		User:        user,
		Message:     message,
		CreatedAt:   now,
		NextTrigger: next,
		Recurrence:  recurrence,
		IsActive:    true,
	}

	r.nextID++
	r.reminders = append(r.reminders, reminder)
	return &r.reminders[len(r.reminders)-1]
}

func (r *InMemoryReminderRepository) CreateMonthlyReminder(daysOfMonth []int, timeStr string, user models.User, message string) *models.Reminder {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	recurrence := models.MonthlyOnDay(daysOfMonth, timeStr)
	next := nextMonthlyTrigger(now, daysOfMonth, timeStr)

	reminder := models.Reminder{
		ID:          r.nextID,
		User:        user,
		Message:     message,
		CreatedAt:   now,
		NextTrigger: next,
		Recurrence:  recurrence,
		IsActive:    true,
	}

	r.nextID++
	r.reminders = append(r.reminders, reminder)
	return &r.reminders[len(r.reminders)-1]
}

func (r *InMemoryReminderRepository) GetReminders() []models.Reminder {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]models.Reminder, len(r.reminders))
	copy(out, r.reminders)
	return out
}

// helpers

func parseHourMinute(timeStr string) (int, int, bool) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0, 0, false
	}
	hour, err1 := strconv.Atoi(parts[0])
	minute, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return 0, 0, false
	}
	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return 0, 0, false
	}
	return hour, minute, true
}

func nextDailyTrigger(from time.Time, timeStr string) time.Time {
	hour, minute, ok := parseHourMinute(timeStr)
	if !ok {
		return from
	}
	candidate := time.Date(from.Year(), from.Month(), from.Day(), hour, minute, 0, 0, from.Location())
	if !candidate.After(from) {
		candidate = candidate.Add(24 * time.Hour)
	}
	return candidate
}

func nextWeeklyTrigger(from time.Time, days []time.Weekday, timeStr string) time.Time {
	if len(days) == 0 {
		return nextDailyTrigger(from, timeStr)
	}
	hour, minute, ok := parseHourMinute(timeStr)
	if !ok {
		return from
	}
	seen := map[time.Weekday]struct{}{}
	uniqueDays := make([]time.Weekday, 0, len(days))
	for _, d := range days {
		if _, exists := seen[d]; !exists {
			seen[d] = struct{}{}
			uniqueDays = append(uniqueDays, d)
		}
	}

	best := time.Time{}
	for i := 0; i < 7; i++ {
		day := from.Add(time.Duration(i) * 24 * time.Hour)
		for _, d := range uniqueDays {
			if day.Weekday() == d {
				candidate := time.Date(day.Year(), day.Month(), day.Day(), hour, minute, 0, 0, from.Location())
				if candidate.After(from) && (best.IsZero() || candidate.Before(best)) {
					best = candidate
				}
			}
		}
	}
	if best.IsZero() {
		return nextWeeklyTrigger(from.Add(7*24*time.Hour), uniqueDays, timeStr)
	}
	return best
}

func daysIn(month time.Month, year int) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func nextMonthlyTrigger(from time.Time, daysOfMonth []int, timeStr string) time.Time {
	if len(daysOfMonth) == 0 {
		return nextDailyTrigger(from, timeStr)
	}
	hour, minute, ok := parseHourMinute(timeStr)
	if !ok {
		return from
	}
	uniq := map[int]struct{}{}
	days := make([]int, 0, len(daysOfMonth))
	for _, d := range daysOfMonth {
		if d >= 1 && d <= 31 {
			if _, exists := uniq[d]; !exists {
				uniq[d] = struct{}{}
				days = append(days, d)
			}
		}
	}
	if len(days) == 0 {
		return nextDailyTrigger(from, timeStr)
	}
	sort.Ints(days)

	best := time.Time{}
	for m := 0; m < 3; m++ {
		t := from.AddDate(0, m, 0)
		dim := daysIn(t.Month(), t.Year())
		for _, d := range days {
			if d > dim {
				continue
			}
			candidate := time.Date(t.Year(), t.Month(), d, hour, minute, 0, 0, from.Location())
			if candidate.After(from) && (best.IsZero() || candidate.Before(best)) {
				best = candidate
			}
		}
		if !best.IsZero() {
			break
		}
	}
	if best.IsZero() {
		best = from.Add(24 * time.Hour)
	}
	return best
}
