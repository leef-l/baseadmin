package dir_rule

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/upload/internal/dao"
	"gbaseadmin/app/upload/internal/logic/shared"
	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/app/upload/internal/service"
	"gbaseadmin/utility/batchutil"
	"gbaseadmin/utility/fieldvalid"
	"gbaseadmin/utility/inpututil"
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
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeDirRuleCreateInput(in)
	if err := validateDirRuleFields(in.DirID, in.Category, in.Status); err != nil {
		return err
	}
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
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeDirRuleUpdateInput(in)
	if err := validateDirRuleFields(in.DirID, in.Category, in.Status); err != nil {
		return err
	}
	if err := s.ensureDirRuleExists(ctx, in.ID); err != nil {
		return err
	}
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
	_, err := dao.UploadDirRule.Ctx(ctx).
		Where(dao.UploadDirRule.Columns().Id, in.ID).
		Where(dao.UploadDirRule.Columns().DeletedAt, nil).
		Data(data).
		Update()
	return err
}

// Delete 软删除文件目录规则
func (s *sDirRule) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := s.ensureDirRuleExists(ctx, id); err != nil {
		return err
	}
	_, err := dao.UploadDirRule.Ctx(ctx).
		Where(dao.UploadDirRule.Columns().Id, id).
		Where(dao.UploadDirRule.Columns().DeletedAt, nil).
		Data(g.Map{
			dao.UploadDirRule.Columns().DeletedAt: gtime.Now(),
		}).
		Update()
	return err
}

// BatchDelete 批量删除文件目录规则
func (s *sDirRule) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	ids = batchutil.CompactIDs(ids)
	if len(ids) == 0 {
		return gerror.New("请选择要删除的目录规则")
	}
	deleteIDs := batchutil.ToInt64s(ids)
	count, err := dao.UploadDirRule.Ctx(ctx).
		WhereIn(dao.UploadDirRule.Columns().Id, deleteIDs).
		Where(dao.UploadDirRule.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if count != len(deleteIDs) {
		return gerror.New("包含不存在或已删除的目录规则")
	}
	return dao.UploadDirRule.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Model(dao.UploadDirRule.Table()).Ctx(ctx).
			WhereIn(dao.UploadDirRule.Columns().Id, deleteIDs).
			Where(dao.UploadDirRule.Columns().DeletedAt, nil).
			Data(g.Map{
				dao.UploadDirRule.Columns().DeletedAt: gtime.Now(),
			}).
			Update()
		return err
	})
}

// Detail 获取文件目录规则详情
func (s *sDirRule) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.DirRuleDetailOutput, err error) {
	if id <= 0 {
		return nil, gerror.New("目录规则不存在或已删除")
	}
	out = &model.DirRuleDetailOutput{}
	err = dao.UploadDirRule.Ctx(ctx).Where(dao.UploadDirRule.Columns().Id, id).Where(dao.UploadDirRule.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("目录规则不存在或已删除")
	}
	out.DirName = shared.LookupDirName(ctx, int64(out.DirID))
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
	dirIDs := make([]int64, 0, len(list))
	for _, item := range list {
		if item.DirID != 0 {
			dirIDs = append(dirIDs, int64(item.DirID))
		}
	}
	dirMap := shared.LoadDirNameMap(ctx, dirIDs)
	for _, item := range list {
		item.DirName = dirMap[int64(item.DirID)]
	}
}

func (s *sDirRule) ensureDirRuleExists(ctx context.Context, id snowflake.JsonInt64) error {
	if id <= 0 {
		return gerror.New("目录规则不存在或已删除")
	}
	count, err := dao.UploadDirRule.Ctx(ctx).
		Where(dao.UploadDirRule.Columns().Id, id).
		Where(dao.UploadDirRule.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return gerror.New("目录规则不存在或已删除")
	}
	return nil
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

func validateDirRuleFields(dirID snowflake.JsonInt64, category, status int) error {
	if dirID <= 0 {
		return gerror.New("目录ID不能为空")
	}
	if err := fieldvalid.Enum("类别", category, 1, 2, 3); err != nil {
		return err
	}
	if err := fieldvalid.Binary("状态", status); err != nil {
		return err
	}
	return nil
}
