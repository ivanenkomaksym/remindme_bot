package bootstrap

import (
	"log"

	"github.com/ivanenkomaksym/remindme_bot/repositories"
	"github.com/spf13/viper"
)

type Env struct {
	AppEnv        string `mapstructure:"APP_ENV"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	BotToken      string `mapstructure:"BOT_TOKEN"`
	WebhookUrl    string `mapstructure:"WEBHOOK_URL"`
	Storage       string `mapstructure:"STORAGE"`
	StorageType   repositories.StorageType
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	// Convert storage string to StorageType
	if env.Storage == "" {
		env.Storage = "inmemory" // Default to in-memory storage
	}

	storageType, err := repositories.ToStorageType(env.Storage)
	if err != nil {
		log.Fatalf("Invalid storage type '%s': %v", env.Storage, err)
	}
	env.StorageType = storageType

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env
}
