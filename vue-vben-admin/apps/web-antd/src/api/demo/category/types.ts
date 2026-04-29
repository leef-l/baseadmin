/** 体验分类类型定义 */

/** 体验分类项 */
export interface CategoryItem {
  id: string;
  parentID?: string;
  categoryName?: string;
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
  children?: CategoryItem[];
}

/** 体验分类列表查询参数 */
export interface CategoryListParams {
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
/** 体验分类树形查询参数 */
export interface CategoryTreeParams {
  startTime?: string;
  endTime?: string;
  name?: string;
  parentID?: string;
  tenantID?: string;
  merchantID?: string;
  status?: number;
}

/** 体验分类创建参数 */
export interface CategoryCreateParams {
  parentID?: string;
  name: string;
  icon?: string;
  sort?: number;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}

/** 体验分类更新参数 */
export interface CategoryUpdateParams {
  id: string;
  parentID?: string;
  name: string;
  icon?: string;
  sort?: number;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}
