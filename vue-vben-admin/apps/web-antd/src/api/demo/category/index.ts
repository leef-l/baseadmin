import { requestClient } from '#/api/request';

import type {
  CategoryItem,
  CategoryListParams,
  CategoryCreateParams,
  CategoryUpdateParams,
  CategoryTreeParams,
} from './types';

/** API 前缀 */
const PREFIX = '/demo/category';

/** 获取体验分类列表 */
export function getCategoryList(params: CategoryListParams) {
  return requestClient.get<{ list: CategoryItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取体验分类详情 */
export function getCategoryDetail(id: string) {
  return requestClient.get<CategoryItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建体验分类 */
export function createCategory(data: CategoryCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新体验分类 */
export function updateCategory(data: CategoryUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除体验分类 */
export function deleteCategory(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除体验分类 */
export function batchDeleteCategory(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出体验分类 */
export function exportCategory(params?: Partial<CategoryListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 获取体验分类树形结构 */
export async function getCategoryTree(params?: CategoryTreeParams) {
  const res = await requestClient.get<{ list: CategoryItem[] }>(`${PREFIX}/tree`, { params });
  return res?.list ?? [];
}

/** 批量编辑体验分类 */
export function batchUpdateCategory(data: { ids: string[]; status?: number; }) {
  return requestClient.put(`${PREFIX}/batch-update`, data);
}
