package config

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/fatih/structs"
	"github.com/jeremywohl/flatten"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const yamlDir = "yaml"
const envDir = "yaml/env"

//go:embed yaml
var configFS embed.FS

var yamlConfigs = []string{
	"base.yaml",
}

type Config struct {
	Env   string
	Addr  string `yaml:"addr"`
	MySql MySql  `yaml:"mysql"`
}

func NewConfig() (*Config, error) {
	var config Config
	viper.SetConfigType("yaml")

	for _, configName := range yamlConfigs {
		configFile, err := configFS.ReadFile(fmt.Sprintf("%s/%s", yamlDir, configName))
		if err != nil {
			return nil, err
		}

		if err = viper.MergeConfig(bytes.NewReader(configFile)); err != nil {
			return nil, err
		}
	}

	err := setupConfigFromEnv()
	if err != nil {
		return nil, err
	}

	fmt.Println(os.Getenv("APP_ENV"))

	if appEnv := os.Getenv("APP_ENV"); appEnv != "" {
		envSpecificConfigFile, envSpecificConfigFileErr := configFS.ReadFile(fmt.Sprintf("%s/%s.yaml", envDir, appEnv))
		if envSpecificConfigFileErr != nil {
			return nil, envSpecificConfigFileErr
		}
		err = viper.MergeConfig(bytes.NewReader(envSpecificConfigFile))
		if err != nil {
			return nil, err
		}
		config.Env = appEnv
	}

	if err = viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func setupConfigFromEnv() error {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	confMap := structs.Map(Config{})

	flat, err := flatten.Flatten(confMap, "", flatten.DotStyle)
	if err != nil {
		return errors.Wrap(err, "Unable to flatten config")
	}

	for key := range flat {
		if err := viper.BindEnv(key); err != nil {
			return errors.Wrapf(err, "Unable to bind env var: %s", key)
		}
	}

	return nil
}
