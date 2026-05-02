// Package teamops 是团队统计与等级晋升的领域服务（基础设施层）。
//
// 设计要点：
//   - 链式更新：注册、激活、营业额变化时，沿 parent_id 链向上递归更新所有上级的 team_count / active_count / team_turnover。
//   - 防环路：迭代深度限制 100 层（等同于 100 代上下级），防止脏数据导致无限循环。
//   - 等级实时检查：每次统计变动后调用 TryUpgrade 看是否能升级，升级写 member_level_log + 更新 user.level_id 和 level_expire_at。
//   - 等级过期：每天定时跑 ScanExpired，对 level_expire_at < now 的会员把 is_qualified 置 0、写 level_log(change_type=3)。
package teamops

import (
	"context"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/member/internal/dao"
	"gbaseadmin/app/member/internal/model/do"
	"gbaseadmin/app/member/internal/model/entity"
	"gbaseadmin/utility/snowflake"
)

// 等级变更类型常量（与 member_level_log.change_type 取值一致）。
const (
	LevelChangeUpgrade = 1 // 自动升级
	LevelChangeAdjust  = 2 // 后台调整
	LevelChangeExpire  = 3 // 过期降级
)

// 链式更新时的最大向上回溯深度，防止脏数据 parent_id 形成环。
const maxAncestorDepth = 100

// IncrAncestorTeamCount 注册成功后调用：把 newUserID 的所有祖先（含直接上级）team_count + 1，
// 直接上级额外 direct_count + 1。
func IncrAncestorTeamCount(ctx context.Context, tx gdb.TX, newUserID int64) error {
	if newUserID <= 0 {
		return nil
	}
	parentID, err := loadParentID(ctx, tx, newUserID)
	if err != nil {
		return err
	}
	if parentID <= 0 {
		return nil
	}

	// 直接上级：direct_count + 1
	if _, err := tx.Model(dao.MemberUser.Table()).Ctx(ctx).
		Where(dao.MemberUser.Columns().Id, parentID).
		Data(g.Map{
			dao.MemberUser.Columns().DirectCount: gdb.Raw("direct_count + 1"),
			dao.MemberUser.Columns().TeamCount:   gdb.Raw("team_count + 1"),
		}).Update(); err != nil {
		return err
	}

	// 链式向上 team_count + 1（直接上级已加过，从其父开始）
	return walkAncestors(ctx, tx, parentID, func(ancestorID int64) error {
		_, err := tx.Model(dao.MemberUser.Table()).Ctx(ctx).
			Where(dao.MemberUser.Columns().Id, ancestorID).
			Data(g.Map{
				dao.MemberUser.Columns().TeamCount: gdb.Raw("team_count + 1"),
			}).Update()
		return err
	})
}

// IncrAncestorActiveCount 当 user.is_active 由 0 → 1 时调用，链式更新所有祖先 active_count + 1。
func IncrAncestorActiveCount(ctx context.Context, tx gdb.TX, userID int64) error {
	if userID <= 0 {
		return nil
	}
	parentID, err := loadParentID(ctx, tx, userID)
	if err != nil {
		return err
	}
	if parentID <= 0 {
		return nil
	}
	return walkAncestorsInclusive(ctx, tx, parentID, func(ancestorID int64) error {
		_, err := tx.Model(dao.MemberUser.Table()).Ctx(ctx).
			Where(dao.MemberUser.Columns().Id, ancestorID).
			Data(g.Map{
				dao.MemberUser.Columns().ActiveCount: gdb.Raw("active_count + 1"),
			}).Update()
		if err != nil {
			return err
		}
		// 每次链式更新都尝试触发升级
		return tryUpgradeOnTx(ctx, tx, ancestorID)
	})
}

// AddAncestorTurnover 当会员 userID 产生 amount（分）的有效消费时，链式累加所有祖先 team_turnover。
//
// amount 必须 > 0。
func AddAncestorTurnover(ctx context.Context, tx gdb.TX, userID int64, amount int64) error {
	if userID <= 0 || amount <= 0 {
		return nil
	}
	parentID, err := loadParentID(ctx, tx, userID)
	if err != nil {
		return err
	}
	if parentID <= 0 {
		return nil
	}
	rawExpr := "team_turnover + " + strconv.FormatInt(amount, 10)
	return walkAncestorsInclusive(ctx, tx, parentID, func(ancestorID int64) error {
		_, err := tx.Model(dao.MemberUser.Table()).Ctx(ctx).
			Where(dao.MemberUser.Columns().Id, ancestorID).
			Data(g.Map{
				dao.MemberUser.Columns().TeamTurnover: gdb.Raw(rawExpr),
			}).Update()
		if err != nil {
			return err
		}
		return tryUpgradeOnTx(ctx, tx, ancestorID)
	})
}

// TryUpgrade 主动检查 userID 是否满足任何更高等级条件，满足则升级。
// 不在事务内调用时，会内部开事务。
func TryUpgrade(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return nil
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		return tryUpgradeOnTx(ctx, tx, userID)
	})
}

// tryUpgradeOnTx 在已有事务里检查并升级。
func tryUpgradeOnTx(ctx context.Context, tx gdb.TX, userID int64) error {
	// 锁会员行
	var u entity.MemberUser
	if err := tx.Model(dao.MemberUser.Table()).Ctx(ctx).
		Where(dao.MemberUser.Columns().Id, userID).
		Where(dao.MemberUser.Columns().DeletedAt, nil).
		LockUpdate().
		Scan(&u); err != nil {
		return err
	}
	if u.Id == 0 {
		return nil
	}

	// 找出当前等级 level_no
	var currentLevelNo uint
	if u.LevelId > 0 {
		var lv entity.MemberLevel
		if err := tx.Model(dao.MemberLevel.Table()).Ctx(ctx).
			Where(dao.MemberLevel.Columns().Id, u.LevelId).
			Where(dao.MemberLevel.Columns().DeletedAt, nil).
			Scan(&lv); err != nil {
			return err
		}
		currentLevelNo = lv.LevelNo
	}

	// 找出最高满足条件的等级（level_no > current 且 active_count/team_turnover 都达标）
	var candidates []*entity.MemberLevel
	if err := tx.Model(dao.MemberLevel.Table()).Ctx(ctx).
		Where(dao.MemberLevel.Columns().Status, 1).
		Where(dao.MemberLevel.Columns().DeletedAt, nil).
		Where(dao.MemberLevel.Columns().LevelNo+" > ?", currentLevelNo).
		Where(dao.MemberLevel.Columns().NeedActiveCount+" <= ?", u.ActiveCount).
		Where(dao.MemberLevel.Columns().NeedTeamTurnover+" <= ?", u.TeamTurnover).
		OrderDesc(dao.MemberLevel.Columns().LevelNo).
		Limit(1).
		Scan(&candidates); err != nil {
		return err
	}
	if len(candidates) == 0 || candidates[0] == nil {
		return nil
	}
	target := candidates[0]

	// 计算到期时间
	now := gtime.Now()
	var expireAt *gtime.Time
	if target.DurationDays > 0 {
		expireAt = gtime.New(now.Time.AddDate(0, 0, int(target.DurationDays)))
	}

	// 写等级日志
	if _, err := tx.Model(dao.MemberLevelLog.Table()).Ctx(ctx).Data(do.MemberLevelLog{
		Id:         snowflake.Generate(),
		UserId:     userID,
		OldLevelId: u.LevelId,
		NewLevelId: target.Id,
		ChangeType: LevelChangeUpgrade,
		ExpireAt:   expireAt,
		Remark:     "满足升级条件，自动升级到 " + target.Name,
		TenantId:   u.TenantId,
		MerchantId: u.MerchantId,
		CreatedBy:  0,
		DeptId:     0,
		CreatedAt:  now,
		UpdatedAt:  now,
	}).Insert(); err != nil {
		return err
	}

	// 更新会员：等级 + 资格 + 每日限购（按目标等级覆盖）
	// daily_purchase_limit 设计上由等级决定，升级即同步；管理员后台仍可单独再调，但下一次升级会被再次覆盖。
	updateData := g.Map{
		dao.MemberUser.Columns().LevelId:            target.Id,
		dao.MemberUser.Columns().IsQualified:        1,
		dao.MemberUser.Columns().DailyPurchaseLimit: target.DailyPurchaseLimit,
	}
	if expireAt != nil {
		updateData[dao.MemberUser.Columns().LevelExpireAt] = expireAt
	}
	_, err := tx.Model(dao.MemberUser.Table()).Ctx(ctx).
		Where(dao.MemberUser.Columns().Id, userID).
		Data(updateData).Update()
	return err
}

// ScanExpiredLevels 每日定时调用：把 level_expire_at < now 且 is_qualified=1 的会员置 is_qualified=0，
// 同时写 level_log(change_type=3)。
//
// 返回处理的会员数。
func ScanExpiredLevels(ctx context.Context) (int, error) {
	now := gtime.Now()
	var users []*entity.MemberUser
	if err := dao.MemberUser.Ctx(ctx).
		Where(dao.MemberUser.Columns().DeletedAt, nil).
		Where(dao.MemberUser.Columns().IsQualified, 1).
		WhereLT(dao.MemberUser.Columns().LevelExpireAt, now).
		WhereNot(dao.MemberUser.Columns().LevelExpireAt, nil).
		Scan(&users); err != nil {
		return 0, err
	}
	if len(users) == 0 {
		return 0, nil
	}

	count := 0
	for _, u := range users {
		if u == nil || u.Id == 0 {
			continue
		}
		err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// 更新 user
			if _, err := tx.Model(dao.MemberUser.Table()).Ctx(ctx).
				Where(dao.MemberUser.Columns().Id, u.Id).
				Where(dao.MemberUser.Columns().IsQualified, 1).
				Data(g.Map{
					dao.MemberUser.Columns().IsQualified: 0,
				}).Update(); err != nil {
				return err
			}
			// 写日志
			if _, err := tx.Model(dao.MemberLevelLog.Table()).Ctx(ctx).Data(do.MemberLevelLog{
				Id:         snowflake.Generate(),
				UserId:     u.Id,
				OldLevelId: u.LevelId,
				NewLevelId: u.LevelId, // 等级不变，仅资格失效
				ChangeType: LevelChangeExpire,
				Remark:     "等级有效期到期，仓库资格失效",
				TenantId:   u.TenantId,
				MerchantId: u.MerchantId,
				CreatedBy:  0,
				DeptId:     0,
				CreatedAt:  now,
				UpdatedAt:  now,
			}).Insert(); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			g.Log().Errorf(ctx, "[teamops] ScanExpiredLevels user=%d err=%v", u.Id, err)
			continue
		}
		count++
	}
	return count, nil
}

// ----- 内部 helpers -----

// loadParentID 读取 user.parent_id（在事务内）。
func loadParentID(ctx context.Context, tx gdb.TX, userID int64) (int64, error) {
	value, err := tx.Model(dao.MemberUser.Table()).Ctx(ctx).
		Where(dao.MemberUser.Columns().Id, userID).
		Where(dao.MemberUser.Columns().DeletedAt, nil).
		Value(dao.MemberUser.Columns().ParentId)
	if err != nil {
		return 0, err
	}
	if value == nil || value.IsNil() || value.IsEmpty() {
		return 0, nil
	}
	return value.Int64(), nil
}

// walkAncestors 从 startParentID 的"父亲"开始向上迭代（不含 startParentID 本身）。
func walkAncestors(ctx context.Context, tx gdb.TX, startParentID int64, fn func(ancestorID int64) error) error {
	cur := startParentID
	visited := map[int64]struct{}{cur: {}}
	for depth := 0; depth < maxAncestorDepth; depth++ {
		next, err := loadParentID(ctx, tx, cur)
		if err != nil {
			return err
		}
		if next <= 0 {
			return nil
		}
		if _, looped := visited[next]; looped {
			return gerror.Newf("会员上级链路存在环：user=%d", next)
		}
		visited[next] = struct{}{}
		if err := fn(next); err != nil {
			return err
		}
		cur = next
	}
	return gerror.New("会员上级链路深度超过 100 层，请检查数据完整性")
}

// walkAncestorsInclusive 从 startUserID 自身开始向上迭代（含 startUserID）。
func walkAncestorsInclusive(ctx context.Context, tx gdb.TX, startUserID int64, fn func(ancestorID int64) error) error {
	if startUserID <= 0 {
		return nil
	}
	if err := fn(startUserID); err != nil {
		return err
	}
	return walkAncestors(ctx, tx, startUserID, fn)
}

// CronInterval 定时任务建议触发频率：每天凌晨 1 点扫一次过期等级。
func CronInterval() time.Duration { return 24 * time.Hour }
