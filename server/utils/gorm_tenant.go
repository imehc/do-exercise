package utils

import "gorm.io/gorm"

func TenantScope(tenantId uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("tenant_id = ?", tenantId)
	}
}

// 使用示例（查询当前租户的用户）：
// func GetUsers(c *gin.Context) {
// 	tenantID := c.MustGet("tenant_id").(uint)
// 	var users []models.User
// 	db.Scopes(utils.TenantScope(tenantID)).Find(&users)
//  ...
// }
