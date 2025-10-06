package bootstrap

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Application struct {
	Env       *Env
	Bot       *tgbotapi.BotAPI
	Container *Container
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	if app.Env.Config.Bot.Enabled {
		app.Bot = NewBot(app.Env)
	}
	app.Container = NewContainer(app)

	return *app
}
