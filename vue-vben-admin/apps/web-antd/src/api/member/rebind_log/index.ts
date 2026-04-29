import { requestClient } from '#/api/request';

import type {
  RebindLogItem,
  RebindLogListParams,
  RebindLogCreateParams,
  RebindLogUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/member/rebind_log';

/** 获取换绑上级日志列表 */
export function getRebindLogList(params: RebindLogListParams) {
  return requestClient.get<{ list: RebindLogItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取换绑上级日志详情 */
export function getRebindLogDetail(id: string) {
  return requestClient.get<RebindLogItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建换绑上级日志 */
export function createRebindLog(data: RebindLogCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新换绑上级日志 */
export function updateRebindLog(data: RebindLogUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除换绑上级日志 */
export function deleteRebindLog(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除换绑上级日志 */
export function batchDeleteRebindLog(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出换绑上级日志 */
export function exportRebindLog(params?: Partial<RebindLogListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 导入换绑上级日志 */
export function importRebindLog(data: FormData) {
  return requestClient.post<{ success: number; fail: number }>(
    `${PREFIX}/import`,
    data,
  );
}

/** 下载换绑上级日志导入模板 */
export function downloadImportTemplateRebindLog() {
  return requestClient.get(`${PREFIX}/import-template`, {
    responseType: 'blob',
  });
}
