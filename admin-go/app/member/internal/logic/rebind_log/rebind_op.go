package rebind_log

// 这里实现"换绑上级"的真正业务（双向重算 team_count / direct_count / active_count / team_turnover）。
//
// 与 rebind_log.go（codegen 自动生成的 CRUD）解耦，由后台 controller 直接调用本文件方法。
// 不进 service 接口（service 是 codegen 维护，避免每次重新生成被覆盖）。

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/member/internal/dao"
	"gbaseadmin/app/member/internal/model/do"
	"gbaseadmin/app/member/internal/model/entity"
	"gbaseadmin/utility/snowflake"
)

// RebindParentInput 换绑入参。
type RebindParentInput struct {
	UserID       int64
	NewParentID  int64 // 0 表示置顶（无上级）
	Reason       string
	OperatorID   int64 // 后台操作人
}

// RebindParent 把 user 的上级改成 newParent，并把 user 自己 + user 所有下级 对原链路 / 新链路的统计贡献做加减。
//
// 前置校验：
//   - user 必须存在且未删除
//   - newParent != user（不能挂自己）
//   - newParent 不能是 user 的下级（防环）
//   - newParent != user.parent_id（无变化直接拒绝）
//
// 步骤（同事务）：
//  1. 锁 user 行，读 oldParentID
//  2. 收集 user + 全部下级 ID（subtree），用于计算贡献量
//  3. 计算贡献量：count = len(subtree)，turnover = 子树成员的"个人营业额贡献"…
//     —— 注意：member_user.team_turnover 是聚合字段，不是个人贡献，所以这里只能：
//     原链路祖先（不含 user 自己）：team_count -= count，
//     新链路祖先（含 newParent 自己）：team_count += count
//     direct_count 只影响原直接父和新直接父
//     active_count 只算子树中 is_active=1 的人数
//     team_turnover 用 user 自身 team_turnover 累加（含子树贡献），近似准确
//  4. 写 user.parent_id = newParentID
//  5. 写 member_rebind_log
func RebindParent(ctx context.Context, in *RebindParentInput) error {
	if in == nil || in.UserID <= 0 {
		return gerror.New("参数不能为空")
	}
	if in.NewParentID == in.UserID {
		return gerror.New("不能将自己设为上级")
	}

	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1. 锁 user
		var user entity.MemberUser
		if err := tx.Model(dao.MemberUser.Table()).Ctx(ctx).
			Where(dao.MemberUser.Columns().Id, in.UserID).
			Where(dao.MemberUser.Columns().DeletedAt, nil).
			LockUpdate().
			Scan(&user); err != nil {
			return err
		}
		if user.Id == 0 {
			return gerror.New("会员不存在")
		}
		oldParentID := int64(user.ParentId)
		if oldParentID == in.NewParentID {
			return gerror.New("新上级与原上级相同，无需换绑")
		}

		// 2. 校验 newParent 存在且非 user 的下级
		if in.NewParentID > 0 {
			var np entity.MemberUser
			if err := tx.Model(dao.MemberUser.Table()).Ctx(ctx).
				Where(dao.MemberUser.Columns().Id, in.NewParentID).
				Where(dao.MemberUser.Columns().DeletedAt, nil).
				Scan(&np); err != nil {
				return err
			}
			if np.Id == 0 {
				return gerror.New("新上级不存在或已删除")
			}
			if np.Status != 1 {
				return gerror.New("新上级账号已被禁用")
			}
		}

		subtree, err := collectSubtreeIDs(ctx, tx, in.UserID)
		if err != nil {
			return err
		}
		if in.NewParentID > 0 {
			for _, id := range subtree {
				if id == in.NewParentID {
					return gerror.Newf("新上级 %d 是当前会员的下级，不能换绑形成环", in.NewParentID)
				}
			}
		}

		// 3. 算贡献：count = len(subtree)；activeCount = 子树 is_active=1 数量；turnover = user.team_turnover
		count := len(subtree)
		activeCount, err := countActiveMembers(ctx, tx, subtree)
		if err != nil {
			return err
		}
		turnover := int64(user.TeamTurnover)

		// 4. 原链路祖先 -= 贡献
		if oldParentID > 0 {
			if err := walkUpAndApply(ctx, tx, oldParentID, -count, -activeCount, -turnover); err != nil {
				return err
			}
			// 原直接父 direct_count - 1
			if _, err := tx.Model(dao.MemberUser.Table()).Ctx(ctx).
				Where(dao.MemberUser.Columns().Id, oldParentID).
				Data(g.Map{
					dao.MemberUser.Columns().DirectCount: gdb.Raw("GREATEST(direct_count - 1, 0)"),
				}).Update(); err != nil {
				return err
			}
		}

		// 5. 新链路祖先 += 贡献
		if in.NewParentID > 0 {
			if err := walkUpAndApply(ctx, tx, in.NewParentID, count, activeCount, turnover); err != nil {
				return err
			}
			// 新直接父 direct_count + 1
			if _, err := tx.Model(dao.MemberUser.Table()).Ctx(ctx).
				Where(dao.MemberUser.Columns().Id, in.NewParentID).
				Data(g.Map{
					dao.MemberUser.Columns().DirectCount: gdb.Raw("direct_count + 1"),
				}).Update(); err != nil {
				return err
			}
		}

		// 6. 改 user.parent_id
		if _, err := tx.Model(dao.MemberUser.Table()).Ctx(ctx).
			Where(dao.MemberUser.Columns().Id, in.UserID).
			Data(g.Map{
				dao.MemberUser.Columns().ParentId: in.NewParentID,
			}).Update(); err != nil {
			return err
		}

		// 7. 写日志
		now := gtime.Now()
		if _, err := tx.Model(dao.MemberRebindLog.Table()).Ctx(ctx).Data(do.MemberRebindLog{
			Id:          snowflake.Generate(),
			UserId:      in.UserID,
			OldParentId: oldParentID,
			NewParentId: in.NewParentID,
			Reason:      strings.TrimSpace(in.Reason),
			OperatorId:  in.OperatorID,
			TenantId:    user.TenantId,
			MerchantId:  user.MerchantId,
			CreatedBy:   in.OperatorID,
			DeptId:      0,
			CreatedAt:   now,
			UpdatedAt:   now,
		}).Insert(); err != nil {
			return err
		}
		return nil
	})
}

// collectSubtreeIDs 用 BFS 收集 root 自身 + 全部下级 ID，最多 5000 个，超过报错。
func collectSubtreeIDs(ctx context.Context, tx gdb.TX, root int64) ([]int64, error) {
	const maxSubtreeSize = 5000
	seen := map[int64]struct{}{root: {}}
	out := []int64{root}
	queue := []int64{root}
	for len(queue) > 0 {
		batch := queue
		queue = nil
		var rows []struct {
			Id int64 `json:"id"`
		}
		if err := tx.Model(dao.MemberUser.Table()).Ctx(ctx).
			Fields(dao.MemberUser.Columns().Id).
			WhereIn(dao.MemberUser.Columns().ParentId, batch).
			Where(dao.MemberUser.Columns().DeletedAt, nil).
			Scan(&rows); err != nil {
			return nil, err
		}
		for _, r := range rows {
			if r.Id <= 0 {
				continue
			}
			if _, ok := seen[r.Id]; ok {
				continue
			}
			seen[r.Id] = struct{}{}
			out = append(out, r.Id)
			queue = append(queue, r.Id)
			if len(out) > maxSubtreeSize {
				return nil, gerror.Newf("子树规模超过 %d 人，无法换绑", maxSubtreeSize)
			}
		}
	}
	return out, nil
}

// countActiveMembers 统计 ids 中 is_active=1 的人数。
func countActiveMembers(ctx context.Context, tx gdb.TX, ids []int64) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	count, err := tx.Model(dao.MemberUser.Table()).Ctx(ctx).
		WhereIn(dao.MemberUser.Columns().Id, ids).
		Where(dao.MemberUser.Columns().IsActive, 1).
		Where(dao.MemberUser.Columns().DeletedAt, nil).
		Count()
	return count, err
}

// walkUpAndApply 沿 startID 向上（含自身）应用 deltaCount/deltaActive/deltaTurnover。
// 不调用 teamops（teamops 设计为只增；这里换绑可能需要负值），所以单独实现。
func walkUpAndApply(ctx context.Context, tx gdb.TX, startID int64, deltaCount, deltaActive int, deltaTurnover int64) error {
	const maxDepth = 100
	cur := startID
	visited := map[int64]struct{}{}
	for depth := 0; depth < maxDepth; depth++ {
		if cur <= 0 {
			return nil
		}
		if _, ok := visited[cur]; ok {
			return gerror.Newf("祖先链路存在环：user=%d", cur)
		}
		visited[cur] = struct{}{}

		// 应用 delta（带 GREATEST 防止 unsigned 列下溢）
		updates := g.Map{}
		if deltaCount != 0 {
			updates[dao.MemberUser.Columns().TeamCount] = gdb.Raw(deltaCountExpr("team_count", deltaCount))
		}
		if deltaActive != 0 {
			updates[dao.MemberUser.Columns().ActiveCount] = gdb.Raw(deltaCountExpr("active_count", deltaActive))
		}
		if deltaTurnover != 0 {
			updates[dao.MemberUser.Columns().TeamTurnover] = gdb.Raw(deltaTurnoverExpr("team_turnover", deltaTurnover))
		}
		if len(updates) > 0 {
			if _, err := tx.Model(dao.MemberUser.Table()).Ctx(ctx).
				Where(dao.MemberUser.Columns().Id, cur).
				Data(updates).Update(); err != nil {
				return err
			}
		}

		next, err := loadAncestorParentID(ctx, tx, cur)
		if err != nil {
			return err
		}
		cur = next
	}
	return gerror.New("祖先链路深度超过 100 层，请检查数据完整性")
}

func loadAncestorParentID(ctx context.Context, tx gdb.TX, userID int64) (int64, error) {
	v, err := tx.Model(dao.MemberUser.Table()).Ctx(ctx).
		Where(dao.MemberUser.Columns().Id, userID).
		Where(dao.MemberUser.Columns().DeletedAt, nil).
		Value(dao.MemberUser.Columns().ParentId)
	if err != nil {
		return 0, err
	}
	if v == nil || v.IsNil() || v.IsEmpty() {
		return 0, nil
	}
	return v.Int64(), nil
}

func deltaCountExpr(column string, delta int) string {
	if delta >= 0 {
		return fmt.Sprintf("%s + %d", column, delta)
	}
	// unsigned 列防下溢：用 GREATEST(col + delta, 0)
	return fmt.Sprintf("GREATEST(CAST(%s AS SIGNED) + %d, 0)", column, delta)
}

func deltaTurnoverExpr(column string, delta int64) string {
	if delta >= 0 {
		return fmt.Sprintf("%s + %s", column, strconv.FormatInt(delta, 10))
	}
	return fmt.Sprintf("GREATEST(CAST(%s AS SIGNED) + %s, 0)", column, strconv.FormatInt(delta, 10))
}
