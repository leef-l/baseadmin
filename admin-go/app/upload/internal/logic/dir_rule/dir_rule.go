package dir_rule

import (
	"context"
	"path"
	"strings"
	"unicode"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"

	"gbaseadmin/app/upload/internal/dao"
	"gbaseadmin/app/upload/internal/logic/shared"
	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/app/upload/internal/service"
	"gbaseadmin/utility/batchutil"
	"gbaseadmin/utility/fieldvalid"
	"gbaseadmin/utility/inpututil"
	"gbaseadmin/utility/pageutil"
	"gbaseadmin/utility/snowflake"
)

func init() {
	service.RegisterDirRule(New())
}

func New() *sDirRule {
	return &sDirRule{}
}

type sDirRule struct{}

type uploadDirRuleSaveData struct {
	DirId        snowflake.JsonInt64 `orm:"dir_id"`
	Category     int                 `orm:"category"`
	FileType     string              `orm:"file_type"`
	StorageTypes string              `orm:"storage_types"`
	SavePath     string              `orm:"save_path"`
	KeepName     int                 `orm:"keep_name"`
	Status       int                 `orm:"status"`
}

type uploadDirRuleCreateData struct {
	Id           snowflake.JsonInt64 `orm:"id"`
	DirId        snowflake.JsonInt64 `orm:"dir_id"`
	Category     int                 `orm:"category"`
	FileType     string              `orm:"file_type"`
	StorageTypes string              `orm:"storage_types"`
	SavePath     string              `orm:"save_path"`
	KeepName     int                 `orm:"keep_name"`
	Status       int                 `orm:"status"`
}

// Create 创建文件目录规则
func (s *sDirRule) Create(ctx context.Context, in *model.DirRuleCreateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeDirRuleCreateInput(in)
	if err := validateDirRuleFields(in.DirID, in.Category, in.Status, in.FileType, in.StorageTypes, in.SavePath, in.KeepName); err != nil {
		return err
	}
	if err := s.ensureDirExists(ctx, in.DirID); err != nil {
		return err
	}
	id := snowflake.Generate()
	_, err := dao.UploadDirRule.Ctx(ctx).Data(uploadDirRuleCreateData{
		Id:           id,
		DirId:        in.DirID,
		Category:     in.Category,
		FileType:     in.FileType,
		StorageTypes: in.StorageTypes,
		SavePath:     in.SavePath,
		KeepName:     in.KeepName,
		Status:       in.Status,
	}).Insert()
	return err
}

// Update 更新文件目录规则
func (s *sDirRule) Update(ctx context.Context, in *model.DirRuleUpdateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeDirRuleUpdateInput(in)
	if err := validateDirRuleFields(in.DirID, in.Category, in.Status, in.FileType, in.StorageTypes, in.SavePath, in.KeepName); err != nil {
		return err
	}
	if err := s.ensureDirRuleExists(ctx, in.ID); err != nil {
		return err
	}
	if err := s.ensureDirExists(ctx, in.DirID); err != nil {
		return err
	}
	data := uploadDirRuleSaveData{
		DirId:        in.DirID,
		Category:     in.Category,
		FileType:     in.FileType,
		StorageTypes: in.StorageTypes,
		SavePath:     in.SavePath,
		KeepName:     in.KeepName,
		Status:       in.Status,
	}
	_, err := dao.UploadDirRule.Ctx(ctx).
		Where(dao.UploadDirRule.Columns().Id, in.ID).
		Where(dao.UploadDirRule.Columns().DeletedAt, nil).
		Data(data).
		Update()
	return err
}

// Delete 软删除文件目录规则
func (s *sDirRule) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := s.ensureDirRuleExists(ctx, id); err != nil {
		return err
	}
	_, err := dao.UploadDirRule.Ctx(ctx).
		Where(dao.UploadDirRule.Columns().Id, id).
		Delete()
	return err
}

// BatchDelete 批量删除文件目录规则
func (s *sDirRule) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	ids = batchutil.CompactIDs(ids)
	if len(ids) == 0 {
		return gerror.New("请选择要删除的目录规则")
	}
	deleteIDs := batchutil.ToInt64s(ids)
	count, err := dao.UploadDirRule.Ctx(ctx).
		WhereIn(dao.UploadDirRule.Columns().Id, deleteIDs).
		Where(dao.UploadDirRule.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if count != len(deleteIDs) {
		return gerror.New("包含不存在或已删除的目录规则")
	}
	return dao.UploadDirRule.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Model(dao.UploadDirRule.Table()).Ctx(ctx).
			WhereIn(dao.UploadDirRule.Columns().Id, deleteIDs).
			Delete()
		return err
	})
}

// Detail 获取文件目录规则详情
func (s *sDirRule) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.DirRuleDetailOutput, err error) {
	if id <= 0 {
		return nil, gerror.New("目录规则不存在或已删除")
	}
	out = &model.DirRuleDetailOutput{}
	err = dao.UploadDirRule.Ctx(ctx).Where(dao.UploadDirRule.Columns().Id, id).Where(dao.UploadDirRule.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("目录规则不存在或已删除")
	}
	out.DirName = shared.LookupDirName(ctx, int64(out.DirID))
	return
}

// List 获取文件目录规则列表
func (s *sDirRule) List(ctx context.Context, in *model.DirRuleListInput) (list []*model.DirRuleListOutput, total int, err error) {
	if in == nil {
		in = &model.DirRuleListInput{}
	}
	normalizeDirRuleListInput(in)
	m := s.buildListModel(ctx, in, true)
	total, err = m.Count()
	if err != nil && shouldRetryDirRuleListWithoutStorageTypes(err, in.Keyword) {
		m = s.buildListModel(ctx, in, false)
		total, err = m.Count()
	}
	if err != nil {
		return
	}
	in.PageNum, in.PageSize = pageutil.Normalize(in.PageNum, in.PageSize)
	err = m.Page(in.PageNum, in.PageSize).OrderAsc(dao.UploadDirRule.Columns().Id).Scan(&list)
	if err != nil {
		return
	}
	s.fillDirNames(ctx, list)
	return
}

func (s *sDirRule) buildListModel(ctx context.Context, in *model.DirRuleListInput, includeStorageTypes bool) *gdb.Model {
	m := dao.UploadDirRule.Ctx(ctx).Where(dao.UploadDirRule.Columns().DeletedAt, nil)
	if in.Keyword != "" {
		keywordBuilder := m.Builder()
		keywordBuilder = keywordBuilder.WhereLike(dao.UploadDirRule.Columns().SavePath, "%"+in.Keyword+"%")
		keywordBuilder = keywordBuilder.WhereOrLike(dao.UploadDirRule.Columns().FileType, "%"+in.Keyword+"%")
		if includeStorageTypes {
			keywordBuilder = keywordBuilder.WhereOrLike(dao.UploadDirRule.Columns().StorageTypes, "%"+in.Keyword+"%")
		}
		m = m.Where(keywordBuilder)
	}
	if in.Category > 0 {
		m = m.Where(dao.UploadDirRule.Columns().Category, in.Category)
	}
	if in.Status != nil {
		m = m.Where(dao.UploadDirRule.Columns().Status, *in.Status)
	}
	return m
}

func (s *sDirRule) fillDirNames(ctx context.Context, list []*model.DirRuleListOutput) {
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

func (s *sDirRule) ensureDirRuleExists(ctx context.Context, id snowflake.JsonInt64) error {
	if id <= 0 {
		return gerror.New("目录规则不存在或已删除")
	}
	count, err := dao.UploadDirRule.Ctx(ctx).
		Where(dao.UploadDirRule.Columns().Id, id).
		Where(dao.UploadDirRule.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return gerror.New("目录规则不存在或已删除")
	}
	return nil
}

func (s *sDirRule) ensureDirExists(ctx context.Context, dirID snowflake.JsonInt64) error {
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

func normalizeDirRuleListInput(in *model.DirRuleListInput) {
	if in == nil {
		return
	}
	in.Keyword = strings.TrimSpace(in.Keyword)
}

func normalizeDirRuleCreateInput(in *model.DirRuleCreateInput) {
	if in == nil {
		return
	}
	in.FileType = normalizeDirRuleMatchValue(in.Category, in.FileType)
	in.StorageTypes = normalizeDirRuleStorageTypes(in.StorageTypes)
	in.SavePath = strings.TrimSpace(in.SavePath)
}

func normalizeDirRuleUpdateInput(in *model.DirRuleUpdateInput) {
	if in == nil {
		return
	}
	in.FileType = normalizeDirRuleMatchValue(in.Category, in.FileType)
	in.StorageTypes = normalizeDirRuleStorageTypes(in.StorageTypes)
	in.SavePath = strings.TrimSpace(in.SavePath)
}

func validateDirRuleFields(dirID snowflake.JsonInt64, category, status int, fileType, storageTypes, savePath string, keepName int) error {
	if dirID <= 0 {
		return gerror.New("目录ID不能为空")
	}
	if err := fieldvalid.Enum("类别", category, 1, 2, 3); err != nil {
		return err
	}
	if storageTypes == "" {
		return gerror.New("适用存储不能为空")
	}
	if len(storageTypes) > 20 {
		return gerror.New("适用存储长度不能超过20个字符")
	}
	for _, item := range splitDirRuleStorageTypes(storageTypes) {
		if err := fieldvalid.Enum("适用存储", toInt(item), 1, 2, 3); err != nil {
			return gerror.New("适用存储值不合法")
		}
	}
	switch category {
	case 2:
		if fileType == "" {
			return gerror.New("文件类型不能为空")
		}
		if len(fileType) > 1000 {
			return gerror.New("文件类型长度不能超过1000个字符")
		}
		for _, item := range strings.Split(fileType, ",") {
			if !isValidDirRuleFileType(item) {
				return gerror.New("文件类型格式不正确")
			}
		}
	case 3:
		if fileType == "" {
			return gerror.New("来源标识不能为空")
		}
		if len(fileType) > 1000 {
			return gerror.New("来源标识长度不能超过1000个字符")
		}
	}
	if hasParentRelativePath(savePath) && !isLocalOnlyStorageTypes(storageTypes) {
		return gerror.New("父级目录规则仅支持本地存储")
	}
	if err := fieldvalid.Binary("保留原文件名", keepName); err != nil {
		return err
	}
	if err := fieldvalid.Binary("状态", status); err != nil {
		return err
	}
	return nil
}

func normalizeDirRuleMatchValue(category int, value string) string {
	switch category {
	case 2:
		return normalizeDirRuleFileType(value)
	case 3:
		return normalizeDirRuleSourceMatchers(value)
	default:
		return ""
	}
}

func normalizeDirRuleFileType(value string) string {
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == '，' || r == ';' || r == '；' || unicode.IsSpace(r)
	})
	normalized := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		part = normalizeDirRuleFileTypeToken(part)
		if part == "" {
			continue
		}
		if _, ok := seen[part]; ok {
			continue
		}
		seen[part] = struct{}{}
		normalized = append(normalized, part)
	}
	return strings.Join(normalized, ",")
}

func normalizeDirRuleSourceMatchers(value string) string {
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == '，' || r == ';' || r == '；' || unicode.IsSpace(r)
	})
	normalized := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		part = normalizeDirRuleSourceMatcher(part)
		if part == "" {
			continue
		}
		if _, ok := seen[part]; ok {
			continue
		}
		seen[part] = struct{}{}
		normalized = append(normalized, part)
	}
	return strings.Join(normalized, ",")
}

func normalizeDirRuleSourceMatcher(value string) string {
	value = strings.TrimSpace(value)
	hasWildcard := strings.HasSuffix(value, "/*")
	if hasWildcard {
		value = strings.TrimSuffix(value, "/*")
	}
	value = shared.NormalizeUploadRuleSource(value)
	if value == "" {
		return ""
	}
	if hasWildcard && value != "/" {
		return value + "/*"
	}
	return value
}

func shouldRetryDirRuleListWithoutStorageTypes(err error, keyword string) bool {
	if err == nil || strings.TrimSpace(keyword) == "" {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "unknown column") && strings.Contains(message, "storage_types")
}

func isValidDirRuleFileType(value string) bool {
	value = normalizeDirRuleFileTypeToken(value)
	if value == "" {
		return false
	}
	if strings.Contains(value, "*") {
		return strings.Count(value, "*") == 1 && strings.HasSuffix(value, "/*") && len(strings.TrimSuffix(value, "/*")) > 0
	}
	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' || r == '+' || r == '.' || r == '/' {
			continue
		}
		return false
	}
	return true
}

func normalizeDirRuleFileTypeToken(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return ""
	}
	if !strings.Contains(value, "/") {
		value = strings.TrimPrefix(value, ".")
	}
	return value
}

func normalizeDirRuleStorageTypes(value string) string {
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == '，' || r == ';' || r == '；' || unicode.IsSpace(r)
	})
	normalized := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if _, ok := seen[part]; ok {
			continue
		}
		seen[part] = struct{}{}
		normalized = append(normalized, part)
	}
	if len(normalized) == 0 {
		return "1,2,3"
	}
	return strings.Join(normalized, ",")
}

func isLocalOnlyStorageTypes(value string) bool {
	items := splitDirRuleStorageTypes(value)
	return len(items) == 1 && items[0] == "1"
}

func splitDirRuleStorageTypes(value string) []string {
	value = normalizeDirRuleStorageTypes(value)
	if value == "" {
		return nil
	}
	return strings.Split(value, ",")
}

func hasParentRelativePath(savePath string) bool {
	savePath = strings.TrimSpace(strings.ReplaceAll(savePath, `\`, "/"))
	if savePath == "" {
		return false
	}
	savePath = normalizeDirRuleSavePathAlias(savePath)
	cleaned := path.Clean(savePath)
	return cleaned == ".." || strings.HasPrefix(cleaned, "../")
}

func normalizeDirRuleSavePathAlias(value string) string {
	value = strings.TrimSpace(strings.ReplaceAll(value, `\`, "/"))
	switch {
	case value == "@up":
		return ".."
	case strings.HasPrefix(value, "@up/"):
		return "../" + strings.TrimPrefix(value, "@up/")
	default:
		return value
	}
}

func toInt(value string) int {
	switch strings.TrimSpace(value) {
	case "1":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	default:
		return 0
	}
}
