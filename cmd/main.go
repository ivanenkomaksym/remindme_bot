package main

import (
	"github.com/ivanenkomaksym/remindme_bot/api/route"
	"github.com/ivanenkomaksym/remindme_bot/bootstrap"
	"github.com/ivanenkomaksym/remindme_bot/repositories"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	bot *tgbotapi.BotAPI

	welcomeMessage = "Welcome to the Reminder Bot!"
	reminderRepo   repositories.ReminderRepository
	userRepo       repositories.UserRepository
)

func main() {
	app := bootstrap.App()

	route.Setup(&app)
}
