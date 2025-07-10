package pkg

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds the application's configuration.
type Config struct {
	ExitCountry             string        `mapstructure:"exit_country"`
	CircuitRotationInterval time.Duration `mapstructure:"circuit_rotation_interval"`
	BootstrapPeers          []string      `mapstructure:"bootstrap_peers"`
	GeoIPDatabasePath       string        `mapstructure:"geoip_database_path"`
}

// LoadConfig loads configuration from file and environment variables.
func LoadConfig() (*Config, error) {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	v.SetDefault("exit_country", "")
	v.SetDefault("circuit_rotation_interval", "15m")
	v.SetDefault("bootstrap_peers", []string{})
	v.SetDefault("geoip_database_path", "GeoLite2-Country.mmdb") // Default GeoIP DB path

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("config.yaml not found, using default values.")
		} else {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	v.AutomaticEnv() // Read environment variables that match.

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}