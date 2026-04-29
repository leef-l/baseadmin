/** 仓库交易记录类型定义 */

/** 仓库交易记录项 */
export interface WarehouseTradeItem {
  id: string;
  tradeNo: string;
  goodsID?: string;
  warehouseGoodsTitle?: string;
  listingID?: string;
  warehouseListingID?: string;
  sellerID?: string;
  userNickname?: string;
  buyerID?: string;
  buyerNickname?: string;
  tradePrice?: string;
  platformFee?: string;
  sellerIncome?: string;
  tradeStatus?: number;
  confirmedAt?: string;
  remark?: string;
  status?: number;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 仓库交易记录列表查询参数 */
export interface WarehouseTradeListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  tradeNo?: string;
  goodsID?: string;
  listingID?: string;
  sellerID?: string;
  buyerID?: string;
  tenantID?: string;
  merchantID?: string;
  tradeStatus?: number;
  status?: number;
  confirmedAtStart?: string;
  confirmedAtEnd?: string;
}

/** 仓库交易记录创建参数 */
export interface WarehouseTradeCreateParams {
  tradeNo: string;
  goodsID?: string;
  listingID?: string;
  sellerID?: string;
  buyerID?: string;
  tradePrice?: string;
  platformFee?: string;
  sellerIncome?: string;
  tradeStatus?: number;
  confirmedAt?: string;
  remark?: string;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}

/** 仓库交易记录更新参数 */
export interface WarehouseTradeUpdateParams {
  id: string;
  tradeNo: string;
  goodsID?: string;
  listingID?: string;
  sellerID?: string;
  buyerID?: string;
  tradePrice?: string;
  platformFee?: string;
  sellerIncome?: string;
  tradeStatus?: number;
  confirmedAt?: string;
  remark?: string;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}
