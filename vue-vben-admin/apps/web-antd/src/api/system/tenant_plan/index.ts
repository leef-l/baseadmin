import { requestClient } from '#/api/request';

import type {
  TenantPlanItem,
  TenantPlanListParams,
  TenantPlanCreateParams,
  TenantPlanUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/system/tenant_plan';

/** 获取租户套餐订阅表列表 */
export function getTenantPlanList(params: TenantPlanListParams) {
  return requestClient.get<{ list: TenantPlanItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取租户套餐订阅表详情 */
export function getTenantPlanDetail(id: string) {
  return requestClient.get<TenantPlanItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建租户套餐订阅表 */
export function createTenantPlan(data: TenantPlanCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新租户套餐订阅表 */
export function updateTenantPlan(data: TenantPlanUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除租户套餐订阅表 */
export function deleteTenantPlan(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除租户套餐订阅表 */
export function batchDeleteTenantPlan(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出租户套餐订阅表 */
export function exportTenantPlan(params?: Partial<TenantPlanListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 导入租户套餐订阅表 */
export function importTenantPlan(data: FormData) {
  return requestClient.post<{ success: number; fail: number }>(
    `${PREFIX}/import`,
    data,
  );
}

/** 下载租户套餐订阅表导入模板 */
export function downloadImportTemplateTenantPlan() {
  return requestClient.get(`${PREFIX}/import-template`, {
    responseType: 'blob',
  });
}

/** 批量编辑租户套餐订阅表 */
export function batchUpdateTenantPlan(data: { ids: string[]; status?: number; }) {
  return requestClient.put(`${PREFIX}/batch-update`, data);
}
