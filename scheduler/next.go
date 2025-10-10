package scheduler

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

func ParseHourMinute(timeStr string) (int, int, bool) {
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

// NextDailyTrigger returns the next occurrence of the provided HH:MM from the reference time.
func NextDailyTrigger(from time.Time, timeOfDay time.Time, location *time.Location) time.Time {
	fromInTz := from.In(location)
	candidate := time.Date(from.Year(), from.Month(), from.Day(), timeOfDay.Hour(), timeOfDay.Minute(), 0, 0, location)
	if !candidate.After(fromInTz) {
		candidate = candidate.Add(24 * time.Hour)
	}

	// Convert back to UTC for storage
	return candidate.UTC()
}

// NextWeeklyTrigger returns the next occurrence on any of the provided weekdays at HH:MM.
func NextWeeklyTrigger(from time.Time, days []time.Weekday, timeOfDay time.Time, location *time.Location) time.Time {
	if len(days) == 0 {
		return NextDailyTrigger(from, timeOfDay, location)
	}
	seen := map[time.Weekday]struct{}{}
	uniqueDays := make([]time.Weekday, 0, len(days))
	for _, d := range days {
		if _, exists := seen[d]; !exists {
			seen[d] = struct{}{}
			uniqueDays = append(uniqueDays, d)
		}
	}

	fromInTz := from.In(location)

	best := time.Time{}
	for i := range 7 {
		day := fromInTz.Add(time.Duration(i) * 24 * time.Hour)
		for _, d := range uniqueDays {
			if day.Weekday() == d {
				candidate := time.Date(day.Year(), day.Month(), day.Day(), timeOfDay.Hour(), timeOfDay.Minute(), 0, 0, location)
				if candidate.After(fromInTz) && (best.IsZero() || candidate.Before(best)) {
					best = candidate
				}
			}
		}
	}
	if best.IsZero() {
		return NextWeeklyTrigger(from.Add(7*24*time.Hour), uniqueDays, timeOfDay, location)
	}

	// Convert back to UTC for storage
	return best.UTC()
}

func daysIn(month time.Month, year int) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// NextMonthlyTrigger returns the next occurrence on any of the provided days-of-month at HH:MM.
func NextMonthlyTrigger(from time.Time, daysOfMonth []int, timeOfDay time.Time, location *time.Location) time.Time {
	if len(daysOfMonth) == 0 {
		return NextDailyTrigger(from, timeOfDay, location)
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
		return NextDailyTrigger(from, timeOfDay, location)
	}
	sort.Ints(days)

	fromInTz := from.In(location)

	best := time.Time{}
	for m := 0; m < 3; m++ {
		t := fromInTz.AddDate(0, m, 0)
		dim := daysIn(t.Month(), t.Year())
		for _, d := range days {
			if d > dim {
				continue
			}
			candidate := time.Date(t.Year(), t.Month(), d, timeOfDay.Hour(), timeOfDay.Minute(), 0, 0, location)
			if candidate.After(fromInTz) && (best.IsZero() || candidate.Before(best)) {
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

	// Convert back to UTC for storage
	return best.UTC()
}

func NextForSpacedBasedRepetition(last time.Time, timeOfDay time.Time, rec *entities.Recurrence) *time.Time {
	if rec.Type != entities.SpacedBasedRepetition || len(rec.SpacedBasedRepetitionDays) == 0 {
		return nil
	}

	next := NextDailyTrigger(last, timeOfDay, rec.GetLocation())

	var nextInterval = rec.SpacedBasedRepetitionDays[0]
	rec.SpacedBasedRepetitionDays = rec.SpacedBasedRepetitionDays[1:]

	// Advance by the next interval, retain time of day
	result := next.Add(time.Duration(nextInterval) * 24 * time.Hour).UTC()
	return &result
}

// NextForRecurrence advances from last trigger according to the recurrence configuration.
func NextForRecurrence(last time.Time, timeOfDay time.Time, rec *entities.Recurrence) *time.Time {
	switch rec.Type {
	case entities.Once:
		return nil
	case entities.Daily:
		// Maintain the same clock time as last trigger, add 24h
		result := NextDailyTrigger(last, timeOfDay, rec.GetLocation())
		return &result
	case entities.Weekly:
		result := NextWeeklyTrigger(last, rec.Weekdays, timeOfDay, rec.GetLocation())
		return &result
	case entities.Monthly:
		result := NextMonthlyTrigger(last, rec.DayOfMonth, timeOfDay, rec.GetLocation())
		return &result
	case entities.Interval:
		// For interval reminders, we need to calculate the next trigger in the user's timezone
		// Convert last trigger to user's timezone
		lastInTz := last.In(rec.GetLocation())

		// Advance by configured number of days, retain time of day
		days := rec.Interval
		if days <= 0 {
			days = 1
		}
		// Convert back to UTC for storage
		result := lastInTz.Add(time.Duration(days) * 24 * time.Hour).UTC()

		return &result
	case entities.SpacedBasedRepetition:
		return NextForSpacedBasedRepetition(last, timeOfDay, rec)
	default:
		result := NextDailyTrigger(last, timeOfDay, rec.GetLocation())
		return &result
	}
}
