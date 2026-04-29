/** 租户类型定义 */

/** 租户项 */
export interface TenantItem {
  id: string;
  name: string;
  code: string;
  contactName?: string;
  contactPhone?: string;
  domain?: string;
  expireAt?: string;
  status?: number;
  remark?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 租户列表查询参数 */
export interface TenantListParams {
  pageNum: number;
  pageSize: number;
  keyword?: string;
  code?: string;
  status?: number;
}

/** 租户创建参数 */
export interface TenantCreateParams {
  name: string;
  code: string;
  contactName?: string;
  contactPhone?: string;
  domain?: string;
  expireAt?: string;
  status?: number;
  remark?: string;
  createAdmin?: number;
  adminUsername?: string;
  adminPassword?: string;
  adminNickname?: string;
  adminEmail?: string;
}

/** 租户更新参数 */
export interface TenantUpdateParams extends TenantCreateParams {
  id: string;
}
