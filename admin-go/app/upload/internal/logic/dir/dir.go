package dir

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/upload/internal/dao"
	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/app/upload/internal/service"
	"gbaseadmin/utility/snowflake"
)

func init() {
	service.RegisterDir(New())
}

func New() *sDir {
	return &sDir{}
}

type sDir struct{}

// Create 创建文件目录
func (s *sDir) Create(ctx context.Context, in *model.DirCreateInput) error {
	if err := s.ensureParentValid(ctx, in.ParentID, 0); err != nil {
		return err
	}
	id := snowflake.Generate()
	_, err := dao.UploadDir.Ctx(ctx).Data(g.Map{
		dao.UploadDir.Columns().Id:        id,
		dao.UploadDir.Columns().ParentId:  in.ParentID,
		dao.UploadDir.Columns().Name:      in.Name,
		dao.UploadDir.Columns().Path:      in.Path,
		dao.UploadDir.Columns().Sort:      in.Sort,
		dao.UploadDir.Columns().Status:    in.Status,
		dao.UploadDir.Columns().CreatedAt: gtime.Now(),
		dao.UploadDir.Columns().UpdatedAt: gtime.Now(),
	}).Insert()
	return err
}

// Update 更新文件目录
func (s *sDir) Update(ctx context.Context, in *model.DirUpdateInput) error {
	if err := s.ensureParentValid(ctx, in.ParentID, in.ID); err != nil {
		return err
	}
	data := g.Map{
		dao.UploadDir.Columns().ParentId:  in.ParentID,
		dao.UploadDir.Columns().Name:      in.Name,
		dao.UploadDir.Columns().Path:      in.Path,
		dao.UploadDir.Columns().Sort:      in.Sort,
		dao.UploadDir.Columns().Status:    in.Status,
		dao.UploadDir.Columns().UpdatedAt: gtime.Now(),
	}
	_, err := dao.UploadDir.Ctx(ctx).Where(dao.UploadDir.Columns().Id, in.ID).Data(data).Update()
	return err
}

// Delete 软删除文件目录
func (s *sDir) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := s.ensureDirDeletable(ctx, id); err != nil {
		return err
	}
	_, err := dao.UploadDir.Ctx(ctx).Where(dao.UploadDir.Columns().Id, id).Data(g.Map{
		dao.UploadDir.Columns().DeletedAt: gtime.Now(),
	}).Update()
	return err
}

// Detail 获取文件目录详情
func (s *sDir) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.DirDetailOutput, err error) {
	out = &model.DirDetailOutput{}
	err = dao.UploadDir.Ctx(ctx).Where(dao.UploadDir.Columns().Id, id).Where(dao.UploadDir.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	// 查询上级目录关联显示
	if out.ParentID != 0 {
		val, err := g.DB().Ctx(ctx).Model("upload_dir").Where("id", out.ParentID).Where("deleted_at", nil).Value("name")
		if err == nil {
			out.DirName = val.String()
		}
	}
	return
}

// List 获取文件目录列表
func (s *sDir) List(ctx context.Context, in *model.DirListInput) (list []*model.DirListOutput, total int, err error) {
	m := dao.UploadDir.Ctx(ctx).Where(dao.UploadDir.Columns().DeletedAt, nil)
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
	if in != nil {
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

	// 使用 map 迭代方式组装树
	nodeMap := make(map[int64]*model.DirTreeOutput, len(list))
	for _, item := range list {
		item.Children = make([]*model.DirTreeOutput, 0)
		nodeMap[int64(item.ID)] = item
	}

	tree = make([]*model.DirTreeOutput, 0)
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

func (s *sDir) fillParentNames(ctx context.Context, list []*model.DirListOutput) {
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
	rows, err := g.DB().Ctx(ctx).Model("upload_dir").
		Fields("id", "name").
		Where("deleted_at", nil).
		WhereIn("id", parentIDs).
		All()
	if err != nil {
		return
	}
	parentMap := make(map[int64]string, len(rows))
	for _, row := range rows {
		parentMap[row["id"].Int64()] = row["name"].String()
	}
	for _, item := range list {
		item.DirName = parentMap[int64(item.ParentID)]
	}
}

func (s *sDir) ensureDirDeletable(ctx context.Context, id snowflake.JsonInt64) error {
	childCount, err := dao.UploadDir.Ctx(ctx).
		Where(dao.UploadDir.Columns().ParentId, id).
		Where(dao.UploadDir.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if childCount > 0 {
		return gerror.New("当前目录下存在子目录，不能直接删除")
	}
	fileCount, err := dao.UploadFile.Ctx(ctx).
		Where(dao.UploadFile.Columns().DirId, id).
		Where(dao.UploadFile.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if fileCount > 0 {
		return gerror.New("当前目录下仍有关联文件，不能直接删除")
	}
	ruleCount, err := dao.UploadDirRule.Ctx(ctx).
		Where(dao.UploadDirRule.Columns().DirId, id).
		Where(dao.UploadDirRule.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if ruleCount > 0 {
		return gerror.New("当前目录仍有关联规则，不能直接删除")
	}
	return nil
}

func (s *sDir) ensureParentValid(ctx context.Context, parentID, currentID snowflake.JsonInt64) error {
	if parentID == 0 {
		return nil
	}
	if currentID != 0 && parentID == currentID {
		return gerror.New("上级目录不能选择自己")
	}
	var parent struct {
		Id       int64 `json:"id"`
		ParentId int64 `json:"parentId"`
	}
	if err := dao.UploadDir.Ctx(ctx).
		Fields(dao.UploadDir.Columns().Id, dao.UploadDir.Columns().ParentId).
		Where(dao.UploadDir.Columns().Id, parentID).
		Where(dao.UploadDir.Columns().DeletedAt, nil).
		Scan(&parent); err != nil {
		return err
	}
	if parent.Id == 0 {
		return gerror.New("上级目录不存在或已删除")
	}
	seen := map[int64]struct{}{int64(parentID): {}}
	for parent.ParentId != 0 {
		if currentID != 0 && parent.ParentId == int64(currentID) {
			return gerror.New("不能将目录移动到自己的子级下")
		}
		if _, ok := seen[parent.ParentId]; ok {
			return gerror.New("目录层级存在循环引用")
		}
		seen[parent.ParentId] = struct{}{}
		next := struct {
			Id       int64 `json:"id"`
			ParentId int64 `json:"parentId"`
		}{}
		if err := dao.UploadDir.Ctx(ctx).
			Fields(dao.UploadDir.Columns().Id, dao.UploadDir.Columns().ParentId).
			Where(dao.UploadDir.Columns().Id, parent.ParentId).
			Where(dao.UploadDir.Columns().DeletedAt, nil).
			Scan(&next); err != nil {
			return err
		}
		if next.Id == 0 {
			return gerror.New("上级目录链路中存在无效节点")
		}
		parent = next
	}
	return nil
}
