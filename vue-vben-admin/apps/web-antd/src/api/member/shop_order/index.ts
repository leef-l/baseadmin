import { requestClient } from '#/api/request';

import type {
  ShopOrderItem,
  ShopOrderListParams,
  ShopOrderCreateParams,
  ShopOrderUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/member/shop_order';

/** 获取商城订单列表 */
export function getShopOrderList(params: ShopOrderListParams) {
  return requestClient.get<{ list: ShopOrderItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取商城订单详情 */
export function getShopOrderDetail(id: string) {
  return requestClient.get<ShopOrderItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建商城订单 */
export function createShopOrder(data: ShopOrderCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新商城订单 */
export function updateShopOrder(data: ShopOrderUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除商城订单 */
export function deleteShopOrder(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除商城订单 */
export function batchDeleteShopOrder(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出商城订单 */
export function exportShopOrder(params?: Partial<ShopOrderListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 导入商城订单 */
export function importShopOrder(data: FormData) {
  return requestClient.post<{ success: number; fail: number }>(
    `${PREFIX}/import`,
    data,
  );
}

/** 下载商城订单导入模板 */
export function downloadImportTemplateShopOrder() {
  return requestClient.get(`${PREFIX}/import-template`, {
    responseType: 'blob',
  });
}

/** 批量编辑商城订单 */
export function batchUpdateShopOrder(data: { ids: string[]; payWallet?: number; orderStatus?: number; status?: number; }) {
  return requestClient.put(`${PREFIX}/batch-update`, data);
}
