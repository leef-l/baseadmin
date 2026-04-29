import { requestClient } from '#/api/request';

import type {
  ShopGoodsItem,
  ShopGoodsListParams,
  ShopGoodsCreateParams,
  ShopGoodsUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/member/shop_goods';

/** 获取商城商品列表 */
export function getShopGoodsList(params: ShopGoodsListParams) {
  return requestClient.get<{ list: ShopGoodsItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取商城商品详情 */
export function getShopGoodsDetail(id: string) {
  return requestClient.get<ShopGoodsItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建商城商品 */
export function createShopGoods(data: ShopGoodsCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新商城商品 */
export function updateShopGoods(data: ShopGoodsUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除商城商品 */
export function deleteShopGoods(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除商城商品 */
export function batchDeleteShopGoods(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出商城商品 */
export function exportShopGoods(params?: Partial<ShopGoodsListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 导入商城商品 */
export function importShopGoods(data: FormData) {
  return requestClient.post<{ success: number; fail: number }>(
    `${PREFIX}/import`,
    data,
  );
}

/** 下载商城商品导入模板 */
export function downloadImportTemplateShopGoods() {
  return requestClient.get(`${PREFIX}/import-template`, {
    responseType: 'blob',
  });
}

/** 批量编辑商城商品 */
export function batchUpdateShopGoods(data: { ids: string[]; isRecommend?: number; status?: number; }) {
  return requestClient.put(`${PREFIX}/batch-update`, data);
}
