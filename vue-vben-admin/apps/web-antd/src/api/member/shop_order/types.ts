/** 商城订单类型定义 */

/** 商城订单项 */
export interface ShopOrderItem {
  id: string;
  orderNo: string;
  userID?: string;
  userNickname?: string;
  goodsID?: string;
  shopGoodsTitle?: string;
  goodsTitle?: string;
  goodsCover?: string;
  quantity?: number;
  totalPrice?: string;
  payWallet?: number;
  orderStatus?: number;
  remark?: string;
  status?: number;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 商城订单列表查询参数 */
export interface ShopOrderListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  keyword?: string;
  orderNo?: string;
  userID?: string;
  goodsID?: string;
  tenantID?: string;
  merchantID?: string;
  payWallet?: number;
  orderStatus?: number;
  status?: number;
  goodsTitle?: string;
}

/** 商城订单创建参数 */
export interface ShopOrderCreateParams {
  orderNo: string;
  userID?: string;
  goodsID?: string;
  goodsTitle?: string;
  goodsCover?: string;
  quantity?: number;
  totalPrice?: string;
  payWallet?: number;
  orderStatus?: number;
  remark?: string;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}

/** 商城订单更新参数 */
export interface ShopOrderUpdateParams {
  id: string;
  orderNo: string;
  userID?: string;
  goodsID?: string;
  goodsTitle?: string;
  goodsCover?: string;
  quantity?: number;
  totalPrice?: string;
  payWallet?: number;
  orderStatus?: number;
  remark?: string;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}
