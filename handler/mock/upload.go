/**
 * Created by zc on 2020/9/4.
 */
package mock

import (
	"github.com/gin-gonic/gin"
	"github.com/pkgms/go/ctr"
	"github.com/zc2638/gotool/utilx"
	"io"
	"net/http"
)

func Upload() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 上传图片
		file, info, err := c.Request.FormFile("file")
		if err != nil {
			ctr.BadRequest(c.Writer, err)
			return
		}

		filename := info.Filename
		imagePath := "uploads/" + filename
		out, err := utilx.CreateFile(imagePath)
		if err != nil {
			ctr.BadRequest(c.Writer, err)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			ctr.BadRequest(c.Writer, err)
			return
		}
		imageUrl := c.Request.Host + "/" + imagePath
		c.String(http.StatusCreated, imageUrl)
	}
}
