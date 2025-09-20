package keyboards

import "github.com/ivanenkomaksym/remindme_bot/domain/entities"

func IsDateSelectionCallback(callbackData string) bool {
	return callbackData == entities.Once.String()
}
