/** 会员钱包类型定义 */

/** 会员钱包项 */
export interface WalletItem {
  id: string;
  userID?: string;
  userUsername?: string;
  walletType?: number;
  balance?: string;
  totalIncome?: string;
  totalExpense?: string;
  frozenAmount?: string;
  status?: number;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 会员钱包列表查询参数 */
export interface WalletListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  userID?: string;
  tenantID?: string;
  merchantID?: string;
  walletType?: number;
  status?: number;
}

/** 会员钱包创建参数 */
export interface WalletCreateParams {
  userID?: string;
  walletType?: number;
  balance?: string;
  totalIncome?: string;
  totalExpense?: string;
  frozenAmount?: string;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}

/** 会员钱包更新参数 */
export interface WalletUpdateParams {
  id: string;
  userID?: string;
  walletType?: number;
  balance?: string;
  totalIncome?: string;
  totalExpense?: string;
  frozenAmount?: string;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}
