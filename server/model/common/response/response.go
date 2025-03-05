package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	OK                    = http.StatusOK
	CREATED               = http.StatusCreated
	NO_CONTENT            = http.StatusNoContent
	BAD_REQUEST           = http.StatusBadRequest
	UNAUTHORIZED          = http.StatusUnauthorized
	FORBIDDEN             = http.StatusForbidden
	NOT_FOUND             = http.StatusNotFound
	NOT_IMPLEMENTED       = http.StatusNotImplemented
	INTERNAL_SERVER_ERROR = http.StatusInternalServerError
)

// Response 响应处理
func Response(c *gin.Context, httpStatus int, data interface{}) {
	// 如果数据为字符串类型，直接返回文本
	if str, ok := data.(string); ok {
		c.String(httpStatus, str)
		return
	}
	// 其他类型返回 JSON
	c.JSON(httpStatus, data)
}

// Success 成功响应
func Success(data interface{}, c *gin.Context) {
	Response(c, OK, data)
}

// Created 创建成功响应
func Created(c *gin.Context) {
	Response(c, CREATED, "Created")
}

// NoContent 无内容响应
func NoContent(c *gin.Context) {
	Response(c, NO_CONTENT, "No Content")
}

// BadRequest 请求错误响应
func BadRequest(message string, c *gin.Context) {
	Response(c, BAD_REQUEST, message)
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context) {
	Response(c, UNAUTHORIZED, "Unauthorized")
}

// Forbidden 禁止访问响应
func Forbidden(c *gin.Context) {
	Response(c, FORBIDDEN, "Forbidden")
}

// NotFound 资源不存在响应
func NotFound(c *gin.Context) {
	Response(c, NOT_FOUND, "Not Found")
}

// NotImplemented 未实现
func NotImplemented(c *gin.Context) {
	Response(c, NOT_IMPLEMENTED, "Not Implemented")
}

// ServerError 服务器错误响应
func ServerError(c *gin.Context) {
	Response(c, INTERNAL_SERVER_ERROR, "Internal Server Error")
}
