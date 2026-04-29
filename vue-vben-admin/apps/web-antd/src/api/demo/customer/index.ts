import { requestClient } from '#/api/request';

import type {
  CustomerItem,
  CustomerListParams,
  CustomerCreateParams,
  CustomerUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/demo/customer';

/** 获取体验客户列表 */
export function getCustomerList(params: CustomerListParams) {
  return requestClient.get<{ list: CustomerItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取体验客户详情 */
export function getCustomerDetail(id: string) {
  return requestClient.get<CustomerItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建体验客户 */
export function createCustomer(data: CustomerCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新体验客户 */
export function updateCustomer(data: CustomerUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除体验客户 */
export function deleteCustomer(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除体验客户 */
export function batchDeleteCustomer(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出体验客户 */
export function exportCustomer(params?: Partial<CustomerListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 导入体验客户 */
export function importCustomer(data: FormData) {
  return requestClient.post<{ success: number; fail: number }>(
    `${PREFIX}/import`,
    data,
  );
}

/** 下载体验客户导入模板 */
export function downloadImportTemplateCustomer() {
  return requestClient.get(`${PREFIX}/import-template`, {
    responseType: 'blob',
  });
}

/** 批量编辑体验客户 */
export function batchUpdateCustomer(data: { ids: string[]; gender?: number; level?: number; sourceType?: number; isVip?: number; status?: number; }) {
  return requestClient.put(`${PREFIX}/batch-update`, data);
}
