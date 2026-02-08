package configloader

import (
	"edge/cmd/config"
	"fmt"

	"github.com/spf13/viper"
)

// Load reads configuration from file or environment variables and
// unmarshals it into the provided Config struct.
func Load() (*config.Config, error) {
	viper.AddConfigPath("./config")
	viper.SetConfigName("default")
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
