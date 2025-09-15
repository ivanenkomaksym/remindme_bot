package keyboards

import (
	"strings"
	"testing"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

func TestHandleRecurrenceTypeSelection(t *testing.T) {
	user := &entities.User{}
	userSelection := &entities.UserSelection{}

	// Once
	res, err := HandleRecurrenceTypeSelection(entities.Once.String(), user, userSelection)
	if err != nil || res == nil {
		t.Fatalf("once should return result and nil error")
	}
	if !strings.Contains(res.Text, "time") {
		t.Fatalf("unexpected result text: %s", res.Text)
	}

	// Daily
	res, err = HandleRecurrenceTypeSelection(entities.Daily.String(), user, userSelection)
	if err != nil || res == nil {
		t.Fatalf("daily should return result and nil error")
	}
	if !strings.Contains(res.Text, "time") {
		t.Fatalf("unexpected result text: %s", res.Text)
	}

	// Weekly
	res, err = HandleRecurrenceTypeSelection(entities.Weekly.String(), user, userSelection)
	if err != nil || res == nil || !userSelection.IsWeekly {
		t.Fatalf("weekly should set IsWeekly and return result")
	}

	// Monthly
	res, err = HandleRecurrenceTypeSelection(entities.Monthly.String(), user, userSelection)
	if err != nil || res == nil {
		t.Fatalf("monthly should return result")
	}

	// Interval
	res, err = HandleRecurrenceTypeSelection(entities.Interval.String(), user, userSelection)
	if err != nil || res == nil {
		t.Fatalf("interval should return result")
	}

	// Custom
	res, err = HandleRecurrenceTypeSelection(entities.Custom.String(), user, userSelection)
	if err != nil || res == nil {
		t.Fatalf("custom should return result")
	}
}
