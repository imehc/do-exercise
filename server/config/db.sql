-- 初始化数据 ;
INSERT INTO sys_api VALUES(1, 'apis.SysApi.Create', '创建api', '/apis/v1/apis', 'BUS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(2, 'apis.SysApi.DELETE', 'api通过id删除', '/apis/v1/apis/:id', 'BUS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(3, 'apis.SysApi.Put', 'api通过id更新', '/apis/v1/apis/:id', 'BUS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(4, 'apis.SysApi.Get', 'api通过id获取', '/apis/v1/apis/:id', 'BUS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(5, 'apis.SysApi.GetPage', '获取api列表', '/apis/v1/apis', 'BUS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(6, 'apis.SysApi.GetAll', '获取所有api', '/apis/v1/apis/all', 'BUS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(7, 'apis.SysUser.Create', '创建用户', '/apis/v1/users', 'BUS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(8, 'apis.SysUser.DELETE', '用户通过id删除', '/apis/v1/users/:id', 'BUS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(9, 'apis.SysUser.Put', '用户通过id更新', '/apis/v1/users/:id', 'BUS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(10, 'apis.SysUser.Get', '用户通过id获取', '/apis/v1/users/:id', 'BUS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_api VALUES(11, 'apis.SysUser.GetPage', '获取用户列表', '/apis/v1/users', 'BUS', 'GET','2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
-- TODO: 更多api ;

INSERT INTO sys_dept VALUES(1, 0, '/0/1', '滑水轨迹', 0, 1, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_dept VALUES(2, 1, '/0/1/2', '研发部', 0, 1, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_dept VALUES(3, 1, '/0/1/3', '设计部', 0, 1, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);

INSERT INTO sys_post VALUES(1, '普通职工', 'normal', 0, 1, NULL, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_post VALUES(2, '设计师', 'designer', 0, 1, NULL, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);

INSERT INTO sys_dept_post VALUES (1, 1);
INSERT INTO sys_dept_post VALUES (2, 1);
INSERT INTO sys_dept_post VALUES (2, 2);

INSERT INTO sys_menu VALUES(1, '系统管理', '系统管理', NULL, '/system', '/0/1', 'M', NULL, NULL, 0, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(2, '用户管理', '用户管理', NULL, 'user', '/0/1/2', 'C', NULL, 'admin:sysUser:list', 1, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(3, '', '新增用户', NULL, NULL, '/0/1/2/3', 'F', 'POST', 'admin:sysUser:add', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(4, '', '删除用户', NULL, NULL, '/0/1/2/4', 'F', 'DELETE', 'admin:sysUser:del', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(5, '', '修改用户', NULL, NULL, '/0/1/2/5', 'F', 'PUT', 'admin:sysUser:update', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_menu VALUES(6, '', '查询用户', NULL, NULL, '/0/1/2/6', 'F', 'GET', 'admin:sysUser:get', 2, true, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
-- TODO: 更多菜单 ;

INSERT INTO sys_menu_api VALUES (3, 7);
INSERT INTO sys_menu_api VALUES (4, 8);
INSERT INTO sys_menu_api VALUES (5, 9);
INSERT INTO sys_menu_api VALUES (6, 10);
INSERT INTO sys_menu_api VALUES (7, 11);

INSERT INTO sys_role VALUES(1, '系统管理员', 'admin', true, 1, 0, NULL, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_role VALUES(2, '普通用户', 'normal', false, 1, 0, NULL, 0, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);

INSERT INTO sys_role_api VALUES (2, 1);

INSERT INTO sys_role_menu VALUES (2, 1);

INSERT INTO sys_user VALUES (1, 'admin', '$2a$10$s6J3WwsH9ghXU07.F1I1huUFX3HOa1dVDevhHvG.mjyvtLutd4Toi', NULL, '管理员', NULL, NULL, 1, 1, 2, 1, NULL, 1, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL );
INSERT INTO sys_user VALUES (2, 'normal', '$2a$10$sI.wC3JJeE8JzuNVzws0sODClvjL5Vbz8bQZ8YeN4P7Y2SmYFISxi', NULL, '普通用户', NULL, NULL, 1, 2, 3, 1, NULL, 1, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL );

INSERT INTO sys_dict VALUES (1, '用户性别', 'sys_user_sex', '用户性别列表', 1, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);

INSERT INTO sys_dict_data VALUES (1, 'sys_user_sex', '男', '1', '性别男', 0, 1, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_dict_data VALUES (2, 'sys_user_sex', '女', '2', '性别女', 0, 1, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
INSERT INTO sys_dict_data VALUES (3, 'sys_user_sex', '未知', '0', '未知性别', 0, 1, '2025-03-06 08:19:00.621', 1, '2025-03-06 08:19:00.621', 1, NULL, NULL);
