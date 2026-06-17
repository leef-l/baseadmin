import { requestClient } from '#/api/request';

import type {
  CronItem,
  CronListParams,
  CronCreateParams,
  CronUpdateParams,
  CronTriggerResult,
  CronLogItem,
  CronLogListParams,
} from './types';

/** API 前缀 */
const PREFIX = '/system/cron';

/** 获取定时任务表列表 */
export function getCronList(params: CronListParams) {
  return requestClient.get<{ list: CronItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取定时任务表详情 */
export function getCronDetail(id: string) {
  return requestClient.get<CronItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建定时任务表 */
export function createCron(data: CronCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新定时任务表 */
export function updateCron(data: CronUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除定时任务表 */
export function deleteCron(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除定时任务表 */
export function batchDeleteCron(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出定时任务表 */
export function exportCron(params?: Partial<CronListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 导入定时任务表 */
export function importCron(data: FormData) {
  return requestClient.post<{ success: number; fail: number }>(
    `${PREFIX}/import`,
    data,
  );
}

/** 下载定时任务表导入模板 */
export function downloadImportTemplateCron() {
  return requestClient.get(`${PREFIX}/import-template`, {
    responseType: 'blob',
  });
}

/** 批量编辑定时任务表 */
export function batchUpdateCron(data: { ids: string[]; status?: number; }) {
  return requestClient.put(`${PREFIX}/batch-update`, data);
}

/** 手动触发定时任务 */
export function triggerCron(name: string) {
  return requestClient.post<CronTriggerResult>(`${PREFIX}/trigger`, { name });
}

/** 查询执行日志 */
export function getCronLogList(params: CronLogListParams) {
  return requestClient.get<{ list: CronLogItem[]; total: number }>(`${PREFIX}/log-list`, { params });
}
