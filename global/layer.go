/**
 * Created by zc on 2020/9/4.
 */
package global

var config *Config

func InitCfg(cfg *Config) {
	config = cfg
}

func Cfg() *Config {
	return config
}
