package ginerrormiddleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

//FormatError generate format error message
func FormatError(httpCode, businessCode int, msg string) (err error) {
	err = fmt.Errorf("%d:%d:%s", httpCode, businessCode, msg)
	return
}

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Errors.Last() == nil {
			return
		}
		err := c.Errors.Last().Err
		c.Errors = c.Errors[:0]
		msg := err.Error()
		strSlice := strings.SplitN(msg, ":", 3)
		businessCode := "100000"
		if len(strSlice) != 3 {
			httpCode := c.Writer.Status()
			out := &AppError{
				Code: businessCode,
				Msg:  err.Error(),
			}
			c.AbortWithStatusJSON(httpCode, out)
			return
		}
		// 格式化错误
		httpCode, err := strconv.Atoi(strSlice[0])
		if err != nil {
			httpCode = http.StatusInternalServerError
		}
		code, _ := strconv.Atoi(strSlice[1])
		if code > 0 {
			businessCode = strconv.Itoa(code)
		}
		out := &AppError{
			Code: businessCode,
			Msg:  strSlice[2],
		}
		c.AbortWithStatusJSON(httpCode, out)

	}
}
