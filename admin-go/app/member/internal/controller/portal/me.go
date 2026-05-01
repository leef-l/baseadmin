package portal

import (
	"context"

	v1 "gbaseadmin/app/member/api/portal/v1"
	"gbaseadmin/app/member/internal/logic/portal"
	"gbaseadmin/app/member/internal/middleware"
)

// Me 控制器（C 端用户中心：个人信息 / 钱包 / 订单 / 团队聚合）。
var Me = cMe{}

type cMe struct{}

// Profile 获取当前会员资料。
func (c *cMe) Profile(ctx context.Context, req *v1.MeProfileReq) (res *v1.MeProfileRes, err error) {
	memberID := int64(middleware.CurrentMemberID(ctx))
	out, err := portal.AuthLogic().GetMyProfile(ctx, memberID)
	if err != nil {
		return nil, err
	}
	return &v1.MeProfileRes{
		MemberID:           out.MemberID,
		Phone:              out.Phone,
		Username:           out.Username,
		Nickname:           out.Nickname,
		Avatar:             out.Avatar,
		RealName:           out.RealName,
		InviteCode:         out.InviteCode,
		ParentID:           out.ParentID,
		LevelID:            out.LevelID,
		LevelName:          out.LevelName,
		LevelExpireAt:      out.LevelExpireAt,
		IsActive:           out.IsActive,
		IsQualified:        out.IsQualified,
		TeamCount:          out.TeamCount,
		DirectCount:        out.DirectCount,
		ActiveCount:        out.ActiveCount,
		TeamTurnover:       out.TeamTurnover,
		InviteURL:          out.InviteURL,
		DailyPurchaseLimit: out.DailyPurchaseLimit,
		TodayPurchaseCount: out.TodayPurchaseCount,
		TotalPurchaseCount: out.TotalPurchaseCount,
	}, nil
}

// Update 修改个人资料。
func (c *cMe) Update(ctx context.Context, req *v1.MeUpdateReq) (res *v1.MeUpdateRes, err error) {
	memberID := int64(middleware.CurrentMemberID(ctx))
	if err = portal.AuthLogic().UpdateMyProfile(ctx, &portal.UpdateMyProfileInput{
		UserID:   memberID,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		RealName: req.RealName,
	}); err != nil {
		return nil, err
	}
	return &v1.MeUpdateRes{}, nil
}

// ChangePassword 修改密码。
func (c *cMe) ChangePassword(ctx context.Context, req *v1.MeChangePasswordReq) (res *v1.MeChangePasswordRes, err error) {
	memberID := int64(middleware.CurrentMemberID(ctx))
	if err = portal.AuthLogic().ChangeMyPassword(ctx, &portal.ChangeMyPasswordInput{
		UserID:      memberID,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}); err != nil {
		return nil, err
	}
	return &v1.MeChangePasswordRes{}, nil
}

// ChangePhone 修改手机号。
func (c *cMe) ChangePhone(ctx context.Context, req *v1.MeChangePhoneReq) (res *v1.MeChangePhoneRes, err error) {
	memberID := int64(middleware.CurrentMemberID(ctx))
	if err = portal.AuthLogic().ChangeMyPhone(ctx, &portal.ChangeMyPhoneInput{
		UserID:   memberID,
		NewPhone: req.NewPhone,
		SmsCode:  req.SmsCode,
	}); err != nil {
		return nil, err
	}
	return &v1.MeChangePhoneRes{}, nil
}

// Wallets 获取三钱包余额。
func (c *cMe) Wallets(ctx context.Context, req *v1.MeWalletsReq) (res *v1.MeWalletsRes, err error) {
	memberID := int64(middleware.CurrentMemberID(ctx))
	out, err := portal.AuthLogic().GetMyWallets(ctx, memberID)
	if err != nil {
		return nil, err
	}
	return &v1.MeWalletsRes{
		Coupon:  toAPIWallet(out.Coupon),
		Reward:  toAPIWallet(out.Reward),
		Promote: toAPIWallet(out.Promote),
	}, nil
}

func toAPIWallet(w portal.WalletInfoData) v1.WalletInfo {
	return v1.WalletInfo{
		Balance:      w.Balance,
		BalanceCent:  w.BalanceCent,
		TotalIncome:  w.TotalIncome,
		TotalExpense: w.TotalExpense,
		FrozenAmount: w.FrozenAmount,
	}
}

// WalletLogs 钱包流水。
func (c *cMe) WalletLogs(ctx context.Context, req *v1.MeWalletLogsReq) (res *v1.MeWalletLogsRes, err error) {
	memberID := int64(middleware.CurrentMemberID(ctx))
	out, err := portal.AuthLogic().ListMyWalletLogs(ctx, &portal.MyWalletLogsInput{
		UserID:     memberID,
		WalletType: req.WalletType,
		PageNum:    req.PageNum,
		PageSize:   req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.MeWalletLogsRes{Total: out.Total, List: make([]*v1.WalletLogRecord, 0, len(out.List))}
	for _, item := range out.List {
		res.List = append(res.List, &v1.WalletLogRecord{
			ID:             item.ID,
			WalletType:     item.WalletType,
			WalletTypeText: item.WalletTypeText,
			ChangeType:     item.ChangeType,
			ChangeTypeText: item.ChangeTypeText,
			ChangeAmount:   item.ChangeAmount,
			BeforeBalance:  item.BeforeBalance,
			AfterBalance:   item.AfterBalance,
			RelatedOrderNo: item.RelatedOrderNo,
			Remark:         item.Remark,
			CreatedAt:      item.CreatedAt,
		})
	}
	return res, nil
}

// Team 获取团队列表。
func (c *cMe) Team(ctx context.Context, req *v1.MeTeamReq) (res *v1.MeTeamRes, err error) {
	memberID := int64(middleware.CurrentMemberID(ctx))
	out, err := portal.AuthLogic().ListMyTeam(ctx, &portal.MyTeamInput{
		UserID:   memberID,
		Scope:    req.Scope,
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.MeTeamRes{Total: out.Total, List: make([]*v1.TeamMemberItem, 0, len(out.List))}
	for _, item := range out.List {
		res.List = append(res.List, &v1.TeamMemberItem{
			MemberID:    item.MemberID,
			Nickname:    item.Nickname,
			Avatar:      item.Avatar,
			Phone:       item.Phone,
			LevelName:   item.LevelName,
			IsQualified: item.IsQualified,
			JoinedAt:    item.JoinedAt,
		})
	}
	return res, nil
}
