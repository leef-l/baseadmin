import { requestClient } from '#/api/request';

import type {
  WarehouseGoodsItem,
  WarehouseGoodsListParams,
  WarehouseGoodsCreateParams,
  WarehouseGoodsUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/member/warehouse_goods';

/** 获取仓库商品列表 */
export function getWarehouseGoodsList(params: WarehouseGoodsListParams) {
  return requestClient.get<{ list: WarehouseGoodsItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取仓库商品详情 */
export function getWarehouseGoodsDetail(id: string) {
  return requestClient.get<WarehouseGoodsItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建仓库商品 */
export function createWarehouseGoods(data: WarehouseGoodsCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新仓库商品 */
export function updateWarehouseGoods(data: WarehouseGoodsUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除仓库商品 */
export function deleteWarehouseGoods(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除仓库商品 */
export function batchDeleteWarehouseGoods(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出仓库商品 */
export function exportWarehouseGoods(params?: Partial<WarehouseGoodsListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 导入仓库商品 */
export function importWarehouseGoods(data: FormData) {
  return requestClient.post<{ success: number; fail: number }>(
    `${PREFIX}/import`,
    data,
  );
}

/** 下载仓库商品导入模板 */
export function downloadImportTemplateWarehouseGoods() {
  return requestClient.get(`${PREFIX}/import-template`, {
    responseType: 'blob',
  });
}

/** 批量编辑仓库商品 */
export function batchUpdateWarehouseGoods(data: { ids: string[]; goodsStatus?: number; status?: number; }) {
  return requestClient.put(`${PREFIX}/batch-update`, data);
}
