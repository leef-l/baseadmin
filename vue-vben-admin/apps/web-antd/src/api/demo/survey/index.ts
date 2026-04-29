import { requestClient } from '#/api/request';

import type {
  SurveyItem,
  SurveyListParams,
  SurveyCreateParams,
  SurveyUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/demo/survey';

/** 获取体验问卷列表 */
export function getSurveyList(params: SurveyListParams) {
  return requestClient.get<{ list: SurveyItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取体验问卷详情 */
export function getSurveyDetail(id: string) {
  return requestClient.get<SurveyItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建体验问卷 */
export function createSurvey(data: SurveyCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新体验问卷 */
export function updateSurvey(data: SurveyUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除体验问卷 */
export function deleteSurvey(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除体验问卷 */
export function batchDeleteSurvey(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出体验问卷 */
export function exportSurvey(params?: Partial<SurveyListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 导入体验问卷 */
export function importSurvey(data: FormData) {
  return requestClient.post<{ success: number; fail: number }>(
    `${PREFIX}/import`,
    data,
  );
}

/** 下载体验问卷导入模板 */
export function downloadImportTemplateSurvey() {
  return requestClient.get(`${PREFIX}/import-template`, {
    responseType: 'blob',
  });
}

/** 批量编辑体验问卷 */
export function batchUpdateSurvey(data: { ids: string[]; isAnonymous?: number; status?: number; }) {
  return requestClient.put(`${PREFIX}/batch-update`, data);
}
