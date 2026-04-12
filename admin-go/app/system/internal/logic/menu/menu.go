package menu

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

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
	normalizeMenuTypeFields(in.Type, &in.Path, &in.Component, &in.LinkURL)
	if err := validateMenuFields(in.Title, in.Type, in.Path, in.Component, in.Permission, in.LinkURL, in.IsShow, in.IsCache, in.Status, in.Sort); err != nil {
		return err
	}
	if err := s.ensureMenuPathUnique(ctx, 0, in.Type, in.Path); err != nil {
		return err
	}
	if err := s.ensureMenuPermissionUnique(ctx, 0, in.Permission); err != nil {
		return err
	}
	if err := s.ensureParentValid(ctx, in.ParentID, 0); err != nil {
		return err
	}
	id := snowflake.Generate()
	_, err := dao.Menu.Ctx(ctx).Data(do.Menu{
		Id:         id,
		ParentId:   in.ParentID,
		Title:      in.Title,
		Type:       in.Type,
		Path:       in.Path,
		Component:  in.Component,
		Permission: in.Permission,
		Icon:       in.Icon,
		Sort:       in.Sort,
		IsShow:     in.IsShow,
		IsCache:    in.IsCache,
		LinkUrl:    in.LinkURL,
		Status:     in.Status,
	}).Insert()
	if err == nil {
		authlogic.ClearUserCaches(ctx, s.getAdminUserIDs(ctx)...)
	}
	return err
}

// Update 更新菜单表
func (s *sMenu) Update(ctx context.Context, in *model.MenuUpdateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeMenuUpdateInput(in)
	normalizeMenuTypeFields(in.Type, &in.Path, &in.Component, &in.LinkURL)
	if err := validateMenuFields(in.Title, in.Type, in.Path, in.Component, in.Permission, in.LinkURL, in.IsShow, in.IsCache, in.Status, in.Sort); err != nil {
		return err
	}
	if err := s.ensureMenuExists(ctx, in.ID); err != nil {
		return err
	}
	if err := s.ensureMenuPathUnique(ctx, in.ID, in.Type, in.Path); err != nil {
		return err
	}
	if err := s.ensureMenuPermissionUnique(ctx, in.ID, in.Permission); err != nil {
		return err
	}
	if err := s.ensureParentValid(ctx, in.ParentID, in.ID); err != nil {
		return err
	}
	data := do.Menu{
		ParentId:   in.ParentID,
		Title:      in.Title,
		Type:       in.Type,
		Path:       in.Path,
		Component:  in.Component,
		Permission: in.Permission,
		Icon:       in.Icon,
		Sort:       in.Sort,
		IsShow:     in.IsShow,
		IsCache:    in.IsCache,
		LinkUrl:    in.LinkURL,
		Status:     in.Status,
	}
	_, err := dao.Menu.Ctx(ctx).
		Where(dao.Menu.Columns().Id, in.ID).
		Where(dao.Menu.Columns().DeletedAt, nil).
		Data(data).
		Update()
	if err == nil {
		authlogic.ClearUserCaches(ctx, s.getMenuAffectedUserIDs(ctx, []int64{int64(in.ID)})...)
	}
	return err
}

// Delete 软删除菜单表
func (s *sMenu) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := s.ensureMenuExists(ctx, id); err != nil {
		return err
	}
	if err := s.ensureMenuDeletable(ctx, id); err != nil {
		return err
	}
	affectedUserIDs := s.getMenuAffectedUserIDs(ctx, []int64{int64(id)})
	err := dao.Menu.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.Menu.Table()).Ctx(ctx).
			Where(dao.Menu.Columns().Id, id).
			Delete(); err != nil {
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
		authlogic.ClearUserCaches(ctx, affectedUserIDs...)
	}
	return err
}

// BatchDelete 批量删除菜单表
func (s *sMenu) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	ids = batchutil.CompactIDs(ids)
	if len(ids) == 0 {
		return gerror.New("请选择要删除的菜单")
	}
	if err := s.ensureMenuIDsExist(ctx, ids); err != nil {
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
		if err := s.ensureMenuBatchDeletable(ctx, id, deleteIDs); err != nil {
			return err
		}
	}
	affectedUserIDs := s.getMenuAffectedUserIDs(ctx, deleteIDs)
	err = dao.Menu.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.Menu.Table()).Ctx(ctx).
			WhereIn(dao.Menu.Columns().Id, deleteIDs).
			Delete(); err != nil {
			return err
		}
		if _, err := tx.Model(dao.RoleMenu.Table()).Ctx(ctx).
			WhereIn(dao.RoleMenu.Columns().MenuId, deleteIDs).
			Delete(); err != nil {
			return err
		}
		return nil
	})
	if err == nil {
		authlogic.ClearUserCaches(ctx, affectedUserIDs...)
	}
	return err
}

// Detail 获取菜单表详情
func (s *sMenu) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.MenuDetailOutput, err error) {
	if id <= 0 {
		return nil, gerror.New("菜单不存在或已删除")
	}
	out = &model.MenuDetailOutput{}
	err = dao.Menu.Ctx(ctx).Where(dao.Menu.Columns().Id, id).Where(dao.Menu.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("菜单不存在或已删除")
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

func normalizeMenuTypeFields(menuType int, path, component, linkURL *string) {
	if path == nil || component == nil || linkURL == nil {
		return
	}
	switch menuType {
	case 1:
		*component = ""
		*linkURL = ""
	case 2:
		*linkURL = ""
	case 3:
		*path = ""
		*component = ""
		*linkURL = ""
	}
}

func validateMenuFields(title string, menuType int, path, component, permission, linkURL string, isShow, isCache, status, sort int) error {
	if title == "" {
		return gerror.New("菜单名称不能为空")
	}
	if err := fieldvalid.Enum("菜单类型", menuType, 1, 2, 3, 4, 5); err != nil {
		return err
	}
	if err := fieldvalid.NonNegative("排序", sort); err != nil {
		return err
	}
	if err := fieldvalid.Binary("是否显示", isShow); err != nil {
		return err
	}
	if err := fieldvalid.Binary("是否缓存", isCache); err != nil {
		return err
	}
	if err := fieldvalid.Binary("状态", status); err != nil {
		return err
	}
	if isRouteMenuType(menuType) {
		if path == "" {
			return gerror.New("当前菜单类型必须填写前端路由路径")
		}
		if !strings.HasPrefix(path, "/") {
			return gerror.New("前端路由路径必须以 / 开头")
		}
	}
	switch menuType {
	case 2:
		if component == "" {
			return gerror.New("菜单类型必须填写前端组件路径")
		}
	case 3:
		if permission == "" {
			return gerror.New("按钮类型必须填写权限标识")
		}
	case 4, 5:
		if linkURL == "" {
			return gerror.New("外链/内链类型必须填写地址")
		}
	}
	return nil
}

func isRouteMenuType(menuType int) bool {
	switch menuType {
	case 1, 2, 4, 5:
		return true
	default:
		return false
	}
}

func (s *sMenu) fillParentTitles(ctx context.Context, list []*model.MenuListOutput) {
	parentIDs := make([]int64, 0, len(list))
	for _, item := range list {
		if item.ParentID != 0 {
			parentIDs = append(parentIDs, int64(item.ParentID))
		}
	}
	parentMap := shared.LoadTitleMap(ctx, "system_menu", parentIDs)
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

func (s *sMenu) ensureMenuExists(ctx context.Context, id snowflake.JsonInt64) error {
	if id <= 0 {
		return gerror.New("菜单不存在或已删除")
	}
	count, err := dao.Menu.Ctx(ctx).
		Where(dao.Menu.Columns().Id, id).
		Where(dao.Menu.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return gerror.New("菜单不存在或已删除")
	}
	return nil
}

func (s *sMenu) ensureMenuIDsExist(ctx context.Context, ids []snowflake.JsonInt64) error {
	dbIDs := batchutil.ToInt64s(ids)
	count, err := dao.Menu.Ctx(ctx).
		WhereIn(dao.Menu.Columns().Id, dbIDs).
		Where(dao.Menu.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if count != len(dbIDs) {
		return gerror.New("包含不存在或已删除的菜单")
	}
	return nil
}

func (s *sMenu) ensureMenuPermissionUnique(ctx context.Context, currentID snowflake.JsonInt64, permission string) error {
	permission = strings.TrimSpace(permission)
	if permission == "" {
		return nil
	}
	m := dao.Menu.Ctx(ctx).
		Where(dao.Menu.Columns().Permission, permission).
		Where(dao.Menu.Columns().DeletedAt, nil)
	if currentID > 0 {
		m = m.WhereNot(dao.Menu.Columns().Id, currentID)
	}
	count, err := m.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New("权限标识已存在")
	}
	return nil
}

func (s *sMenu) ensureMenuPathUnique(ctx context.Context, currentID snowflake.JsonInt64, menuType int, path string) error {
	path = strings.TrimSpace(path)
	if !isRouteMenuType(menuType) || path == "" {
		return nil
	}
	m := dao.Menu.Ctx(ctx).
		Where(dao.Menu.Columns().Path, path).
		WhereIn(dao.Menu.Columns().Type, []int{1, 2, 4, 5}).
		Where(dao.Menu.Columns().DeletedAt, nil)
	if currentID > 0 {
		m = m.WhereNot(dao.Menu.Columns().Id, currentID)
	}
	count, err := m.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New("前端路由路径已存在")
	}
	return nil
}

func (s *sMenu) ensureMenuBatchDeletable(ctx context.Context, id snowflake.JsonInt64, deleteIDs []int64) error {
	childModel := dao.Menu.Ctx(ctx).
		Where(dao.Menu.Columns().ParentId, id).
		Where(dao.Menu.Columns().DeletedAt, nil)
	if len(deleteIDs) > 0 {
		childModel = childModel.WhereNotIn(dao.Menu.Columns().Id, deleteIDs)
	}
	childCount, err := childModel.Count()
	if err != nil {
		return err
	}
	if childCount > 0 {
		return gerror.New("当前菜单下存在子菜单，不能直接删除")
	}
	return nil
}

func (s *sMenu) collectBatchDeleteOrder(ctx context.Context, ids []snowflake.JsonInt64) ([]snowflake.JsonInt64, error) {
	var rows []struct {
		Id       int64 `json:"id"`
		ParentId int64 `json:"parentId"`
	}
	if err := dao.Menu.Ctx(ctx).
		Fields(dao.Menu.Columns().Id, dao.Menu.Columns().ParentId).
		Where(dao.Menu.Columns().DeletedAt, nil).
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

func (s *sMenu) getMenuAffectedUserIDs(ctx context.Context, menuIDs []int64) []int64 {
	menuIDs = compactMenuIDs(menuIDs)
	userIDs := make([]int64, 0)
	userIDs = append(userIDs, s.getAdminUserIDs(ctx)...)
	if len(menuIDs) == 0 {
		return compactMenuIDs(userIDs)
	}

	var rows []struct {
		UserId int64 `json:"userId"`
	}
	if err := g.DB().Ctx(ctx).
		Model(dao.UserRole.Table()+" ur").
		LeftJoin(dao.Role.Table()+" r", "r.id = ur.role_id").
		LeftJoin(dao.RoleMenu.Table()+" rm", "rm.role_id = ur.role_id").
		LeftJoin(dao.Users.Table()+" u", "u.id = ur.user_id").
		Fields("ur.user_id AS userId").
		WhereIn("rm.menu_id", menuIDs).
		Where("r.deleted_at", nil).
		Where("r.status", 1).
		Where("u.deleted_at", nil).
		Scan(&rows); err != nil {
		return compactMenuIDs(userIDs)
	}
	for _, row := range rows {
		userIDs = append(userIDs, row.UserId)
	}
	return compactMenuIDs(userIDs)
}

func (s *sMenu) getAdminUserIDs(ctx context.Context) []int64 {
	var rows []struct {
		UserId int64 `json:"userId"`
	}
	if err := g.DB().Ctx(ctx).
		Model(dao.UserRole.Table()+" ur").
		LeftJoin(dao.Role.Table()+" r", "r.id = ur.role_id").
		LeftJoin(dao.Users.Table()+" u", "u.id = ur.user_id").
		Fields("ur.user_id AS userId").
		Where("r.is_admin", 1).
		Where("r.deleted_at", nil).
		Where("r.status", 1).
		Where("u.deleted_at", nil).
		Scan(&rows); err != nil {
		return nil
	}
	userIDs := make([]int64, 0, len(rows))
	for _, row := range rows {
		userIDs = append(userIDs, row.UserId)
	}
	return compactMenuIDs(userIDs)
}

func compactMenuIDs(values []int64) []int64 {
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
