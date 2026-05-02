import { requestClient } from '#/api/request';

import type {
  ContractTemplateItem,
  ContractTemplateListParams,
  ContractTemplateCreateParams,
  ContractTemplateUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/member/contract_template';

/** 获取会员合同模板列表 */
export function getContractTemplateList(params: ContractTemplateListParams) {
  return requestClient.get<{ list: ContractTemplateItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取会员合同模板详情 */
export function getContractTemplateDetail(id: string) {
  return requestClient.get<ContractTemplateItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建会员合同模板 */
export function createContractTemplate(data: ContractTemplateCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新会员合同模板 */
export function updateContractTemplate(data: ContractTemplateUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除会员合同模板 */
export function deleteContractTemplate(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除会员合同模板 */
export function batchDeleteContractTemplate(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出会员合同模板 */
export function exportContractTemplate(params?: Partial<ContractTemplateListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 导入会员合同模板 */
export function importContractTemplate(data: FormData) {
  return requestClient.post<{ success: number; fail: number }>(
    `${PREFIX}/import`,
    data,
  );
}

/** 下载会员合同模板导入模板 */
export function downloadImportTemplateContractTemplate() {
  return requestClient.get(`${PREFIX}/import-template`, {
    responseType: 'blob',
  });
}

/** 批量编辑会员合同模板 */
export function batchUpdateContractTemplate(data: { ids: string[]; isDefault?: number; status?: number; }) {
  return requestClient.put(`${PREFIX}/batch-update`, data);
}
