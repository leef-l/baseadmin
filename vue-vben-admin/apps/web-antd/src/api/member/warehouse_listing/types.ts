/** 仓库挂卖记录类型定义 */

/** 仓库挂卖记录项 */
export interface WarehouseListingItem {
  id: string;
  goodsID?: string;
  warehouseGoodsTitle?: string;
  sellerID?: string;
  userNickname?: string;
  listingPrice?: string;
  listingStatus?: number;
  listedAt?: string;
  soldAt?: string;
  remark?: string;
  status?: number;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 仓库挂卖记录列表查询参数 */
export interface WarehouseListingListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  goodsID?: string;
  sellerID?: string;
  tenantID?: string;
  merchantID?: string;
  listingStatus?: number;
  status?: number;
  listedAtStart?: string;
  listedAtEnd?: string;
  soldAtStart?: string;
  soldAtEnd?: string;
}

/** 仓库挂卖记录创建参数 */
export interface WarehouseListingCreateParams {
  goodsID?: string;
  sellerID?: string;
  listingPrice?: string;
  listingStatus?: number;
  listedAt?: string;
  soldAt?: string;
  remark?: string;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}

/** 仓库挂卖记录更新参数 */
export interface WarehouseListingUpdateParams {
  id: string;
  goodsID?: string;
  sellerID?: string;
  listingPrice?: string;
  listingStatus?: number;
  listedAt?: string;
  soldAt?: string;
  remark?: string;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}
