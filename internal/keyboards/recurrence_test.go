package keyboards

import (
	"strings"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/internal/models"
	"github.com/ivanenkomaksym/remindme_bot/internal/types"
)

func TestHandleRecurrenceTypeSelection(t *testing.T) {
	var msg tgbotapi.EditMessageTextConfig
	state := &types.UserSelectionState{}

	// Daily
	mk, err := HandleRecurrenceTypeSelection(models.Daily.String(), &msg, state)
	if err != nil || mk == nil {
		t.Fatalf("daily should return markup and nil error")
	}
	if !strings.Contains(msg.Text, "daily") {
		t.Fatalf("unexpected msg text: %s", msg.Text)
	}

	// Weekly
	mk, err = HandleRecurrenceTypeSelection(models.Weekly.String(), &msg, state)
	if err != nil || mk == nil || !state.IsWeekly {
		t.Fatalf("weekly should set IsWeekly and return markup")
	}

	// Monthly
	mk, err = HandleRecurrenceTypeSelection(models.Monthly.String(), &msg, state)
	if err != nil || mk == nil {
		t.Fatalf("monthly should return markup")
	}

	// Interval
	mk, err = HandleRecurrenceTypeSelection(models.Interval.String(), &msg, state)
	if err != nil || mk == nil {
		t.Fatalf("interval should return markup")
	}

	// Custom
	mk, err = HandleRecurrenceTypeSelection(models.Custom.String(), &msg, state)
	if err != nil || mk == nil {
		t.Fatalf("custom should return markup")
	}
}
