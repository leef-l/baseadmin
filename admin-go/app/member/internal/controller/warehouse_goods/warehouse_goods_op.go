package warehouse_goods

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	v1 "gbaseadmin/app/member/api/member/v1"
	"gbaseadmin/app/member/internal/dao"
	"gbaseadmin/app/member/internal/model/entity"
)

// Assign 后台分配仓库商品给指定会员。
//
// 直接写在 controller 层调用 dao（不重新走 service.WarehouseGoods 的 Update，因为 Update 校验过严：
// 比如它会带上数据权限校验，而管理员后台账号期望走通；codegen Update 还会强制校验所有字段）。
func (c *cWarehouseGoods) Assign(ctx context.Context, req *v1.WarehouseGoodsAssignReq) (res *v1.WarehouseGoodsAssignRes, err error) {
	if int64(req.GoodsID) <= 0 || int64(req.OwnerID) <= 0 {
		return nil, gerror.New("商品 ID / 目标会员 ID 不能为空")
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 锁商品行
		var goods entity.MemberWarehouseGoods
		if err := tx.Model(dao.MemberWarehouseGoods.Table()).Ctx(ctx).
			Where(dao.MemberWarehouseGoods.Columns().Id, req.GoodsID).
			Where(dao.MemberWarehouseGoods.Columns().DeletedAt, nil).
			LockUpdate().
			Scan(&goods); err != nil {
			return err
		}
		if goods.Id == 0 {
			return gerror.New("商品不存在或已删除")
		}
		// 状态约束：仅持有中可换持有人
		if goods.GoodsStatus != 1 {
			return gerror.New("仅持有中状态的商品可重新分配")
		}

		// 校验目标会员
		var owner entity.MemberUser
		if err := tx.Model(dao.MemberUser.Table()).Ctx(ctx).
			Where(dao.MemberUser.Columns().Id, req.OwnerID).
			Where(dao.MemberUser.Columns().DeletedAt, nil).
			Scan(&owner); err != nil {
			return err
		}
		if owner.Id == 0 {
			return gerror.New("目标会员不存在或已删除")
		}
		if owner.Status != 1 {
			return gerror.New("目标会员账号已被禁用")
		}

		// 写入 owner_id；其他字段不动
		if _, err := tx.Model(dao.MemberWarehouseGoods.Table()).Ctx(ctx).
			Where(dao.MemberWarehouseGoods.Columns().Id, goods.Id).
			Data(g.Map{
				dao.MemberWarehouseGoods.Columns().OwnerId: req.OwnerID,
				dao.MemberWarehouseGoods.Columns().Remark:  req.Remark,
			}).Update(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &v1.WarehouseGoodsAssignRes{}, nil
}
