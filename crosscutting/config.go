package crosscutting

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/yzx9/otodo/infrastructure/config"
)

func LoadConfig(dir string) (*viper.Viper, error) {
	ins := viper.New()
	ins.SetConfigType("yaml")
	ins.AddConfigPath(dir)

	ins.SetConfigName("config.yaml")
	if err := ins.ReadInConfig(); err != nil {
		return ins, fmt.Errorf("fails to load config.yaml: %w", err)
	}

	ins.SetConfigName("secret.yaml")
	if err := ins.MergeInConfig(); err != nil {
		return ins, fmt.Errorf("fails to load secret.yaml: %w", err)
	}

	config.SetConfig(ins)

	return ins, nil
}

func LoadAndWatchConfig(dir string) (*viper.Viper, error) {
	ins, err := LoadConfig(dir)
	if err != nil {
		return nil, err
	}

	ins.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed: ", e.Name)
		config.SetConfig(ins)
	})

	ins.WatchConfig()

	return ins, nil
}
