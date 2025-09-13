package bootstrap

import (
	"github.com/ivanenkomaksym/remindme_bot/internal/repositories"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Application struct {
	Env          *Env
	ReminderRepo repositories.ReminderRepository
	Bot          *tgbotapi.BotAPI
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.Bot = NewBot(app.Env)
	return *app
}
