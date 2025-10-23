package bootstrap

import (
	"fmt"
	"log"

	"github.com/ivanenkomaksym/remindme_bot/api/controllers"
	"github.com/ivanenkomaksym/remindme_bot/config"
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/repositories"
	"github.com/ivanenkomaksym/remindme_bot/domain/services"
	"github.com/ivanenkomaksym/remindme_bot/domain/usecases"
	"github.com/ivanenkomaksym/remindme_bot/repositories/inmemory"
	"github.com/ivanenkomaksym/remindme_bot/repositories/persistent"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Container holds all the dependencies
type Container struct {
	Config config.Config

	// Repositories
	UserRepo          repositories.UserRepository
	ReminderRepo      repositories.ReminderRepository
	UserSelectionRepo repositories.UserSelectionRepository
	PremiumUsageRepo  repositories.PremiumUsageRepository

	// Services
	NLPService services.NLPService

	// Use Cases
	UserUseCase         usecases.UserUseCase
	ReminderUseCase     usecases.ReminderUseCase
	PremiumUsageUseCase usecases.PremiumUsageUseCase
	BotUseCase          usecases.BotUseCase
	DateUseCase         usecases.DateUseCase

	// Controllers
	BotController          *controllers.BotController
	UserController         *controllers.UserController
	ReminderController     *controllers.ReminderController
	TimezoneController     *controllers.TimezoneController
	PremiumUsageController *controllers.PremiumUsageController
}

// NewContainer creates a new dependency injection container
func NewContainer(app *Application) *Container {
	container := &Container{Config: *app.Env.Config}

	// Initialize repositories
	container.initRepositories(app.Env)

	// Initialize services
	container.initServices()

	// Initialize use cases
	container.initUseCases()

	// Initialize controllers
	container.initControllers(app.Bot)

	return container
}

// initRepositories initializes all repositories
func (c *Container) initRepositories(env *Env) {
	switch env.StorageType {
	case repositories.InMemory:
		c.UserRepo = inmemory.NewInMemoryUserRepository()
		c.ReminderRepo = inmemory.NewInMemoryReminderRepository()
		c.UserSelectionRepo = inmemory.NewInMemoryUserSelectionRepository()
		c.PremiumUsageRepo = inmemory.NewInMemoryPremiumUsageRepository()
	case repositories.Mongo:
		// Expect connection string and database name from config
		conn := env.Config.Database.ConnectionString
		if conn == "" {
			log.Fatalf("Missing database connection string")
		}
		dbName := env.Config.Database.Database
		userRepo, err := persistent.NewMongoUserRepository(conn, dbName)
		if err != nil {
			log.Fatalf("Failed to init Mongo user repo: %v", err)
		}
		remRepo, err := persistent.NewMongoReminderRepository(conn, dbName)
		if err != nil {
			log.Fatalf("Failed to init Mongo reminder repo: %v", err)
		}
		premiumRepo, err := persistent.NewMongoPremiumUsageRepository(conn, dbName)
		if err != nil {
			log.Fatalf("Failed to init Mongo premium usage repo: %v", err)
		}
		c.UserRepo = userRepo
		c.ReminderRepo = remRepo
		c.PremiumUsageRepo = premiumRepo
		// User selections still in-memory for now
		c.UserSelectionRepo = inmemory.NewInMemoryUserSelectionRepository()
	default:
		log.Fatalf("Unsupported storage type: %v", env.StorageType)
	}
}

// initServices initializes all services
func (c *Container) initServices() {
	c.NLPService = &noOpNLPService{}
	if !c.Config.OpenAI.Enabled {
		return
	}

	var err error
	c.NLPService, err = services.NewNLPService(&c.Config, c.PremiumUsageRepo)
	if err != nil {
		log.Printf("Warning: NLP service initialization failed: %v", err)
	}
}

// initUseCases initializes all use cases
func (c *Container) initUseCases() {
	c.UserUseCase = usecases.NewUserUseCase(c.UserRepo, c.UserSelectionRepo)
	c.ReminderUseCase = usecases.NewReminderUseCase(c.ReminderRepo, c.UserRepo)
	c.PremiumUsageUseCase = usecases.NewPremiumUsageUseCase(c.PremiumUsageRepo)
}

// noOpNLPService is a no-op implementation when OpenAI is not configured
type noOpNLPService struct{}

func (n *noOpNLPService) ParseReminderText(userID int64, text string, userTimezone string, userLanguage string) (*entities.UserSelection, error) {
	return nil, fmt.Errorf("NLP service is not configured - OpenAI API key required")
}

// initControllers initializes all controllers
func (c *Container) initControllers(bot *tgbotapi.BotAPI) {
	c.UserController = controllers.NewUserController(c.UserUseCase)
	c.ReminderController = controllers.NewReminderController(c.ReminderUseCase, c.NLPService, c.UserUseCase)
	c.PremiumUsageController = controllers.NewPremiumUsageController(c.PremiumUsageRepo, c.UserRepo)
	if bot != nil {
		c.DateUseCase = usecases.NewDateUseCase(c.UserUseCase, bot)
		c.BotUseCase = usecases.NewBotUseCase(c.UserUseCase, c.ReminderUseCase, c.DateUseCase, c.Config, bot, c.PremiumUsageUseCase, c.NLPService)
		c.BotController = controllers.NewBotController(c.BotUseCase, c.UserUseCase, c.DateUseCase, bot)
		c.TimezoneController = controllers.NewTimezoneController(c.UserUseCase, bot, c.Config)
	}
}
