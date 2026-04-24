package uploader

import (
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
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
	"gbaseadmin/utility/uploadticket"
)

func init() {
	service.RegisterUploader(&sUploader{})
}

type sUploader struct{}

type uploadContentInfo struct {
	Mime     string
	ImageExt string
}

type uploadPolicy struct {
	ConfigID            int64
	Source              string
	PathPrefix          string
	MaxSize             int64
	AllowedExts         map[string]struct{}
	AllowedMimePrefixes []string
}

const maxInt64AsUint64 = ^uint64(0) >> 1
const defaultLocalStoragePath = shared.DefaultLocalStoragePath

var imageExts = map[string]bool{
	"jpg": true, "jpeg": true, "png": true, "gif": true,
	"webp": true, "bmp": true, "heic": true, "heif": true,
}

var blockedUploadExts = map[string]bool{
	"svg":   true,
	"html":  true,
	"htm":   true,
	"xhtml": true,
	"shtml": true,
	"xml":   true,
	"js":    true,
	"mjs":   true,
	"cjs":   true,
	"css":   true,
	"php":   true,
	"phtml": true,
	"jsp":   true,
	"asp":   true,
	"aspx":  true,
}

func (s *sUploader) Upload(ctx context.Context) (*model.UploadOutput, error) {
	return s.upload(ctx, nil)
}

func (s *sUploader) UploadByTicket(ctx context.Context, ticket string) (*model.UploadOutput, error) {
	secret, err := uploadTicketSecret(ctx)
	if err != nil {
		return nil, err
	}
	claims, err := uploadticket.Verify(ticket, secret)
	if err != nil {
		return nil, err
	}

	policy := &uploadPolicy{
		ConfigID:            claims.ConfigID,
		Source:              strings.TrimSpace(claims.Scene),
		PathPrefix:          sanitizeUploadSubdir(claims.Dir),
		MaxSize:             claims.MaxSize,
		AllowedExts:         sliceToSet(claims.AllowedExts),
		AllowedMimePrefixes: normalizeMimePrefixes(claims.AllowedMimePrefixes),
	}
	return s.upload(ctx, policy)
}

func (s *sUploader) upload(ctx context.Context, policy *uploadPolicy) (*model.UploadOutput, error) {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return nil, fmt.Errorf("请求上下文无效")
	}

	// 获取上传文件
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请选择要上传的文件")
	}

	configID := r.Get("configId").Int64()
	if configID == 0 {
		configID = r.Get("configID").Int64()
	}
	if policy != nil && policy.ConfigID > 0 {
		configID = policy.ConfigID
	}

	cfg, err := loadUploadStorageConfig(ctx, configID)
	if err != nil {
		return nil, err
	}

	// 验证文件大小
	if policy != nil && policy.MaxSize > 0 && policy.MaxSize < cfg.MaxSize {
		cfg.MaxSize = policy.MaxSize
	}
	if file.Size > cfg.MaxSize {
		return nil, fmt.Errorf("文件大小超过限制（最大 %dMB）", cfg.MaxSize/1024/1024)
	}

	// 解析文件信息
	fileMeta, err := buildUploadFileMeta(file, policy)
	if err != nil {
		return nil, err
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
	uniqueName := buildUniqueName(now, randomSuffix(10000), fileMeta.Ext)
	source := strings.TrimSpace(r.Get("source").String())
	if source == "" {
		source = strings.TrimSpace(r.Get("scene").String())
	}
	if policy != nil && policy.Source != "" {
		source = policy.Source
	}
	if source == "" {
		source = strings.TrimSpace(r.Header.Get("X-Upload-Source"))
	}
	if source == "" {
		source = strings.TrimSpace(r.Header.Get("Referer"))
	}
	systemUserID := currentUploadSystemUserID(ctx)

	// 始终先保存到本地临时目录，云存储场景下上传后再清理
	resolvedDirID := dirId
	var (
		relativeDir string
		savePath    string
		keepName    bool
	)
	fallbackRelativeDir := ""
	if policy != nil {
		fallbackRelativeDir = policy.PathPrefix
	}
	resolvedDirID, relativeDir, savePath, keepName, err = resolveUploadSavePath(ctx, cfg, dirId, fileMeta.Ext, fileMeta.Mime, source, now, fallbackRelativeDir, systemUserID)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(savePath, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %v", err)
	}

	saveName := uniqueName
	if keepName {
		if preservedName := buildPreservedUploadName(file.Filename, fileMeta.Ext); preservedName != "" {
			saveName = preservedName
		}
	}
	saveName = ensureUniqueUploadName(savePath, saveName)
	file.Filename = saveName
	fullPath := filepath.Join(savePath, saveName)
	_, err = file.Save(savePath)
	if err != nil {
		return nil, fmt.Errorf("保存文件失败: %v", err)
	}

	objectKey := buildObjectKey(relativeDir, saveName)
	fileURL := buildLocalFileURL(cfg.LocalPath, relativeDir, saveName)
	storeResult, err := newStorageProvider(cfg).Store(ctx, storeRequest{
		BaseLocalPath: cfg.LocalPath,
		FileURL:       fileURL,
		RelativeDir:   relativeDir,
		LocalFilePath: fullPath,
		ObjectKey:     objectKey,
		UniqueName:    saveName,
	})
	if err != nil {
		return nil, err
	}

	// 生成ID并写入数据库
	id := snowflake.Generate()
	_, err = dao.UploadFile.Ctx(ctx).Data(do.UploadFile{
		Id:      id,
		DirId:   resolvedDirID,
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

func currentUploadSystemUserID(ctx context.Context) int64 {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return 0
	}
	return r.GetCtxVar("jwt_user_id").Int64()
}

type uploadFileMeta struct {
	Ext     string
	IsImage int
	Mime    string
}

func buildUploadFileMeta(file *ghttp.UploadFile, policy *uploadPolicy) (uploadFileMeta, error) {
	ext := normalizeExt(filepath.Ext(file.Filename))
	if blockedUploadExts[ext] {
		return uploadFileMeta{}, fmt.Errorf("不允许上传 .%s 文件", ext)
	}
	if policy != nil && len(policy.AllowedExts) > 0 {
		if _, ok := policy.AllowedExts[ext]; !ok {
			return uploadFileMeta{}, fmt.Errorf("当前场景不允许上传 .%s 文件", ext)
		}
	}

	contentInfo, err := inspectUploadContent(file)
	if err != nil {
		return uploadFileMeta{}, err
	}
	mimeType := contentInfo.Mime
	if contentInfo.ImageExt == "" && (ext == "heic" || ext == "heif") && mimeType == "application/octet-stream" {
		contentInfo.ImageExt = ext
		mimeType = imageMimeByExt(ext)
	}
	if isBlockedDetectedMime(mimeType) {
		return uploadFileMeta{}, fmt.Errorf("不允许上传 %s 类型文件", mimeType)
	}
	if imageExts[ext] && contentInfo.ImageExt != "" && contentInfo.ImageExt != ext {
		if !(ext == "jpg" && contentInfo.ImageExt == "jpeg") && !(ext == "jpeg" && contentInfo.ImageExt == "jpg") {
			return uploadFileMeta{}, fmt.Errorf("文件内容与扩展名 .%s 不匹配", ext)
		}
	}
	if imageExts[ext] && contentInfo.ImageExt == "" {
		return uploadFileMeta{}, fmt.Errorf("文件内容与扩展名 .%s 不匹配", ext)
	}
	if policy != nil && len(policy.AllowedMimePrefixes) > 0 {
		if !matchMimePrefixes(mimeType, policy.AllowedMimePrefixes) {
			return uploadFileMeta{}, fmt.Errorf("当前场景不允许上传 %s 类型文件", mimeType)
		}
	}

	isImage := 0
	if imageExts[ext] {
		isImage = 1
	}
	return uploadFileMeta{
		Ext:     ext,
		IsImage: isImage,
		Mime:    mimeType,
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

func buildPreservedUploadName(filename string, ext string) string {
	base := strings.TrimSpace(filename)
	if base == "" {
		return ""
	}
	base = filepath.Base(strings.ReplaceAll(base, "\\", "/"))
	base = strings.ReplaceAll(base, string(os.PathSeparator), "_")
	base = strings.ReplaceAll(base, "/", "_")
	base = strings.TrimSpace(base)
	if base == "" || base == "." || base == ".." {
		return ""
	}
	if ext != "" {
		expectedExt := "." + strings.TrimPrefix(strings.ToLower(ext), ".")
		if strings.ToLower(filepath.Ext(base)) != expectedExt {
			base = strings.TrimSuffix(base, filepath.Ext(base)) + expectedExt
		}
	}
	return base
}

func ensureUniqueUploadName(dir string, name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return name
	}
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	candidate := name
	for index := 1; ; index++ {
		if _, err := os.Stat(filepath.Join(dir, candidate)); os.IsNotExist(err) {
			return candidate
		}
		candidate = fmt.Sprintf("%s_%d%s", base, index, ext)
	}
}

func normalizeLocalStoragePath(path string) string { return shared.NormalizeLocalStoragePath(path) }

func buildLocalFileURL(baseLocalPath, relativeDir, fileName string) string {
	return shared.BuildLocalFileURLWithBase(baseLocalPath, relativeDir, fileName)
}

func localStoragePhysicalPath(fileURL string) string { return shared.LocalStoragePhysicalPath(fileURL) }

func uploadTicketSecret(ctx context.Context) (string, error) {
	secret := strings.TrimSpace(g.Cfg().MustGet(ctx, "uploadTicket.secret").String())
	if secret != "" {
		return secret, nil
	}
	return "", fmt.Errorf("未配置 uploadTicket.secret")
}

func sanitizeUploadSubdir(value string) string {
	value = strings.TrimSpace(strings.ReplaceAll(value, "\\", "/"))
	value = strings.Trim(value, "/")
	if value == "" {
		return ""
	}
	cleaned := path.Clean("/" + value)
	cleaned = strings.TrimPrefix(cleaned, "/")
	if cleaned == "." || cleaned == "" || strings.Contains(cleaned, "..") {
		return ""
	}
	return cleaned
}

func sliceToSet(values []string) map[string]struct{} {
	if len(values) == 0 {
		return nil
	}
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		value = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(value)), ".")
		if value == "" {
			continue
		}
		result[value] = struct{}{}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func normalizeMimePrefixes(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.ToLower(strings.TrimSpace(value))
		if value == "" {
			continue
		}
		result = append(result, value)
	}
	return result
}

func matchMimePrefixes(value string, prefixes []string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	for _, prefix := range prefixes {
		if strings.HasPrefix(value, prefix) {
			return true
		}
	}
	return false
}

func inspectUploadContent(file *ghttp.UploadFile) (*uploadContentInfo, error) {
	if file == nil {
		return nil, fmt.Errorf("请选择要上传的文件")
	}

	reader, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("读取上传文件失败: %v", err)
	}
	defer reader.Close()

	head := make([]byte, 512)
	n, err := reader.Read(head)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("读取上传文件失败: %v", err)
	}
	head = head[:n]

	mimeType := strings.ToLower(strings.TrimSpace(http.DetectContentType(head)))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	imageExt := sniffImageExt(head)
	if imageExt != "" && !strings.HasPrefix(mimeType, "image/") {
		mimeType = imageMimeByExt(imageExt)
	}
	return &uploadContentInfo{
		Mime:     mimeType,
		ImageExt: imageExt,
	}, nil
}

func sniffImageExt(head []byte) string {
	switch {
	case len(head) >= 3 && head[0] == 0xFF && head[1] == 0xD8 && head[2] == 0xFF:
		return "jpg"
	case len(head) >= 8 &&
		head[0] == 0x89 &&
		head[1] == 0x50 &&
		head[2] == 0x4E &&
		head[3] == 0x47 &&
		head[4] == 0x0D &&
		head[5] == 0x0A &&
		head[6] == 0x1A &&
		head[7] == 0x0A:
		return "png"
	case len(head) >= 6 && (string(head[:6]) == "GIF87a" || string(head[:6]) == "GIF89a"):
		return "gif"
	case len(head) >= 12 && string(head[:4]) == "RIFF" && string(head[8:12]) == "WEBP":
		return "webp"
	case len(head) >= 2 && string(head[:2]) == "BM":
		return "bmp"
	case len(head) >= 12 &&
		string(head[4:8]) == "ftyp" &&
		(string(head[8:12]) == "heic" ||
			string(head[8:12]) == "heix" ||
			string(head[8:12]) == "hevc" ||
			string(head[8:12]) == "hevx" ||
			string(head[8:12]) == "mif1" ||
			string(head[8:12]) == "msf1"):
		return "heic"
	default:
		return ""
	}
}

func imageMimeByExt(ext string) string {
	switch normalizeExt(ext) {
	case "jpg", "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "gif":
		return "image/gif"
	case "webp":
		return "image/webp"
	case "bmp":
		return "image/bmp"
	case "heic":
		return "image/heic"
	case "heif":
		return "image/heif"
	default:
		return "application/octet-stream"
	}
}

func isBlockedDetectedMime(value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	baseValue := value
	if index := strings.Index(baseValue, ";"); index >= 0 {
		baseValue = strings.TrimSpace(baseValue[:index])
	}
	switch baseValue {
	case "image/svg+xml", "text/html", "application/xhtml+xml", "text/xml", "application/xml":
		return true
	default:
		return false
	}
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
