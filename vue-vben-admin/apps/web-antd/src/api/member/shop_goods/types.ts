/** 商城商品类型定义 */

/** 商城商品项 */
export interface ShopGoodsItem {
  id: string;
  categoryID?: string;
  shopCategoryName?: string;
  title: string;
  cover?: string;
  images?: string;
  price?: string;
  originalPrice?: string;
  stock?: number;
  sales?: number;
  content?: string;
  sort?: number;
  isRecommend?: number;
  status?: number;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 商城商品列表查询参数 */
export interface ShopGoodsListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  title?: string;
  categoryID?: string;
  tenantID?: string;
  merchantID?: string;
  isRecommend?: number;
  status?: number;
}

/** 商城商品创建参数 */
export interface ShopGoodsCreateParams {
  categoryID?: string;
  title: string;
  cover?: string;
  images?: string;
  price?: string;
  originalPrice?: string;
  stock?: number;
  sales?: number;
  content?: string;
  sort?: number;
  isRecommend?: number;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}

/** 商城商品更新参数 */
export interface ShopGoodsUpdateParams {
  id: string;
  categoryID?: string;
  title: string;
  cover?: string;
  images?: string;
  price?: string;
  originalPrice?: string;
  stock?: number;
  sales?: number;
  content?: string;
  sort?: number;
  isRecommend?: number;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}
