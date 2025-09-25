package bootstrap

import (
	"log"
	"strconv"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/domain/repositories"
	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Bot      BotConfig
	App      AppConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Address         string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Type             repositories.StorageType
	Host             string
	Port             int
	Username         string
	Password         string
	Database         string
	SSLMode          string
	ConnectionString string
}

// BotConfig holds bot-related configuration
type BotConfig struct {
	Token      string
	Debug      bool
	WebhookURL string
}

// AppConfig holds application-related configuration
type AppConfig struct {
	Environment string
	LogLevel    string
	Timezone    string
	APIKey      string
}

// LoadConfig loads configuration from environment variables and config files
func LoadConfig() *Config {
	config := &Config{}

	// Set default values
	config.setDefaults()

	// Load from .env file
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read .env file: %v", err)
	}

	// Load from environment variables
	viper.AutomaticEnv()

	// Load configuration
	config.loadServerConfig()
	config.loadDatabaseConfig()
	config.loadBotConfig()
	config.loadAppConfig()

	// Validate configuration
	config.validate()

	return config
}

// setDefaults sets default configuration values
func (c *Config) setDefaults() {
	c.Server = ServerConfig{
		Address:         "0.0.0.0",
		Port:            "8080",
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
		IdleTimeout:     120 * time.Second,
		ShutdownTimeout: 10 * time.Second,
	}

	c.Database = DatabaseConfig{
		Type:             repositories.InMemory,
		Host:             "localhost",
		Port:             5432,
		Username:         "postgres",
		Password:         "",
		Database:         "remindme_bot",
		SSLMode:          "disable",
		ConnectionString: "",
	}

	c.Bot = BotConfig{
		Token:      "",
		Debug:      false,
		WebhookURL: "",
	}

	c.App = AppConfig{
		Environment: "development",
		LogLevel:    "info",
		Timezone:    "UTC",
		APIKey:      "",
	}
}

// loadServerConfig loads server configuration
func (c *Config) loadServerConfig() {
	if addr := viper.GetString("SERVER_ADDRESS"); addr != "" {
		c.Server.Address = addr
	}
	if port := viper.GetString("PORT"); port != "" {
		c.Server.Port = port
	}
	if readTimeout := viper.GetString("READ_TIMEOUT"); readTimeout != "" {
		if duration, err := time.ParseDuration(readTimeout); err == nil {
			c.Server.ReadTimeout = duration
		}
	}
	if writeTimeout := viper.GetString("WRITE_TIMEOUT"); writeTimeout != "" {
		if duration, err := time.ParseDuration(writeTimeout); err == nil {
			c.Server.WriteTimeout = duration
		}
	}
	if idleTimeout := viper.GetString("IDLE_TIMEOUT"); idleTimeout != "" {
		if duration, err := time.ParseDuration(idleTimeout); err == nil {
			c.Server.IdleTimeout = duration
		}
	}
	if shutdownTimeout := viper.GetString("SHUTDOWN_TIMEOUT"); shutdownTimeout != "" {
		if duration, err := time.ParseDuration(shutdownTimeout); err == nil {
			c.Server.ShutdownTimeout = duration
		}
	}
}

// loadDatabaseConfig loads database configuration
func (c *Config) loadDatabaseConfig() {
	if dbType := viper.GetString("STORAGE"); dbType != "" {
		if storageType, err := repositories.ToStorageType(dbType); err == nil {
			c.Database.Type = storageType
		}
	}
	if conn := viper.GetString("DB_CONNECTION_STRING"); conn != "" {
		c.Database.ConnectionString = conn
	}
	if host := viper.GetString("DB_HOST"); host != "" {
		c.Database.Host = host
	}
	if port := viper.GetString("DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			c.Database.Port = p
		}
	}
	if username := viper.GetString("DB_USERNAME"); username != "" {
		c.Database.Username = username
	}
	if password := viper.GetString("DB_PASSWORD"); password != "" {
		c.Database.Password = password
	}
	if database := viper.GetString("DB_DATABASE"); database != "" {
		c.Database.Database = database
	}
	if sslMode := viper.GetString("DB_SSL_MODE"); sslMode != "" {
		c.Database.SSLMode = sslMode
	}
}

// loadBotConfig loads bot configuration
func (c *Config) loadBotConfig() {
	if token := viper.GetString("BOT_TOKEN"); token != "" {
		c.Bot.Token = token
	}
	if debug := viper.GetString("BOT_DEBUG"); debug != "" {
		if d, err := strconv.ParseBool(debug); err == nil {
			c.Bot.Debug = d
		}
	}
	if webhookURL := viper.GetString("WEBHOOK_URL"); webhookURL != "" {
		c.Bot.WebhookURL = webhookURL
	}
}

// loadAppConfig loads application configuration
func (c *Config) loadAppConfig() {
	if env := viper.GetString("APP_ENV"); env != "" {
		c.App.Environment = env
	}
	if logLevel := viper.GetString("LOG_LEVEL"); logLevel != "" {
		c.App.LogLevel = logLevel
	}
	if timezone := viper.GetString("TIMEZONE"); timezone != "" {
		c.App.Timezone = timezone
	}
	if apiKey := viper.GetString("API_KEY"); apiKey != "" {
		c.App.APIKey = apiKey
	}
}

// validate validates the configuration
func (c *Config) validate() {
	if c.Bot.Token == "" {
		log.Fatal("BOT_TOKEN is required")
	}
	if c.Bot.WebhookURL == "" {
		log.Fatal("WEBHOOK_URL is required")
	}
	if c.Server.Port == "" {
		log.Fatal("Server port is required")
	}
}

// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	return c.Server.Address + ":" + c.Server.Port
}

// IsDevelopment returns true if the application is in development mode
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsProduction returns true if the application is in production mode
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// GetLogLevel returns the log level
func (c *Config) GetLogLevel() string {
	return c.App.LogLevel
}

// GetTimezone returns the timezone
func (c *Config) GetTimezone() string {
	return c.App.Timezone
}
