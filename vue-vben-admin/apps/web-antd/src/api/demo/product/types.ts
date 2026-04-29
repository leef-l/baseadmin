/** 体验商品类型定义 */

/** 体验商品项 */
export interface ProductItem {
  id: string;
  categoryID?: string;
  categoryName?: string;
  skuNo: string;
  name: string;
  cover?: string;
  manualFile?: string;
  detailContent?: string;
  specJSON?: string;
  websiteURL?: string;
  type?: number;
  isRecommend?: number;
  salePrice?: number;
  stockNum?: number;
  weightNum?: number;
  sort?: number;
  icon?: string;
  status?: number;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 体验商品列表查询参数 */
export interface ProductListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  skuNo?: string;
  name?: string;
  categoryID?: string;
  tenantID?: string;
  merchantID?: string;
  type?: number;
  isRecommend?: number;
  status?: number;
}

/** 体验商品创建参数 */
export interface ProductCreateParams {
  categoryID?: string;
  skuNo: string;
  name: string;
  cover?: string;
  manualFile?: string;
  detailContent?: string;
  specJSON?: string;
  websiteURL?: string;
  type?: number;
  isRecommend?: number;
  salePrice?: number;
  stockNum?: number;
  weightNum?: number;
  sort?: number;
  icon?: string;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}

/** 体验商品更新参数 */
export interface ProductUpdateParams {
  id: string;
  categoryID?: string;
  skuNo: string;
  name: string;
  cover?: string;
  manualFile?: string;
  detailContent?: string;
  specJSON?: string;
  websiteURL?: string;
  type?: number;
  isRecommend?: number;
  salePrice?: number;
  stockNum?: number;
  weightNum?: number;
  sort?: number;
  icon?: string;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}
