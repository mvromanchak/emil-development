package sdk

import (
	"flag"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// LoadConfig loads config in order  flag >  file
func LoadConfig() (*viper.Viper, error) {
	flag.String("config", "./config/settings.json", "config file location")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return nil, errors.Wrap(err, "LoadConfig: unable to read flags")
	}
	configPath := viper.GetString("config")
	viper.SetConfigFile(configPath)

	err = viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "LoadConfig: unable to read config")
	}

	v := viper.GetViper()
	return v, nil
}
