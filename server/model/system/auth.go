package system

type Claims struct {
	Username  string // 用户名
	ID        uint   // 用户ID
	DeptId    int    // 归属部门
	Status    int    // 状态
	PostId    int    // 岗位
	RoleId    int    // 角色
	IsAdmin   bool   // 是否为管理员
	DataScope uint   // 数据权限
}
