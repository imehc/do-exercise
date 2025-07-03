package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common/status"
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
func Success(c *gin.Context, data interface{}) {
	Response(c, http.StatusOK, data)
}

// Created 创建成功响应
func Created(c *gin.Context) {
	Response(c, http.StatusCreated, "Created")
}

// NoContent 无内容响应
func NoContent(c *gin.Context) {
	Response(c, http.StatusNoContent, "No Content")
}

// BadRequest 请求错误响应
func BadRequest(c *gin.Context, value string) {
	lang := c.GetString("lang")
	Response(c, http.StatusBadRequest, ValidationError{
		Type:    status.BAD_REQUEST_MSG,
		Message: global.I18.Translate(value, lang),
	})
}

// BadRequest 请求错误响应带详情
func BadRequestDetails(c *gin.Context, value string, details []ValidationDetail) {
	lang := c.GetString("lang")
	Response(c, http.StatusBadRequest, ValidationError{
		Type:    status.BAD_REQUEST_MSG,
		Message: global.I18.Translate(value, lang),
		Details: details,
	})
}

// StatusTooManyRequests 访问过于频繁
func StatusTooManyRequests(c *gin.Context) {
	lang := c.GetString("lang")
	text := global.I18.Translate("tooManyRequests", lang)
	Response(c, http.StatusTooManyRequests, text)
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context) {
	Response(c, http.StatusUnauthorized, "Unauthorized")
}

// Forbidden 禁止访问响应
func Forbidden(c *gin.Context, message string) {
	if message != "" {
		Response(c, http.StatusForbidden, message)
		return
	}
	Response(c, http.StatusForbidden, "Forbidden")
}

// NotFound 资源不存在响应
func NotFound(c *gin.Context) {
	Response(c, http.StatusNotFound, "Not Found")
}

// StatusUnprocessableEntity 请求参数错误
func StatusUnprocessableEntity(c *gin.Context) {
	Response(c, http.StatusUnprocessableEntity, "Unprocessable Entity")
}

// NotImplemented 未实现
func NotImplemented(c *gin.Context) {
	Response(c, http.StatusNotImplemented, "Not Implemented")
}

// ServerError 服务器错误响应
func ServerError(c *gin.Context) {
	Response(c, http.StatusInternalServerError, "Internal Server Error")
}
