/** 域名绑定类型定义 */

export interface DomainItem {
  id: string;
  domain: string;
  ownerType: number;
  tenantId?: string;
  tenantName?: string;
  merchantId?: string;
  merchantName?: string;
  appCode: string;
  verifyToken?: string;
  verifyStatus?: number;
  sslStatus?: number;
  nginxStatus?: number;
  status?: number;
  remark?: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface DomainListParams {
  pageNum: number;
  pageSize: number;
  keyword?: string;
  domain?: string;
  ownerType?: number;
  tenantId?: string;
  merchantId?: string;
  appCode?: string;
  status?: number;
}

export interface DomainCreateParams {
  domain: string;
  ownerType?: number;
  tenantId?: string;
  merchantId?: string;
  appCode?: string;
  verifyStatus?: number;
  sslStatus?: number;
  status?: number;
  remark?: string;
}

export interface DomainUpdateParams extends DomainCreateParams {
  id: string;
}

export interface DomainApplyNginxResult {
  configPath: string;
  nginxStatus: number;
  sslStatus: number;
}

export interface DomainApplySSLResult extends DomainApplyNginxResult {
  certPath: string;
}
