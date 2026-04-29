import { requestClient } from '#/api/request';

import type {
  TenantCreateParams,
  TenantItem,
  TenantListParams,
  TenantUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/system/tenant';

/** 获取租户列表 */
export function getTenantList(params: TenantListParams) {
  return requestClient.get<{ list: TenantItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取租户详情 */
export function getTenantDetail(id: string) {
  return requestClient.get<TenantItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建租户 */
export function createTenant(data: TenantCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新租户 */
export function updateTenant(data: TenantUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除租户 */
export function deleteTenant(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除租户 */
export function batchDeleteTenant(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}
