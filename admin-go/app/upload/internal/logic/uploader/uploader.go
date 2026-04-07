package uploader

import (
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/upload/internal/dao"
	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/app/upload/internal/service"
	"gbaseadmin/utility/snowflake"
)

func init() {
	service.RegisterUploader(&sUploader{})
}

type sUploader struct{}

const maxInt64AsUint64 = ^uint64(0) >> 1
const defaultLocalStoragePath = "resource/upload"

var imageExts = map[string]bool{
	"jpg": true, "jpeg": true, "png": true, "gif": true,
	"webp": true, "svg": true, "bmp": true,
}

func (s *sUploader) Upload(ctx context.Context) (*model.UploadOutput, error) {
	r := g.RequestFromCtx(ctx)

	// 获取上传文件
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请选择要上传的文件")
	}

	// 读取默认上传配置
	maxSize := int64(10 * 1024 * 1024) // 默认10MB
	storageType := 1                   // 默认本地
	localPath := defaultLocalStoragePath

	var configRecord map[string]interface{}
	err := dao.UploadConfig.Ctx(ctx).
		Where("is_default", 1).
		Where("status", 1).
		Where(dao.UploadConfig.Columns().DeletedAt, nil).
		Scan(&configRecord)
	if err == nil && configRecord != nil {
		if v := getInt64(configRecord, "max_size"); v > 0 {
			maxSize = v * 1024 * 1024
		}
		if v := getInt64(configRecord, "storage"); v > 0 {
			storageType = int(v)
		}
		if v := getString(configRecord, "local_path"); v != "" {
			localPath = v
		}
	}
	localPath = normalizeLocalStoragePath(localPath)

	// 验证文件大小
	if file.Size > maxSize {
		return nil, fmt.Errorf("文件大小超过限制（最大 %dMB）", maxSize/1024/1024)
	}

	// 解析文件信息
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != "" {
		ext = ext[1:] // 去掉点号
	}
	isImage := 0
	if imageExts[ext] {
		isImage = 1
	}

	// 获取目录ID
	dirId := r.Get("dirId").Int64()
	if dirId > 0 {
		dirCount, err := dao.UploadDir.Ctx(ctx).
			Where(dao.UploadDir.Columns().Id, dirId).
			Where(dao.UploadDir.Columns().DeletedAt, nil).
			Count()
		if err != nil {
			return nil, err
		}
		if dirCount == 0 {
			return nil, fmt.Errorf("所选目录不存在或已删除")
		}
	}

	// 生成唯一文件名和对象路径
	now := time.Now()
	dateDir := now.Format("2006-01-02")
	uniqueName := buildUniqueName(now, randomSuffix(10000), ext)

	// 始终先保存到本地临时目录，云存储场景下上传后再清理
	savePath := filepath.Join(localPath, dateDir)
	if err := os.MkdirAll(savePath, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %v", err)
	}

	file.Filename = uniqueName
	fullPath := filepath.Join(savePath, uniqueName)
	_, err = file.Save(savePath)
	if err != nil {
		return nil, fmt.Errorf("保存文件失败: %v", err)
	}

	// objectKey：云存储对象路径，本地存储时复用为相对路径
	objectKey := dateDir + "/" + uniqueName
	var fileURL string

	switch storageType {
	case 2: // 阿里云OSS
		cfg := ossConfig{
			Endpoint:  getString(configRecord, "oss_endpoint"),
			Bucket:    getString(configRecord, "oss_bucket"),
			AccessKey: getString(configRecord, "oss_access_key"),
			SecretKey: getString(configRecord, "oss_secret_key"),
		}
		fileURL, err = uploadToOSS(cfg, fullPath, objectKey)
		if err != nil {
			_ = os.Remove(fullPath)
			return nil, fmt.Errorf("上传至OSS失败: %v", err)
		}
		// 上传成功后删除本地临时文件
		_ = os.Remove(fullPath)

	case 3: // 腾讯云COS
		cfg := cosConfig{
			Region:    getString(configRecord, "cos_region"),
			Bucket:    getString(configRecord, "cos_bucket"),
			SecretId:  getString(configRecord, "cos_secret_id"),
			SecretKey: getString(configRecord, "cos_secret_key"),
		}
		fileURL, err = uploadToCOS(cfg, fullPath, objectKey)
		if err != nil {
			_ = os.Remove(fullPath)
			return nil, fmt.Errorf("上传至COS失败: %v", err)
		}
		// 上传成功后删除本地临时文件
		_ = os.Remove(fullPath)

	default: // case 1: 本地存储
		fileURL = buildLocalFileURL(dateDir, uniqueName)
	}

	// 生成ID并写入数据库
	id := snowflake.Generate()
	_, err = dao.UploadFile.Ctx(ctx).Data(g.Map{
		"id":         id,
		"dir_id":     dirId,
		"name":       file.Filename,
		"url":        fileURL,
		"ext":        ext,
		"size":       file.Size,
		"mime":       file.FileHeader.Header.Get("Content-Type"),
		"storage":    storageType,
		"is_image":   isImage,
		"created_at": gtime.Now(),
		"updated_at": gtime.Now(),
	}).Insert()
	if err != nil {
		// 本地存储时回滚物理文件
		if storageType == 1 {
			_ = os.Remove(fullPath)
		}
		if storageType == 2 {
			cfg := ossConfig{
				Endpoint:  getString(configRecord, "oss_endpoint"),
				Bucket:    getString(configRecord, "oss_bucket"),
				AccessKey: getString(configRecord, "oss_access_key"),
				SecretKey: getString(configRecord, "oss_secret_key"),
			}
			if delErr := deleteFromOSS(cfg, objectKey); delErr != nil {
				g.Log().Warningf(ctx, "回滚OSS文件失败: objectKey=%s, err=%v", objectKey, delErr)
			}
		}
		if storageType == 3 {
			cfg := cosConfig{
				Region:    getString(configRecord, "cos_region"),
				Bucket:    getString(configRecord, "cos_bucket"),
				SecretId:  getString(configRecord, "cos_secret_id"),
				SecretKey: getString(configRecord, "cos_secret_key"),
			}
			if delErr := deleteFromCOS(cfg, objectKey); delErr != nil {
				g.Log().Warningf(ctx, "回滚COS文件失败: objectKey=%s, err=%v", objectKey, delErr)
			}
		}
		return nil, fmt.Errorf("保存文件记录失败: %v", err)
	}

	return &model.UploadOutput{
		ID:      snowflake.JsonInt64(id),
		URL:     fileURL,
		Name:    file.Filename,
		Size:    file.Size,
		Ext:     ext,
		Mime:    file.FileHeader.Header.Get("Content-Type"),
		IsImage: isImage,
	}, nil
}

// getString 安全地从 map[string]interface{} 中取字符串值
func getString(m map[string]interface{}, key string) string {
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

func getInt64(m map[string]interface{}, key string) int64 {
	if m == nil {
		return 0
	}
	v, ok := m[key]
	if !ok {
		return 0
	}
	switch value := v.(type) {
	case int:
		return int64(value)
	case int8:
		return int64(value)
	case int16:
		return int64(value)
	case int32:
		return int64(value)
	case int64:
		return value
	case uint:
		return int64(value)
	case uint8:
		return int64(value)
	case uint16:
		return int64(value)
	case uint32:
		return int64(value)
	case uint64:
		if value > maxInt64AsUint64 {
			return 0
		}
		return int64(value)
	case float32:
		return int64(value)
	case float64:
		return int64(value)
	case json.Number:
		n, _ := value.Int64()
		return n
	case string:
		n, _ := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
		return n
	case []byte:
		n, _ := strconv.ParseInt(strings.TrimSpace(string(value)), 10, 64)
		return n
	default:
		return 0
	}
}

func buildUniqueName(now time.Time, randSuffix int, ext string) string {
	ext = normalizeExt(ext)
	base := fmt.Sprintf("%d%d%04d", now.UnixMilli(), now.UnixNano()%1000, randSuffix)
	if ext == "" {
		return base
	}
	return base + "." + ext
}

func normalizeExt(ext string) string {
	ext = strings.TrimSpace(ext)
	ext = strings.TrimPrefix(ext, ".")
	return strings.ToLower(ext)
}

func normalizeLocalStoragePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return defaultLocalStoragePath
	}
	cleaned := filepath.Clean(path)
	if cleaned == "." {
		return defaultLocalStoragePath
	}
	return cleaned
}

func buildLocalFileURL(parts ...string) string {
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.Trim(strings.TrimSpace(part), `/\`)
		if part != "" {
			filtered = append(filtered, part)
		}
	}
	if len(filtered) == 0 {
		return "/upload"
	}
	return "/upload/" + strings.Join(filtered, "/")
}

func localStoragePhysicalPath(fileURL string) string {
	fileURL = strings.TrimSpace(fileURL)
	if fileURL == "" {
		return defaultLocalStoragePath
	}
	trimmed := strings.TrimPrefix(fileURL, "/upload")
	trimmed = strings.TrimPrefix(trimmed, "/")
	if trimmed == "" {
		return defaultLocalStoragePath
	}
	parts := strings.Split(trimmed, "/")
	all := make([]string, 0, len(parts)+1)
	all = append(all, defaultLocalStoragePath)
	for _, part := range parts {
		if part = strings.TrimSpace(part); part != "" {
			all = append(all, part)
		}
	}
	return filepath.Join(all...)
}

func randomSuffix(max int) int {
	if max <= 1 {
		return 0
	}
	var buf [8]byte
	if _, err := crand.Read(buf[:]); err == nil {
		return int(binary.BigEndian.Uint64(buf[:]) % uint64(max))
	}
	return int(time.Now().UnixNano() % int64(max))
}
