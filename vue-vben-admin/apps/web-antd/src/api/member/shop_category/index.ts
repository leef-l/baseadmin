import { requestClient } from '#/api/request';

import type {
  ShopCategoryItem,
  ShopCategoryListParams,
  ShopCategoryCreateParams,
  ShopCategoryUpdateParams,
  ShopCategoryTreeParams,
} from './types';

/** API 前缀 */
const PREFIX = '/member/shop_category';

/** 获取商城商品分类列表 */
export function getShopCategoryList(params: ShopCategoryListParams) {
  return requestClient.get<{ list: ShopCategoryItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取商城商品分类详情 */
export function getShopCategoryDetail(id: string) {
  return requestClient.get<ShopCategoryItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建商城商品分类 */
export function createShopCategory(data: ShopCategoryCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新商城商品分类 */
export function updateShopCategory(data: ShopCategoryUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除商城商品分类 */
export function deleteShopCategory(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除商城商品分类 */
export function batchDeleteShopCategory(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出商城商品分类 */
export function exportShopCategory(params?: Partial<ShopCategoryListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 获取商城商品分类树形结构 */
export async function getShopCategoryTree(params?: ShopCategoryTreeParams) {
  const res = await requestClient.get<{ list: ShopCategoryItem[] }>(`${PREFIX}/tree`, { params });
  return res?.list ?? [];
}

/** 批量编辑商城商品分类 */
export function batchUpdateShopCategory(data: { ids: string[]; status?: number; }) {
  return requestClient.put(`${PREFIX}/batch-update`, data);
}
