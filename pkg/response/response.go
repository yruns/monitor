package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

const (
	ERROR   = -1
	SUCCESS = 0
	INVALID = 1
)

func Result(response *Response, c *gin.Context) {
	//开始时间
	c.JSON(http.StatusOK, response)
}

func R(code int, data interface{}, msg string) *Response {
	return &Response{
		code,
		data,
		msg,
	}
}

func Ok() *Response {
	return R(SUCCESS, nil, "success")
}

func OkWithMessage(message string) *Response {
	return R(SUCCESS, nil, message)
}

func OkWithData(data interface{}) *Response {
	return R(SUCCESS, data, "success")
}

func OkWithDetailed(data interface{}, message string) *Response {
	return R(SUCCESS, data, message)
}

func Fail() *Response {
	return R(ERROR, nil, "failed")
}

func FailWithMessage(message string) *Response {
	return R(ERROR, nil, message)
}

func FailWithError(err error) *Response {
	return R(ERROR, nil, err.Error())
}

func FailWithDetailed(data interface{}, message string) *Response {
	return R(ERROR, data, message)
}

func FailWithStatusCode(status int, message string, c *gin.Context) {
	c.JSON(status, Response{
		-1,
		nil,
		message,
	})
}
