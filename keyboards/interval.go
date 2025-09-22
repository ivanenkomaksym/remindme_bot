package keyboards

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

const (
	CallbackIntervalStart = "interval_start"
)

func IsIntervalCallback(callbackData string) bool {
	return callbackData == CallbackIntervalStart
}

// GetIntervalPrompt returns a text asking user to enter N for interval days. No keyboard.
func GetIntervalPrompt(userSelection *entities.UserSelection, lang string) *tgbotapi.InlineKeyboardMarkup {
	userSelection.StartCustomInterval()
	// no keyboard for interval; user should type a number
	return nil
}

func HandleCustomIntervalInput(text string, userEntity *entities.User, selection *entities.UserSelection) *SelectionResult {
	s := T(userEntity.Language)
	i, err := strconv.Atoi(text)
	if err != nil || i < 0 || i > 7 {
		return &SelectionResult{Text: s.MsgInvalidIntervalFormat, Markup: nil}
	}

	selection.SetCustomInterval(i)

	outputText := s.MsgSelectTime
	markup := GetHourRangeMarkup(userEntity.Language)

	return &SelectionResult{Text: outputText, Markup: markup}
}
