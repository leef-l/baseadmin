package file

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	cos "github.com/tencentyun/cos-go-sdk-v5"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/upload/internal/dao"
	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/app/upload/internal/model/entity"
	"gbaseadmin/app/upload/internal/service"
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

// Create 创建文件记录
func (s *sFile) Create(ctx context.Context, in *model.FileCreateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeFileCreateInput(in)
	if err := s.ensureDirExists(ctx, in.DirID); err != nil {
		return err
	}
	id := snowflake.Generate()
	_, err := dao.UploadFile.Ctx(ctx).Data(g.Map{
		dao.UploadFile.Columns().Id:        id,
		dao.UploadFile.Columns().DirId:     in.DirID,
		dao.UploadFile.Columns().Name:      in.Name,
		dao.UploadFile.Columns().Url:       in.URL,
		dao.UploadFile.Columns().Ext:       in.Ext,
		dao.UploadFile.Columns().Size:      in.Size,
		dao.UploadFile.Columns().Mime:      in.Mime,
		dao.UploadFile.Columns().Storage:   in.Storage,
		dao.UploadFile.Columns().IsImage:   in.IsImage,
		dao.UploadFile.Columns().CreatedAt: gtime.Now(),
		dao.UploadFile.Columns().UpdatedAt: gtime.Now(),
	}).Insert()
	return err
}

// Update 更新文件记录
func (s *sFile) Update(ctx context.Context, in *model.FileUpdateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeFileUpdateInput(in)
	if err := s.ensureDirExists(ctx, in.DirID); err != nil {
		return err
	}
	data := g.Map{
		dao.UploadFile.Columns().DirId:     in.DirID,
		dao.UploadFile.Columns().Name:      in.Name,
		dao.UploadFile.Columns().Url:       in.URL,
		dao.UploadFile.Columns().Ext:       in.Ext,
		dao.UploadFile.Columns().Size:      in.Size,
		dao.UploadFile.Columns().Mime:      in.Mime,
		dao.UploadFile.Columns().Storage:   in.Storage,
		dao.UploadFile.Columns().IsImage:   in.IsImage,
		dao.UploadFile.Columns().UpdatedAt: gtime.Now(),
	}
	_, err := dao.UploadFile.Ctx(ctx).Where(dao.UploadFile.Columns().Id, in.ID).Data(data).Update()
	return err
}

// Delete 删除文件记录并物理删除文件
func (s *sFile) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	// 先查询文件信息，用于物理删除
	var fileInfo struct {
		Url     string `orm:"url"`
		Storage int    `orm:"storage"`
	}
	err := dao.UploadFile.Ctx(ctx).Where(dao.UploadFile.Columns().Id, id).
		Where(dao.UploadFile.Columns().DeletedAt, nil).Scan(&fileInfo)
	if err != nil {
		return err
	}

	// 软删除记录
	_, err = dao.UploadFile.Ctx(ctx).Where(dao.UploadFile.Columns().Id, id).Data(g.Map{
		dao.UploadFile.Columns().DeletedAt: gtime.Now(),
	}).Update()
	if err != nil {
		return err
	}

	// 物理删除文件
	if fileInfo.Url != "" {
		switch fileInfo.Storage {
		case 1: // 本地存储: URL /upload/xxx -> 物理路径 resource/upload/xxx
			localPath := localStoragePhysicalPath(fileInfo.Url)
			_ = os.Remove(localPath)
		case 2: // 阿里云OSS
			if delErr := deleteCloudFileOSS(ctx, fileInfo.Url); delErr != nil {
				g.Log().Warningf(ctx, "OSS删除文件失败: url=%s, err=%v", fileInfo.Url, delErr)
			}
		case 3: // 腾讯云COS
			if delErr := deleteCloudFileCOS(ctx, fileInfo.Url); delErr != nil {
				g.Log().Warningf(ctx, "COS删除文件失败: url=%s, err=%v", fileInfo.Url, delErr)
			}
		}
	}
	return nil
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
	out = &model.FileDetailOutput{}
	err = dao.UploadFile.Ctx(ctx).Where(dao.UploadFile.Columns().Id, id).Where(dao.UploadFile.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	// 查询所属目录关联显示
	if out.DirID != 0 {
		val, err := g.DB().Ctx(ctx).Model("upload_dir").Where("id", out.DirID).Where("deleted_at", nil).Value("name")
		if err == nil {
			out.DirName = val.String()
		}
	}
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
	dirSet := make(map[int64]struct{})
	for _, item := range list {
		if item.DirID != 0 {
			dirSet[int64(item.DirID)] = struct{}{}
		}
	}
	if len(dirSet) == 0 {
		return
	}
	dirIDs := make([]int64, 0, len(dirSet))
	for id := range dirSet {
		dirIDs = append(dirIDs, id)
	}
	rows, err := g.DB().Ctx(ctx).Model("upload_dir").
		Fields("id", "name").
		Where("deleted_at", nil).
		WhereIn("id", dirIDs).
		All()
	if err != nil {
		return
	}
	dirMap := make(map[int64]string, len(rows))
	for _, row := range rows {
		dirMap[row["id"].Int64()] = row["name"].String()
	}
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

func loadUploadConfigByURL(ctx context.Context, storage int, fileURL string) (*entity.UploadConfig, string, error) {
	var configs []*entity.UploadConfig
	if err := dao.UploadConfig.Ctx(ctx).
		Where(dao.UploadConfig.Columns().Storage, storage).
		Where(dao.UploadConfig.Columns().DeletedAt, nil).
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
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return nil, ""
	}
	for _, config := range configs {
		if config == nil {
			continue
		}
		switch storage {
		case 2:
			if objectKey, ok := matchOSSObjectKeyParsed(parsedURL, config); ok {
				return config, objectKey
			}
		case 3:
			if objectKey, ok := matchCOSObjectKeyParsed(parsedURL, config); ok {
				return config, objectKey
			}
		}
	}
	return nil, ""
}

func matchOSSObjectKey(fileURL string, config *entity.UploadConfig) (string, bool) {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return "", false
	}
	return matchOSSObjectKeyParsed(parsedURL, config)
}

func matchOSSObjectKeyParsed(parsedURL *url.URL, config *entity.UploadConfig) (string, bool) {
	if config == nil || config.OssBucket == "" || config.OssEndpoint == "" {
		return "", false
	}
	if parsedURL == nil {
		return "", false
	}
	expectedHost := strings.ToLower(fmt.Sprintf("%s.%s", normalizeHostPart(config.OssBucket), normalizeHostPart(config.OssEndpoint)))
	if strings.ToLower(parsedURL.Hostname()) != expectedHost {
		return "", false
	}
	objectKey := objectKeyFromPath(parsedURL.Path)
	if objectKey == "" {
		return "", false
	}
	return objectKey, true
}

func matchCOSObjectKey(fileURL string, config *entity.UploadConfig) (string, bool) {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return "", false
	}
	return matchCOSObjectKeyParsed(parsedURL, config)
}

func matchCOSObjectKeyParsed(parsedURL *url.URL, config *entity.UploadConfig) (string, bool) {
	if config == nil || config.CosBucket == "" || config.CosRegion == "" {
		return "", false
	}
	if parsedURL == nil {
		return "", false
	}
	expectedHost := strings.ToLower(fmt.Sprintf("%s.cos.%s.myqcloud.com", normalizeHostPart(config.CosBucket), normalizeHostPart(config.CosRegion)))
	if strings.ToLower(parsedURL.Hostname()) != expectedHost {
		return "", false
	}
	objectKey := objectKeyFromPath(parsedURL.Path)
	if objectKey == "" {
		return "", false
	}
	return objectKey, true
}

func normalizeHostPart(value string) string {
	return strings.TrimSpace(value)
}

func objectKeyFromPath(path string) string {
	objectKey := strings.TrimPrefix(path, "/")
	if objectKey == "" {
		return ""
	}
	if decoded, err := url.PathUnescape(objectKey); err == nil && decoded != "" {
		return decoded
	}
	return objectKey
}
