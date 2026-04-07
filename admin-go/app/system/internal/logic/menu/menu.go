package menu

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/system/internal/dao"
	authlogic "gbaseadmin/app/system/internal/logic/auth"
	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/app/system/internal/service"
	"gbaseadmin/utility/snowflake"
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
	// 查询上级菜单ID，0 表示顶级菜单关联显示
	if out.ParentID != 0 {
		val, _ := g.DB().Ctx(ctx).Model("system_menu").Where("id", out.ParentID).Where("deleted_at", nil).Value("title")
		out.MenuTitle = val.String()
	}
	return
}

// List 获取菜单表列表
func (s *sMenu) List(ctx context.Context, in *model.MenuListInput) (list []*model.MenuListOutput, total int, err error) {
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

	// 使用 map 迭代方式组装树
	nodeMap := make(map[int64]*model.MenuTreeOutput, len(list))
	for _, item := range list {
		item.Children = make([]*model.MenuTreeOutput, 0)
		nodeMap[int64(item.ID)] = item
	}

	tree = make([]*model.MenuTreeOutput, 0)
	for _, item := range list {
		if int64(item.ParentID) == 0 {
			tree = append(tree, item)
		} else if parent, ok := nodeMap[int64(item.ParentID)]; ok {
			parent.Children = append(parent.Children, item)
		} else {
			tree = append(tree, item)
		}
	}
	return
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
	if parentID == 0 {
		return nil
	}
	if currentID != 0 && parentID == currentID {
		return gerror.New("上级菜单不能选择自己")
	}
	var parent struct {
		Id       int64 `json:"id"`
		ParentId int64 `json:"parentId"`
	}
	if err := dao.Menu.Ctx(ctx).
		Fields(dao.Menu.Columns().Id, dao.Menu.Columns().ParentId).
		Where(dao.Menu.Columns().Id, parentID).
		Where(dao.Menu.Columns().DeletedAt, nil).
		Scan(&parent); err != nil {
		return err
	}
	if parent.Id == 0 {
		return gerror.New("上级菜单不存在或已删除")
	}
	seen := map[int64]struct{}{int64(parentID): {}}
	for parent.ParentId != 0 {
		if currentID != 0 && parent.ParentId == int64(currentID) {
			return gerror.New("不能将菜单移动到自己的子级下")
		}
		if _, ok := seen[parent.ParentId]; ok {
			return gerror.New("菜单层级存在循环引用")
		}
		seen[parent.ParentId] = struct{}{}
		next := struct {
			Id       int64 `json:"id"`
			ParentId int64 `json:"parentId"`
		}{}
		if err := dao.Menu.Ctx(ctx).
			Fields(dao.Menu.Columns().Id, dao.Menu.Columns().ParentId).
			Where(dao.Menu.Columns().Id, parent.ParentId).
			Where(dao.Menu.Columns().DeletedAt, nil).
			Scan(&next); err != nil {
			return err
		}
		if next.Id == 0 {
			return gerror.New("上级菜单链路中存在无效节点")
		}
		parent = next
	}
	return nil
}
