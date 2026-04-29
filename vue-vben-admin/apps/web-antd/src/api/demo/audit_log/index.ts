import { requestClient } from '#/api/request';

import type {
  AuditLogItem,
  AuditLogListParams,
  AuditLogCreateParams,
  AuditLogUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/demo/audit_log';

/** 获取体验审计日志列表 */
export function getAuditLogList(params: AuditLogListParams) {
  return requestClient.get<{ list: AuditLogItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取体验审计日志详情 */
export function getAuditLogDetail(id: string) {
  return requestClient.get<AuditLogItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建体验审计日志 */
export function createAuditLog(data: AuditLogCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新体验审计日志 */
export function updateAuditLog(data: AuditLogUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除体验审计日志 */
export function deleteAuditLog(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除体验审计日志 */
export function batchDeleteAuditLog(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出体验审计日志 */
export function exportAuditLog(params?: Partial<AuditLogListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 导入体验审计日志 */
export function importAuditLog(data: FormData) {
  return requestClient.post<{ success: number; fail: number }>(
    `${PREFIX}/import`,
    data,
  );
}

/** 下载体验审计日志导入模板 */
export function downloadImportTemplateAuditLog() {
  return requestClient.get(`${PREFIX}/import-template`, {
    responseType: 'blob',
  });
}
