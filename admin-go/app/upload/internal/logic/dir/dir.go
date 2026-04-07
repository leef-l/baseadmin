package dir

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/upload/internal/dao"
	"gbaseadmin/app/upload/internal/logic/shared"
	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/app/upload/internal/service"
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

// Create 创建文件目录
func (s *sDir) Create(ctx context.Context, in *model.DirCreateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeDirCreateInput(in)
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
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeDirUpdateInput(in)
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
	_, err := dao.UploadDir.Ctx(ctx).
		Where(dao.UploadDir.Columns().Id, in.ID).
		Where(dao.UploadDir.Columns().DeletedAt, nil).
		Data(data).
		Update()
	return err
}

// Delete 软删除文件目录
func (s *sDir) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := s.ensureDirDeletable(ctx, id); err != nil {
		return err
	}
	_, err := dao.UploadDir.Ctx(ctx).
		Where(dao.UploadDir.Columns().Id, id).
		Where(dao.UploadDir.Columns().DeletedAt, nil).
		Data(g.Map{
			dao.UploadDir.Columns().DeletedAt: gtime.Now(),
		}).
		Update()
	return err
}

// Detail 获取文件目录详情
func (s *sDir) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.DirDetailOutput, err error) {
	out = &model.DirDetailOutput{}
	err = dao.UploadDir.Ctx(ctx).Where(dao.UploadDir.Columns().Id, id).Where(dao.UploadDir.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
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
	return treeutil.ValidateParent(parentID, currentID, func(id int64) (int64, int64, error) {
		var parent struct {
			Id       int64 `json:"id"`
			ParentId int64 `json:"parentId"`
		}
		if err := dao.UploadDir.Ctx(ctx).
			Fields(dao.UploadDir.Columns().Id, dao.UploadDir.Columns().ParentId).
			Where(dao.UploadDir.Columns().Id, id).
			Where(dao.UploadDir.Columns().DeletedAt, nil).
			Scan(&parent); err != nil {
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
