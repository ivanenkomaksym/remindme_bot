package config

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
	OpenAI   OpenAIConfig
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
	Enabled                 bool
	Token                   string
	Debug                   bool
	PublicURL               string
	MonitorPendingUpdates   bool
	PendingUpdatesThreshold int
	AutoClearPendingUpdates bool
}

// AppConfig holds application-related configuration
type AppConfig struct {
	Environment     string
	LogLevel        string
	Timezone        string
	APIKey          string
	NotifierTimeout time.Duration
}

// OpenAIConfig holds OpenAI-related configuration
type OpenAIConfig struct {
	Enabled bool
	APIKey  string
	Model   string
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
	config.loadOpenAIConfig()

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
		Enabled:                 false,
		Token:                   "",
		Debug:                   false,
		PublicURL:               "",
		MonitorPendingUpdates:   true,
		PendingUpdatesThreshold: 100,
		AutoClearPendingUpdates: true,
	}

	c.App = AppConfig{
		Environment:     "development",
		LogLevel:        "info",
		Timezone:        "UTC",
		APIKey:          "",
		NotifierTimeout: 1 * time.Minute,
	}

	c.OpenAI = OpenAIConfig{
		Enabled: false,
		APIKey:  "",
		Model:   "gpt-4o-mini",
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
	c.Bot.Enabled = viper.GetBool("BOT_ENABLED")
	if !c.Bot.Enabled {
		return
	}

	if token := viper.GetString("BOT_TOKEN"); token != "" {
		c.Bot.Token = token
	}
	if debug := viper.GetString("BOT_DEBUG"); debug != "" {
		if d, err := strconv.ParseBool(debug); err == nil {
			c.Bot.Debug = d
		}
	}
	if publicURL := viper.GetString("PUBLIC_URL"); publicURL != "" {
		c.Bot.PublicURL = publicURL
	}
	if monitorUpdates := viper.GetString("BOT_MONITOR_PENDING_UPDATES"); monitorUpdates != "" {
		if m, err := strconv.ParseBool(monitorUpdates); err == nil {
			c.Bot.MonitorPendingUpdates = m
		}
	}
	if threshold := viper.GetString("BOT_PENDING_UPDATES_THRESHOLD"); threshold != "" {
		if t, err := strconv.Atoi(threshold); err == nil {
			c.Bot.PendingUpdatesThreshold = t
		}
	}
	if autoClear := viper.GetString("BOT_AUTO_CLEAR_PENDING_UPDATES"); autoClear != "" {
		if a, err := strconv.ParseBool(autoClear); err == nil {
			c.Bot.AutoClearPendingUpdates = a
		}
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
	if notifierTimeout := viper.GetString("NOTIFIER_TIMEOUT"); notifierTimeout != "" {
		if duration, err := time.ParseDuration(notifierTimeout); err == nil {
			c.App.NotifierTimeout = duration
		}
	}
}

// loadOpenAIConfig loads OpenAI configuration
func (c *Config) loadOpenAIConfig() {
	c.OpenAI.Enabled = viper.GetBool("OPENAI_ENABLED")
	if apiKey := viper.GetString("OPENAI_API_KEY"); apiKey != "" {
		c.OpenAI.APIKey = apiKey
	}
	if model := viper.GetString("OPENAI_MODEL"); model != "" {
		c.OpenAI.Model = model
	}
}

// validate validates the configuration
func (c *Config) validate() {
	if c.Bot.Enabled {
		if c.Bot.Token == "" {
			log.Fatal("BOT_TOKEN is required")
		}
		if c.Bot.PublicURL == "" {
			log.Fatal("PUBLIC_URL is required")
		}
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
