import { requestClient } from '#/api/request';

import type {
  LevelItem,
  LevelListParams,
  LevelCreateParams,
  LevelUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/member/level';

/** 获取会员等级配置列表 */
export function getLevelList(params: LevelListParams) {
  return requestClient.get<{ list: LevelItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取会员等级配置详情 */
export function getLevelDetail(id: string) {
  return requestClient.get<LevelItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建会员等级配置 */
export function createLevel(data: LevelCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新会员等级配置 */
export function updateLevel(data: LevelUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除会员等级配置 */
export function deleteLevel(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除会员等级配置 */
export function batchDeleteLevel(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出会员等级配置 */
export function exportLevel(params?: Partial<LevelListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 导入会员等级配置 */
export function importLevel(data: FormData) {
  return requestClient.post<{ success: number; fail: number }>(
    `${PREFIX}/import`,
    data,
  );
}

/** 下载会员等级配置导入模板 */
export function downloadImportTemplateLevel() {
  return requestClient.get(`${PREFIX}/import-template`, {
    responseType: 'blob',
  });
}

/** 批量编辑会员等级配置 */
export function batchUpdateLevel(data: { ids: string[]; isTop?: number; autoDeploy?: number; status?: number; }) {
  return requestClient.put(`${PREFIX}/batch-update`, data);
}
