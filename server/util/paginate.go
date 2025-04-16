package util

import "gorm.io/gorm"

// Paginate 分页
func Paginate(pageSize, pageIndex int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := max((pageIndex-1)*pageSize, 0)
		return db.Offset(offset).Limit(pageSize)
	}
}
