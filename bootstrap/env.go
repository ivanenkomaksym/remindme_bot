package bootstrap

import (
	"log"

	"github.com/ivanenkomaksym/remindme_bot/config"
	"github.com/ivanenkomaksym/remindme_bot/domain/repositories"
)

type Env struct {
	Config      *config.Config
	StorageType repositories.StorageType
}

func NewEnv() *Env {
	config := config.LoadConfig()

	env := &Env{
		Config:      config,
		StorageType: config.Database.Type,
	}

	if env.Config.IsDevelopment() {
		log.Println("The App is running in development env")
	}

	return env
}
