/**
 * Created by zc on 2020/9/4.
 */
package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pkgms/go/ctr"
	"github.com/zc2638/mock/handler/mock"
)

func Init(e *gin.Engine) {
	e.GET("/", func(c *gin.Context) {
		ctr.OK(c.Writer, "Hello World!")
	})
	e.Any("/mock/any", mock.Any())
}
