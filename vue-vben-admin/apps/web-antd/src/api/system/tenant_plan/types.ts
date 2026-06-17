/** 租户套餐订阅表类型定义 */

/** 租户套餐订阅表项 */
export interface TenantPlanItem {
  id: string;
  tenantID: string;
  tenantName?: string;
  planID: string;
  planName?: string;
  startAt?: string;
  expireAt?: string;
  status?: number;
  remark?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 租户套餐订阅表列表查询参数 */
export interface TenantPlanListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  tenantID?: string;
  planID?: string;
  merchantID?: string;
  status?: number;
  remark?: string;
  startAtStart?: string;
  startAtEnd?: string;
  expireAtStart?: string;
  expireAtEnd?: string;
}

/** 租户套餐订阅表创建参数 */
export interface TenantPlanCreateParams {
  tenantID: string;
  planID: string;
  startAt?: string;
  expireAt?: string;
  status?: number;
  remark?: string;
  merchantID?: string;
}

/** 租户套餐订阅表更新参数 */
export interface TenantPlanUpdateParams {
  id: string;
  tenantID: string;
  planID: string;
  startAt?: string;
  expireAt?: string;
  status?: number;
  remark?: string;
  merchantID?: string;
}
