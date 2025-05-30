package util

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// TranslateDBError 将数据库错误转换为业务错误码
func TranslateDBError(err error, notFoundTxt string) string {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return "operationTimeout"
	case errors.Is(err, context.Canceled):
		return "requestCanceled"
	case errors.Is(err, gorm.ErrRecordNotFound):
		return notFoundTxt
	default:
		return "queryFailed"
	}
}
