/** 仓库商品类型定义 */

/** 仓库商品项 */
export interface WarehouseGoodsItem {
  id: string;
  goodsNo: string;
  title: string;
  cover?: string;
  initPrice?: string;
  currentPrice?: string;
  priceRiseRate?: number;
  platformFeeRate?: number;
  ownerID?: string;
  userNickname?: string;
  tradeCount?: number;
  goodsStatus?: number;
  remark?: string;
  sort?: number;
  status?: number;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 仓库商品列表查询参数 */
export interface WarehouseGoodsListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  keyword?: string;
  goodsNo?: string;
  title?: string;
  ownerID?: string;
  tenantID?: string;
  merchantID?: string;
  goodsStatus?: number;
  status?: number;
}

/** 仓库商品创建参数 */
export interface WarehouseGoodsCreateParams {
  goodsNo: string;
  title: string;
  cover?: string;
  initPrice?: string;
  currentPrice?: string;
  priceRiseRate?: number;
  platformFeeRate?: number;
  ownerID?: string;
  tradeCount?: number;
  goodsStatus?: number;
  remark?: string;
  sort?: number;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}

/** 仓库商品更新参数 */
export interface WarehouseGoodsUpdateParams {
  id: string;
  goodsNo: string;
  title: string;
  cover?: string;
  initPrice?: string;
  currentPrice?: string;
  priceRiseRate?: number;
  platformFeeRate?: number;
  ownerID?: string;
  tradeCount?: number;
  goodsStatus?: number;
  remark?: string;
  sort?: number;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}
