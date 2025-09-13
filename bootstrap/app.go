package bootstrap

import (
	"github.com/ivanenkomaksym/remindme_bot/repositories"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Application struct {
	Env          *Env
	Bot          *tgbotapi.BotAPI
	ReminderRepo repositories.ReminderRepository
	UserRepo     repositories.UserRepository
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.Bot = NewBot(app.Env)

	reminderFactory := repositories.NewReminderRepositoryFactory()
	app.ReminderRepo = reminderFactory.CreateRepository(app.Env.StorageType)

	userFactory := repositories.NewUserRepositoryFactory()
	app.UserRepo = userFactory.CreateRepository(app.Env.StorageType)

	return *app
}
