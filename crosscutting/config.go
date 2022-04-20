package crosscutting

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/yzx9/otodo/infrastructure/config"
)

var configManager *viper.Viper

func LoadConfig(dir string) error {
	configManager = viper.New()
	configManager.SetConfigType("yaml")
	configManager.AddConfigPath(dir)

	configManager.SetConfigName("config.yaml")
	if err := configManager.ReadInConfig(); err != nil {
		return fmt.Errorf("fails to load config.yaml: %w", err)
	}

	configManager.SetConfigName("secret.yaml")
	if err := configManager.MergeInConfig(); err != nil {
		return fmt.Errorf("fails to load secret.yaml: %w", err)
	}

	config.SetConfig(configManager)

	return nil
}

func LoadAndWatchConfig(dir string) (<-chan time.Time, error) {
	if err := LoadConfig(dir); err != nil {
		return nil, err
	}

	changed := make(chan time.Time, 1)
	configManager.OnConfigChange(func(e fsnotify.Event) {
		config.SetConfig(configManager)

		select {
		case changed <- time.Now():
		default:
		}
	})

	configManager.WatchConfig()

	return changed, nil
}
