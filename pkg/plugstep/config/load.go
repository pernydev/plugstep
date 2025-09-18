package config

import (
	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/log"
)

func LoadPlugstepConfig(configLocation string) (*PlugstepConfig, error) {
	var config PlugstepConfig
	log.Debug("loading config", "configLocation", configLocation)

	if _, err := toml.DecodeFile(configLocation, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
