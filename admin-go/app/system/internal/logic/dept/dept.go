package dept

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"

	"gbaseadmin/app/system/internal/dao"
	authlogic "gbaseadmin/app/system/internal/logic/auth"
	"gbaseadmin/app/system/internal/logic/shared"
	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/app/system/internal/model/do"
	"gbaseadmin/app/system/internal/service"
	"gbaseadmin/utility/batchutil"
	"gbaseadmin/utility/fieldvalid"
	"gbaseadmin/utility/inpututil"
	"gbaseadmin/utility/pageutil"
	"gbaseadmin/utility/snowflake"
	"gbaseadmin/utility/treeutil"
)

func init() {
	service.RegisterDept(New())
}

func New() *sDept {
	return &sDept{}
}

type sDept struct{}

// Create 创建部门表
func (s *sDept) Create(ctx context.Context, in *model.DeptCreateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeDeptCreateInput(in)
	if err := validateDeptFields(in.Title, in.Sort, in.Status); err != nil {
		return err
	}
	if err := s.ensureDeptTitleUnique(ctx, 0, in.ParentID, in.Title); err != nil {
		return err
	}
	if err := s.ensureParentValid(ctx, in.ParentID, 0); err != nil {
		return err
	}
	if err := s.ensureParentAccessible(ctx, in.ParentID); err != nil {
		return err
	}
	id := snowflake.Generate()
	_, err := dao.Dept.Ctx(ctx).Data(do.Dept{
		Id:       id,
		ParentId: in.ParentID,
		Title:    in.Title,
		Username: in.Username,
		Email:    in.Email,
		Sort:     in.Sort,
		Status:   in.Status,
	}).Insert()
	return err
}

// Update 更新部门表
func (s *sDept) Update(ctx context.Context, in *model.DeptUpdateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeDeptUpdateInput(in)
	if err := validateDeptFields(in.Title, in.Sort, in.Status); err != nil {
		return err
	}
	if err := s.ensureDeptExists(ctx, in.ID); err != nil {
		return err
	}
	if err := s.ensureDeptAccessible(ctx, in.ID); err != nil {
		return err
	}
	if err := s.ensureDeptTitleUnique(ctx, in.ID, in.ParentID, in.Title); err != nil {
		return err
	}
	if err := s.ensureParentValid(ctx, in.ParentID, in.ID); err != nil {
		return err
	}
	if err := s.ensureParentAccessible(ctx, in.ParentID); err != nil {
		return err
	}
	data := do.Dept{
		ParentId: in.ParentID,
		Title:    in.Title,
		Username: in.Username,
		Email:    in.Email,
		Sort:     in.Sort,
		Status:   in.Status,
	}
	_, err := dao.Dept.Ctx(ctx).
		Where(dao.Dept.Columns().Id, in.ID).
		Where(dao.Dept.Columns().DeletedAt, nil).
		Data(data).
		Update()
	if err == nil {
		authlogic.ClearUserCaches(ctx, s.getDeptUserIDs(ctx, []int64{int64(in.ID)})...)
	}
	return err
}

// Delete 软删除部门表
func (s *sDept) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := s.ensureDeptExists(ctx, id); err != nil {
		return err
	}
	if err := s.ensureDeptAccessible(ctx, id); err != nil {
		return err
	}
	if err := s.ensureDeptDeletable(ctx, id); err != nil {
		return err
	}
	_, err := dao.Dept.Ctx(ctx).
		Where(dao.Dept.Columns().Id, id).
		Delete()
	return err
}

// BatchDelete 批量删除部门表
func (s *sDept) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	ids = batchutil.CompactIDs(ids)
	if len(ids) == 0 {
		return gerror.New("请选择要删除的部门")
	}
	if err := s.ensureDeptIDsExist(ctx, ids); err != nil {
		return err
	}
	order, err := s.collectBatchDeleteOrder(ctx, ids)
	if err != nil {
		return err
	}
	if len(order) == 0 {
		return nil
	}
	if err := s.ensureDeptIDsAccessible(ctx, order); err != nil {
		return err
	}
	deleteIDs := batchutil.ToInt64s(order)
	for _, id := range order {
		if err := s.ensureDeptBatchDeletable(ctx, id, deleteIDs); err != nil {
			return err
		}
	}
	err = dao.Dept.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Model(dao.Dept.Table()).Ctx(ctx).
			WhereIn(dao.Dept.Columns().Id, deleteIDs).
			Delete()
		return err
	})
	return err
}

// Detail 获取部门表详情
func (s *sDept) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.DeptDetailOutput, err error) {
	if id <= 0 {
		return nil, gerror.New("部门不存在或已删除")
	}
	out = &model.DeptDetailOutput{}
	err = dao.Dept.Ctx(ctx).Where(dao.Dept.Columns().Id, id).Where(dao.Dept.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("部门不存在或已删除")
	}
	allowed, err := shared.CanAccessDept(ctx, int64(out.ID))
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, gerror.New("部门不存在或已删除")
	}
	out.DeptTitle = shared.LookupTitle(ctx, "system_dept", int64(out.ParentID))
	return
}

// List 获取部门表列表
func (s *sDept) List(ctx context.Context, in *model.DeptListInput) (list []*model.DeptListOutput, total int, err error) {
	if in == nil {
		in = &model.DeptListInput{}
	}
	normalizeDeptListInput(in)
	m := dao.Dept.Ctx(ctx).Where(dao.Dept.Columns().DeletedAt, nil)
	m, err = shared.ApplyDeptScopeToDeptModel(ctx, m, dao.Dept.Columns().Id)
	if err != nil {
		return nil, 0, err
	}
	if in.Keyword != "" {
		keywordBuilder := m.Builder().
			WhereLike(dao.Dept.Columns().Title, "%"+in.Keyword+"%").
			WhereOrLike(dao.Dept.Columns().Username, "%"+in.Keyword+"%").
			WhereOrLike(dao.Dept.Columns().Email, "%"+in.Keyword+"%")
		m = m.Where(keywordBuilder)
	}
	if in.Status != nil {
		m = m.Where(dao.Dept.Columns().Status, *in.Status)
	}
	total, err = m.Count()
	if err != nil {
		return
	}
	in.PageNum, in.PageSize = pageutil.Normalize(in.PageNum, in.PageSize)
	err = m.Page(in.PageNum, in.PageSize).OrderAsc(dao.Dept.Columns().Id).Scan(&list)
	if err != nil {
		return
	}
	s.fillParentTitles(ctx, list)
	return
}

// Tree 获取部门表树形结构
func (s *sDept) Tree(ctx context.Context, in *model.DeptTreeInput) (tree []*model.DeptTreeOutput, err error) {
	var list []*model.DeptTreeOutput
	m := dao.Dept.Ctx(ctx).Where(dao.Dept.Columns().DeletedAt, nil)
	m, err = shared.ApplyDeptScopeToDeptModel(ctx, m, dao.Dept.Columns().Id)
	if err != nil {
		return nil, err
	}
	if in != nil {
		normalizeDeptTreeInput(in)
		if in.Keyword != "" {
			keywordBuilder := m.Builder().
				WhereLike(dao.Dept.Columns().Title, "%"+in.Keyword+"%").
				WhereOrLike(dao.Dept.Columns().Username, "%"+in.Keyword+"%").
				WhereOrLike(dao.Dept.Columns().Email, "%"+in.Keyword+"%")
			m = m.Where(keywordBuilder)
		}
		if in.Status != nil {
			m = m.Where(dao.Dept.Columns().Status, *in.Status)
		}
	}
	err = m.OrderAsc(dao.Dept.Columns().Sort).Scan(&list)
	if err != nil {
		return
	}

	tree = treeutil.BuildForest(list, treeutil.TreeNodeAccessor[*model.DeptTreeOutput]{
		ID:       func(item *model.DeptTreeOutput) int64 { return int64(item.ID) },
		ParentID: func(item *model.DeptTreeOutput) int64 { return int64(item.ParentID) },
		Init: func(item *model.DeptTreeOutput) {
			item.Children = make([]*model.DeptTreeOutput, 0)
		},
		Append: func(parent *model.DeptTreeOutput, child *model.DeptTreeOutput) {
			parent.Children = append(parent.Children, child)
		},
	})
	return
}

func normalizeDeptCreateInput(in *model.DeptCreateInput) {
	if in == nil {
		return
	}
	in.Title = strings.TrimSpace(in.Title)
	in.Username = strings.TrimSpace(in.Username)
	in.Email = strings.TrimSpace(in.Email)
}

func normalizeDeptUpdateInput(in *model.DeptUpdateInput) {
	if in == nil {
		return
	}
	in.Title = strings.TrimSpace(in.Title)
	in.Username = strings.TrimSpace(in.Username)
	in.Email = strings.TrimSpace(in.Email)
}

func normalizeDeptListInput(in *model.DeptListInput) {
	if in == nil {
		return
	}
	in.Keyword = strings.TrimSpace(in.Keyword)
}

func normalizeDeptTreeInput(in *model.DeptTreeInput) {
	if in == nil {
		return
	}
	in.Keyword = strings.TrimSpace(in.Keyword)
}

func validateDeptFields(title string, sort, status int) error {
	if title == "" {
		return gerror.New("部门名称不能为空")
	}
	if err := fieldvalid.NonNegative("排序", sort); err != nil {
		return err
	}
	if err := fieldvalid.Binary("状态", status); err != nil {
		return err
	}
	return nil
}

func (s *sDept) fillParentTitles(ctx context.Context, list []*model.DeptListOutput) {
	parentIDs := make([]int64, 0, len(list))
	for _, item := range list {
		if item.ParentID != 0 {
			parentIDs = append(parentIDs, int64(item.ParentID))
		}
	}
	parentMap := shared.LoadTitleMap(ctx, "system_dept", parentIDs)
	for _, item := range list {
		item.DeptTitle = parentMap[int64(item.ParentID)]
	}
}

func (s *sDept) ensureDeptDeletable(ctx context.Context, id snowflake.JsonInt64) error {
	childCount, err := dao.Dept.Ctx(ctx).
		Where(dao.Dept.Columns().ParentId, id).
		Where(dao.Dept.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if childCount > 0 {
		return gerror.New("当前部门下存在子部门，不能直接删除")
	}
	userCount, err := dao.Users.Ctx(ctx).
		Where(dao.Users.Columns().DeptId, id).
		Where(dao.Users.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if userCount > 0 {
		return gerror.New("当前部门下仍有关联用户，不能直接删除")
	}
	roleDeptCount, err := dao.RoleDept.Ctx(ctx).
		Where(dao.RoleDept.Columns().DeptId, id).
		Count()
	if err != nil {
		return err
	}
	if roleDeptCount > 0 {
		return gerror.New("当前部门仍被角色数据权限引用，不能直接删除")
	}
	return nil
}

func (s *sDept) ensureDeptExists(ctx context.Context, id snowflake.JsonInt64) error {
	if id <= 0 {
		return gerror.New("部门不存在或已删除")
	}
	count, err := dao.Dept.Ctx(ctx).
		Where(dao.Dept.Columns().Id, id).
		Where(dao.Dept.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return gerror.New("部门不存在或已删除")
	}
	return nil
}

func (s *sDept) ensureDeptAccessible(ctx context.Context, id snowflake.JsonInt64) error {
	allowed, err := shared.CanAccessDept(ctx, int64(id))
	if err != nil {
		return err
	}
	if !allowed {
		return gerror.New("部门不存在或已删除")
	}
	return nil
}

func (s *sDept) ensureDeptIDsAccessible(ctx context.Context, ids []snowflake.JsonInt64) error {
	for _, id := range ids {
		if err := s.ensureDeptAccessible(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (s *sDept) ensureDeptIDsExist(ctx context.Context, ids []snowflake.JsonInt64) error {
	dbIDs := batchutil.ToInt64s(ids)
	count, err := dao.Dept.Ctx(ctx).
		WhereIn(dao.Dept.Columns().Id, dbIDs).
		Where(dao.Dept.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if count != len(dbIDs) {
		return gerror.New("包含不存在或已删除的部门")
	}
	return nil
}

func (s *sDept) ensureDeptTitleUnique(ctx context.Context, currentID, parentID snowflake.JsonInt64, title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil
	}
	m := dao.Dept.Ctx(ctx).
		Where(dao.Dept.Columns().ParentId, parentID).
		Where(dao.Dept.Columns().Title, title).
		Where(dao.Dept.Columns().DeletedAt, nil)
	if currentID > 0 {
		m = m.WhereNot(dao.Dept.Columns().Id, currentID)
	}
	count, err := m.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New("同级部门名称已存在")
	}
	return nil
}

func (s *sDept) ensureDeptBatchDeletable(ctx context.Context, id snowflake.JsonInt64, deleteIDs []int64) error {
	childModel := dao.Dept.Ctx(ctx).
		Where(dao.Dept.Columns().ParentId, id).
		Where(dao.Dept.Columns().DeletedAt, nil)
	if len(deleteIDs) > 0 {
		childModel = childModel.WhereNotIn(dao.Dept.Columns().Id, deleteIDs)
	}
	childCount, err := childModel.Count()
	if err != nil {
		return err
	}
	if childCount > 0 {
		return gerror.New("当前部门下存在子部门，不能直接删除")
	}
	userCount, err := dao.Users.Ctx(ctx).
		Where(dao.Users.Columns().DeptId, id).
		Where(dao.Users.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if userCount > 0 {
		return gerror.New("当前部门下仍有关联用户，不能直接删除")
	}
	roleDeptCount, err := dao.RoleDept.Ctx(ctx).
		Where(dao.RoleDept.Columns().DeptId, id).
		Count()
	if err != nil {
		return err
	}
	if roleDeptCount > 0 {
		return gerror.New("当前部门仍被角色数据权限引用，不能直接删除")
	}
	return nil
}

func (s *sDept) collectBatchDeleteOrder(ctx context.Context, ids []snowflake.JsonInt64) ([]snowflake.JsonInt64, error) {
	var rows []struct {
		Id       int64 `json:"id"`
		ParentId int64 `json:"parentId"`
	}
	if err := dao.Dept.Ctx(ctx).
		Fields(dao.Dept.Columns().Id, dao.Dept.Columns().ParentId).
		Where(dao.Dept.Columns().DeletedAt, nil).
		Scan(&rows); err != nil {
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

func (s *sDept) ensureParentValid(ctx context.Context, parentID, currentID snowflake.JsonInt64) error {
	return treeutil.ValidateParent(parentID, currentID, func(id int64) (int64, int64, error) {
		var parent struct {
			Id       int64 `json:"id"`
			ParentId int64 `json:"parentId"`
		}
		if err := dao.Dept.Ctx(ctx).
			Fields(dao.Dept.Columns().Id, dao.Dept.Columns().ParentId).
			Where(dao.Dept.Columns().Id, id).
			Where(dao.Dept.Columns().DeletedAt, nil).
			Scan(&parent); err != nil {
			return 0, 0, err
		}
		return parent.Id, parent.ParentId, nil
	}, treeutil.Messages{
		Self:         "上级部门不能选择自己",
		Missing:      "上级部门不存在或已删除",
		ChildLoop:    "不能将部门移动到自己的子级下",
		Cycle:        "部门层级存在循环引用",
		InvalidChain: "上级部门链路中存在无效节点",
	})
}

func (s *sDept) ensureParentAccessible(ctx context.Context, parentID snowflake.JsonInt64) error {
	if parentID == 0 {
		scope, err := shared.ResolveDataAccessScope(ctx)
		if err != nil {
			return err
		}
		if scope.All {
			return nil
		}
		return gerror.New("上级部门不存在或无权操作")
	}
	allowed, err := shared.CanAccessDept(ctx, int64(parentID))
	if err != nil {
		return err
	}
	if !allowed {
		return gerror.New("上级部门不存在或无权操作")
	}
	return nil
}

func (s *sDept) getDeptUserIDs(ctx context.Context, deptIDs []int64) []int64 {
	deptIDs = compactDeptIDs(deptIDs)
	if len(deptIDs) == 0 {
		return nil
	}
	var rows []struct {
		Id int64 `json:"id"`
	}
	if err := dao.Users.Ctx(ctx).
		Fields(dao.Users.Columns().Id).
		WhereIn(dao.Users.Columns().DeptId, deptIDs).
		Where(dao.Users.Columns().DeletedAt, nil).
		Scan(&rows); err != nil {
		return nil
	}
	userIDs := make([]int64, 0, len(rows))
	for _, row := range rows {
		userIDs = append(userIDs, row.Id)
	}
	return userIDs
}

func compactDeptIDs(values []int64) []int64 {
	if len(values) == 0 {
		return nil
	}
	seen := make(map[int64]struct{}, len(values))
	ids := make([]int64, 0, len(values))
	for _, value := range values {
		if value <= 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		ids = append(ids, value)
	}
	if len(ids) == 0 {
		return nil
	}
	return ids
}
