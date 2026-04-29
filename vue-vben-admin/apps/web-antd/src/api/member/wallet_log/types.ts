/** 钱包流水记录类型定义 */

/** 钱包流水记录项 */
export interface WalletLogItem {
  id: string;
  userID?: string;
  userNickname?: string;
  walletType?: number;
  changeType?: number;
  changeAmount?: string;
  beforeBalance?: string;
  afterBalance?: string;
  relatedOrderNo?: string;
  remark?: string;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 钱包流水记录列表查询参数 */
export interface WalletLogListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  relatedOrderNo?: string;
  userID?: string;
  tenantID?: string;
  merchantID?: string;
  walletType?: number;
  changeType?: number;
}

/** 钱包流水记录创建参数 */
export interface WalletLogCreateParams {
  userID?: string;
  walletType?: number;
  changeType?: number;
  changeAmount?: string;
  beforeBalance?: string;
  afterBalance?: string;
  relatedOrderNo?: string;
  remark?: string;
  tenantID?: string;
  merchantID?: string;
}

/** 钱包流水记录更新参数 */
export interface WalletLogUpdateParams {
  id: string;
  userID?: string;
  walletType?: number;
  changeType?: number;
  changeAmount?: string;
  beforeBalance?: string;
  afterBalance?: string;
  relatedOrderNo?: string;
  remark?: string;
  tenantID?: string;
  merchantID?: string;
}
