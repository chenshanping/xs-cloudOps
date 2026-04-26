package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PageResponse struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

const (
	SUCCESS      = 200
	ERROR        = 500
	UNAUTHORIZED = 401
	FORBIDDEN    = 403
	NOT_FOUND    = 404
	BAD_REQUEST  = 400
)

func Result(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func Ok(c *gin.Context) {
	Result(c, SUCCESS, "success", nil)
}

func OkWithMessage(c *gin.Context, message string) {
	Result(c, SUCCESS, message, nil)
}

func OkWithData(c *gin.Context, data interface{}) {
	Result(c, SUCCESS, "success", data)
}

func OkWithPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	Result(c, SUCCESS, "success", PageResponse{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func Fail(c *gin.Context, message string) {
	Result(c, ERROR, message, nil)
}

func FailWithCode(c *gin.Context, code int, message string) {
	Result(c, code, message, nil)
}

func Unauthorized(c *gin.Context, message string) {
	Result(c, UNAUTHORIZED, message, nil)
}

func Forbidden(c *gin.Context, message string) {
	Result(c, FORBIDDEN, message, nil)
}

func NotFound(c *gin.Context, message string) {
	Result(c, NOT_FOUND, message, nil)
}

func BadRequest(c *gin.Context, message string) {
	Result(c, BAD_REQUEST, message, nil)
}
