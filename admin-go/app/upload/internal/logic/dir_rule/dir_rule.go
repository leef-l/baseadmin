package dir_rule

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/upload/internal/dao"
	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/app/upload/internal/service"
	"gbaseadmin/utility/pageutil"
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
	normalizeDirRuleCreateInput(in)
	if err := s.ensureDirExists(ctx, in.DirID); err != nil {
		return err
	}
	id := snowflake.Generate()
	_, err := dao.UploadDirRule.Ctx(ctx).Data(g.Map{
		dao.UploadDirRule.Columns().Id:        id,
		dao.UploadDirRule.Columns().DirId:     in.DirID,
		dao.UploadDirRule.Columns().Category:  in.Category,
		dao.UploadDirRule.Columns().SavePath:  in.SavePath,
		dao.UploadDirRule.Columns().Status:    in.Status,
		dao.UploadDirRule.Columns().CreatedAt: gtime.Now(),
		dao.UploadDirRule.Columns().UpdatedAt: gtime.Now(),
	}).Insert()
	return err
}

// Update 更新文件目录规则
func (s *sDirRule) Update(ctx context.Context, in *model.DirRuleUpdateInput) error {
	normalizeDirRuleUpdateInput(in)
	if err := s.ensureDirExists(ctx, in.DirID); err != nil {
		return err
	}
	data := g.Map{
		dao.UploadDirRule.Columns().DirId:     in.DirID,
		dao.UploadDirRule.Columns().Category:  in.Category,
		dao.UploadDirRule.Columns().SavePath:  in.SavePath,
		dao.UploadDirRule.Columns().Status:    in.Status,
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
	if in == nil {
		in = &model.DirRuleListInput{}
	}
	normalizeDirRuleListInput(in)
	m := dao.UploadDirRule.Ctx(ctx).Where(dao.UploadDirRule.Columns().DeletedAt, nil)
	if in.Keyword != "" {
		m = m.WhereLike(dao.UploadDirRule.Columns().SavePath, "%"+in.Keyword+"%")
	}
	if in.Category > 0 {
		m = m.Where(dao.UploadDirRule.Columns().Category, in.Category)
	}
	if in.Status != nil {
		m = m.Where(dao.UploadDirRule.Columns().Status, *in.Status)
	}
	total, err = m.Count()
	if err != nil {
		return
	}
	in.PageNum, in.PageSize = pageutil.Normalize(in.PageNum, in.PageSize)
	err = m.Page(in.PageNum, in.PageSize).OrderAsc(dao.UploadDirRule.Columns().Id).Scan(&list)
	if err != nil {
		return
	}
	s.fillDirNames(ctx, list)
	return
}

func (s *sDirRule) fillDirNames(ctx context.Context, list []*model.DirRuleListOutput) {
	dirSet := make(map[int64]struct{})
	for _, item := range list {
		if item.DirID != 0 {
			dirSet[int64(item.DirID)] = struct{}{}
		}
	}
	if len(dirSet) == 0 {
		return
	}
	dirIDs := make([]int64, 0, len(dirSet))
	for id := range dirSet {
		dirIDs = append(dirIDs, id)
	}
	rows, err := g.DB().Ctx(ctx).Model("upload_dir").
		Fields("id", "name").
		Where("deleted_at", nil).
		WhereIn("id", dirIDs).
		All()
	if err != nil {
		return
	}
	dirMap := make(map[int64]string, len(rows))
	for _, row := range rows {
		dirMap[row["id"].Int64()] = row["name"].String()
	}
	for _, item := range list {
		item.DirName = dirMap[int64(item.DirID)]
	}
}

func (s *sDirRule) ensureDirExists(ctx context.Context, dirID snowflake.JsonInt64) error {
	if dirID == 0 {
		return nil
	}
	count, err := dao.UploadDir.Ctx(ctx).
		Where(dao.UploadDir.Columns().Id, dirID).
		Where(dao.UploadDir.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return gerror.New("所选目录不存在或已删除")
	}
	return nil
}

func normalizeDirRuleListInput(in *model.DirRuleListInput) {
	if in == nil {
		return
	}
	in.Keyword = strings.TrimSpace(in.Keyword)
}

func normalizeDirRuleCreateInput(in *model.DirRuleCreateInput) {
	if in == nil {
		return
	}
	in.SavePath = strings.TrimSpace(in.SavePath)
}

func normalizeDirRuleUpdateInput(in *model.DirRuleUpdateInput) {
	if in == nil {
		return
	}
	in.SavePath = strings.TrimSpace(in.SavePath)
}
