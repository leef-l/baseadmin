/** 体验客户类型定义 */

/** 体验客户项 */
export interface CustomerItem {
  id: string;
  avatar?: string;
  name: string;
  customerNo: string;
  phone?: string;
  email?: string;
  gender?: number;
  level?: number;
  sourceType?: number;
  isVip?: number;
  registeredAt?: string;
  remark?: string;
  status?: number;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 体验客户列表查询参数 */
export interface CustomerListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  keyword?: string;
  customerNo?: string;
  name?: string;
  phone?: string;
  email?: string;
  tenantID?: string;
  merchantID?: string;
  gender?: number;
  level?: number;
  sourceType?: number;
  isVip?: number;
  status?: number;
  registeredAtStart?: string;
  registeredAtEnd?: string;
}

/** 体验客户创建参数 */
export interface CustomerCreateParams {
  avatar?: string;
  name: string;
  customerNo: string;
  phone?: string;
  email?: string;
  gender?: number;
  level?: number;
  sourceType?: number;
  isVip?: number;
  registeredAt?: string;
  remark?: string;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}

/** 体验客户更新参数 */
export interface CustomerUpdateParams {
  id: string;
  avatar?: string;
  name: string;
  customerNo: string;
  phone?: string;
  email?: string;
  gender?: number;
  level?: number;
  sourceType?: number;
  isVip?: number;
  registeredAt?: string;
  remark?: string;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}
