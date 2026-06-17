/** 套餐表类型定义 */

/** 套餐表项 */
export interface PlanItem {
  id: string;
  name: string;
  code: string;
  description?: string;
  userLimit?: number;
  storageMb?: number;
  features?: string;
  price?: number;
  sort?: number;
  status?: number;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 套餐表列表查询参数 */
export interface PlanListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  keyword?: string;
  code?: string;
  name?: string;
  tenantID?: string;
  merchantID?: string;
  status?: number;
  description?: string;
}

/** 套餐表创建参数 */
export interface PlanCreateParams {
  name: string;
  code: string;
  description?: string;
  userLimit?: number;
  storageMb?: number;
  features?: string;
  price?: number;
  sort?: number;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}

/** 套餐表更新参数 */
export interface PlanUpdateParams {
  id: string;
  name: string;
  code: string;
  description?: string;
  userLimit?: number;
  storageMb?: number;
  features?: string;
  price?: number;
  sort?: number;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}
