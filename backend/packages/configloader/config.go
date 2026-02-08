package configloader

import (
	"CRM/cmd/config" // Import the new config package
	"fmt"

	"github.com/spf13/viper"
)

// Load reads configuration from file or environment variables and
// unmarshals it into the provided Config struct.
func Load() (*config.Config, error) {
	// Set the path to look for the config file in.
	viper.AddConfigPath("./config")
	viper.SetConfigName("default")
	viper.SetConfigType("toml")

	// Attempt to read the config file.
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg config.Config
	// Unmarshal the config into the struct.
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
