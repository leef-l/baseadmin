import { requestClient } from '#/api/request';

import type {
  OrderItem,
  OrderListParams,
  OrderCreateParams,
  OrderUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/demo/order';

/** 获取体验订单列表 */
export function getOrderList(params: OrderListParams) {
  return requestClient.get<{ list: OrderItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取体验订单详情 */
export function getOrderDetail(id: string) {
  return requestClient.get<OrderItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建体验订单 */
export function createOrder(data: OrderCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新体验订单 */
export function updateOrder(data: OrderUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除体验订单 */
export function deleteOrder(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除体验订单 */
export function batchDeleteOrder(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出体验订单 */
export function exportOrder(params?: Partial<OrderListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 导入体验订单 */
export function importOrder(data: FormData) {
  return requestClient.post<{ success: number; fail: number }>(
    `${PREFIX}/import`,
    data,
  );
}

/** 下载体验订单导入模板 */
export function downloadImportTemplateOrder() {
  return requestClient.get(`${PREFIX}/import-template`, {
    responseType: 'blob',
  });
}

/** 批量编辑体验订单 */
export function batchUpdateOrder(data: { ids: string[]; payStatus?: number; deliverStatus?: number; status?: number; }) {
  return requestClient.put(`${PREFIX}/batch-update`, data);
}
