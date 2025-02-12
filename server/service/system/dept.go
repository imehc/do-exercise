package system

import (
	"errors"
	"fmt"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/model/system/response"
	"gorm.io/gorm"
)

type DeptService struct{}

// 创建部门
func (d *DeptService) CreateDept(request request.DeptRequest, createdBy uint) (err error) {
	db := global.DB

	dept := system.Dept{
		ParentId: request.ParentId,
		Name:     request.Name,
		Leader:   request.Leader,
		ControlBy: global.ControlBy{
			CreatedBy: createdBy,
		},
	}

	if request.Sort != 0 {
		dept.Sort = request.Sort
	}
	if request.Phone != "" {
		dept.Phone = request.Phone
	}
	if request.Email != "" {
		dept.Email = request.Email
	}
	if request.Status != 0 {
		dept.Status = request.Status
	}

	err = db.Create(&dept).Error
	return
}

// 删除部门
func (d *DeptService) DeleteDept(param request.DeptParam, deletedBy uint) (err error) {
	db := global.DB

	var dept system.Dept
	result := db.Unscoped().First(&dept, param.DeptId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("部门不存在")
		}
		return result.Error
	}

	if !dept.DeletedAt.Time.IsZero() {
		return errors.New("部门已删除")
	}

	// 检查该部门下是否有其它用户
	var user system.User
	result = db.Model(system.User{}).Where("dept_id = ?", param.DeptId).Find(&user)
	if result.RowsAffected > 0 {
		return errors.New("该部门下存在用户，无法删除")
	}

	db.Model(system.Dept{}).Where("id = ?", param.DeptId).Update("deleted_by", deletedBy).Delete(&dept)
	return nil
}

// 更新部门
func (d *DeptService) UpdateDept(param request.DeptParam, request request.DeptRequest, updatedBy uint) (err error) {
	db := global.DB

	var dept system.Dept
	result := db.First(&dept, param.DeptId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("部门不存在")
		}
		return result.Error
	}

	dept.Email = request.Email
	dept.Leader = request.Leader
	dept.Name = request.Name
	dept.Phone = request.Phone
	dept.Sort = request.Sort
	dept.Status = request.Status
	dept.ParentId = request.ParentId
	dept.ControlBy = global.ControlBy{
		UpdatedBy: updatedBy,
	}

	db.Model(system.Dept{}).Where("id = ?", param.DeptId).Updates(&dept).Omit("id", "created_at")

	return nil
}

// 查询单个部门
func (d *DeptService) GetDep(param request.DeptParam) (response.DeptItem, error) {
	db := global.DB

	var deptRep response.DeptItem
	var dept system.Dept
	result := db.First(&dept, param.DeptId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return deptRep, errors.New("部门不存在")
		}
		return deptRep, result.Error
	}

	return response.DeptItem{
		Model: dept.Model,
		DeptRequest: request.DeptRequest{
			Email:    dept.Email,
			Leader:   dept.Leader,
			Name:     dept.Name,
			ParentId: dept.ParentId,
			Phone:    dept.Phone,
			Sort:     dept.Sort,
			Status:   dept.Status,
		},
	}, nil
}

// 查询部门
func (d *DeptService) GetDeptList(query request.DeptQueryParams) (response.DeptResponse, error) {
	db := global.DB

	var total int64
	var originDepts []system.Dept
	offset := (query.Page - 1) * query.PageSize
	depts := response.DeptResponse{}
	err := db.Model(&system.Dept{}).Where(fmt.Sprintf("name LIKE '%%%s%%'", query.Name)).Count(&total).Offset(offset).Limit(query.PageSize).Find(&originDepts).Error
	if err != nil {
		return depts, err
	}
	depts.Meta.Page = query.Page
	depts.Meta.PageSize = query.PageSize
	depts.Meta.Total = total

	depts.Data = make([]response.DeptItem, len(originDepts))
	for i, dept := range originDepts {
		depts.Data[i].Model = dept.Model
		depts.Data[i].ParentId = dept.ParentId
		depts.Data[i].Name = dept.Name
		depts.Data[i].Sort = dept.Sort
		depts.Data[i].Leader = dept.Leader
		depts.Data[i].Phone = dept.Phone
		depts.Data[i].Email = dept.Email
		depts.Data[i].Status = dept.Status
	}

	return depts, nil
}

// 部门树
func (d *DeptService) GetDeptTree() ([]response.DeptTree, error) {
	db := global.DB
	var depts []system.Dept
	if err := db.Find(&depts).Error; err != nil {
		return []response.DeptTree{}, err
	}

	// 创建一个映射来存储部门及其子部门
	deptMap := make(map[int][]*system.Dept)
	for _, dept := range depts {
		if dept.ParentId != nil {
			deptMap[*dept.ParentId] = append(deptMap[*dept.ParentId], &dept)
		} else {
			deptMap[0] = append(deptMap[0], &dept)
		}
	}

	// 递归构建部门树
	var buildTree func(parentId int) []response.DeptTree
	buildTree = func(parentId int) []response.DeptTree {
		var tree []response.DeptTree
		for _, dept := range deptMap[parentId] {
			// 递归构建子部门
			children := buildTree(int(dept.ID))
			if children == nil {
				// 如果没有子部门，则设置 children 为空数组
				children = []response.DeptTree{}
			}
			tree = append(tree, response.DeptTree{
				ID:       int(dept.ID),
				Label:    dept.Name,
				Children: children,
			})
		}
		return tree
	}

	// 从 parentId 为 0 的部门开始构建树
	return buildTree(0), nil
}
