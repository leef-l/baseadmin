package user

import (
	"context"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"

	v1 "gbaseadmin/app/member/api/member/v1"
	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/app/member/internal/service"
)

func csvSafeUser(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == 0 || strings.ContainsAny(s[:1], "=+-@\t\r") {
		return "'" + s
	}
	return s
}

var User = cUser{}

type cUser struct{}

// Create 创建会员用户
func (c *cUser) Create(ctx context.Context, req *v1.UserCreateReq) (res *v1.UserCreateRes, err error) {
	err = service.User().Create(ctx, &model.UserCreateInput{
		ParentID: req.ParentID,
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Phone: req.Phone,
		Avatar: req.Avatar,
		RealName: req.RealName,
		LevelID: req.LevelID,
		LevelExpireAt: req.LevelExpireAt,
		TeamCount: req.TeamCount,
		DirectCount: req.DirectCount,
		ActiveCount: req.ActiveCount,
		TeamTurnover: req.TeamTurnover,
		IsActive: req.IsActive,
		IsQualified: req.IsQualified,
		InviteCode: req.InviteCode,
		RegisterIP: req.RegisterIP,
		LastLoginAt: req.LastLoginAt,
		Remark: req.Remark,
		Sort: req.Sort,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新会员用户
func (c *cUser) Update(ctx context.Context, req *v1.UserUpdateReq) (res *v1.UserUpdateRes, err error) {
	err = service.User().Update(ctx, &model.UserUpdateInput{
		ID: req.ID,
		ParentID: req.ParentID,
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Phone: req.Phone,
		Avatar: req.Avatar,
		RealName: req.RealName,
		LevelID: req.LevelID,
		LevelExpireAt: req.LevelExpireAt,
		TeamCount: req.TeamCount,
		DirectCount: req.DirectCount,
		ActiveCount: req.ActiveCount,
		TeamTurnover: req.TeamTurnover,
		IsActive: req.IsActive,
		IsQualified: req.IsQualified,
		InviteCode: req.InviteCode,
		RegisterIP: req.RegisterIP,
		LastLoginAt: req.LastLoginAt,
		Remark: req.Remark,
		Sort: req.Sort,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除会员用户
func (c *cUser) Delete(ctx context.Context, req *v1.UserDeleteReq) (res *v1.UserDeleteRes, err error) {
	err = service.User().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除会员用户
func (c *cUser) BatchDelete(ctx context.Context, req *v1.UserBatchDeleteReq) (res *v1.UserBatchDeleteRes, err error) {
	err = service.User().BatchDelete(ctx, req.IDs)
	return
}

// BatchUpdate 批量编辑会员用户
func (c *cUser) BatchUpdate(ctx context.Context, req *v1.UserBatchUpdateReq) (res *v1.UserBatchUpdateRes, err error) {
	err = service.User().BatchUpdate(ctx, &model.UserBatchUpdateInput{
		IDs: req.IDs,
		IsActive: req.IsActive,
		IsQualified: req.IsQualified,
		Status: req.Status,
	})
	return
}

// Detail 获取会员用户详情
func (c *cUser) Detail(ctx context.Context, req *v1.UserDetailReq) (res *v1.UserDetailRes, err error) {
	res = &v1.UserDetailRes{}
	res.UserDetailOutput, err = service.User().Detail(ctx, req.ID)
	return
}

// List 获取会员用户列表
func (c *cUser) List(ctx context.Context, req *v1.UserListReq) (res *v1.UserListRes, err error) {
	res = &v1.UserListRes{}
	res.List, res.Total, err = service.User().List(ctx, &model.UserListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Keyword: req.Keyword,
		Username: req.Username,
		InviteCode: req.InviteCode,
		Nickname: req.Nickname,
		RealName: req.RealName,
		Phone: req.Phone,
		ParentID: req.ParentID,
		LevelID: req.LevelID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		IsActive: req.IsActive,
		IsQualified: req.IsQualified,
		Status: req.Status,
		LevelExpireAtStart: req.LevelExpireAtStart,
		LevelExpireAtEnd: req.LevelExpireAtEnd,
		LastLoginAtStart: req.LastLoginAtStart,
		LastLoginAtEnd: req.LastLoginAtEnd,
	})
	return
}
// Export 导出会员用户
func (c *cUser) Export(ctx context.Context, req *v1.UserExportReq) (res *v1.UserExportRes, err error) {
	list, err := service.User().Export(ctx, &model.UserListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Keyword: req.Keyword,
		Username: req.Username,
		InviteCode: req.InviteCode,
		Nickname: req.Nickname,
		RealName: req.RealName,
		Phone: req.Phone,
		ParentID: req.ParentID,
		LevelID: req.LevelID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		IsActive: req.IsActive,
		IsQualified: req.IsQualified,
		Status: req.Status,
		LevelExpireAtStart: req.LevelExpireAtStart,
		LevelExpireAtEnd: req.LevelExpireAtEnd,
		LastLoginAtStart: req.LastLoginAtStart,
		LastLoginAtEnd: req.LastLoginAtEnd,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="user.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头（与导入模板列对齐，末尾追加只读列）
	_ = w.Write([]string{"上级会员", "用户名", "昵称", "手机号", "头像", "真实姓名", "当前等级", "团队总人数", "直推人数", "有效用户数", "团队总营业额", "是否激活", "仓库资格", "邀请码", "注册IP", "备注", "排序", "状态", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			csvSafeUser(item.UserUsername),
			csvSafeUser(item.Username),
			csvSafeUser(item.Nickname),
			csvSafeUser(item.Phone),
			csvSafeUser(item.Avatar),
			csvSafeUser(item.RealName),
			csvSafeUser(item.LevelName),
			fmt.Sprintf("%v", item.TeamCount),
			fmt.Sprintf("%v", item.DirectCount),
			fmt.Sprintf("%v", item.ActiveCount),
			fmt.Sprintf("%v", item.TeamTurnover),
			fmt.Sprintf("%v", item.IsActive),
			fmt.Sprintf("%v", item.IsQualified),
			csvSafeUser(item.InviteCode),
			csvSafeUser(item.RegisterIP),
			csvSafeUser(item.Remark),
			fmt.Sprintf("%v", item.Sort),
			fmt.Sprintf("%v", item.Status),
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Tree 获取会员用户树形结构
func (c *cUser) Tree(ctx context.Context, req *v1.UserTreeReq) (res *v1.UserTreeRes, err error) {
	res = &v1.UserTreeRes{}
	res.List, err = service.User().Tree(ctx, &model.UserTreeInput{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Keyword: req.Keyword,
		Username: req.Username,
		InviteCode: req.InviteCode,
		Nickname: req.Nickname,
		RealName: req.RealName,
		Phone: req.Phone,
		ParentID: req.ParentID,
		LevelID: req.LevelID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		IsActive: req.IsActive,
		IsQualified: req.IsQualified,
		Status: req.Status,
		LevelExpireAtStart: req.LevelExpireAtStart,
		LevelExpireAtEnd: req.LevelExpireAtEnd,
		LastLoginAtStart: req.LastLoginAtStart,
		LastLoginAtEnd: req.LastLoginAtEnd,
	})
	return
}

