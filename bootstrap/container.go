package bootstrap

import (
	"log"

	"github.com/ivanenkomaksym/remindme_bot/api/controllers"
	"github.com/ivanenkomaksym/remindme_bot/domain/repositories"
	"github.com/ivanenkomaksym/remindme_bot/domain/usecases"
	"github.com/ivanenkomaksym/remindme_bot/repositories/inmemory"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Container holds all the dependencies
type Container struct {
	// Repositories
	UserRepo     repositories.UserRepository
	ReminderRepo repositories.ReminderRepository

	// Use Cases
	UserUseCase     usecases.UserUseCase
	ReminderUseCase usecases.ReminderUseCase
	BotUseCase      usecases.BotUseCase

	// Controllers
	BotController      *controllers.BotController
	UserController     *controllers.UserController
	ReminderController *controllers.ReminderController
}

// NewContainer creates a new dependency injection container
func NewContainer(app *Application) *Container {
	container := &Container{}

	// Initialize repositories
	container.initRepositories(app.Env.StorageType)

	// Initialize use cases
	container.initUseCases()

	// Initialize controllers
	container.initControllers(app.Bot)

	return container
}

// initRepositories initializes all repositories
func (c *Container) initRepositories(storageType repositories.StorageType) {
	switch storageType {
	case repositories.InMemory:
		c.UserRepo = inmemory.NewInMemoryUserRepository()
		c.ReminderRepo = inmemory.NewInMemoryReminderRepository()
	default:
		log.Fatalf("Unsupported storage type: %v", storageType)
	}
}

// initUseCases initializes all use cases
func (c *Container) initUseCases() {
	c.UserUseCase = usecases.NewUserUseCase(c.UserRepo)
	c.ReminderUseCase = usecases.NewReminderUseCase(c.ReminderRepo, c.UserRepo)
}

// initControllers initializes all controllers
func (c *Container) initControllers(bot *tgbotapi.BotAPI) {
	c.BotUseCase = usecases.NewBotUseCase(c.UserUseCase, c.ReminderUseCase, bot)
	c.BotController = controllers.NewBotController(c.BotUseCase, bot)
	c.UserController = controllers.NewUserController(c.UserUseCase)
	c.ReminderController = controllers.NewReminderController(c.ReminderUseCase)
}
