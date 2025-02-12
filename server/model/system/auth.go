package system

type Claims struct {
	Username string // 用户名
	ID       uint   // 用户ID
	DeptId   int    // 归属部门
	Status   int    // 状态
	PostId   int    // 岗位
	RoleId   int    // 角色
}
