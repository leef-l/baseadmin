import { requestClient } from '#/api/request';

import type {
  LevelLogItem,
  LevelLogListParams,
  LevelLogCreateParams,
  LevelLogUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/member/level_log';

/** 获取等级变更日志列表 */
export function getLevelLogList(params: LevelLogListParams) {
  return requestClient.get<{ list: LevelLogItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取等级变更日志详情 */
export function getLevelLogDetail(id: string) {
  return requestClient.get<LevelLogItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建等级变更日志 */
export function createLevelLog(data: LevelLogCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新等级变更日志 */
export function updateLevelLog(data: LevelLogUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除等级变更日志 */
export function deleteLevelLog(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除等级变更日志 */
export function batchDeleteLevelLog(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出等级变更日志 */
export function exportLevelLog(params?: Partial<LevelLogListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 导入等级变更日志 */
export function importLevelLog(data: FormData) {
  return requestClient.post<{ success: number; fail: number }>(
    `${PREFIX}/import`,
    data,
  );
}

/** 下载等级变更日志导入模板 */
export function downloadImportTemplateLevelLog() {
  return requestClient.get(`${PREFIX}/import-template`, {
    responseType: 'blob',
  });
}
