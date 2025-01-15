package middleware

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 定义一个responseWriterWrapper类型，用于包裹gin.ResponseWriter，以扩展其功能
type responseWriterWrapper struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

// 重写Write方法，实现在响应体内容被写入时同时缓存这些内容
func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	return w.buf.Write(b)
}

func ResponseError() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建一个responseWriterWrapper实例，用于替换当前的ResponseWriter
		writer := &responseWriterWrapper{
			ResponseWriter: c.Writer,        // 使用原ResponseWriter初始化
			buf:            &bytes.Buffer{}, // 初始化一个空的缓冲区
		}
		c.Writer = writer // 将上下文中的Writer替换为我们自定义的writer
		// 继续执行后续的请求处理链
		c.Next()

		if c.Writer.Status() >= http.StatusBadRequest {
			data := writer.buf.String()              // 获取缓冲区中的内容
			text := strings.TrimSuffix(data, "{}\n") // 去掉末尾的{}和换行符
			writer.ResponseWriter.WriteString(text)  // 将处理后的内容写入实际的响应体
			writer.buf.Reset()                       // 清空缓冲区
		} else {
			data := writer.buf.Bytes()
			writer.ResponseWriter.Write([]byte(data)) // 将处理后的内容写入实际的响应体
		}
	}
}
