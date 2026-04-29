package dir

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"

	"gbaseadmin/app/upload/internal/dao"
	"gbaseadmin/app/upload/internal/logic/shared"
	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/app/upload/internal/service"
	"gbaseadmin/utility/batchutil"
	"gbaseadmin/utility/fieldvalid"
	"gbaseadmin/utility/inpututil"
	"gbaseadmin/utility/pageutil"
	"gbaseadmin/utility/snowflake"
	"gbaseadmin/utility/treeutil"
)

func init() {
	service.RegisterDir(New())
}

func New() *sDir {
	return &sDir{}
}

type sDir struct{}

type uploadDirSaveData struct {
	ParentId snowflake.JsonInt64 `orm:"parent_id"`
	Name     string              `orm:"name"`
	Path     string              `orm:"path"`
	KeepName int                 `orm:"keep_name"`
	Sort     int                 `orm:"sort"`
	Status   int                 `orm:"status"`
}

type uploadDirCreateData struct {
	Id         snowflake.JsonInt64 `orm:"id"`
	ParentId   snowflake.JsonInt64 `orm:"parent_id"`
	Name       string              `orm:"name"`
	Path       string              `orm:"path"`
	KeepName   int                 `orm:"keep_name"`
	Sort       int                 `orm:"sort"`
	Status     int                 `orm:"status"`
	TenantId   snowflake.JsonInt64 `orm:"tenant_id"`
	MerchantId snowflake.JsonInt64 `orm:"merchant_id"`
	CreatedBy  snowflake.JsonInt64 `orm:"created_by"`
	DeptId     snowflake.JsonInt64 `orm:"dept_id"`
}

// Create 创建文件目录
func (s *sDir) Create(ctx context.Context, in *model.DirCreateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeDirCreateInput(in)
	if err := validateDirFields(in.Name, in.Path, in.KeepName, in.Sort, in.Status); err != nil {
		return err
	}
	if err := s.ensureDirPathUnique(ctx, 0, in.Path); err != nil {
		return err
	}
	if err := s.ensureParentValid(ctx, in.ParentID, 0); err != nil {
		return err
	}
	id := snowflake.Generate()
	var tenantID, merchantID, createdBy, deptID snowflake.JsonInt64
	shared.ApplyWriteScope(ctx, &tenantID, &merchantID, &createdBy, &deptID)
	_, err := dao.UploadDir.Ctx(ctx).Data(uploadDirCreateData{
		Id:         id,
		ParentId:   in.ParentID,
		Name:       in.Name,
		Path:       in.Path,
		KeepName:   in.KeepName,
		Sort:       in.Sort,
		Status:     in.Status,
		TenantId:   tenantID,
		MerchantId: merchantID,
		CreatedBy:  createdBy,
		DeptId:     deptID,
	}).Insert()
	return err
}

// Update 更新文件目录
func (s *sDir) Update(ctx context.Context, in *model.DirUpdateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeDirUpdateInput(in)
	if err := validateDirFields(in.Name, in.Path, in.KeepName, in.Sort, in.Status); err != nil {
		return err
	}
	if err := s.ensureDirExists(ctx, in.ID); err != nil {
		return err
	}
	if err := s.ensureDirPathUnique(ctx, in.ID, in.Path); err != nil {
		return err
	}
	if err := s.ensureParentValid(ctx, in.ParentID, in.ID); err != nil {
		return err
	}
	data := uploadDirSaveData{
		ParentId: in.ParentID,
		Name:     in.Name,
		Path:     in.Path,
		KeepName: in.KeepName,
		Sort:     in.Sort,
		Status:   in.Status,
	}
	m := dao.UploadDir.Ctx(ctx).
		Where(dao.UploadDir.Columns().Id, in.ID).
		Where(dao.UploadDir.Columns().DeletedAt, nil).
		Data(data)
	m = shared.ApplyAccessScope(ctx, m)
	_, err := m.Update()
	return err
}

// Delete 软删除文件目录
func (s *sDir) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := s.ensureDirExists(ctx, id); err != nil {
		return err
	}
	if err := s.ensureDirDeletable(ctx, id); err != nil {
		return err
	}
	_, err := dao.UploadDir.Ctx(ctx).
		Where(dao.UploadDir.Columns().Id, id).
		Delete()
	return err
}

// BatchDelete 批量删除文件目录
func (s *sDir) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	ids = batchutil.CompactIDs(ids)
	if len(ids) == 0 {
		return gerror.New("请选择要删除的目录")
	}
	if err := s.ensureDirIDsExist(ctx, ids); err != nil {
		return err
	}
	order, err := s.collectBatchDeleteOrder(ctx, ids)
	if err != nil {
		return err
	}
	if len(order) == 0 {
		return nil
	}
	deleteIDs := batchutil.ToInt64s(order)
	for _, id := range order {
		if err := s.ensureDirBatchDeletable(ctx, id, deleteIDs); err != nil {
			return err
		}
	}
	return dao.UploadDir.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Model(dao.UploadDir.Table()).Ctx(ctx).
			WhereIn(dao.UploadDir.Columns().Id, deleteIDs).
			Delete()
		return err
	})
}

// Detail 获取文件目录详情
func (s *sDir) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.DirDetailOutput, err error) {
	if id <= 0 {
		return nil, gerror.New("目录不存在或已删除")
	}
	out = &model.DirDetailOutput{}
	m := dao.UploadDir.Ctx(ctx).Where(dao.UploadDir.Columns().Id, id).Where(dao.UploadDir.Columns().DeletedAt, nil)
	m = shared.ApplyAccessScope(ctx, m)
	err = m.Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("目录不存在或已删除")
	}
	out.DirName = shared.LookupDirName(ctx, int64(out.ParentID))
	return
}

// List 获取文件目录列表
func (s *sDir) List(ctx context.Context, in *model.DirListInput) (list []*model.DirListOutput, total int, err error) {
	if in == nil {
		in = &model.DirListInput{}
	}
	normalizeDirListInput(in)
	m := dao.UploadDir.Ctx(ctx).Where(dao.UploadDir.Columns().DeletedAt, nil)
	m = shared.ApplyAccessScope(ctx, m)
	if in.Keyword != "" {
		keywordBuilder := m.Builder().
			WhereLike(dao.UploadDir.Columns().Name, "%"+in.Keyword+"%").
			WhereOrLike(dao.UploadDir.Columns().Path, "%"+in.Keyword+"%")
		m = m.Where(keywordBuilder)
	}
	if in.Status != nil {
		m = m.Where(dao.UploadDir.Columns().Status, *in.Status)
	}
	total, err = m.Count()
	if err != nil {
		return
	}
	in.PageNum, in.PageSize = pageutil.Normalize(in.PageNum, in.PageSize)
	err = m.Page(in.PageNum, in.PageSize).OrderAsc(dao.UploadDir.Columns().Id).Scan(&list)
	if err != nil {
		return
	}
	s.fillParentNames(ctx, list)
	return
}

// Tree 获取文件目录树形结构
func (s *sDir) Tree(ctx context.Context, in *model.DirTreeInput) (tree []*model.DirTreeOutput, err error) {
	var list []*model.DirTreeOutput
	m := dao.UploadDir.Ctx(ctx).Where(dao.UploadDir.Columns().DeletedAt, nil)
	m = shared.ApplyAccessScope(ctx, m)
	if in != nil {
		normalizeDirTreeInput(in)
		if in.Keyword != "" {
			keywordBuilder := m.Builder().
				WhereLike(dao.UploadDir.Columns().Name, "%"+in.Keyword+"%").
				WhereOrLike(dao.UploadDir.Columns().Path, "%"+in.Keyword+"%")
			m = m.Where(keywordBuilder)
		}
		if in.Status != nil {
			m = m.Where(dao.UploadDir.Columns().Status, *in.Status)
		}
	}
	err = m.OrderAsc(dao.UploadDir.Columns().Sort).Scan(&list)
	if err != nil {
		return
	}

	tree = treeutil.BuildForest(list, treeutil.TreeNodeAccessor[*model.DirTreeOutput]{
		ID:       func(item *model.DirTreeOutput) int64 { return int64(item.ID) },
		ParentID: func(item *model.DirTreeOutput) int64 { return int64(item.ParentID) },
		Init: func(item *model.DirTreeOutput) {
			item.Children = make([]*model.DirTreeOutput, 0)
		},
		Append: func(parent *model.DirTreeOutput, child *model.DirTreeOutput) {
			parent.Children = append(parent.Children, child)
		},
	})
	return
}

func normalizeDirListInput(in *model.DirListInput) {
	if in == nil {
		return
	}
	in.Keyword = strings.TrimSpace(in.Keyword)
}

func normalizeDirTreeInput(in *model.DirTreeInput) {
	if in == nil {
		return
	}
	in.Keyword = strings.TrimSpace(in.Keyword)
}

func normalizeDirCreateInput(in *model.DirCreateInput) {
	if in == nil {
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Path = strings.TrimSpace(in.Path)
}

func normalizeDirUpdateInput(in *model.DirUpdateInput) {
	if in == nil {
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Path = strings.TrimSpace(in.Path)
}

func validateDirFields(name, path string, keepName, sort, status int) error {
	if name == "" {
		return gerror.New("目录名称不能为空")
	}
	if path == "" {
		return gerror.New("目录路径不能为空")
	}
	if err := fieldvalid.NonNegative("排序", sort); err != nil {
		return err
	}
	if err := fieldvalid.Binary("状态", status); err != nil {
		return err
	}
	if err := fieldvalid.Binary("保留原文件名", keepName); err != nil {
		return err
	}
	return nil
}

func (s *sDir) fillParentNames(ctx context.Context, list []*model.DirListOutput) {
	parentIDs := make([]int64, 0, len(list))
	for _, item := range list {
		if item.ParentID != 0 {
			parentIDs = append(parentIDs, int64(item.ParentID))
		}
	}
	parentMap := shared.LoadDirNameMap(ctx, parentIDs)
	for _, item := range list {
		item.DirName = parentMap[int64(item.ParentID)]
	}
}

func (s *sDir) ensureDirDeletable(ctx context.Context, id snowflake.JsonInt64) error {
	childModel := dao.UploadDir.Ctx(ctx).
		Where(dao.UploadDir.Columns().ParentId, id).
		Where(dao.UploadDir.Columns().DeletedAt, nil)
	childModel = shared.ApplyAccessScope(ctx, childModel)
	childCount, err := childModel.Count()
	if err != nil {
		return err
	}
	if childCount > 0 {
		return gerror.New("当前目录下存在子目录，不能直接删除")
	}
	fileModel := dao.UploadFile.Ctx(ctx).
		Where(dao.UploadFile.Columns().DirId, id).
		Where(dao.UploadFile.Columns().DeletedAt, nil)
	fileModel = shared.ApplyAccessScope(ctx, fileModel)
	fileCount, err := fileModel.Count()
	if err != nil {
		return err
	}
	if fileCount > 0 {
		return gerror.New("当前目录下仍有关联文件，不能直接删除")
	}
	ruleModel := dao.UploadDirRule.Ctx(ctx).
		Where(dao.UploadDirRule.Columns().DirId, id).
		Where(dao.UploadDirRule.Columns().DeletedAt, nil)
	ruleModel = shared.ApplyAccessScope(ctx, ruleModel)
	ruleCount, err := ruleModel.Count()
	if err != nil {
		return err
	}
	if ruleCount > 0 {
		return gerror.New("当前目录仍有关联规则，不能直接删除")
	}
	return nil
}

func (s *sDir) ensureDirExists(ctx context.Context, id snowflake.JsonInt64) error {
	if id <= 0 {
		return gerror.New("目录不存在或已删除")
	}
	m := dao.UploadDir.Ctx(ctx).
		Where(dao.UploadDir.Columns().Id, id).
		Where(dao.UploadDir.Columns().DeletedAt, nil)
	m = shared.ApplyAccessScope(ctx, m)
	count, err := m.Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return gerror.New("目录不存在或已删除")
	}
	return nil
}

func (s *sDir) ensureDirIDsExist(ctx context.Context, ids []snowflake.JsonInt64) error {
	dbIDs := batchutil.ToInt64s(ids)
	m := dao.UploadDir.Ctx(ctx).
		WhereIn(dao.UploadDir.Columns().Id, dbIDs).
		Where(dao.UploadDir.Columns().DeletedAt, nil)
	m = shared.ApplyAccessScope(ctx, m)
	count, err := m.Count()
	if err != nil {
		return err
	}
	if count != len(dbIDs) {
		return gerror.New("包含不存在或已删除的目录")
	}
	return nil
}

func (s *sDir) ensureDirPathUnique(ctx context.Context, currentID snowflake.JsonInt64, path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}
	m := dao.UploadDir.Ctx(ctx).
		Where(dao.UploadDir.Columns().Path, path).
		Where(dao.UploadDir.Columns().DeletedAt, nil)
	m = shared.ApplyTenantScopeToModel(ctx, m, shared.ColumnTenantID, shared.ColumnMerchantID)
	if currentID > 0 {
		m = m.WhereNot(dao.UploadDir.Columns().Id, currentID)
	}
	count, err := m.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New("目录路径已存在")
	}
	return nil
}

func (s *sDir) ensureDirBatchDeletable(ctx context.Context, id snowflake.JsonInt64, deleteIDs []int64) error {
	childModel := dao.UploadDir.Ctx(ctx).
		Where(dao.UploadDir.Columns().ParentId, id).
		Where(dao.UploadDir.Columns().DeletedAt, nil)
	childModel = shared.ApplyAccessScope(ctx, childModel)
	if len(deleteIDs) > 0 {
		childModel = childModel.WhereNotIn(dao.UploadDir.Columns().Id, deleteIDs)
	}
	childCount, err := childModel.Count()
	if err != nil {
		return err
	}
	if childCount > 0 {
		return gerror.New("当前目录下存在子目录，不能直接删除")
	}
	fileModel := dao.UploadFile.Ctx(ctx).
		Where(dao.UploadFile.Columns().DirId, id).
		Where(dao.UploadFile.Columns().DeletedAt, nil)
	fileModel = shared.ApplyAccessScope(ctx, fileModel)
	fileCount, err := fileModel.Count()
	if err != nil {
		return err
	}
	if fileCount > 0 {
		return gerror.New("当前目录下仍有关联文件，不能直接删除")
	}
	ruleModel := dao.UploadDirRule.Ctx(ctx).
		Where(dao.UploadDirRule.Columns().DirId, id).
		Where(dao.UploadDirRule.Columns().DeletedAt, nil)
	ruleModel = shared.ApplyAccessScope(ctx, ruleModel)
	ruleCount, err := ruleModel.Count()
	if err != nil {
		return err
	}
	if ruleCount > 0 {
		return gerror.New("当前目录仍有关联规则，不能直接删除")
	}
	return nil
}

func (s *sDir) collectBatchDeleteOrder(ctx context.Context, ids []snowflake.JsonInt64) ([]snowflake.JsonInt64, error) {
	var rows []struct {
		Id       int64 `json:"id"`
		ParentId int64 `json:"parentId"`
	}
	m := dao.UploadDir.Ctx(ctx).
		Fields(dao.UploadDir.Columns().Id, dao.UploadDir.Columns().ParentId).
		Where(dao.UploadDir.Columns().DeletedAt, nil)
	m = shared.ApplyAccessScope(ctx, m)
	if err := m.Scan(&rows); err != nil {
		return nil, err
	}
	treeRows := make([]batchutil.TreeRow, 0, len(rows))
	for _, row := range rows {
		treeRows = append(treeRows, batchutil.TreeRow{
			ID:       row.Id,
			ParentID: row.ParentId,
		})
	}
	return batchutil.ExpandTreeDeleteOrder(ids, treeRows), nil
}

func (s *sDir) ensureParentValid(ctx context.Context, parentID, currentID snowflake.JsonInt64) error {
	return treeutil.ValidateParent(parentID, currentID, func(id int64) (int64, int64, error) {
		var parent struct {
			Id       int64 `json:"id"`
			ParentId int64 `json:"parentId"`
		}
		m := dao.UploadDir.Ctx(ctx).
			Fields(dao.UploadDir.Columns().Id, dao.UploadDir.Columns().ParentId).
			Where(dao.UploadDir.Columns().Id, id).
			Where(dao.UploadDir.Columns().DeletedAt, nil)
		m = shared.ApplyAccessScope(ctx, m)
		if err := m.Scan(&parent); err != nil {
			return 0, 0, err
		}
		return parent.Id, parent.ParentId, nil
	}, treeutil.Messages{
		Self:         "上级目录不能选择自己",
		Missing:      "上级目录不存在或已删除",
		ChildLoop:    "不能将目录移动到自己的子级下",
		Cycle:        "目录层级存在循环引用",
		InvalidChain: "上级目录链路中存在无效节点",
	})
}
