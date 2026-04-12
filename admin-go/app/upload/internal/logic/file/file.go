package file

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	cos "github.com/tencentyun/cos-go-sdk-v5"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"gbaseadmin/app/upload/internal/dao"
	"gbaseadmin/app/upload/internal/logic/shared"
	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/app/upload/internal/model/do"
	"gbaseadmin/app/upload/internal/model/entity"
	"gbaseadmin/app/upload/internal/service"
	"gbaseadmin/utility/batchutil"
	"gbaseadmin/utility/fieldvalid"
	"gbaseadmin/utility/inpututil"
	"gbaseadmin/utility/pageutil"
	"gbaseadmin/utility/snowflake"
)

func init() {
	service.RegisterFile(New())
}

func New() *sFile {
	return &sFile{}
}

type sFile struct{}

type fileDeleteTarget struct {
	ID      int64  `json:"id"`
	URL     string `json:"url"`
	Storage int    `json:"storage"`
}

type stagedLocalDelete struct {
	RecordID     int64
	OriginalPath string
	StagedPath   string
}

// Create 创建文件记录
func (s *sFile) Create(ctx context.Context, in *model.FileCreateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeFileCreateInput(in)
	if err := validateFileFields(in.Name, in.URL, in.Storage, in.IsImage, in.Size); err != nil {
		return err
	}
	if err := s.ensureDirExists(ctx, in.DirID); err != nil {
		return err
	}
	id := snowflake.Generate()
	_, err := dao.UploadFile.Ctx(ctx).Data(do.UploadFile{
		Id:      id,
		DirId:   in.DirID,
		Name:    in.Name,
		Url:     in.URL,
		Ext:     in.Ext,
		Size:    in.Size,
		Mime:    in.Mime,
		Storage: in.Storage,
		IsImage: in.IsImage,
	}).Insert()
	return err
}

// Update 更新文件记录
func (s *sFile) Update(ctx context.Context, in *model.FileUpdateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeFileUpdateInput(in)
	if err := validateFileFields(in.Name, in.URL, in.Storage, in.IsImage, in.Size); err != nil {
		return err
	}
	if err := s.ensureFileExists(ctx, in.ID); err != nil {
		return err
	}
	if err := s.ensureDirExists(ctx, in.DirID); err != nil {
		return err
	}
	data := do.UploadFile{
		DirId:   in.DirID,
		Name:    in.Name,
		Url:     in.URL,
		Ext:     in.Ext,
		Size:    in.Size,
		Mime:    in.Mime,
		Storage: in.Storage,
		IsImage: in.IsImage,
	}
	_, err := dao.UploadFile.Ctx(ctx).
		Where(dao.UploadFile.Columns().Id, in.ID).
		Where(dao.UploadFile.Columns().DeletedAt, nil).
		Data(data).
		Update()
	return err
}

// Delete 删除文件记录并物理删除文件
func (s *sFile) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	targets, err := s.loadDeleteTargets(ctx, []snowflake.JsonInt64{id})
	if err != nil {
		return err
	}
	if err := precheckDeleteTargets(ctx, targets); err != nil {
		return err
	}
	return s.deleteTargets(ctx, targets)
}

// BatchDelete 批量删除文件记录
func (s *sFile) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	ids = batchutil.CompactIDs(ids)
	if len(ids) == 0 {
		return gerror.New("请选择要删除的文件")
	}
	targets, err := s.loadDeleteTargets(ctx, ids)
	if err != nil {
		return err
	}
	if err := precheckDeleteTargets(ctx, targets); err != nil {
		return err
	}
	return s.deleteTargets(ctx, targets)
}

// deleteCloudFileOSS 从阿里云 OSS 删除文件
func deleteCloudFileOSS(ctx context.Context, fileURL string) error {
	configRecord, objectKey, err := loadUploadConfigByURL(ctx, 2, fileURL)
	if err != nil {
		return err
	}

	client, err := oss.New(configRecord.OssEndpoint, configRecord.OssAccessKey, configRecord.OssSecretKey)
	if err != nil {
		return fmt.Errorf("创建OSS客户端失败: %w", err)
	}
	b, err := client.Bucket(configRecord.OssBucket)
	if err != nil {
		return fmt.Errorf("获取OSS Bucket失败: %w", err)
	}
	return b.DeleteObject(objectKey)
}

// deleteCloudFileCOS 从腾讯云 COS 删除文件
func deleteCloudFileCOS(ctx context.Context, fileURL string) error {
	configRecord, objectKey, err := loadUploadConfigByURL(ctx, 3, fileURL)
	if err != nil {
		return err
	}

	bucketURL := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", configRecord.CosBucket, configRecord.CosRegion)
	u, err := url.Parse(bucketURL)
	if err != nil {
		return fmt.Errorf("解析COS URL失败: %w", err)
	}
	cosClient := cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  configRecord.CosSecretId,
			SecretKey: configRecord.CosSecretKey,
		},
	})
	_, err = cosClient.Object.Delete(ctx, objectKey, nil)
	return err
}

func deleteStoredFile(ctx context.Context, storage int, fileURL string) error {
	if strings.TrimSpace(fileURL) == "" {
		return nil
	}
	switch storage {
	case 1:
		localPath := shared.LocalStoragePhysicalPath(fileURL)
		if err := os.Remove(localPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("删除本地文件失败: %w", err)
		}
		return nil
	case 2:
		return deleteCloudFileOSS(ctx, fileURL)
	case 3:
		return deleteCloudFileCOS(ctx, fileURL)
	default:
		return nil
	}
}

func precheckDeleteTargets(ctx context.Context, targets []fileDeleteTarget) error {
	for _, target := range targets {
		switch target.Storage {
		case 2, 3:
			if _, _, err := loadUploadConfigByURL(ctx, target.Storage, target.URL); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *sFile) deleteTargets(ctx context.Context, targets []fileDeleteTarget) error {
	stagedLocalDeletes, err := stageLocalDeleteTargets(targets)
	if err != nil {
		return err
	}
	deleteIDs := make([]int64, 0, len(targets))
	for _, target := range targets {
		deleteIDs = append(deleteIDs, target.ID)
	}
	if err := s.softDeleteTargets(ctx, deleteIDs); err != nil {
		if restoreErr := restoreStagedLocalDeletes(stagedLocalDeletes); restoreErr != nil {
			return gerror.Wrapf(err, "删除文件记录失败，且恢复本地文件失败: %v", restoreErr)
		}
		return err
	}
	return s.finalizeDeleteTargets(ctx, targets, stagedLocalDeletes)
}

func (s *sFile) softDeleteTargets(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return dao.UploadFile.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Model(dao.UploadFile.Table()).Ctx(ctx).
			WhereIn(dao.UploadFile.Columns().Id, ids).
			Where(dao.UploadFile.Columns().DeletedAt, nil).
			Delete()
		return err
	})
}

func (s *sFile) restoreDeletedTargets(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	_, err := dao.UploadFile.Ctx(ctx).
		Unscoped().
		WhereIn(dao.UploadFile.Columns().Id, ids).
		Update(dao.UploadFile.Columns().DeletedAt + "=NULL")
	return err
}

func (s *sFile) finalizeDeleteTargets(ctx context.Context, targets []fileDeleteTarget, stagedLocalDeletes []stagedLocalDelete) error {
	stagedMap := make(map[int64]stagedLocalDelete, len(stagedLocalDeletes))
	for _, item := range stagedLocalDeletes {
		stagedMap[item.RecordID] = item
	}

	var deleteErrors []string
	for _, target := range targets {
		if target.Storage == 1 {
			if staged, ok := stagedMap[target.ID]; ok {
				if err := finalizeStagedLocalDelete(staged); err != nil {
					g.Log().Warningf(ctx, "failed to cleanup staged local file %d: %v", target.ID, err)
				}
			}
			continue
		}
		if err := deleteStoredFile(ctx, target.Storage, target.URL); err != nil {
			if restoreErr := s.restoreDeletedTargets(ctx, []int64{target.ID}); restoreErr != nil {
				deleteErrors = append(deleteErrors, fmt.Sprintf("文件 %d 删除失败，且记录恢复失败: %v / %v", target.ID, err, restoreErr))
				continue
			}
			deleteErrors = append(deleteErrors, fmt.Sprintf("文件 %d 删除失败，记录已恢复: %v", target.ID, err))
		}
	}

	if len(deleteErrors) > 0 {
		return gerror.New(strings.Join(deleteErrors, "; "))
	}
	return nil
}

func stageLocalDeleteTargets(targets []fileDeleteTarget) ([]stagedLocalDelete, error) {
	staged := make([]stagedLocalDelete, 0, len(targets))
	for _, target := range targets {
		if target.Storage != 1 {
			continue
		}
		item, err := stageLocalDeleteTarget(target)
		if err != nil {
			if restoreErr := restoreStagedLocalDeletes(staged); restoreErr != nil {
				return nil, gerror.Newf("准备删除本地文件失败: %v；恢复已暂存文件失败: %v", err, restoreErr)
			}
			return nil, err
		}
		if item.RecordID > 0 {
			staged = append(staged, item)
		}
	}
	return staged, nil
}

func stageLocalDeleteTarget(target fileDeleteTarget) (stagedLocalDelete, error) {
	if strings.TrimSpace(target.URL) == "" {
		return stagedLocalDelete{}, nil
	}
	localPath := shared.LocalStoragePhysicalPath(target.URL)
	if _, err := os.Stat(localPath); err != nil {
		if os.IsNotExist(err) {
			return stagedLocalDelete{}, nil
		}
		return stagedLocalDelete{}, fmt.Errorf("读取本地文件失败: %w", err)
	}

	stagedPath := buildLocalDeleteStagePath(target.ID, localPath)
	if err := os.MkdirAll(filepath.Dir(stagedPath), 0o755); err != nil {
		return stagedLocalDelete{}, fmt.Errorf("创建本地删除暂存目录失败: %w", err)
	}
	if err := os.Rename(localPath, stagedPath); err != nil {
		return stagedLocalDelete{}, fmt.Errorf("暂存本地待删文件失败: %w", err)
	}
	return stagedLocalDelete{
		RecordID:     target.ID,
		OriginalPath: localPath,
		StagedPath:   stagedPath,
	}, nil
}

func buildLocalDeleteStagePath(recordID int64, localPath string) string {
	return filepath.Join(
		shared.DefaultLocalStoragePath,
		".trash",
		strconv.FormatInt(time.Now().UnixNano(), 10),
		fmt.Sprintf("%d-%s", recordID, filepath.Base(localPath)),
	)
}

func restoreStagedLocalDeletes(staged []stagedLocalDelete) error {
	for i := len(staged) - 1; i >= 0; i-- {
		if err := restoreStagedLocalDelete(staged[i]); err != nil {
			return err
		}
	}
	return nil
}

func restoreStagedLocalDelete(staged stagedLocalDelete) error {
	if staged.OriginalPath == "" || staged.StagedPath == "" {
		return nil
	}
	if _, err := os.Stat(staged.StagedPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("读取本地暂存文件失败: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(staged.OriginalPath), 0o755); err != nil {
		return fmt.Errorf("恢复本地文件目录失败: %w", err)
	}
	if err := os.Rename(staged.StagedPath, staged.OriginalPath); err != nil {
		return fmt.Errorf("恢复本地文件失败: %w", err)
	}
	return nil
}

func finalizeStagedLocalDelete(staged stagedLocalDelete) error {
	if staged.StagedPath == "" {
		return nil
	}
	if err := os.Remove(staged.StagedPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("清理本地暂存文件失败: %w", err)
	}
	return nil
}

func (s *sFile) loadDeleteTargets(ctx context.Context, ids []snowflake.JsonInt64) ([]fileDeleteTarget, error) {
	ids = batchutil.CompactIDs(ids)
	if len(ids) == 0 {
		return nil, gerror.New("请选择要删除的文件")
	}
	dbIDs := batchutil.ToInt64s(ids)
	var rows []fileDeleteTarget
	if err := dao.UploadFile.Ctx(ctx).
		Fields(dao.UploadFile.Columns().Id, dao.UploadFile.Columns().Url, dao.UploadFile.Columns().Storage).
		WhereIn(dao.UploadFile.Columns().Id, dbIDs).
		Where(dao.UploadFile.Columns().DeletedAt, nil).
		Scan(&rows); err != nil {
		return nil, err
	}
	if len(rows) != len(dbIDs) {
		if len(dbIDs) == 1 {
			return nil, gerror.New("文件记录不存在或已删除")
		}
		return nil, gerror.New("包含不存在或已删除的文件")
	}
	return orderDeleteTargets(rows, dbIDs), nil
}

func orderDeleteTargets(rows []fileDeleteTarget, ids []int64) []fileDeleteTarget {
	if len(rows) == 0 || len(ids) == 0 {
		return nil
	}
	rowMap := make(map[int64]fileDeleteTarget, len(rows))
	for _, row := range rows {
		rowMap[row.ID] = row
	}
	ordered := make([]fileDeleteTarget, 0, len(ids))
	for _, id := range ids {
		if row, ok := rowMap[id]; ok {
			ordered = append(ordered, row)
		}
	}
	return ordered
}

// getStr 安全地从 map[string]interface{} 中取字符串值
func getStr(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	v, ok := m[key]
	if !ok {
		return ""
	}
	switch value := v.(type) {
	case nil:
		return ""
	case string:
		return value
	case []byte:
		return string(value)
	case fmt.Stringer:
		return value.String()
	default:
		return fmt.Sprintf("%v", value)
	}
}

// Detail 获取文件记录详情
func (s *sFile) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.FileDetailOutput, err error) {
	if id <= 0 {
		return nil, gerror.New("文件记录不存在或已删除")
	}
	out = &model.FileDetailOutput{}
	err = dao.UploadFile.Ctx(ctx).Where(dao.UploadFile.Columns().Id, id).Where(dao.UploadFile.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("文件记录不存在或已删除")
	}
	out.DirName = shared.LookupDirName(ctx, int64(out.DirID))
	return
}

// List 获取文件记录列表
func (s *sFile) List(ctx context.Context, in *model.FileListInput) (list []*model.FileListOutput, total int, err error) {
	if in == nil {
		in = &model.FileListInput{}
	}
	normalizeFileListInput(in)
	m := dao.UploadFile.Ctx(ctx).Where(dao.UploadFile.Columns().DeletedAt, nil)
	if in.Keyword != "" {
		keywordBuilder := m.Builder().
			WhereLike(dao.UploadFile.Columns().Name, "%"+in.Keyword+"%").
			WhereOrLike(dao.UploadFile.Columns().Url, "%"+in.Keyword+"%").
			WhereOrLike(dao.UploadFile.Columns().Ext, "%"+in.Keyword+"%").
			WhereOrLike(dao.UploadFile.Columns().Mime, "%"+in.Keyword+"%")
		m = m.Where(keywordBuilder)
	}
	if in.DirID > 0 {
		m = m.Where(dao.UploadFile.Columns().DirId, in.DirID)
	}
	if in.Name != "" {
		m = m.WhereLike(dao.UploadFile.Columns().Name, "%"+in.Name+"%")
	}
	if in.Storage > 0 {
		m = m.Where(dao.UploadFile.Columns().Storage, in.Storage)
	}
	if in.IsImage != nil {
		m = m.Where(dao.UploadFile.Columns().IsImage, *in.IsImage)
	}
	total, err = m.Count()
	if err != nil {
		return
	}
	in.PageNum, in.PageSize = pageutil.Normalize(in.PageNum, in.PageSize)
	err = m.Page(in.PageNum, in.PageSize).OrderDesc(dao.UploadFile.Columns().Id).Scan(&list)
	if err != nil {
		return
	}
	s.fillDirNames(ctx, list)
	return
}

func (s *sFile) fillDirNames(ctx context.Context, list []*model.FileListOutput) {
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

func normalizeFileCreateInput(in *model.FileCreateInput) {
	if in == nil {
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.URL = strings.TrimSpace(in.URL)
	in.Ext = strings.TrimSpace(in.Ext)
	in.Mime = strings.TrimSpace(in.Mime)
}

func normalizeFileUpdateInput(in *model.FileUpdateInput) {
	if in == nil {
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.URL = strings.TrimSpace(in.URL)
	in.Ext = strings.TrimSpace(in.Ext)
	in.Mime = strings.TrimSpace(in.Mime)
}

func normalizeFileListInput(in *model.FileListInput) {
	if in == nil {
		return
	}
	in.Keyword = strings.TrimSpace(in.Keyword)
	in.Name = strings.TrimSpace(in.Name)
}

func validateFileFields(name, fileURL string, storage, isImage int, size int64) error {
	if name == "" {
		return gerror.New("文件名称不能为空")
	}
	if fileURL == "" {
		return gerror.New("文件地址不能为空")
	}
	if err := fieldvalid.Enum("存储类型", storage, 1, 2, 3); err != nil {
		return err
	}
	if err := fieldvalid.Binary("是否图片", isImage); err != nil {
		return err
	}
	if err := fieldvalid.NonNegative64("文件大小", size); err != nil {
		return err
	}
	return nil
}

func (s *sFile) ensureDirExists(ctx context.Context, dirID snowflake.JsonInt64) error {
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

func (s *sFile) ensureFileExists(ctx context.Context, id snowflake.JsonInt64) error {
	if id <= 0 {
		return gerror.New("文件记录不存在或已删除")
	}
	count, err := dao.UploadFile.Ctx(ctx).
		Where(dao.UploadFile.Columns().Id, id).
		Where(dao.UploadFile.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return gerror.New("文件记录不存在或已删除")
	}
	return nil
}

func loadUploadConfigByURL(ctx context.Context, storage int, fileURL string) (*entity.UploadConfig, string, error) {
	var configs []*entity.UploadConfig
	if err := dao.UploadConfig.Ctx(ctx).
		Where(dao.UploadConfig.Columns().Storage, storage).
		OrderAsc(dao.UploadConfig.Columns().DeletedAt).
		OrderDesc(dao.UploadConfig.Columns().Id).
		Scan(&configs); err != nil {
		return nil, "", fmt.Errorf("读取上传配置失败: %w", err)
	}
	config, objectKey := matchUploadConfigByURL(configs, storage, fileURL)
	if config == nil {
		return nil, "", fmt.Errorf("未找到与文件地址匹配的上传配置: %s", fileURL)
	}
	return config, objectKey, nil
}

func matchUploadConfigByURL(configs []*entity.UploadConfig, storage int, fileURL string) (*entity.UploadConfig, string) {
	return shared.MatchUploadConfigByURL(configs, storage, fileURL)
}

func matchOSSObjectKey(fileURL string, config *entity.UploadConfig) (string, bool) {
	return shared.MatchOSSObjectKey(fileURL, config)
}

func matchOSSObjectKeyParsed(parsedURL *url.URL, config *entity.UploadConfig) (string, bool) {
	return shared.MatchOSSObjectKeyParsed(parsedURL, config)
}

func matchCOSObjectKey(fileURL string, config *entity.UploadConfig) (string, bool) {
	return shared.MatchCOSObjectKey(fileURL, config)
}

func matchCOSObjectKeyParsed(parsedURL *url.URL, config *entity.UploadConfig) (string, bool) {
	return shared.MatchCOSObjectKeyParsed(parsedURL, config)
}
