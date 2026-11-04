package request

import "slices"

type Email struct {
	Email string `json:"email" form:"email" binding:"required,email"`
}

// EmailCache 一次邮箱验证码的缓存内容。
type EmailCache struct {
	Code string `json:"code"`
	// UserIds 本次验证码可用于哪些账号。
	//
	// 邮箱在多租户下不唯一（同一个人在 N 个租户里是 N 行 sys_user），一个验证码
	// 只能证明「请求者拥有这个邮箱」，证明不了他要进哪个租户。因此这里存候选集合，
	// 由后续的租户选择收敛到具体账号。已登录的绑定/换绑/改密流程只会有一个元素。
	UserIds []string `json:"user_ids"`
}

// Allows 该验证码是否可用于指定账号
func (c EmailCache) Allows(userId string) bool {
	return slices.Contains(c.UserIds, userId)
}

// SharesUser 两个候选集合是否有交集。用于判断「这封验证码是不是发给同一个人的」。
func (c EmailCache) SharesUser(userIds []string) bool {
	return slices.ContainsFunc(userIds, c.Allows)
}

type BindEmailReq struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

type UserResetPasswordReq struct {
	Email     string `json:"email" binding:"required,email"`
	Code      string `json:"code" binding:"required"`
	Password  string `json:"password" binding:"required"`
	PublicKey string `json:"public_key" binding:"required"`
	// TenantId 要重置哪个租户下的账号。同一邮箱只归属一个租户时可省略；
	// 归属多个租户时服务端会返回候选列表要求显式指定，避免误改到别的租户的口令。
	TenantId string `json:"tenant_id"`
}

type UserModifyPasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	Password    string `json:"password" binding:"required"`
	PublicKey   string `json:"public_key" binding:"required"`
}

type UserModifyProfileReq struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}
