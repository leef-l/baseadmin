/** 商城商品分类类型定义 */

/** 商城商品分类项 */
export interface ShopCategoryItem {
  id: string;
  parentID?: string;
  shopCategoryName?: string;
  name: string;
  icon?: string;
  sort?: number;
  status?: number;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
  children?: ShopCategoryItem[];
}

/** 商城商品分类列表查询参数 */
export interface ShopCategoryListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  name?: string;
  parentID?: string;
  tenantID?: string;
  merchantID?: string;
  status?: number;
}
/** 商城商品分类树形查询参数 */
export interface ShopCategoryTreeParams {
  startTime?: string;
  endTime?: string;
  name?: string;
  parentID?: string;
  tenantID?: string;
  merchantID?: string;
  status?: number;
}

/** 商城商品分类创建参数 */
export interface ShopCategoryCreateParams {
  parentID?: string;
  name: string;
  icon?: string;
  sort?: number;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}

/** 商城商品分类更新参数 */
export interface ShopCategoryUpdateParams {
  id: string;
  parentID?: string;
  name: string;
  icon?: string;
  sort?: number;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}
