import { requestClient } from '#/api/request';

import type {
  ProductItem,
  ProductListParams,
  ProductCreateParams,
  ProductUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/demo/product';

/** 获取体验商品列表 */
export function getProductList(params: ProductListParams) {
  return requestClient.get<{ list: ProductItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取体验商品详情 */
export function getProductDetail(id: string) {
  return requestClient.get<ProductItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建体验商品 */
export function createProduct(data: ProductCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新体验商品 */
export function updateProduct(data: ProductUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除体验商品 */
export function deleteProduct(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除体验商品 */
export function batchDeleteProduct(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出体验商品 */
export function exportProduct(params?: Partial<ProductListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 导入体验商品 */
export function importProduct(data: FormData) {
  return requestClient.post<{ success: number; fail: number }>(
    `${PREFIX}/import`,
    data,
  );
}

/** 下载体验商品导入模板 */
export function downloadImportTemplateProduct() {
  return requestClient.get(`${PREFIX}/import-template`, {
    responseType: 'blob',
  });
}

/** 批量编辑体验商品 */
export function batchUpdateProduct(data: { ids: string[]; type?: number; isRecommend?: number; status?: number; }) {
  return requestClient.put(`${PREFIX}/batch-update`, data);
}
