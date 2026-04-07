package menu

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
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
	service.RegisterMenu(New())
}

func New() *sMenu {
	return &sMenu{}
}

type sMenu struct{}

// Create 创建菜单表
func (s *sMenu) Create(ctx context.Context, in *model.MenuCreateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeMenuCreateInput(in)
	if err := s.ensureParentValid(ctx, in.ParentID, 0); err != nil {
		return err
	}
	id := snowflake.Generate()
	_, err := dao.Menu.Ctx(ctx).Data(g.Map{
		dao.Menu.Columns().Id:         id,
		dao.Menu.Columns().ParentId:   in.ParentID,
		dao.Menu.Columns().Title:      in.Title,
		dao.Menu.Columns().Type:       in.Type,
		dao.Menu.Columns().Path:       in.Path,
		dao.Menu.Columns().Component:  in.Component,
		dao.Menu.Columns().Permission: in.Permission,
		dao.Menu.Columns().Icon:       in.Icon,
		dao.Menu.Columns().Sort:       in.Sort,
		dao.Menu.Columns().IsShow:     in.IsShow,
		dao.Menu.Columns().IsCache:    in.IsCache,
		dao.Menu.Columns().LinkUrl:    in.LinkURL,
		dao.Menu.Columns().Status:     in.Status,
		dao.Menu.Columns().CreatedAt:  gtime.Now(),
		dao.Menu.Columns().UpdatedAt:  gtime.Now(),
	}).Insert()
	if err == nil {
		authlogic.ClearAllUserCaches(ctx)
	}
	return err
}

// Update 更新菜单表
func (s *sMenu) Update(ctx context.Context, in *model.MenuUpdateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeMenuUpdateInput(in)
	if err := s.ensureParentValid(ctx, in.ParentID, in.ID); err != nil {
		return err
	}
	data := g.Map{
		dao.Menu.Columns().ParentId:   in.ParentID,
		dao.Menu.Columns().Title:      in.Title,
		dao.Menu.Columns().Type:       in.Type,
		dao.Menu.Columns().Path:       in.Path,
		dao.Menu.Columns().Component:  in.Component,
		dao.Menu.Columns().Permission: in.Permission,
		dao.Menu.Columns().Icon:       in.Icon,
		dao.Menu.Columns().Sort:       in.Sort,
		dao.Menu.Columns().IsShow:     in.IsShow,
		dao.Menu.Columns().IsCache:    in.IsCache,
		dao.Menu.Columns().LinkUrl:    in.LinkURL,
		dao.Menu.Columns().Status:     in.Status,
		dao.Menu.Columns().UpdatedAt:  gtime.Now(),
	}
	_, err := dao.Menu.Ctx(ctx).Where(dao.Menu.Columns().Id, in.ID).Data(data).Update()
	if err == nil {
		authlogic.ClearAllUserCaches(ctx)
	}
	return err
}

// Delete 软删除菜单表
func (s *sMenu) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := s.ensureMenuDeletable(ctx, id); err != nil {
		return err
	}
	err := dao.Menu.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.Menu.Table()).Ctx(ctx).
			Where(dao.Menu.Columns().Id, id).
			Where(dao.Menu.Columns().DeletedAt, nil).
			Data(g.Map{
				dao.Menu.Columns().DeletedAt: gtime.Now(),
			}).
			Update(); err != nil {
			return err
		}
		if _, err := tx.Model(dao.RoleMenu.Table()).Ctx(ctx).
			Where(dao.RoleMenu.Columns().MenuId, id).
			Delete(); err != nil {
			return err
		}
		return nil
	})
	if err == nil {
		authlogic.ClearAllUserCaches(ctx)
	}
	return err
}

// Detail 获取菜单表详情
func (s *sMenu) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.MenuDetailOutput, err error) {
	out = &model.MenuDetailOutput{}
	err = dao.Menu.Ctx(ctx).Where(dao.Menu.Columns().Id, id).Where(dao.Menu.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	out.MenuTitle = shared.LookupTitle(ctx, "system_menu", int64(out.ParentID))
	return
}

// List 获取菜单表列表
func (s *sMenu) List(ctx context.Context, in *model.MenuListInput) (list []*model.MenuListOutput, total int, err error) {
	if in == nil {
		in = &model.MenuListInput{}
	}
	normalizeMenuListInput(in)
	m := dao.Menu.Ctx(ctx).Where(dao.Menu.Columns().DeletedAt, nil)
	if in.Keyword != "" {
		keywordBuilder := m.Builder().
			WhereLike(dao.Menu.Columns().Title, "%"+in.Keyword+"%").
			WhereOrLike(dao.Menu.Columns().Path, "%"+in.Keyword+"%").
			WhereOrLike(dao.Menu.Columns().Component, "%"+in.Keyword+"%").
			WhereOrLike(dao.Menu.Columns().Permission, "%"+in.Keyword+"%")
		m = m.Where(keywordBuilder)
	}
	if in.Type > 0 {
		m = m.Where(dao.Menu.Columns().Type, in.Type)
	}
	if in.IsShow != nil {
		m = m.Where(dao.Menu.Columns().IsShow, *in.IsShow)
	}
	if in.IsCache != nil {
		m = m.Where(dao.Menu.Columns().IsCache, *in.IsCache)
	}
	if in.Status != nil {
		m = m.Where(dao.Menu.Columns().Status, *in.Status)
	}
	total, err = m.Count()
	if err != nil {
		return
	}
	in.PageNum, in.PageSize = pageutil.Normalize(in.PageNum, in.PageSize)
	err = m.Page(in.PageNum, in.PageSize).OrderAsc(dao.Menu.Columns().Id).Scan(&list)
	if err != nil {
		return
	}
	s.fillParentTitles(ctx, list)
	return
}

// Tree 获取菜单表树形结构
func (s *sMenu) Tree(ctx context.Context, in *model.MenuTreeInput) (tree []*model.MenuTreeOutput, err error) {
	var list []*model.MenuTreeOutput
	m := dao.Menu.Ctx(ctx).Where(dao.Menu.Columns().DeletedAt, nil)
	if in != nil {
		normalizeMenuTreeInput(in)
		if in.Keyword != "" {
			keywordBuilder := m.Builder().
				WhereLike(dao.Menu.Columns().Title, "%"+in.Keyword+"%").
				WhereOrLike(dao.Menu.Columns().Path, "%"+in.Keyword+"%").
				WhereOrLike(dao.Menu.Columns().Component, "%"+in.Keyword+"%").
				WhereOrLike(dao.Menu.Columns().Permission, "%"+in.Keyword+"%")
			m = m.Where(keywordBuilder)
		}
		if in.Type > 0 {
			m = m.Where(dao.Menu.Columns().Type, in.Type)
		}
		if in.IsShow != nil {
			m = m.Where(dao.Menu.Columns().IsShow, *in.IsShow)
		}
		if in.IsCache != nil {
			m = m.Where(dao.Menu.Columns().IsCache, *in.IsCache)
		}
		if in.Status != nil {
			m = m.Where(dao.Menu.Columns().Status, *in.Status)
		}
	}
	err = m.OrderAsc(dao.Menu.Columns().Sort).Scan(&list)
	if err != nil {
		return
	}

	tree = treeutil.BuildForest(list, treeutil.TreeNodeAccessor[*model.MenuTreeOutput]{
		ID:       func(item *model.MenuTreeOutput) int64 { return int64(item.ID) },
		ParentID: func(item *model.MenuTreeOutput) int64 { return int64(item.ParentID) },
		Init: func(item *model.MenuTreeOutput) {
			item.Children = make([]*model.MenuTreeOutput, 0)
		},
		Append: func(parent *model.MenuTreeOutput, child *model.MenuTreeOutput) {
			parent.Children = append(parent.Children, child)
		},
	})
	return
}

func normalizeMenuCreateInput(in *model.MenuCreateInput) {
	if in == nil {
		return
	}
	in.Title = strings.TrimSpace(in.Title)
	in.Path = strings.TrimSpace(in.Path)
	in.Component = strings.TrimSpace(in.Component)
	in.Permission = strings.TrimSpace(in.Permission)
	in.Icon = strings.TrimSpace(in.Icon)
	in.LinkURL = strings.TrimSpace(in.LinkURL)
}

func normalizeMenuUpdateInput(in *model.MenuUpdateInput) {
	if in == nil {
		return
	}
	in.Title = strings.TrimSpace(in.Title)
	in.Path = strings.TrimSpace(in.Path)
	in.Component = strings.TrimSpace(in.Component)
	in.Permission = strings.TrimSpace(in.Permission)
	in.Icon = strings.TrimSpace(in.Icon)
	in.LinkURL = strings.TrimSpace(in.LinkURL)
}

func normalizeMenuListInput(in *model.MenuListInput) {
	if in == nil {
		return
	}
	in.Keyword = strings.TrimSpace(in.Keyword)
}

func normalizeMenuTreeInput(in *model.MenuTreeInput) {
	if in == nil {
		return
	}
	in.Keyword = strings.TrimSpace(in.Keyword)
}

func (s *sMenu) fillParentTitles(ctx context.Context, list []*model.MenuListOutput) {
	parentSet := make(map[int64]struct{})
	for _, item := range list {
		if item.ParentID != 0 {
			parentSet[int64(item.ParentID)] = struct{}{}
		}
	}
	if len(parentSet) == 0 {
		return
	}
	parentIDs := make([]int64, 0, len(parentSet))
	for id := range parentSet {
		parentIDs = append(parentIDs, id)
	}
	rows, err := g.DB().Ctx(ctx).Model("system_menu").
		Fields("id", "title").
		Where("deleted_at", nil).
		WhereIn("id", parentIDs).
		All()
	if err != nil {
		return
	}
	parentMap := make(map[int64]string, len(rows))
	for _, row := range rows {
		parentMap[row["id"].Int64()] = row["title"].String()
	}
	for _, item := range list {
		item.MenuTitle = parentMap[int64(item.ParentID)]
	}
}

func (s *sMenu) ensureMenuDeletable(ctx context.Context, id snowflake.JsonInt64) error {
	childCount, err := dao.Menu.Ctx(ctx).
		Where(dao.Menu.Columns().ParentId, id).
		Where(dao.Menu.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if childCount > 0 {
		return gerror.New("当前菜单下存在子菜单，不能直接删除")
	}
	return nil
}

func (s *sMenu) ensureParentValid(ctx context.Context, parentID, currentID snowflake.JsonInt64) error {
	return treeutil.ValidateParent(parentID, currentID, func(id int64) (int64, int64, error) {
		var parent struct {
			Id       int64 `json:"id"`
			ParentId int64 `json:"parentId"`
		}
		if err := dao.Menu.Ctx(ctx).
			Fields(dao.Menu.Columns().Id, dao.Menu.Columns().ParentId).
			Where(dao.Menu.Columns().Id, id).
			Where(dao.Menu.Columns().DeletedAt, nil).
			Scan(&parent); err != nil {
			return 0, 0, err
		}
		return parent.Id, parent.ParentId, nil
	}, treeutil.Messages{
		Self:         "上级菜单不能选择自己",
		Missing:      "上级菜单不存在或已删除",
		ChildLoop:    "不能将菜单移动到自己的子级下",
		Cycle:        "菜单层级存在循环引用",
		InvalidChain: "上级菜单链路中存在无效节点",
	})
}
