/** 商户类型定义 */

/** 商户项 */
export interface MerchantItem {
  id: string;
  tenantId: string;
  tenantName?: string;
  name: string;
  code: string;
  contactName?: string;
  contactPhone?: string;
  address?: string;
  status?: number;
  remark?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 商户列表查询参数 */
export interface MerchantListParams {
  pageNum: number;
  pageSize: number;
  keyword?: string;
  tenantId?: string;
  code?: string;
  status?: number;
}

/** 商户创建参数 */
export interface MerchantCreateParams {
  tenantId?: string;
  name: string;
  code: string;
  contactName?: string;
  contactPhone?: string;
  address?: string;
  status?: number;
  remark?: string;
  createAdmin?: number;
  adminUsername?: string;
  adminPassword?: string;
  adminNickname?: string;
  adminEmail?: string;
}

/** 商户更新参数 */
export interface MerchantUpdateParams extends MerchantCreateParams {
  id: string;
}
