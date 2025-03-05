package system

type Claims struct {
	Username  string // 用户名
	UserId    uint   // 用户ID
	DeptId    uint   // 归属部门
	Status    uint   // 状态
	PostId    uint   // 岗位
	RoleId    uint   // 角色
	IsAdmin   bool   // 是否为管理员
	DataScope uint   // 数据权限
}
