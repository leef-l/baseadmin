package dept

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/system/internal/dao"
	authlogic "gbaseadmin/app/system/internal/logic/auth"
	"gbaseadmin/app/system/internal/logic/shared"
	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/app/system/internal/service"
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
	if err := s.ensureParentValid(ctx, in.ParentID, 0); err != nil {
		return err
	}
	id := snowflake.Generate()
	_, err := dao.Dept.Ctx(ctx).Data(g.Map{
		dao.Dept.Columns().Id:        id,
		dao.Dept.Columns().ParentId:  in.ParentID,
		dao.Dept.Columns().Title:     in.Title,
		dao.Dept.Columns().Username:  in.Username,
		dao.Dept.Columns().Email:     in.Email,
		dao.Dept.Columns().Sort:      in.Sort,
		dao.Dept.Columns().Status:    in.Status,
		dao.Dept.Columns().CreatedAt: gtime.Now(),
		dao.Dept.Columns().UpdatedAt: gtime.Now(),
	}).Insert()
	if err == nil {
		authlogic.ClearAllUserCaches(ctx)
	}
	return err
}

// Update 更新部门表
func (s *sDept) Update(ctx context.Context, in *model.DeptUpdateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeDeptUpdateInput(in)
	if err := s.ensureParentValid(ctx, in.ParentID, in.ID); err != nil {
		return err
	}
	data := g.Map{
		dao.Dept.Columns().ParentId:  in.ParentID,
		dao.Dept.Columns().Title:     in.Title,
		dao.Dept.Columns().Username:  in.Username,
		dao.Dept.Columns().Email:     in.Email,
		dao.Dept.Columns().Sort:      in.Sort,
		dao.Dept.Columns().Status:    in.Status,
		dao.Dept.Columns().UpdatedAt: gtime.Now(),
	}
	_, err := dao.Dept.Ctx(ctx).
		Where(dao.Dept.Columns().Id, in.ID).
		Where(dao.Dept.Columns().DeletedAt, nil).
		Data(data).
		Update()
	if err == nil {
		authlogic.ClearAllUserCaches(ctx)
	}
	return err
}

// Delete 软删除部门表
func (s *sDept) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := s.ensureDeptDeletable(ctx, id); err != nil {
		return err
	}
	_, err := dao.Dept.Ctx(ctx).
		Where(dao.Dept.Columns().Id, id).
		Where(dao.Dept.Columns().DeletedAt, nil).
		Data(g.Map{
			dao.Dept.Columns().DeletedAt: gtime.Now(),
		}).
		Update()
	if err == nil {
		authlogic.ClearAllUserCaches(ctx)
	}
	return err
}

// Detail 获取部门表详情
func (s *sDept) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.DeptDetailOutput, err error) {
	out = &model.DeptDetailOutput{}
	err = dao.Dept.Ctx(ctx).Where(dao.Dept.Columns().Id, id).Where(dao.Dept.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
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
