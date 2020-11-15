package lib

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	name         string
	fileName     string
	filePathList []string
	data         *viper.Viper
}

func NewConfig(name string, fileName string, filePathList []string) (*Config, error) {
	config := Config{
		name:         name,
		fileName:     fileName,
		filePathList: filePathList,
	}

	viper.SetConfigName(name)
	for _, path := range filePathList {
		viper.AddConfigPath(path)
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Msg(err.Error())
	}

	config.data = viper.GetViper()
	return &config, nil
}

func (c *Config) GetSubConfig(name string) *Config {
	subConfig := Config{
		name:         fmt.Sprintf("%s.%s", c.name, name),
		fileName:     c.fileName,
		filePathList: c.filePathList,
	}

	subConfig.data = c.data.Sub(subConfig.name)
	return &subConfig
}

func (c *Config) Unmarshal(rawVal interface{}) error {
	return c.data.Unmarshal(rawVal)
}

func (c *Config) Save() error {
	return c.data.WriteConfig()
}

func (c *Config) GetHardwareNameList() []string {
	hardwares := c.data.GetStringMapString("hardwares")
	names := make([]string, 0, len(hardwares))
	for name := range hardwares {
		names = append(names, name)
	}

	return names
}
