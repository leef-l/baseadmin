/** 体验合同类型定义 */

/** 体验合同项 */
export interface ContractItem {
  id: string;
  contractNo: string;
  customerID?: string;
  customerName?: string;
  orderID?: string;
  orderOrderNo?: string;
  title: string;
  contractFile?: string;
  signImage?: string;
  contractAmount?: number;
  signedAt?: string;
  expiresAt?: string;
  status?: number;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 体验合同列表查询参数 */
export interface ContractListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  contractNo?: string;
  title?: string;
  customerID?: string;
  orderID?: string;
  tenantID?: string;
  merchantID?: string;
  status?: number;
  signedAtStart?: string;
  signedAtEnd?: string;
  expiresAtStart?: string;
  expiresAtEnd?: string;
}

/** 体验合同创建参数 */
export interface ContractCreateParams {
  contractNo: string;
  customerID?: string;
  orderID?: string;
  title: string;
  contractFile?: string;
  signImage?: string;
  contractAmount?: number;
  signPassword?: string;
  signedAt?: string;
  expiresAt?: string;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}

/** 体验合同更新参数 */
export interface ContractUpdateParams {
  id: string;
  contractNo: string;
  customerID?: string;
  orderID?: string;
  title: string;
  contractFile?: string;
  signImage?: string;
  contractAmount?: number;
  signPassword?: string;
  signedAt?: string;
  expiresAt?: string;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}
