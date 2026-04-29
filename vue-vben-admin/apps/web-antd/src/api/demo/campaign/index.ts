import { requestClient } from '#/api/request';

import type {
  CampaignItem,
  CampaignListParams,
  CampaignCreateParams,
  CampaignUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/demo/campaign';

/** 获取体验活动列表 */
export function getCampaignList(params: CampaignListParams) {
  return requestClient.get<{ list: CampaignItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取体验活动详情 */
export function getCampaignDetail(id: string) {
  return requestClient.get<CampaignItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建体验活动 */
export function createCampaign(data: CampaignCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新体验活动 */
export function updateCampaign(data: CampaignUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除体验活动 */
export function deleteCampaign(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除体验活动 */
export function batchDeleteCampaign(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出体验活动 */
export function exportCampaign(params?: Partial<CampaignListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 导入体验活动 */
export function importCampaign(data: FormData) {
  return requestClient.post<{ success: number; fail: number }>(
    `${PREFIX}/import`,
    data,
  );
}

/** 下载体验活动导入模板 */
export function downloadImportTemplateCampaign() {
  return requestClient.get(`${PREFIX}/import-template`, {
    responseType: 'blob',
  });
}

/** 批量编辑体验活动 */
export function batchUpdateCampaign(data: { ids: string[]; type?: number; channel?: number; isPublic?: number; status?: number; }) {
  return requestClient.put(`${PREFIX}/batch-update`, data);
}
