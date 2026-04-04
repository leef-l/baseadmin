package dir_rule

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/upload/internal/dao"
	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/app/upload/internal/service"
	"gbaseadmin/utility/snowflake"
)

func init() {
	service.RegisterDirRule(New())
}

func New() *sDirRule {
	return &sDirRule{}
}

type sDirRule struct{}

// Create 创建文件目录规则
func (s *sDirRule) Create(ctx context.Context, in *model.DirRuleCreateInput) error {
	id := snowflake.Generate()
	_, err := dao.UploadDirRule.Ctx(ctx).Data(g.Map{
		dao.UploadDirRule.Columns().Id:        id,
		dao.UploadDirRule.Columns().DirId: in.DirID,
		dao.UploadDirRule.Columns().Category: in.Category,
		dao.UploadDirRule.Columns().SavePath: in.SavePath,
		dao.UploadDirRule.Columns().Status: in.Status,
		dao.UploadDirRule.Columns().CreatedAt: gtime.Now(),
		dao.UploadDirRule.Columns().UpdatedAt: gtime.Now(),
	}).Insert()
	return err
}

// Update 更新文件目录规则
func (s *sDirRule) Update(ctx context.Context, in *model.DirRuleUpdateInput) error {
	data := g.Map{
		dao.UploadDirRule.Columns().DirId: in.DirID,
		dao.UploadDirRule.Columns().Category: in.Category,
		dao.UploadDirRule.Columns().SavePath: in.SavePath,
		dao.UploadDirRule.Columns().Status: in.Status,
		dao.UploadDirRule.Columns().UpdatedAt: gtime.Now(),
	}
	_, err := dao.UploadDirRule.Ctx(ctx).Where(dao.UploadDirRule.Columns().Id, in.ID).Data(data).Update()
	return err
}

// Delete 软删除文件目录规则
func (s *sDirRule) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	_, err := dao.UploadDirRule.Ctx(ctx).Where(dao.UploadDirRule.Columns().Id, id).Data(g.Map{
		dao.UploadDirRule.Columns().DeletedAt: gtime.Now(),
	}).Update()
	return err
}

// Detail 获取文件目录规则详情
func (s *sDirRule) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.DirRuleDetailOutput, err error) {
	out = &model.DirRuleDetailOutput{}
	err = dao.UploadDirRule.Ctx(ctx).Where(dao.UploadDirRule.Columns().Id, id).Where(dao.UploadDirRule.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	// 查询目录ID关联显示
	if out.DirID != 0 {
		val, err := g.DB().Ctx(ctx).Model("upload_dir").Where("id", out.DirID).Where("deleted_at", nil).Value("name")
		if err == nil {
			out.DirName = val.String()
		}
	}
	return
}

// List 获取文件目录规则列表
func (s *sDirRule) List(ctx context.Context, in *model.DirRuleListInput) (list []*model.DirRuleListOutput, total int, err error) {
	m := dao.UploadDirRule.Ctx(ctx).Where(dao.UploadDirRule.Columns().DeletedAt, nil)
	if in.Category > 0 {
		m = m.Where(dao.UploadDirRule.Columns().Category, in.Category)
	}
	if in.Status > 0 {
		m = m.Where(dao.UploadDirRule.Columns().Status, in.Status)
	}
	total, err = m.Count()
	if err != nil {
		return
	}
	err = m.Page(in.PageNum, in.PageSize).OrderAsc(dao.UploadDirRule.Columns().Id).Scan(&list)
	if err != nil {
		return
	}
	// 填充关联显示字段
	for _, item := range list {
		if item.DirID != 0 {
			val, err := g.DB().Ctx(ctx).Model("upload_dir").Where("id", item.DirID).Where("deleted_at", nil).Value("name")
			if err == nil {
				item.DirName = val.String()
			}
		}
	}
	return
}

