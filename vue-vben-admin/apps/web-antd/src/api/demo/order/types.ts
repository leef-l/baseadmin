/** 体验订单类型定义 */

/** 体验订单项 */
export interface OrderItem {
  id: string;
  orderNo: string;
  customerID?: string;
  customerName?: string;
  productID?: string;
  productSkuNo?: string;
  quantity?: number;
  amount?: number;
  payStatus?: number;
  deliverStatus?: number;
  paidAt?: string;
  deliverAt?: string;
  receiverPhone?: string;
  address?: string;
  remark?: string;
  status?: number;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 体验订单列表查询参数 */
export interface OrderListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  keyword?: string;
  orderNo?: string;
  receiverPhone?: string;
  customerID?: string;
  productID?: string;
  tenantID?: string;
  merchantID?: string;
  payStatus?: number;
  deliverStatus?: number;
  status?: number;
  paidAtStart?: string;
  paidAtEnd?: string;
  deliverAtStart?: string;
  deliverAtEnd?: string;
}

/** 体验订单创建参数 */
export interface OrderCreateParams {
  orderNo: string;
  customerID?: string;
  productID?: string;
  quantity?: number;
  amount?: number;
  payStatus?: number;
  deliverStatus?: number;
  paidAt?: string;
  deliverAt?: string;
  receiverPhone?: string;
  address?: string;
  remark?: string;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}

/** 体验订单更新参数 */
export interface OrderUpdateParams {
  id: string;
  orderNo: string;
  customerID?: string;
  productID?: string;
  quantity?: number;
  amount?: number;
  payStatus?: number;
  deliverStatus?: number;
  paidAt?: string;
  deliverAt?: string;
  receiverPhone?: string;
  address?: string;
  remark?: string;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}
