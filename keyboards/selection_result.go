package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SelectionResult is a unified return type for selection handlers.
type SelectionResult struct {
	Text   string
	Markup *tgbotapi.InlineKeyboardMarkup
}
