/**
 * Created by zc on 2020/9/4.
 */
package global

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	"github.com/pkgms/go/server"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Config struct {
	Server server.Config `json:"server"`
	Logger bool          `json:"logger"`
}

func Environ() *Config {
	cfg := &Config{}
	cfg.Logger = true
	return cfg
}

func ParseConfig(cfgPath string) (*Config, error) {
	if cfgPath != "" {
		viper.SetConfigFile(cfgPath)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			return nil, err
		}
		viper.AddConfigPath(home)
		viper.SetConfigName("config.yaml")
	}
	viper.SetEnvPrefix("MOCK")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	cfg := Environ()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(*os.PathError); ok {
			fmt.Println("Warning: not find config file, use default.")
			return cfg, nil
		}
		return nil, err
	}
	fmt.Println("Using config file:", viper.ConfigFileUsed())
	err := viper.Unmarshal(cfg, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "json"
	})
	return cfg, err
}
