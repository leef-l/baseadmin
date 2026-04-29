import { requestClient } from '#/api/request';

import type {
  ContractItem,
  ContractListParams,
  ContractCreateParams,
  ContractUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/demo/contract';

/** 获取体验合同列表 */
export function getContractList(params: ContractListParams) {
  return requestClient.get<{ list: ContractItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取体验合同详情 */
export function getContractDetail(id: string) {
  return requestClient.get<ContractItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建体验合同 */
export function createContract(data: ContractCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新体验合同 */
export function updateContract(data: ContractUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除体验合同 */
export function deleteContract(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除体验合同 */
export function batchDeleteContract(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出体验合同 */
export function exportContract(params?: Partial<ContractListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 导入体验合同 */
export function importContract(data: FormData) {
  return requestClient.post<{ success: number; fail: number }>(
    `${PREFIX}/import`,
    data,
  );
}

/** 下载体验合同导入模板 */
export function downloadImportTemplateContract() {
  return requestClient.get(`${PREFIX}/import-template`, {
    responseType: 'blob',
  });
}

/** 批量编辑体验合同 */
export function batchUpdateContract(data: { ids: string[]; status?: number; }) {
  return requestClient.put(`${PREFIX}/batch-update`, data);
}
