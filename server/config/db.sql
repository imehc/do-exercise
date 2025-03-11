-- 初始化数据 ;
INSERT INTO sys_api VALUES(1, 'apis.SysApi.Create', '创建接口', '/apis/v1/apis', 'SYS', 'POST','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(2, 'apis.SysApi.Del', '接口通过id删除', '/apis/v1/apis/:id', 'SYS', 'DELETE','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(3, 'apis.SysApi.Put', '接口通过id更新', '/apis/v1/apis/:id', 'SYS', 'PUT','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(4, 'apis.SysApi.Get', '接口通过id获取', '/apis/v1/apis/:id', 'SYS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(5, 'apis.SysApi.GetPage', '获取接口列表', '/apis/v1/apis', 'SYS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(6, 'apis.SysApi.GetAll', '获取所有接口', '/apis/v1/apis/all', 'SYS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(7, 'apis.SysPost.Create', '创建岗位', '/apis/v1/posts', 'SYS', 'POST','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(8, 'apis.SysPost.Del', '岗位通过id删除', '/apis/v1/posts/:id', 'SYS', 'DELETE','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(9, 'apis.SysPost.Put', '岗位通过id更新', '/apis/v1/posts/:id', 'SYS', 'PUT','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(10, 'apis.SysPost.Get', '岗位通过id获取', '/apis/v1/posts/:id', 'SYS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(11, 'apis.SysPost.GetPage', '获取岗位列表', '/apis/v1/posts', 'SYS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(12, 'apis.SysDept.Create', '创建部门', '/apis/v1/depts', 'SYS', 'POST','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(13, 'apis.SysDept.Del', '部门通过id删除', '/apis/v1/depts/:id', 'SYS', 'DELETE','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(14, 'apis.SysDept.Put', '部门通过id更新', '/apis/v1/depts/:id', 'SYS', 'PUT','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(15, 'apis.SysDept.Get', '部门通过id获取', '/apis/v1/depts/:id', 'SYS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(16, 'apis.SysDept.GetPage', '获取部门列表', '/apis/v1/depts', 'SYS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(17, 'apis.SysDept.GetTree', '获取部门树', '/apis/v1/depts/tree', 'SYS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(18, 'apis.SysMenu.Create', '创建菜单', '/apis/v1/menus', 'SYS', 'POST','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(19, 'apis.SysMenu.Del', '菜单通过id删除', '/apis/v1/menus/:id', 'SYS', 'DELETE','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(20, 'apis.SysMenu.Put', '菜单通过id更新', '/apis/v1/menus/:id', 'SYS', 'PUT','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(21, 'apis.SysMenu.Get', '菜单通过id获取', '/apis/v1/menus/:id', 'SYS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(22, 'apis.SysMenu.GetPage', '获取菜单列表', '/apis/v1/menus', 'SYS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(23, 'apis.SysMenu.GetTree', '获取菜单树', '/apis/v1/menus/tree', 'SYS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(24, 'apis.SysRole.Create', '创建角色', '/apis/v1/roles', 'SYS', 'POST','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(25, 'apis.SysRole.Del', '角色通过id删除', '/apis/v1/roles/:id', 'SYS', 'DELETE','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(26, 'apis.SysRole.Put', '角色通过id更新', '/apis/v1/roles/:id', 'SYS', 'PUT','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(27, 'apis.SysRole.Get', '角色通过id获取', '/apis/v1/roles/:id', 'SYS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(28, 'apis.SysRole.GetPage', '获取角色列表', '/apis/v1/roles', 'SYS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(29, 'apis.SysRole.PutData', '角色更新数据权限', '/apis/v1/roles/:id/data-scope', 'SYS', 'PUT','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(30, 'apis.SysRole.PutMenu', '角色更新菜单权限', '/apis/v1/roles/:id/menu-scope', 'SYS', 'PUT','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(31, 'apis.SysRole.PutApi', '角色更新api权限', '/apis/v1/roles/:id/api-scope', 'SYS', 'PUT','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(32, 'apis.SysUser.Create', '创建用户', '/apis/v1/users', 'SYS', 'POST','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(33, 'apis.SysUser.Del', '用户通过id删除', '/apis/v1/users/:id', 'SYS', 'DELETE','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(34, 'apis.SysUser.Put', '用户通过id更新', '/apis/v1/users/:id', 'SYS', 'PUT','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(35, 'apis.SysUser.Get', '用户通过id获取', '/apis/v1/users/:id', 'SYS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(36, 'apis.SysUser.GetPage', '获取用户列表', '/apis/v1/users', 'SYS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(37, 'apis.SysDict.Create', '创建字典', '/apis/v1/dicts', 'SYS', 'POST','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(38, 'apis.SysDict.Del', '字典通过id删除', '/apis/v1/dicts/:id', 'SYS', 'DELETE','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(39, 'apis.SysDict.Put', '字典通过id更新', '/apis/v1/dicts/:id', 'SYS', 'PUT','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(40, 'apis.SysDict.Get', '字典通过id获取', '/apis/v1/dicts/:id', 'SYS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(41, 'apis.SysDict.GetPage', '获取字典列表', '/apis/v1/dicts', 'SYS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(42, 'apis.SysDictData.Create', '创建字典数据', '/apis/v1/dict-data', 'SYS', 'POST','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(43, 'apis.SysDictData.Del', '字典数据通过id删除', '/apis/v1/dict-data/:id', 'SYS', 'DELETE','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(44, 'apis.SysDictData.Put', '字典数据通过id更新', '/apis/v1/dict-data/:id', 'SYS', 'PUT','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(45, 'apis.SysDictData.Get', '字典数据通过id获取', '/apis/v1/dict-data/:id', 'SYS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(46, 'apis.SysDictData.GetPage', '获取字典数据列表', '/apis/v1/dict-data', 'SYS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
-- 业务
INSERT INTO sys_api VALUES(47, 'apis.User.Put', '用户更新信息', '/apis/v1/user', 'BUS', 'PUT','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(48, 'apis.User.Get', '用户获取信息', '/apis/v1/user', 'BUS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(49, 'apis.Menu.GetCompact', '获取精简菜单', '/apis/v1/menu/compact', 'BUS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);

INSERT INTO sys_dept VALUES(1, 0, '/0/1', '滑水轨迹', 0, 1, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_dept VALUES(2, 1, '/0/1/2', '研发部', 0, 1, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_dept VALUES(3, 1, '/0/1/3', '设计部', 0, 1, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);

INSERT INTO sys_post VALUES(1, '普通职工', 'normal', 0, 1, NULL, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_post VALUES(2, '设计师', 'designer', 0, 1, NULL, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);

INSERT INTO sys_dept_post VALUES (1, 1);
INSERT INTO sys_dept_post VALUES (2, 1);
INSERT INTO sys_dept_post VALUES (2, 2);

INSERT INTO sys_menu VALUES(1, '系统管理', '系统管理', NULL, '/system', '/0/1', 'M', NULL, NULL, 0, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(2, '接口管理', '接口管理', 'material-symbols:api', 'api', '/0/1/2', 'C', NULL, 'admin:sysApi:list', 1, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(3, '', '新增接口', '', NULL, '/0/1/2/3', 'F', 'POST', 'admin:sysApi:add', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(4, '', '删除接口', NULL, NULL, '/0/1/2/4', 'F', 'DELETE', 'admin:sysApi:del', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(5, '', '修改接口', NULL, NULL, '/0/1/2/5', 'F', 'PUT', 'admin:sysApi:update', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(6, '', '查询接口', NULL, NULL, '/0/1/2/6', 'F', 'GET', 'admin:sysApi:get', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
-- INSERT INTO sys_menu VALUES(7, '', '查询所有接口', NULL, NULL, '/0/1/2/7', 'F', 'GET', 'admin:sysApi:getAll', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(8, '用户管理', '用户管理', 'material-symbols:person', 'user', '/0/1/8', 'C', NULL, 'admin:sysUser:list', 1, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(9, '', '新增用户', NULL, NULL, '/0/1/8/9', 'F', 'POST', 'admin:sysUser:add', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(10, '', '删除用户', NULL, NULL, '/0/1/8/10', 'F', 'DELETE', 'admin:sysUser:del', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(11, '', '修改用户', NULL, NULL, '/0/1/8/11', 'F', 'PUT', 'admin:sysUser:update', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(12, '', '查询用户', NULL, NULL, '/0/1/8/12', 'F', 'GET', 'admin:sysUser:get', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(13, '角色管理', '角色管理', 'material-symbols:admin-panel-settings', 'role', '/0/1/13', 'C', NULL, 'admin:sysRole:list', 1, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(14, '', '新增角色', NULL, NULL, '/0/1/13/14', 'F', 'POST', 'admin:sysRole:add', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(15, '', '删除角色', NULL, NULL, '/0/1/13/15', 'F', 'DELETE', 'admin:sysRole:del', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(16, '', '修改角色', NULL, NULL, '/0/1/13/16', 'F', 'PUT', 'admin:sysRole:update', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(17, '', '查询角色', NULL, NULL, '/0/1/13/17', 'F', 'GET', 'admin:sysRole:get', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(18, '', '修改角色数据范围', NULL, NULL, '/0/1/13/18', 'F', 'PUT', 'admin:sysRole:putData', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(19, '', '修改角色菜单范围', NULL, NULL, '/0/1/13/19', 'F', 'PUT', 'admin:sysRole:putMenu', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(20, '', '修改角色接口范围', NULL, NULL, '/0/1/13/20', 'F', 'PUT', 'admin:sysRole:putApi', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(21, '菜单管理', '菜单管理', 'material-symbols:menu', 'role', '/0/1/21', 'C', NULL, 'admin:sysMenu:list', 1, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(22, '', '新增菜单', NULL, NULL, '/0/1/21/22', 'F', 'POST', 'admin:sysMenu:add', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(23, '', '删除菜单', NULL, NULL, '/0/1/21/23', 'F', 'DELETE', 'admin:sysMenu:del', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(24, '', '修改菜单', NULL, NULL, '/0/1/21/24', 'F', 'PUT', 'admin:sysMenu:update', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(25, '', '查询菜单', NULL, NULL, '/0/1/21/25', 'F', 'GET', 'admin:sysMenu:get', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
-- INSERT INTO sys_menu VALUES(26, '', '查询菜单树', NULL, NULL, '/0/1/21/26', 'F', 'GET', 'admin:sysMenu:getTree', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
-- TODO: 更多菜单 ;
INSERT INTO sys_menu VALUES(27, '部门管理', '部门管理', 'material-symbols:corporate-fare', 'dept', '/0/1/27', 'C', NULL, 'admin:sysDept:list', 1, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(28, '', '新增部门', NULL, NULL, '/0/1/27/28', 'F', 'POST', 'admin:sysDept:add', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(29, '', '删除部门', NULL, NULL, '/0/1/27/29', 'F', 'DELETE', 'admin:sysDept:del', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(30, '', '修改部门', NULL, NULL, '/0/1/27/30', 'F', 'PUT', 'admin:sysDept:update', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(31, '', '查询部门', NULL, NULL, '/0/1/27/31', 'F', 'GET', 'admin:sysDept:get', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(32, '岗位管理', '岗位管理', 'material-symbols:work', 'post', '/0/1/32', 'C', NULL, 'admin:sysPost:list', 1, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(33, '', '新增岗位', NULL, NULL, '/0/1/32/33', 'F', 'POST', 'admin:sysPost:add', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(34, '', '删除岗位', NULL, NULL, '/0/1/32/34', 'F', 'DELETE', 'admin:sysPost:del', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(35, '', '修改岗位', NULL, NULL, '/0/1/32/35', 'F', 'PUT', 'admin:sysPost:update', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(36, '', '查询岗位', NULL, NULL, '/0/1/32/36', 'F', 'GET', 'admin:sysPost:get', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(37, '字典管理', '字典管理', 'material-symbols:database', 'dict', '/0/1/37', 'C', NULL, 'admin:sysDict:list', 1, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(38, '', '新增字典', NULL, NULL, '/0/1/37/38', 'F', 'POST', 'admin:sysDict:add', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(39, '', '删除字典', NULL, NULL, '/0/1/37/39', 'F', 'DELETE', 'admin:sysDict:del', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(40, '', '修改字典', NULL, NULL, '/0/1/37/40', 'F', 'PUT', 'admin:sysDict:update', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(41, '', '查询字典', NULL, NULL, '/0/1/37/41', 'F', 'GET', 'admin:sysDict:get', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);

INSERT INTO sys_menu_api VALUES (2, 5); -- 接口管理
INSERT INTO sys_menu_api VALUES (3, 1); -- 添加
INSERT INTO sys_menu_api VALUES (4, 2); -- 删除
INSERT INTO sys_menu_api VALUES (5, 3); -- 修改
INSERT INTO sys_menu_api VALUES (5, 4); -- 修改
INSERT INTO sys_menu_api VALUES (6, 5); -- 查询
INSERT INTO sys_menu_api VALUES (8, 36); -- 用户管理
INSERT INTO sys_menu_api VALUES (9, 32); -- 添加
INSERT INTO sys_menu_api VALUES (10, 33); -- 删除
INSERT INTO sys_menu_api VALUES (11, 34); -- 修改
INSERT INTO sys_menu_api VALUES (11, 35); -- 修改
INSERT INTO sys_menu_api VALUES (12, 36); -- 查询
INSERT INTO sys_menu_api VALUES (13, 28); -- 角色管理
INSERT INTO sys_menu_api VALUES (14, 24); -- 添加
INSERT INTO sys_menu_api VALUES (15, 25); -- 删除
INSERT INTO sys_menu_api VALUES (16, 26); -- 修改
INSERT INTO sys_menu_api VALUES (16, 27); -- 修改
INSERT INTO sys_menu_api VALUES (17, 28); -- 查询
INSERT INTO sys_menu_api VALUES (18, 27); -- 修改数据范围
INSERT INTO sys_menu_api VALUES (18, 29); -- 修改数据范围
INSERT INTO sys_menu_api VALUES (19, 27); -- 修改菜单范围
INSERT INTO sys_menu_api VALUES (19, 30); -- 修改菜单范围
INSERT INTO sys_menu_api VALUES (20, 27); -- 修改接口范围
INSERT INTO sys_menu_api VALUES (20, 31); -- 修改接口范围
INSERT INTO sys_menu_api VALUES (21, 22); -- 菜单管理
INSERT INTO sys_menu_api VALUES (22, 18); -- 添加
INSERT INTO sys_menu_api VALUES (23, 19); -- 删除
INSERT INTO sys_menu_api VALUES (24, 20); -- 修改
INSERT INTO sys_menu_api VALUES (24, 21); -- 修改
INSERT INTO sys_menu_api VALUES (25, 22); -- 查询
INSERT INTO sys_menu_api VALUES (27, 16); -- 部门管理
INSERT INTO sys_menu_api VALUES (28, 12); -- 添加
INSERT INTO sys_menu_api VALUES (29, 13); -- 删除
INSERT INTO sys_menu_api VALUES (30, 14); -- 修改
INSERT INTO sys_menu_api VALUES (30, 15); -- 修改
INSERT INTO sys_menu_api VALUES (31, 16); -- 查询
INSERT INTO sys_menu_api VALUES (32, 11); -- 岗位管理
INSERT INTO sys_menu_api VALUES (33, 7); -- 添加
INSERT INTO sys_menu_api VALUES (34, 8); -- 删除
INSERT INTO sys_menu_api VALUES (35, 9); -- 修改
INSERT INTO sys_menu_api VALUES (35, 10); -- 修改
INSERT INTO sys_menu_api VALUES (36, 11); -- 查询
INSERT INTO sys_menu_api VALUES (37, 41); -- 字典管理
INSERT INTO sys_menu_api VALUES (38, 37); -- 添加
INSERT INTO sys_menu_api VALUES (39, 38); -- 删除
INSERT INTO sys_menu_api VALUES (40, 39); -- 修改
INSERT INTO sys_menu_api VALUES (40, 40); -- 修改
INSERT INTO sys_menu_api VALUES (41, 41); -- 查询


INSERT INTO sys_role VALUES(1, '系统管理员', 'admin', true, 1, 0, NULL, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_role VALUES(2, '普通用户', 'normal', false, 1, 0, NULL, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);

INSERT INTO sys_role_api VALUES (2, 1);
INSERT INTO sys_role_api VALUES (2, 3);
INSERT INTO sys_role_api VALUES (2, 4);
INSERT INTO sys_role_api VALUES (2, 5);

INSERT INTO sys_role_menu VALUES (2, 1);
INSERT INTO sys_role_menu VALUES (2, 2);
INSERT INTO sys_role_menu VALUES (2, 3);
INSERT INTO sys_role_menu VALUES (2, 5);
INSERT INTO sys_role_menu VALUES (2, 6);

INSERT INTO sys_user VALUES (1, 'admin', '$2a$10$s6J3WwsH9ghXU07.F1I1huUFX3HOa1dVDevhHvG.mjyvtLutd4Toi', NULL, '管理员', NULL, NULL, 1, 1, 2, 1, NULL, 1, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL );
INSERT INTO sys_user VALUES (2, 'normal', '$2a$10$sI.wC3JJeE8JzuNVzws0sODClvjL5Vbz8bQZ8YeN4P7Y2SmYFISxi', NULL, '普通用户', NULL, NULL, 1, 2, 3, 1, NULL, 1, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL );

INSERT INTO sys_dict VALUES (1, '用户性别', 'sys_user_sex', '用户性别列表', 1, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);

INSERT INTO sys_dict_data VALUES (1, 'sys_user_sex', '男', '1', '性别男', 0, 1, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_dict_data VALUES (2, 'sys_user_sex', '女', '2', '性别女', 0, 1, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_dict_data VALUES (3, 'sys_user_sex', '未知', '0', '未知性别', 0, 1, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
