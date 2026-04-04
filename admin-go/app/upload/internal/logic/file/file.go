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
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/upload/internal/dao"
	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/app/upload/internal/service"
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
			localPath := "resource" + fileInfo.Url
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
	var configRecord map[string]interface{}
	if err := dao.UploadConfig.Ctx(ctx).
		Where("is_default", 1).Where("status", 1).Where("deleted_at IS NULL").
		Scan(&configRecord); err != nil || configRecord == nil {
		return fmt.Errorf("读取OSS配置失败")
	}

	endpoint := getStr(configRecord, "oss_endpoint")
	bucket := getStr(configRecord, "oss_bucket")
	accessKey := getStr(configRecord, "oss_access_key")
	secretKey := getStr(configRecord, "oss_secret_key")

	if endpoint == "" || bucket == "" || accessKey == "" || secretKey == "" {
		return fmt.Errorf("阿里云OSS配置不完整")
	}

	// 从 URL 提取 objectKey：https://{bucket}.{endpoint}/{objectKey}
	prefix := fmt.Sprintf("https://%s.%s/", bucket, endpoint)
	objectKey := strings.TrimPrefix(fileURL, prefix)
	if objectKey == fileURL {
		return fmt.Errorf("无法从URL解析objectKey: %s", fileURL)
	}

	client, err := oss.New(endpoint, accessKey, secretKey)
	if err != nil {
		return fmt.Errorf("创建OSS客户端失败: %w", err)
	}
	b, err := client.Bucket(bucket)
	if err != nil {
		return fmt.Errorf("获取OSS Bucket失败: %w", err)
	}
	return b.DeleteObject(objectKey)
}

// deleteCloudFileCOS 从腾讯云 COS 删除文件
func deleteCloudFileCOS(ctx context.Context, fileURL string) error {
	var configRecord map[string]interface{}
	if err := dao.UploadConfig.Ctx(ctx).
		Where("is_default", 1).Where("status", 1).Where("deleted_at IS NULL").
		Scan(&configRecord); err != nil || configRecord == nil {
		return fmt.Errorf("读取COS配置失败")
	}

	region := getStr(configRecord, "cos_region")
	bucket := getStr(configRecord, "cos_bucket")
	secretId := getStr(configRecord, "cos_secret_id")
	secretKey := getStr(configRecord, "cos_secret_key")

	if region == "" || bucket == "" || secretId == "" || secretKey == "" {
		return fmt.Errorf("腾讯云COS配置不完整")
	}

	// 从 URL 提取 objectKey：https://{bucket}.cos.{region}.myqcloud.com/{objectKey}
	prefix := fmt.Sprintf("https://%s.cos.%s.myqcloud.com/", bucket, region)
	objectKey := strings.TrimPrefix(fileURL, prefix)
	if objectKey == fileURL {
		return fmt.Errorf("无法从URL解析objectKey: %s", fileURL)
	}

	bucketURL := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", bucket, region)
	u, err := url.Parse(bucketURL)
	if err != nil {
		return fmt.Errorf("解析COS URL失败: %w", err)
	}
	cosClient := cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretId,
			SecretKey: secretKey,
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
	s, _ := v.(string)
	return s
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
	m := dao.UploadFile.Ctx(ctx).Where(dao.UploadFile.Columns().DeletedAt, nil)
	if in.DirID > 0 {
		m = m.Where(dao.UploadFile.Columns().DirId, in.DirID)
	}
	if in.Name != "" {
		m = m.WhereLike(dao.UploadFile.Columns().Name, "%"+in.Name+"%")
	}
	if in.Storage > 0 {
		m = m.Where(dao.UploadFile.Columns().Storage, in.Storage)
	}
	if in.IsImage > 0 {
		m = m.Where(dao.UploadFile.Columns().IsImage, in.IsImage)
	}
	total, err = m.Count()
	if err != nil {
		return
	}
	err = m.Page(in.PageNum, in.PageSize).OrderDesc(dao.UploadFile.Columns().Id).Scan(&list)
	if err != nil {
		return
	}
	// 填充关联显示字段
	for _, item := range list {
		if item.DirID != 0 {
			val, err := g.DB().Ctx(ctx).Model("upload_dir").Where("id", item.DirID).Where("deleted_at", nil).Value("name")
			if err == nil {
				item.DirName = val.String()
			}
		}
	}
	return
}
