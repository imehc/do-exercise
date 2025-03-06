package system

import (
	"errors"
	"sort"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	sysRes "github.com/imehc/do-exercise/server/model/system/response"
	"github.com/imehc/do-exercise/server/pkg/utils/scope"
	"github.com/imehc/do-exercise/server/utils"
	"gorm.io/gorm"
)

type DeptService struct{}

// 创建部门
func (d *DeptService) CreateDept(request request.DeptRequest, createBy uint) (err error) {
	db := global.DB

	// 检查是否有效的父级部门
	if *request.ParentId != 0 {
		var parentDept system.Dept
		result := db.
			First(&parentDept, request.ParentId)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				return errors.New("父级部门不存在")
			}
			return result.Error
		}
	}

	dept := system.Dept{
		ParentId: request.ParentId,
		Name:     request.Name,
		ControlWrapper: model.ControlWrapper{
			CreateBy: createBy,
		},
	}

	if request.Sort != 0 {
		dept.Sort = request.Sort
	}
	if request.Status != 0 {
		dept.Status = request.Status
	}

	tx := db.Begin()
	err = tx.
		Create(&dept).
		Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if *dept.ParentId == 0 {
		dept.Path = utils.FormatFullpath(uint(*dept.ParentId), dept.DeptId, "")
	} else {
		var parentDept system.Dept
		err = tx.First(&parentDept, *request.ParentId).Error
		if err != nil {
			return err
		}
		dept.Path = utils.FormatFullpath(uint(*dept.ParentId), dept.DeptId, parentDept.Path)
	}

	if err = tx.Model(&dept).Update("path", dept.Path).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return
}

// 删除部门
func (d *DeptService) DeleteDept(param request.DeptParam, deleteBy uint) (err error) {
	db := global.DB

	var dept system.Dept
	result := db.
		Unscoped().
		First(&dept, param.DeptId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("部门不存在")
		}
		return result.Error
	}

	if !dept.DeleteAt.Time.IsZero() {
		return errors.New("部门已删除")
	}

	// 检查该部门下是否有其它用户
	var user system.User
	result = db.
		Model(system.User{}).
		Where("dept_id = ?", param.DeptId).
		Find(&user)
	if result.RowsAffected > 0 {
		return errors.New("该部门下存在用户，无法删除")
	}

	db.
		Model(system.Dept{}).
		Where("dept_id = ?", param.DeptId).
		Update("deleted_by", deleteBy).
		Delete(&dept)
	return nil
}

// 更新部门
func (d *DeptService) UpdateDept(param request.DeptParam, request request.DeptRequest, updateBy uint) (err error) {
	db := global.DB

	var dept system.Dept
	result := db.
		First(&dept, param.DeptId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("部门不存在")
		}
		return result.Error
	}

	dept.Name = request.Name
	dept.Sort = request.Sort
	dept.Status = request.Status
	dept.ParentId = request.ParentId
	dept.ControlWrapper = model.ControlWrapper{
		UpdateBy: updateBy,
	}

	if *dept.ParentId == 0 {
		dept.Path = utils.FormatFullpath(uint(*dept.ParentId), dept.DeptId, "")
	} else {
		var parentDept system.Dept
		err = db.First(&parentDept, *request.ParentId).Error
		if err != nil {
			return err
		}
		dept.Path = utils.FormatFullpath(uint(*dept.ParentId), dept.DeptId, parentDept.Path)
	}

	db.
		Model(system.Dept{}).
		Where("dept_id = ?", param.DeptId).
		Updates(&dept).
		Omit("dept_id", "created_at")

	return nil
}

// 查询单个部门
func (d *DeptService) GetDep(param request.DeptParam) (response sysRes.DeptItem, err error) {
	db := global.DB

	var dept system.Dept
	result := db.
		First(&dept, param.DeptId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return response, errors.New("部门不存在")
		}
		return response, result.Error
	}

	response.ID = dept.DeptId
	response.DeptRequest = request.DeptRequest{
		Name:        dept.Name,
		ParentId:    dept.ParentId,
		SortWrapper: dept.SortWrapper,
		StatusWrapper: model.StatusWrapper{
			Status: dept.Status,
		},
	}
	response.ControlWrapper = dept.ControlWrapper
	return
}

// 查询部门树列表
func (d *DeptService) GetDeptTreeList(s common.ScopeData) (response []sysRes.DeptResponse, err error) {
	db := global.DB
	// 应用数据权限过滤
	db = scope.GetDataScope(db, &s, "sys_dept")

	var depts []system.Dept
	if err := db.
		Order("sort DESC").
		Order("dept_id ASC").
		Find(&depts).Error; err != nil {
		return []sysRes.DeptResponse{}, err
	}

	// 构建ID到部门的映射
	deptMap := make(map[uint]*system.Dept)
	for i := range depts {
		deptMap[depts[i].DeptId] = &depts[i]
	}

	// 构建树形结构
	response = make([]sysRes.DeptResponse, 0)
	for _, dept := range depts {
		if *dept.ParentId == 0 {
			// 根节点
			root := sysRes.DeptResponse{
				DeptItem: sysRes.DeptItem{
					IDWrapper: common.IDWrapper{
						ID: dept.DeptId,
					},
					ControlWrapper: dept.ControlWrapper,
					DeptRequest: request.DeptRequest{
						ParentId:    dept.ParentId,
						Name:        dept.Name,
						SortWrapper: dept.SortWrapper,
						StatusWrapper: model.StatusWrapper{
							Status: dept.Status,
						},
					},
				},
				Children: make([]sysRes.DeptResponse, 0), // 初始化为空数组
			}
			response = append(response, root)
		}
	}

	// 递归构建子树
	for i := range response {
		response[i].Children = d.buildDeptResponseSubTree(deptMap, uint(response[i].ID))
	}

	return
}

// 部门树
func (d *DeptService) GetDeptTree(s *common.ScopeData) (response []sysRes.DeptTree, err error) {
	db := global.DB
	db = scope.GetDataScope(db, s, "sys_dept")

	var depts []system.Dept
	if err := db.
		Order("sort DESC").
		Order("dept_id ASC").
		Find(&depts).Error; err != nil {
		return []sysRes.DeptTree{}, err
	}

	// 构建ID到部门的映射
	deptMap := make(map[uint]*system.Dept)
	for i := range depts {
		deptMap[depts[i].DeptId] = &depts[i]
	}

	// 构建树形结构
	roots := make([]sysRes.DeptTree, 0)
	for _, dept := range depts {
		if *dept.ParentId == 0 {
			// 根节点
			root := sysRes.DeptTree{
				ID:       int(dept.DeptId),
				Label:    dept.Name,
				Children: make([]sysRes.DeptTree, 0), // 初始化为空数组
			}
			roots = append(roots, root)
		}
	}

	// 递归构建子树
	for i := range roots {
		roots[i].Children = d.buildDeptSubTree(deptMap, uint(roots[i].ID))
	}

	return roots, nil
}

// 递归构建子树
func (d *DeptService) buildDeptSubTree(deptMap map[uint]*system.Dept, parentId uint) []sysRes.DeptTree {
	children := make([]sysRes.DeptTree, 0)

	// 遍历所有部门，找到当前父节点的直接子节点
	for id, dept := range deptMap {
		if dept != nil && *dept.ParentId == parentId {
			child := sysRes.DeptTree{
				ID:       int(id),
				Label:    dept.Name,
				Children: make([]sysRes.DeptTree, 0), // 初始化为空数组
			}
			// 递归构建该子节点的子树
			child.Children = d.buildDeptSubTree(deptMap, id)
			children = append(children, child)
		}
	}

	// 按照Sort和ID排序
	if len(children) > 1 {
		sort.Slice(children, func(i, j int) bool {
			deptI := deptMap[uint(children[i].ID)]
			deptJ := deptMap[uint(children[j].ID)]
			// 优先按Sort降序
			if deptI.Sort != deptJ.Sort {
				return deptI.Sort > deptJ.Sort
			}
			// Sort相同时按ID升序
			return children[i].ID < children[j].ID
		})
	}

	return children
}

// 递归构建子树
func (d *DeptService) buildDeptResponseSubTree(deptMap map[uint]*system.Dept, parentId uint) []sysRes.DeptResponse {
	children := make([]sysRes.DeptResponse, 0)

	// 遍历所有部门，找到当前父节点的直接子节点
	for id, dept := range deptMap {
		if dept != nil && *dept.ParentId == parentId {
			child := sysRes.DeptResponse{
				DeptItem: sysRes.DeptItem{
					IDWrapper: common.IDWrapper{
						ID: dept.DeptId,
					},
					ControlWrapper: dept.ControlWrapper,
					DeptRequest: request.DeptRequest{
						ParentId:    dept.ParentId,
						Name:        dept.Name,
						SortWrapper: dept.SortWrapper,
						StatusWrapper: model.StatusWrapper{
							Status: dept.Status,
						},
					},
				},
				Children: make([]sysRes.DeptResponse, 0), // 初始化为空数组
			}
			// 递归构建该子节点的子树
			child.Children = d.buildDeptResponseSubTree(deptMap, id)
			children = append(children, child)
		}
	}

	// 按照Sort和ID排序
	if len(children) > 1 {
		sort.Slice(children, func(i, j int) bool {
			deptI := deptMap[uint(children[i].ID)]
			deptJ := deptMap[uint(children[j].ID)]
			// 优先按Sort降序
			if deptI.Sort != deptJ.Sort {
				return deptI.Sort > deptJ.Sort
			}
			// Sort相同时按ID升序
			return children[i].ID < children[j].ID
		})
	}

	return children
}
