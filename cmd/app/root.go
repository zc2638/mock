/*
Copyright © 2020 zc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package app

import (
	"github.com/gin-gonic/gin"
	"github.com/pkgms/go/server"
	"github.com/spf13/cobra"
	"github.com/zc2638/mock/global"
	"github.com/zc2638/mock/handler"
	"net/http"
	"os"
)

var cfgFile string

func NewServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mock",
		Short: "mock service",
		Long:  `Mock Service.`,
		RunE:  Run,
	}
	cfgFilePath := os.Getenv("LUBAN_CONFIG")
	if cfgFilePath == "" {
		cfgFilePath = global.DefaultConfigPath
	}
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", cfgFilePath, "config file (default is $HOME/config.yaml)")
	return cmd
}

func Run(cmd *cobra.Command, args []string) error {
	cfg, err := global.ParseConfig(cfgFile)
	if err != nil {
		return err
	}
	global.InitCfg(cfg)

	engine := gin.New()
	engine.Use(gin.Recovery(), Cors())
	if cfg.Logger {
		engine.Use(gin.Logger())
	}
	handler.Init(engine)
	s := server.New(&cfg.Server)
	s.Handler = engine
	return s.Run()
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, X-Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
