import { requestClient } from '#/api/request';

import type {
  WorkOrderItem,
  WorkOrderListParams,
  WorkOrderCreateParams,
  WorkOrderUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/demo/work_order';

/** 获取体验工单列表 */
export function getWorkOrderList(params: WorkOrderListParams) {
  return requestClient.get<{ list: WorkOrderItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取体验工单详情 */
export function getWorkOrderDetail(id: string) {
  return requestClient.get<WorkOrderItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建体验工单 */
export function createWorkOrder(data: WorkOrderCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新体验工单 */
export function updateWorkOrder(data: WorkOrderUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除体验工单 */
export function deleteWorkOrder(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除体验工单 */
export function batchDeleteWorkOrder(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出体验工单 */
export function exportWorkOrder(params?: Partial<WorkOrderListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 导入体验工单 */
export function importWorkOrder(data: FormData) {
  return requestClient.post<{ success: number; fail: number }>(
    `${PREFIX}/import`,
    data,
  );
}

/** 下载体验工单导入模板 */
export function downloadImportTemplateWorkOrder() {
  return requestClient.get(`${PREFIX}/import-template`, {
    responseType: 'blob',
  });
}

/** 批量编辑体验工单 */
export function batchUpdateWorkOrder(data: { ids: string[]; priority?: number; sourceType?: number; status?: number; }) {
  return requestClient.put(`${PREFIX}/batch-update`, data);
}
