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
func NextDailyTrigger(from time.Time, timeStr string) time.Time {
	hour, minute, ok := ParseHourMinute(timeStr)
	if !ok {
		return from
	}
	candidate := time.Date(from.Year(), from.Month(), from.Day(), hour, minute, 0, 0, from.Location())
	if !candidate.After(from) {
		candidate = candidate.Add(24 * time.Hour)
	}
	return candidate
}

// NextWeeklyTrigger returns the next occurrence on any of the provided weekdays at HH:MM.
func NextWeeklyTrigger(from time.Time, days []time.Weekday, timeStr string) time.Time {
	if len(days) == 0 {
		return NextDailyTrigger(from, timeStr)
	}
	hour, minute, ok := ParseHourMinute(timeStr)
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
	for i := range 7 {
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
		return NextWeeklyTrigger(from.Add(7*24*time.Hour), uniqueDays, timeStr)
	}
	return best
}

func daysIn(month time.Month, year int) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// NextMonthlyTrigger returns the next occurrence on any of the provided days-of-month at HH:MM.
func NextMonthlyTrigger(from time.Time, daysOfMonth []int, timeStr string) time.Time {
	if len(daysOfMonth) == 0 {
		return NextDailyTrigger(from, timeStr)
	}
	hour, minute, ok := ParseHourMinute(timeStr)
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
		return NextDailyTrigger(from, timeStr)
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

// NextForRecurrence advances from last trigger according to the recurrence configuration.
func NextForRecurrence(last time.Time, rec *entities.Recurrence) *time.Time {
	switch rec.Type {
	case entities.Once:
		return nil
	case entities.Daily:
		// Maintain the same clock time as last trigger, add 24h
		result := last.Add(24 * time.Hour)
		return &result
	case entities.Weekly:
		result := NextWeeklyTrigger(last, rec.Weekdays, rec.GetTimeOfDay())
		return &result
	case entities.Monthly:
		result := NextMonthlyTrigger(last, rec.DayOfMonth, rec.GetTimeOfDay())
		return &result
	case entities.Interval:
		// Advance by configured number of days, retain time of day
		days := rec.Interval
		if days <= 0 {
			days = 1
		}
		result := last.Add(time.Duration(days) * 24 * time.Hour)
		return &result
	case entities.Custom:
		// Fallback simple daily for now
		result := last.Add(24 * time.Hour)
		return &result
	default:
		result := last.Add(24 * time.Hour)
		return &result
	}
}
