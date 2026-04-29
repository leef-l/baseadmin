/** 体验工单类型定义 */

/** 体验工单项 */
export interface WorkOrderItem {
  id: string;
  ticketNo: string;
  customerID?: string;
  customerName?: string;
  productID?: string;
  productSkuNo?: string;
  orderID?: string;
  orderOrderNo?: string;
  title: string;
  priority?: number;
  sourceType?: number;
  description?: string;
  attachmentFile?: string;
  dueAt?: string;
  status?: number;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 体验工单列表查询参数 */
export interface WorkOrderListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  keyword?: string;
  ticketNo?: string;
  title?: string;
  customerID?: string;
  productID?: string;
  orderID?: string;
  tenantID?: string;
  merchantID?: string;
  priority?: number;
  sourceType?: number;
  status?: number;
  dueAtStart?: string;
  dueAtEnd?: string;
}

/** 体验工单创建参数 */
export interface WorkOrderCreateParams {
  ticketNo: string;
  customerID?: string;
  productID?: string;
  orderID?: string;
  title: string;
  priority?: number;
  sourceType?: number;
  description?: string;
  attachmentFile?: string;
  dueAt?: string;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}

/** 体验工单更新参数 */
export interface WorkOrderUpdateParams {
  id: string;
  ticketNo: string;
  customerID?: string;
  productID?: string;
  orderID?: string;
  title: string;
  priority?: number;
  sourceType?: number;
  description?: string;
  attachmentFile?: string;
  dueAt?: string;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}
