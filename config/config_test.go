package config

import (
	"testing"
	"time"

	"github.com/spf13/viper"
)

func resetViper() {
	// Best-effort reset to avoid cross-test contamination
	viper.Reset()
}

func TestLoadConfig_WithDefaultsRequiresMandatoryEnv(t *testing.T) {
	resetViper()

	var expectedApiKey = "73086a67-5eac-4dff-9f34-71335d5a6244"

	// Provide mandatory envs so validate() does not fatal
	t.Setenv("BOT_TOKEN", "test-token")
	t.Setenv("PUBLIC_URL", "https://example.com")
	// PORT is taken from PORT env variable
	t.Setenv("PORT", "8081")
	t.Setenv("API_KEY", expectedApiKey)
	t.Setenv("NOTIFIER_TIMEOUT", "15m")

	cfg := LoadConfig()
	if cfg == nil {
		t.Fatalf("expected config, got nil")
	}

	if cfg.Server.Port != "8081" {
		t.Errorf("expected port 8081, got %s", cfg.Server.Port)
	}
	if got := cfg.GetServerAddress(); got != "0.0.0.0:8081" {
		t.Errorf("unexpected server address: %s", got)
	}

	if cfg.GetLogLevel() == "" {
		t.Errorf("expected default log level to be set")
	}
	if cfg.GetTimezone() == "" {
		t.Errorf("expected default timezone to be set")
	}
	if cfg.App.APIKey != expectedApiKey {
		t.Errorf("expected %s api key", expectedApiKey)
	}
	if cfg.App.NotifierTimeout != 15*time.Minute {
		t.Errorf("expected NotifierTimeout 15m, got %v", cfg.App.NotifierTimeout)
	}
}

func TestLoadConfig_OverridesDurations(t *testing.T) {
	resetViper()

	// Mandatory to pass validation
	t.Setenv("BOT_TOKEN", "token")
	t.Setenv("PUBLIC_URL", "https://example.com")
	t.Setenv("PORT", "9090")

	// Duration overrides
	t.Setenv("READ_TIMEOUT", "45s")
	t.Setenv("WRITE_TIMEOUT", "1m")
	t.Setenv("IDLE_TIMEOUT", "2m")
	t.Setenv("SHUTDOWN_TIMEOUT", "5s")

	cfg := LoadConfig()

	if cfg.Server.ReadTimeout != 45*time.Second {
		t.Errorf("expected ReadTimeout 45s, got %v", cfg.Server.ReadTimeout)
	}
	if cfg.Server.WriteTimeout != 1*time.Minute {
		t.Errorf("expected WriteTimeout 1m, got %v", cfg.Server.WriteTimeout)
	}
	if cfg.Server.IdleTimeout != 2*time.Minute {
		t.Errorf("expected IdleTimeout 2m, got %v", cfg.Server.IdleTimeout)
	}
	if cfg.Server.ShutdownTimeout != 5*time.Second {
		t.Errorf("expected ShutdownTimeout 5s, got %v", cfg.Server.ShutdownTimeout)
	}
}
