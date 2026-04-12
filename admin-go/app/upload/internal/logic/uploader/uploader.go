package uploader

import (
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"gbaseadmin/app/upload/internal/dao"
	"gbaseadmin/app/upload/internal/logic/shared"
	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/app/upload/internal/model/do"
	"gbaseadmin/app/upload/internal/service"
	"gbaseadmin/utility/snowflake"
)

func init() {
	service.RegisterUploader(&sUploader{})
}

type sUploader struct{}

const maxInt64AsUint64 = ^uint64(0) >> 1
const defaultLocalStoragePath = shared.DefaultLocalStoragePath

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
	cfg, err := loadUploadStorageConfig(ctx)
	if err != nil {
		return nil, err
	}

	// 验证文件大小
	if file.Size > cfg.MaxSize {
		return nil, fmt.Errorf("文件大小超过限制（最大 %dMB）", cfg.MaxSize/1024/1024)
	}

	// 解析文件信息
	fileMeta := buildUploadFileMeta(file)

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
	uniqueName := buildUniqueName(now, randomSuffix(10000), fileMeta.Ext)

	// 始终先保存到本地临时目录，云存储场景下上传后再清理
	relativeDir, savePath, err := resolveUploadSavePath(ctx, cfg, dirId, fileMeta.Ext, now)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(savePath, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %v", err)
	}

	file.Filename = uniqueName
	fullPath := filepath.Join(savePath, uniqueName)
	_, err = file.Save(savePath)
	if err != nil {
		return nil, fmt.Errorf("保存文件失败: %v", err)
	}

	objectKey := buildObjectKey(relativeDir, uniqueName)
	storeResult, err := newStorageProvider(cfg).Store(ctx, storeRequest{
		RelativeDir:   relativeDir,
		LocalFilePath: fullPath,
		ObjectKey:     objectKey,
		UniqueName:    uniqueName,
	})
	if err != nil {
		return nil, err
	}

	// 生成ID并写入数据库
	id := snowflake.Generate()
	_, err = dao.UploadFile.Ctx(ctx).Data(do.UploadFile{
		Id:      id,
		DirId:   dirId,
		Name:    file.Filename,
		Url:     storeResult.FileURL,
		Ext:     fileMeta.Ext,
		Size:    file.Size,
		Mime:    fileMeta.Mime,
		Storage: cfg.StorageType,
		IsImage: fileMeta.IsImage,
	}).Insert()
	if err != nil {
		runCleanupHook(ctx, "rollback", storeResult.OnRollback)
		return nil, fmt.Errorf("保存文件记录失败: %v", err)
	}
	runCleanupHook(ctx, "commit", storeResult.OnCommit)

	return &model.UploadOutput{
		ID:      snowflake.JsonInt64(id),
		URL:     storeResult.FileURL,
		Name:    file.Filename,
		Size:    file.Size,
		Ext:     fileMeta.Ext,
		Mime:    fileMeta.Mime,
		IsImage: fileMeta.IsImage,
	}, nil
}

type uploadFileMeta struct {
	Ext     string
	IsImage int
	Mime    string
}

func buildUploadFileMeta(file *ghttp.UploadFile) uploadFileMeta {
	ext := normalizeExt(filepath.Ext(file.Filename))
	mimeType := strings.TrimSpace(file.FileHeader.Header.Get("Content-Type"))
	if mimeType == "" {
		mimeType = mime.TypeByExtension("." + ext)
	}
	isImage := 0
	if imageExts[ext] || strings.HasPrefix(mimeType, "image/") {
		isImage = 1
	}
	return uploadFileMeta{
		Ext:     ext,
		IsImage: isImage,
		Mime:    mimeType,
	}
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

func normalizeLocalStoragePath(path string) string { return shared.NormalizeLocalStoragePath(path) }

func buildLocalFileURL(parts ...string) string { return shared.BuildLocalFileURL(parts...) }

func localStoragePhysicalPath(fileURL string) string { return shared.LocalStoragePhysicalPath(fileURL) }

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
