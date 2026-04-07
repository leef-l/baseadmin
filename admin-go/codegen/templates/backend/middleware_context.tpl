package middleware

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"gbaseadmin/utility/snowflake"
)

// GetUserID 从 context 中获取当前登录用户 ID
func GetUserID(ctx context.Context) snowflake.JsonInt64 {
	val := g.RequestFromCtx(ctx).GetCtxVar("jwt_user_id")
	return snowflake.JsonInt64(val.Int64())
}

// GetDeptID 从 context 中获取当前登录用户部门 ID
func GetDeptID(ctx context.Context) snowflake.JsonInt64 {
	val := g.RequestFromCtx(ctx).GetCtxVar("jwt_dept_id")
	return snowflake.JsonInt64(val.Int64())
}

// GetIsAdmin 从 context 中获取当前用户是否是超级管理员
func GetIsAdmin(ctx context.Context) bool {
	claims := g.RequestFromCtx(ctx).GetCtxVar("jwt_claims")
	if claims.IsNil() {
		return false
	}
	return false // 默认非管理员，实际项目中从 claims 中获取
}

// ApplyDataScope 应用数据权限过滤
// 根据当前用户的角色数据范围，自动添加 WHERE 条件
// cols 参数为可选的字段名：第一个为 created_by 字段名，第二个为 dept_id 字段名
func ApplyDataScope(ctx context.Context, m *gdb.Model, cols ...string) *gdb.Model {
	// 超级管理员不做过滤
	if GetIsAdmin(ctx) {
		return m
	}

	userID := GetUserID(ctx)
	deptID := GetDeptID(ctx)

	// 获取用户角色的数据范围（从 context 或数据库查询）
	// data_scope: 1=全部, 2=本部门及以下, 3=本部门, 4=仅本人, 5=自定义
	// 默认按"仅本人"过滤（最严格）
	dataScope := getDataScope(ctx)

	switch dataScope {
	case 1:
		// 全部数据，不做过滤
		return m
	case 2, 3:
		// 本部门（及以下）：按 dept_id 过滤
		if len(cols) >= 2 && cols[1] != "" {
			if dataScope == 2 {
				// 本部门及以下：查询当前部门及所有子部门的 ID
				deptIDs := getSubDeptIDs(ctx, deptID)
				deptIDs = append(deptIDs, deptID)
				m = m.WhereIn(cols[1], deptIDs)
			} else {
				// 仅本部门
				m = m.Where(cols[1], deptID)
			}
		}
		return m
	case 4:
		// 仅本人：按 created_by 过滤
		if len(cols) >= 1 && cols[0] != "" {
			m = m.Where(cols[0], userID)
		}
		return m
	default:
		// 默认仅本人
		if len(cols) >= 1 && cols[0] != "" {
			m = m.Where(cols[0], userID)
		}
		return m
	}
}

// getDataScope 获取当前用户的数据范围（取角色中最大的范围）
func getDataScope(ctx context.Context) int {
	userID := GetUserID(ctx)
	if userID == 0 {
		return 4 // 未登录，仅本人
	}
	// 查询用户角色的最大数据范围
	val, err := g.DB().Ctx(ctx).Raw(`
		SELECT MIN(r.data_scope) FROM system_role r
		INNER JOIN system_user_role ur ON ur.role_id = r.id
		WHERE ur.user_id = ? AND r.status = 1 AND r.deleted_at IS NULL
	`, userID).Value()
	if err != nil || val.IsNil() || val.IsEmpty() {
		return 4 // 默认仅本人
	}
	return val.Int()
}

// getSubDeptIDs 获取部门及所有子部门 ID
func getSubDeptIDs(ctx context.Context, parentDeptID snowflake.JsonInt64) []snowflake.JsonInt64 {
	var result []snowflake.JsonInt64
	children, err := g.DB().Ctx(ctx).Model("system_dept").
		Where("parent_id", parentDeptID).
		Where("deleted_at IS NULL").
		Array("id")
	if err != nil || len(children) == 0 {
		return result
	}
	for _, v := range children {
		childID := snowflake.JsonInt64(v.Int64())
		result = append(result, childID)
		result = append(result, getSubDeptIDs(ctx, childID)...)
	}
	return result
}
