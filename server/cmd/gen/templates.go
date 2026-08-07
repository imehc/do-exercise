package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const bt = "`"

// goTag 生成字段的 gorm 标签片段
func (f Field) goTag() string {
	segs := []string{}
	switch f.Kind {
	case "string":
		segs = append(segs, "size:255")
	case "text":
		segs = append(segs, "type:text")
	case "bool":
		segs = append(segs, "default:false")
	}
	if f.Comment != "" {
		segs = append(segs, "comment:"+f.Comment)
	}
	return strings.Join(segs, ";")
}

func (f Field) commentSuffix() string {
	if f.Comment == "" {
		return ""
	}
	return " // " + f.Comment
}

// modelField 生成结构体字段行，含 json/gorm 标签
func (f Field) modelField() string {
	return fmt.Sprintf("\t%s %s %sjson:%q gorm:%q%s%s\n",
		f.Pascal, f.GoType, bt, f.Name, f.goTag(), bt, f.commentSuffix())
}

// genModel 生成 model/system/{file}.go
func (spec *Spec) genModel() string {
	var b strings.Builder
	b.WriteString("package system\n\n")
	b.WriteString("import (\n\t\"time\"\n\n\t\"github.com/imehc/do-exercise/server/model\"\n\t\"gorm.io/gorm\"\n)\n\n")
	fmt.Fprintf(&b, "// %s %s\ntype %s struct {\n\tmodel.IdWrapper\n", spec.Name, spec.Desc, spec.Name)
	for _, f := range spec.Fields {
		b.WriteString(f.modelField())
	}
	b.WriteString("\tCreatedAt time.Time      " + bt + `json:"created_at" gorm:"column:created_at;autoCreateTime;comment:创建时间"` + bt + "\n")
	b.WriteString("\tUpdatedAt time.Time      " + bt + `json:"updated_at" gorm:"column:updated_at;autoUpdateTime;comment:更新时间"` + bt + "\n")
	b.WriteString("\tDeletedAt gorm.DeletedAt " + bt + `json:"-" gorm:"column:deleted_at;index;comment:删除时间"` + bt + "\n")
	b.WriteString("}\n")
	return b.String()
}

// genRequest 生成 model/system/request/{file}.go
func (spec *Spec) genRequest() string {
	var b strings.Builder
	b.WriteString("package request\n\n")
	fmt.Fprintf(&b, "type Create%sReq struct {\n", spec.Name)
	for _, f := range spec.Fields {
		fmt.Fprintf(&b, "\t%s %s %sjson:%q%s%s\n", f.Pascal, f.GoType, bt, f.Name, bt, f.commentSuffix())
	}
	b.WriteString("}\n\n")
	fmt.Fprintf(&b, "type Update%sReq struct {\n\tId uint %sjson:\"id\" binding:\"required\"%s\n", spec.Name, bt, bt)
	for _, f := range spec.Fields {
		fmt.Fprintf(&b, "\t%s %s %sjson:%q%s%s\n", f.Pascal, f.GoType, bt, f.Name, bt, f.commentSuffix())
	}
	b.WriteString("}\n")
	return b.String()
}

// genResponse 生成 model/system/response/{file}.go
func (spec *Spec) genResponse() string {
	var b strings.Builder
	b.WriteString("package response\n\nimport \"time\"\n\n")
	fmt.Fprintf(&b, "type %sResp struct {\n\tId uint %sjson:\"id\"%s\n", spec.Name, bt, bt)
	for _, f := range spec.Fields {
		fmt.Fprintf(&b, "\t%s %s %sjson:%q,omitzero%s\n", f.Pascal, f.GoType, bt, f.Name, bt)
	}
	b.WriteString("\tCreatedAt time.Time " + bt + `json:"created_at,omitzero"` + bt + "\n")
	b.WriteString("\tUpdatedAt time.Time " + bt + `json:"updated_at,omitzero"` + bt + "\n")
	b.WriteString("}\n")
	return b.String()
}

// genService 生成 service/system/{file}.go（遵循现有 db 首参约定）
func (spec *Spec) genService() string {
	name := spec.Name
	var b strings.Builder
	b.WriteString("package system\n\n")
	b.WriteString("import (\n\t\"errors\"\n\n\t\"github.com/imehc/do-exercise/server/global\"\n")
	b.WriteString("\t\"github.com/imehc/do-exercise/server/model/common\"\n\t\"github.com/imehc/do-exercise/server/model/system\"\n")
	b.WriteString("\t\"github.com/imehc/do-exercise/server/model/system/request\"\n\t\"github.com/imehc/do-exercise/server/model/system/response\"\n")
	b.WriteString("\t\"github.com/imehc/do-exercise/server/util\"\n\t\"go.uber.org/zap\"\n\t\"gorm.io/gorm\"\n)\n\n")
	fmt.Fprintf(&b, "type %sService struct{}\n\n", name)

	// Create
	fmt.Fprintf(&b, "// Create 创建%s\n", name)
	fmt.Fprintf(&b, "func (s *%sService) Create(db *gorm.DB, req request.Create%sReq) (uint, error) {\n\tentity := system.%s{\n", name, name, name)
	for _, f := range spec.Fields {
		fmt.Fprintf(&b, "\t\t%s: req.%s,\n", f.Pascal, f.Pascal)
	}
	fmt.Fprintf(&b, "\t}\n\tif err := db.Create(&entity).Error; err != nil {\n\t\tglobal.Log.Error(\"创建%s失败\", zap.Error(err))\n\t\treturn 0, errors.New(\"create%sFailed\")\n\t}\n\treturn entity.Id, nil\n}\n\n", name, name)

	// Update
	fmt.Fprintf(&b, "// Update 更新%s\n", name)
	fmt.Fprintf(&b, "func (s *%sService) Update(db *gorm.DB, req request.Update%sReq) error {\n", name, name)
	fmt.Fprintf(&b, "\tentity := &system.%s{}\n\tif err := db.Where(\"id = ?\", req.Id).First(entity).Error; err != nil {\n\t\treturn errors.New(\"%sNotFound\")\n\t}\n", name, name)
	for _, f := range spec.Fields {
		fmt.Fprintf(&b, "\tentity.%s = req.%s\n", f.Pascal, f.Pascal)
	}
	fmt.Fprintf(&b, "\tif err := db.Model(entity).Updates(entity).Error; err != nil {\n\t\tglobal.Log.Error(\"更新%s失败\", zap.Error(err))\n\t\treturn errors.New(\"update%sFailed\")\n\t}\n\treturn nil\n}\n\n", name, name)

	// Get
	fmt.Fprintf(&b, "// Get 获取单个%s\n", name)
	fmt.Fprintf(&b, "func (s *%sService) Get(db *gorm.DB, id uint) (*response.%sResp, error) {\n", name, name)
	fmt.Fprintf(&b, "\tentity := system.%s{}\n\tif err := db.Where(\"id = ?\", id).First(&entity).Error; err != nil {\n\t\treturn nil, errors.New(\"%sNotFound\")\n\t}\n\treturn &response.%sResp{\n\t\tId: entity.Id,\n", name, name, name)
	for _, f := range spec.Fields {
		fmt.Fprintf(&b, "\t\t%s: entity.%s,\n", f.Pascal, f.Pascal)
	}
	fmt.Fprintf(&b, "\t\tCreatedAt: entity.CreatedAt,\n\t\tUpdatedAt: entity.UpdatedAt,\n\t}, nil\n}\n\n")

	// GetList
	fmt.Fprintf(&b, "// GetList 获取%s列表\n", name)
	fmt.Fprintf(&b, "func (s *%sService) GetList(db *gorm.DB, req common.Pagination) (*common.PageResult[response.%sResp], error) {\n", name, name)
	fmt.Fprintf(&b, "\tvar list []system.%s\n\tvar total int64\n\tdb = db.Model(&system.%s{})\n\tdb.Count(&total)\n\treq.Normalize()\n\tdb = db.Scopes(util.Paginate(req.PageSize, req.Page)).Order(\"id DESC\")\n\tif err := db.Find(&list).Error; err != nil {\n\t\treturn nil, errors.New(\"get%sListFailed\")\n\t}\n", name, name, name)
	fmt.Fprintf(&b, "\tdata := make([]response.%sResp, 0, len(list))\n\tfor _, item := range list {\n\t\tdata = append(data, response.%sResp{\n\t\t\tId: item.Id,\n", name, name)
	for _, f := range spec.Fields {
		fmt.Fprintf(&b, "\t\t\t%s: item.%s,\n", f.Pascal, f.Pascal)
	}
	fmt.Fprintf(&b, "\t\t\tCreatedAt: item.CreatedAt,\n\t\t\tUpdatedAt: item.UpdatedAt,\n\t\t})\n\t}\n")
	fmt.Fprintf(&b, "\treturn &common.PageResult[response.%sResp]{\n\t\tData: data,\n\t\tMeta: common.PageMeta{Page: req.Page, PageSize: req.PageSize, Total: total},\n\t}, nil\n}\n\n", name)

	// Delete
	fmt.Fprintf(&b, "// Delete 删除%s\n", name)
	fmt.Fprintf(&b, "func (s *%sService) Delete(db *gorm.DB, id uint) error {\n\tif err := db.Delete(&system.%s{}, id).Error; err != nil {\n\t\tglobal.Log.Error(\"删除%s失败\", zap.Error(err))\n\t\treturn errors.New(\"delete%sFailed\")\n\t}\n\treturn nil\n}\n", name, name, name, name)
	return b.String()
}

// genApi 生成 api/v1/system/{file}.go
func (spec *Spec) genApi() string {
	name := spec.Name
	svc := spec.serviceVar()
	var b strings.Builder
	b.WriteString("package system\n\n")
	b.WriteString("import (\n\t\"github.com/gin-gonic/gin\"\n\t\"github.com/imehc/do-exercise/server/model/common\"\n")
	b.WriteString("\t\"github.com/imehc/do-exercise/server/model/common/response\"\n\t\"github.com/imehc/do-exercise/server/model/system/request\"\n")
	b.WriteString("\t\"github.com/imehc/do-exercise/server/util\"\n\t\"github.com/spf13/cast\"\n)\n\n")
	fmt.Fprintf(&b, "type %sApi struct{}\n\n", name)

	fmt.Fprintf(&b, "// Create 创建%s\n", name)
	fmt.Fprintf(&b, "func (s %sApi) Create(ctx *gin.Context) {\n\tvar req request.Create%sReq\n\tif err := ctx.ShouldBindJSON(&req); err != nil {\n\t\tctx.Error(err)\n\t\treturn\n\t}\n\tif _, err := %s.Create(util.DB(ctx), req); err != nil {\n\t\tresponse.BadRequest(ctx, err.Error())\n\t\treturn\n\t}\n\tresponse.NoContent(ctx)\n}\n\n", name, name, svc)

	fmt.Fprintf(&b, "// Update 更新%s\n", name)
	fmt.Fprintf(&b, "func (s %sApi) Update(ctx *gin.Context) {\n\tid := cast.ToUint(ctx.Param(\"id\"))\n\tif id == 0 {\n\t\tresponse.BadRequest(ctx, \"idCannotBeEmpty\")\n\t\treturn\n\t}\n\tvar req request.Update%sReq\n\tif err := ctx.ShouldBindJSON(&req); err != nil {\n\t\tctx.Error(err)\n\t\treturn\n\t}\n\treq.Id = id\n\tif err := %s.Update(util.DB(ctx), req); err != nil {\n\t\tresponse.BadRequest(ctx, err.Error())\n\t\treturn\n\t}\n\tresponse.NoContent(ctx)\n}\n\n", name, name, svc)

	fmt.Fprintf(&b, "// Get 获取%s详情\n", name)
	fmt.Fprintf(&b, "func (s %sApi) Get(ctx *gin.Context) {\n\tid := cast.ToUint(ctx.Param(\"id\"))\n\tif id == 0 {\n\t\tresponse.BadRequest(ctx, \"idCannotBeEmpty\")\n\t\treturn\n\t}\n\tdata, err := %s.Get(util.DB(ctx), id)\n\tif err != nil {\n\t\tresponse.BadRequest(ctx, err.Error())\n\t\treturn\n\t}\n\tresponse.Success(ctx, data)\n}\n\n", name, svc)

	fmt.Fprintf(&b, "// GetList 获取%s列表\n", name)
	fmt.Fprintf(&b, "func (s %sApi) GetList(ctx *gin.Context) {\n\tvar req common.Pagination\n\tif err := ctx.ShouldBindQuery(&req); err != nil {\n\t\tctx.Error(err)\n\t\treturn\n\t}\n\tdata, err := %s.GetList(util.DB(ctx), req)\n\tif err != nil {\n\t\tresponse.BadRequest(ctx, err.Error())\n\t\treturn\n\t}\n\tresponse.Success(ctx, data)\n}\n\n", name, svc)

	fmt.Fprintf(&b, "// Delete 删除%s\n", name)
	fmt.Fprintf(&b, "func (s %sApi) Delete(ctx *gin.Context) {\n\tid := cast.ToUint(ctx.Param(\"id\"))\n\tif id == 0 {\n\t\tresponse.BadRequest(ctx, \"idCannotBeEmpty\")\n\t\treturn\n\t}\n\tif err := %s.Delete(util.DB(ctx), id); err != nil {\n\t\tresponse.BadRequest(ctx, err.Error())\n\t\treturn\n\t}\n\tresponse.NoContent(ctx)\n}\n", name, svc)
	return b.String()
}

// genRouter 生成 router/system/{file}.go
func (spec *Spec) genRouter() string {
	name := spec.Name
	apiVar := spec.apiVar()
	var b strings.Builder
	b.WriteString("package system\n\nimport \"github.com/gin-gonic/gin\"\n\n")
	fmt.Fprintf(&b, "type %sRouter struct{}\n\n", name)
	fmt.Fprintf(&b, "func (s %sRouter) Init%sRouter(r *gin.RouterGroup) gin.IRoutes {\n\trouter := r.Group(%q)\n\t{\n", name, name, spec.Path)
	fmt.Fprintf(&b, "\t\trouter.POST(\"\", %s.Create)     // 创建%s\n", apiVar, name)
	fmt.Fprintf(&b, "\t\trouter.PUT(\":id\", %s.Update)   // 更新%s\n", apiVar, name)
	fmt.Fprintf(&b, "\t\trouter.GET(\":id\", %s.Get)      // 获取单个%s\n", apiVar, name)
	fmt.Fprintf(&b, "\t\trouter.GET(\"\", %s.GetList)    // 获取%s列表\n", apiVar, name)
	fmt.Fprintf(&b, "\t\trouter.DELETE(\":id\", %s.Delete) // 删除%s\n", apiVar, name)
	b.WriteString("\t}\n\treturn router\n}\n")
	return b.String()
}

// writeFile 写入生成文件
func (spec *Spec) writeFiles(dir string) error {
	files := []struct {
		rel string
		gen string
	}{
		{filepath.Join("model", "system", spec.fileName()), spec.genModel()},
		{filepath.Join("model", "system", "request", spec.fileName()), spec.genRequest()},
		{filepath.Join("model", "system", "response", spec.fileName()), spec.genResponse()},
		{filepath.Join("service", "system", spec.fileName()), spec.genService()},
		{filepath.Join("api", "v1", "system", spec.fileName()), spec.genApi()},
		{filepath.Join("router", "system", spec.fileName()), spec.genRouter()},
	}
	for _, f := range files {
		abs := filepath.Join(dir, f.rel)
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(abs, []byte(f.gen), 0o644); err != nil {
			return fmt.Errorf("写入 %s 失败: %w", f.rel, err)
		}
		fmt.Printf("已生成 %s\n", f.rel)
	}
	return nil
}
