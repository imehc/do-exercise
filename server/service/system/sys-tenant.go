package system

import (
	"errors"
	"strings"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common"
	commonResponse "github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/model/system/response"
	"github.com/imehc/do-exercise/server/util"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SysTenantService struct{}

// TenantAdminRoleCode 租户管理员角色编码（各租户域内独立，互不冲突）
const TenantAdminRoleCode = "tenant_admin"

// checkCodeDuplicate 检查租户编码是否重复
func (s *SysTenantService) checkCodeDuplicate(db *gorm.DB, code string) error {
	var count int64
	if err := db.Model(&system.SysTenant{}).
		Where("code = ?", code).
		Count(&count).
		Error; err != nil || count > 0 {
		return errors.New("tenantCodeDuplicated")
	}
	return nil
}

// Create 创建租户并自动供应租户管理员（角色 + 全量菜单 + 管理员账号 + Casbin 策略）。
// 平台层调用，经 BypassTenantDB 显式控制 tenant_id。
func (s *SysTenantService) Create(db *gorm.DB, req request.CreateSysTenantReq) (*response.SysTenantResp, error) {
	if err := s.checkCodeDuplicate(db, req.Code); err != nil {
		return nil, err
	}

	tenantId := util.NextID()
	tenant := &system.SysTenant{
		TenantId: tenantId,
		Name:     strings.TrimSpace(req.Name),
		Code:     strings.ToLower(strings.TrimSpace(req.Code)),
		Status:   true,
		Remark:   req.Remark,
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(tenant).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("createTenantFailed")
	}

	// 供应租户管理员角色：挂全量菜单（排除平台独有菜单，租户管理仅平台超管可用）
	adminRole := &system.SysRole{
		Name:     "租户管理员",
		Code:     TenantAdminRoleCode,
		TenantId: tenantId,
	}
	if err := tx.Create(adminRole).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("createTenantFailed")
	}

	var menus []system.SysMenu
	if err := tx.Preload("Apis").Where("id NOT IN ?", global.PlatformOnlyMenuIDs).Find(&menus).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("createTenantFailed")
	}
	if len(menus) > 0 {
		if err := tx.Model(adminRole).Association("Menus").Replace(menus); err != nil {
			tx.Rollback()
			return nil, errors.New("createTenantFailed")
		}
	}

	// 创建/选定租户管理员账号
	// admin_mode=new：新建账号（密码传明文，SysUser.BeforeCreate 会自动 bcrypt，避免二次哈希）
	// admin_mode=existing：复用现有用户（复制到本租户，密码哈希原样复用，避免二次 bcrypt）
	adminUser, err := s.provisionAdminUser(tx, adminRole, tenantId, req)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("createTenantFailed")
	}

	// Casbin 策略在事务提交成功之后同步，保持与角色/用户服务的 fail-closed 语义一致
	if err := s.syncProvisionPolicy(adminRole, menus, adminUser); err != nil {
		global.Log.Error("租户管理员策略同步失败", zap.Error(err))
		return nil, errors.New("createTenantFailed")
	}

	return &response.SysTenantResp{
		TenantId:  tenant.TenantId,
		Name:      tenant.Name,
		Code:      tenant.Code,
		Status:    tenant.Status,
		Remark:    tenant.Remark,
		CreatedAt: tenant.CreatedAt,
		CreatedBy: tenant.CreatedBy,
	}, nil
}

// provisionAdminUser 在事务内创建或复用租户管理员账号并挂载管理员角色。
// admin_mode=new：新建账号，复用传入的明文密码（BeforeCreate 自动 bcrypt）。
// admin_mode=existing：将现有用户复制到本租户（密码哈希原样复用，跳过 BeforeCreate 的二次哈希），
// 并校验其非平台超级管理员、未归属目标租户。
func (s *SysTenantService) provisionAdminUser(
	tx *gorm.DB,
	adminRole *system.SysRole,
	tenantId string,
	req request.CreateSysTenantReq,
) (*system.SysUser, error) {
	var adminUser *system.SysUser

	if req.AdminMode == request.AdminModeExisting {
		if req.AdminUserId == "" {
			return nil, errors.New("adminUserRequired")
		}
		var src system.SysUser
		if err := tx.
			Where("id = ?", req.AdminUserId).
			Where("is_super_admin = ?", false).
			Where("tenant_id != ?", global.PlatformTenantID).
			First(&src).Error; err != nil {
			return nil, errors.New("adminUserNotFound")
		}
		// 校验该用户未归属目标租户（同 username 已有记录）
		var count int64
		if err := tx.Model(&system.SysUser{}).
			Where("tenant_id = ?", tenantId).
			Where("username = ?", src.Username).
			Count(&count).Error; err != nil {
			return nil, errors.New("createTenantFailed")
		}
		if count > 0 {
			return nil, errors.New("adminUserAlreadyAssigned")
		}
		adminUser = &system.SysUser{
			UserId:   util.NextID(),
			Username: src.Username,
			Nickname: src.Nickname,
			Email:    src.Email,
			Avatar:   src.Avatar,
			Password: src.Password,
			TenantId: tenantId,
		}
		// 跳过 BeforeCreate 的密码二次哈希，密码哈希原样复用
		if err := tx.Session(&gorm.Session{SkipHooks: true}).Create(adminUser).Error; err != nil {
			return nil, errors.New("createTenantFailed")
		}
	} else {
		adminUser = &system.SysUser{
			UserId:   util.NextID(),
			Username: req.AdminUsername,
			Nickname: "租户管理员",
			Password: req.AdminPassword,
			TenantId: tenantId,
		}
		if err := tx.Create(adminUser).Error; err != nil {
			return nil, errors.New("createTenantFailed")
		}
	}

	// 挂载租户管理员角色
	if err := tx.Model(adminUser).Association("Roles").Replace([]system.SysRole{*adminRole}); err != nil {
		return nil, errors.New("createTenantFailed")
	}
	return adminUser, nil
}

// syncProvisionPolicy 为租户管理员同步 Casbin p 策略与 g 规则
func (s *SysTenantService) syncProvisionPolicy(role *system.SysRole, menus []system.SysMenu, user *system.SysUser) error {
	enforcer := global.Enforcer
	apis := lo.FlatMap(menus, func(menu system.SysMenu, _ int) []system.SysApi {
		return menu.Apis
	})
	if len(apis) > 0 {
		policies := lo.Map(apis, func(item system.SysApi, _ int) []string {
			return []string{role.Code, role.TenantId, item.Path, item.Method}
		})
		if _, err := enforcer.AddPolicies(policies); err != nil {
			return err
		}
	}
	if _, err := enforcer.AddRoleForUser(user.UserId, role.Code, role.TenantId); err != nil {
		return err
	}
	return nil
}

// Update 更新租户（名称/状态/备注）
func (s *SysTenantService) Update(db *gorm.DB, id string, req request.UpdateSysTenantReq) error {
	tenant, err := s.checkTenantExist(db, id)
	if err != nil {
		return err
	}
	updates := map[string]interface{}{
		"Name":   strings.TrimSpace(req.Name),
		"Remark": req.Remark,
	}
	if req.Status != nil {
		updates["Status"] = *req.Status
	}
	if err := db.Model(tenant).Updates(updates).Error; err != nil {
		return errors.New("updateTenantFailed")
	}
	return nil
}

// Delete 删除租户（软删除 + 停用，禁止后续登录）
// 平台至少保留一个租户，仅剩一个时禁止删除，避免平台无可用租户。
func (s *SysTenantService) Delete(db *gorm.DB, id string) error {
	if _, err := s.checkTenantExist(db, id); err != nil {
		return err
	}
	var count int64
	if err := db.Model(&system.SysTenant{}).Count(&count).Error; err != nil {
		return errors.New("deleteTenantFailed")
	}
	if count <= 1 {
		return errors.New("lastTenantNotDeletable")
	}
	if err := db.Model(&system.SysTenant{}).Where("tenant_id = ?", id).Delete(nil).Error; err != nil {
		return errors.New("deleteTenantFailed")
	}
	return nil
}

// Get 获取租户详情
func (s *SysTenantService) Get(db *gorm.DB, id string) (*response.SysTenantResp, error) {
	tenant, err := s.checkTenantExist(db, id)
	if err != nil {
		return nil, err
	}
	return &response.SysTenantResp{
		TenantId:   tenant.TenantId,
		Name:       tenant.Name,
		Code:       tenant.Code,
		Status:     tenant.Status,
		ExpireTime: tenant.ExpireTime,
		Remark:     tenant.Remark,
		CreatedAt:  tenant.CreatedAt,
		CreatedBy:  tenant.CreatedBy,
		UpdatedAt:  tenant.UpdatedAt,
		UpdatedBy:  tenant.UpdatedBy,
	}, nil
}

// GetList 分页查询租户列表
func (s *SysTenantService) GetList(db *gorm.DB, req request.QuerySysTenantReq) (common.PageResult[response.SysTenantResp], error) {
	var total int64
	// Count 用独立 builder，避免污染后续 Find 的状态
	countDB := db.Model(&system.SysTenant{})
	if req.Name != "" {
		countDB = countDB.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Code != "" {
		countDB = countDB.Where("code LIKE ?", "%"+req.Code+"%")
	}
	if err := countDB.Count(&total).Error; err != nil {
		return common.PageResult[response.SysTenantResp]{}, errors.New("getTenantListFailed")
	}

	pagination := common.Pagination{Page: req.Page, PageSize: req.PageSize}
	pagination.Normalize()

	listDB := db.Model(&system.SysTenant{})
	if req.Name != "" {
		listDB = listDB.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Code != "" {
		listDB = listDB.Where("code LIKE ?", "%"+req.Code+"%")
	}

	var tenants []system.SysTenant
	if err := listDB.
		Scopes(util.Paginate(pagination.PageSize, pagination.Page)).
		Order("created_at DESC").
		Find(&tenants).Error; err != nil {
		return common.PageResult[response.SysTenantResp]{}, errors.New("getTenantListFailed")
	}

	data := make([]response.SysTenantResp, len(tenants))
	for i, t := range tenants {
		data[i] = response.SysTenantResp{
			TenantId:   t.TenantId,
			Name:       t.Name,
			Code:       t.Code,
			Status:     t.Status,
			ExpireTime: t.ExpireTime,
			Remark:     t.Remark,
			CreatedAt:  t.CreatedAt,
			CreatedBy:  t.CreatedBy,
			UpdatedAt:  t.UpdatedAt,
			UpdatedBy:  t.UpdatedBy,
		}
	}
	return common.PageResult[response.SysTenantResp]{
		Data: data,
		Meta: common.PageMeta{Page: pagination.Page, PageSize: pagination.PageSize, Total: total},
	}, nil
}

// checkTenantExist 检查租户是否存在且未被删除
func (s *SysTenantService) checkTenantExist(db *gorm.DB, tenantId string) (*system.SysTenant, error) {
	if tenantId == "" {
		return nil, errors.New("idCannotBeEmpty")
	}
	var tenant system.SysTenant
	if err := db.Where("tenant_id = ?", tenantId).First(&tenant).Error; err != nil {
		return nil, errors.New("tenantNotFound")
	}
	return &tenant, nil
}

// ListEnabled 查询启用的租户（登录页选择器用）
func (s *SysTenantService) ListEnabled(db *gorm.DB) ([]commonResponse.TenantOption, error) {
	var tenants []system.SysTenant
	if err := db.Model(&system.SysTenant{}).
		Where("status = ?", true).
		Order("created_at ASC").
		Find(&tenants).Error; err != nil {
		return nil, errors.New("getTenantListFailed")
	}
	return lo.Map(tenants, func(t system.SysTenant, _ int) commonResponse.TenantOption {
		return commonResponse.TenantOption{TenantId: t.TenantId, Name: t.Name, Code: t.Code}
	}), nil
}

// ListEnabledForUsername 查询指定用户名归属的可用租户（租户切换器、多租户登录候选）。
// 仅返回启用中的业务租户；平台保留租户（platform）不出现在任何租户选择器中，
// 平台管理员直接登录平台域，无需也无法切换。
func (s *SysTenantService) ListEnabledForUsername(db *gorm.DB, username string) ([]commonResponse.TenantOption, error) {
	var tenantIds []string
	if err := db.Model(&system.SysUser{}).
		Where("username = ?", username).
		Where("tenant_id != ?", global.PlatformTenantID).
		Distinct().
		Pluck("tenant_id", &tenantIds).Error; err != nil {
		return nil, errors.New("getTenantListFailed")
	}

	options := make([]commonResponse.TenantOption, 0, len(tenantIds))
	for _, tid := range tenantIds {
		var tenant system.SysTenant
		if err := db.Where("tenant_id = ? AND status = ?", tid, true).First(&tenant).Error; err != nil {
			continue
		}
		options = append(options, commonResponse.TenantOption{
			TenantId: tenant.TenantId,
			Name:     tenant.Name,
			Code:     tenant.Code,
		})
	}
	return options, nil
}

// ListAssignableAdmins 查询可被选作租户管理员的现有用户（创建租户时用）。
// 排除平台超级管理员（is_super_admin=true）与平台保留租户（tenant_id=platform）。
func (s *SysTenantService) ListAssignableAdmins(db *gorm.DB) ([]response.AssignableUser, error) {
	var users []system.SysUser
	if err := db.
		Where("is_super_admin = ?", false).
		Where("tenant_id != ?", global.PlatformTenantID).
		Order("username ASC").
		Find(&users).Error; err != nil {
		return nil, errors.New("getUserListFailed")
	}
	return lo.Map(users, func(u system.SysUser, _ int) response.AssignableUser {
		return response.AssignableUser{
			Id:       u.UserId,
			Username: u.Username,
			Nickname: u.Nickname,
			Email:    u.Email,
			TenantId: u.TenantId,
		}
	}), nil
}

// ListAssignableUsers 查询可分配给指定租户的现有用户（平台层，跨租户）。
// 排除平台超级管理员（is_super_admin=true）、平台保留租户（tenant_id=platform）
// 以及已归属目标租户（目标租户下已存在同名 username）的用户。
func (s *SysTenantService) ListAssignableUsers(db *gorm.DB, tenantId string) ([]response.AssignableUser, error) {
	users, err := s.queryAssignableUsers(db, tenantId)
	if err != nil {
		return nil, err
	}
	return lo.Map(users, func(u system.SysUser, _ int) response.AssignableUser {
		return response.AssignableUser{
			Id:       u.UserId,
			Username: u.Username,
			Nickname: u.Nickname,
			Email:    u.Email,
			TenantId: u.TenantId,
		}
	}), nil
}

// queryAssignableUsers 返回当前可分配的用户实体列表。
func (s *SysTenantService) queryAssignableUsers(db *gorm.DB, tenantId string) ([]system.SysUser, error) {
	var users []system.SysUser
	sub := db.Model(&system.SysUser{}).
		Where("tenant_id = ?", tenantId).
		Select("username")
	if err := db.
		Where("is_super_admin = ?", false).
		Where("tenant_id != ?", global.PlatformTenantID).
		Where("username NOT IN (?)", sub).
		Order("username ASC").
		Find(&users).Error; err != nil {
		return nil, errors.New("getUserListFailed")
	}
	return users, nil
}

// AssignUsers 把选中的现有用户复制到目标租户下：复用原密码哈希与昵称邮箱，
// 并沿用其角色（按角色 code 在目标租户下匹配，目标租户无同名角色则不沿用）。
// 平台超级管理员不可分配；已归属目标租户的用户会被跳过。
func (s *SysTenantService) AssignUsers(db *gorm.DB, tenantId string, userIds []string) error {
	if _, err := s.checkTenantExist(db, tenantId); err != nil {
		return err
	}

	// 目标租户已有角色（按 code 索引），用于沿用角色匹配
	var targetRoles []system.SysRole
	if err := db.Where("tenant_id = ?", tenantId).Find(&targetRoles).Error; err != nil {
		return errors.New("assignUsersFailed")
	}
	roleByCode := lo.SliceToMap(targetRoles, func(r system.SysRole) (string, system.SysRole) {
		return r.Code, r
	})

	// 目标租户已归属的用户名，分配时跳过
	var targetUsernames []string
	if err := db.Model(&system.SysUser{}).
		Where("tenant_id = ?", tenantId).
		Pluck("username", &targetUsernames).Error; err != nil {
		return errors.New("assignUsersFailed")
	}
	existingName := make(map[string]struct{}, len(targetUsernames))
	for _, u := range targetUsernames {
		existingName[u] = struct{}{}
	}

	// 取出可分配的源用户（含角色），排除超管与平台租户
	ids := lo.Uniq(userIds)
	var sourceUsers []system.SysUser
	if err := db.
		Preload("Roles").
		Where("id IN ?", ids).
		Where("is_super_admin = ?", false).
		Where("tenant_id != ?", global.PlatformTenantID).
		Find(&sourceUsers).Error; err != nil {
		return errors.New("assignUsersFailed")
	}
	if len(sourceUsers) == 0 {
		return nil
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	type assignedUser struct {
		user  *system.SysUser
		roles []system.SysRole
	}
	assigned := make([]assignedUser, 0, len(sourceUsers))

	for _, src := range sourceUsers {
		if _, ok := existingName[src.Username]; ok {
			continue
		}
		newUser := &system.SysUser{
			UserId:   util.NextID(),
			Username: src.Username,
			Nickname: src.Nickname,
			Email:    src.Email,
			Avatar:   src.Avatar,
			Password: src.Password, // 复用原密码哈希，避免二次 bcrypt
			TenantId: tenantId,
		}
		// 跳过 BeforeCreate 的密码哈希，密码哈希直接原样复用
		if err := tx.Session(&gorm.Session{SkipHooks: true}).Create(newUser).Error; err != nil {
			tx.Rollback()
			return errors.New("assignUsersFailed")
		}

		// 沿用角色：按 code 匹配目标租户下的同名角色
		roles := make([]system.SysRole, 0, len(src.Roles))
		for _, r := range src.Roles {
			if tr, ok := roleByCode[r.Code]; ok {
				roles = append(roles, tr)
			}
		}
		if len(roles) > 0 {
			if err := tx.Model(newUser).Association("Roles").Replace(roles); err != nil {
				tx.Rollback()
				return errors.New("assignUsersFailed")
			}
		}
		assigned = append(assigned, assignedUser{user: newUser, roles: roles})
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.New("assignUsersFailed")
	}

	// 事务提交成功后同步 Casbin g 规则（global.Enforcer 走独立 adapter）
	for _, a := range assigned {
		if err := s.syncAssignedUserPolicy(a.user, a.roles); err != nil {
			return err
		}
	}
	return nil
}

// syncAssignedUserPolicy 为复制到目标租户的用户同步 Casbin 角色绑定（限定目标租户域）。
func (s *SysTenantService) syncAssignedUserPolicy(user *system.SysUser, roles []system.SysRole) error {
	enforcer := global.Enforcer
	if _, err := enforcer.DeleteRolesForUser(user.UserId, user.TenantId); err != nil {
		return errors.New("assignUsersFailed")
	}
	for _, r := range roles {
		if _, err := enforcer.AddRoleForUser(user.UserId, r.Code, user.TenantId); err != nil {
			return errors.New("assignUsersFailed")
		}
	}
	return nil
}
