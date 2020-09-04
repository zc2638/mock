/**
 * Created by zc on 2020/9/4.
 */
package mock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkgms/go/ctr"
	"github.com/zc2638/gotool/curlx"
	"github.com/zc2638/mock/global"
	"github.com/zc2638/mock/pkg/network"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

type Call struct {
	Sleep    int64   `json:"sleep"`     // ms
	BodySize int     `json:"body_size"` // KB
	Fault    *Fault  `json:"fault"`
	Remote   *Remote `json:"remote"`
}

type Fault struct {
	Percent int `json:"percent"` // 百分比
	Code    int `json:"code"`    // http code
}

type Remote struct {
	Address string `json:"address"`
	Method  string `json:"method"`
	Extend  Extend `json:"extend"`
}

type Extend struct {
	Headers []string `json:"headers"`
}

func Any() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 增加异常模拟方法
		var calls []Call
		if c.Request.Method != http.MethodPost &&
			c.Request.Method != http.MethodPatch &&
			c.Request.Method != http.MethodDelete &&
			c.Request.Method != http.MethodPut {
			c.JSON(http.StatusOK, BuildResult(c, nil, 0, nil))
			return
		}

		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			ctr.BadRequest(c.Writer, err)
			return
		}
		defer c.Request.Body.Close()
		if err := json.Unmarshal(body, &calls); err != nil {
			ctr.BadRequest(c.Writer, err)
			return
		}

		callsLen := len(calls)
		if callsLen == 0 {
			c.JSON(http.StatusOK, BuildResult(c, body, 0, nil))
			return
		}
		call := calls[0]
		// sleep
		CallSleep(call.Sleep)
		// fault
		if code := CallFault(c, call.Fault); code > 0 {
			c.JSON(code, BuildResult(c, nil, 0, nil))
			return
		}
		// bodySize
		if res := CallBodySize(call.BodySize); res != "" {
			c.String(http.StatusOK, res)
			return
		}
		// remote
		resp, err := CallRemote(c, call, calls)
		if err != nil {
			c.JSON(http.StatusOK, BuildResult(c, body, 0, gin.H{
				"status":  "error",
				"message": err.Error(),
			}))
			return
		}
		if resp == nil {
			c.JSON(http.StatusOK, BuildResult(c, body, 0, nil))
			return
		}
		respBody, err := resp.ParseBody()
		if err != nil {
			c.JSON(http.StatusOK, BuildResult(c, body, 0, gin.H{
				"status":  "error",
				"message": err.Error(),
			}))
		}
		c.JSON(http.StatusOK, BuildResult(c, body, 0, string(respBody)))
	}
}

func BuildResult(c *gin.Context, body []byte, nextCode int, nextResponse interface{}) map[string]interface{} {
	return gin.H{
		"timestamp":     time.Now().UnixNano() / 1e3,
		"Request-Host":  c.Request.Host,
		"URL":           c.Request.URL.String(),
		"RequestURI":    c.Request.RequestURI,
		"RemoteAddr":    c.Request.RemoteAddr,
		"Method":        c.Request.Method,
		"Header":        c.Request.Header,
		"Server-Host":   network.Hostname(),
		"Server-Ip":     network.IP(),
		"Body":          string(body),
		"Next-Code":     nextCode,
		"Next-Response": nextResponse,
	}
}

func CallSleep(ms int64) {
	if ms > 0 {
		time.Sleep(time.Millisecond * time.Duration(ms))
	}
}

func CallFault(c *gin.Context, fault *Fault) int {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(100)
	fmt.Printf("%d < %d\n", i, fault.Percent)
	if i > fault.Percent {
		return 0
	}
	return fault.Code
}

func CallBodySize(size int) string {
	if size > 0 {
		var buffer bytes.Buffer
		currentSize := size * 1024
		for i := 0; i < currentSize; i++ {
			buffer.WriteString("a")
		}
		return buffer.String()
	}
	return ""
}

func CallRemote(c *gin.Context, call Call, calls []Call) (*curlx.Response, error) {
	if call.Remote == nil {
		return nil, nil
	}
	r := curlx.NewRequest()
	r.Url = call.Remote.Address
	r.Method = call.Remote.Method
	if len(calls) > 1 {
		newCalls := calls[1:]
		nb, err := json.Marshal(newCalls)
		if err != nil {
			return nil, err
		}
		r.Body = nb
	}
	for _, h := range global.OpenTracingHeaders {
		if hv := c.GetHeader(h); hv != "" {
			r.Header.Set(h, hv)
		}
	}
	for _, h := range call.Remote.Extend.Headers {
		if hv := c.GetHeader(h); hv != "" {
			r.Header.Set(h, hv)
		}
	}
	return r.Do()
}
