import { requestClient } from '#/api/request';

import type {
  TeamExportItem,
  TeamExportListParams,
  TeamExportCreateParams,
  TeamExportUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/member/team_export';

/** 获取团队数据导出列表 */
export function getTeamExportList(params: TeamExportListParams) {
  return requestClient.get<{ list: TeamExportItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取团队数据导出详情 */
export function getTeamExportDetail(id: string) {
  return requestClient.get<TeamExportItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建团队数据导出 */
export function createTeamExport(data: TeamExportCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新团队数据导出 */
export function updateTeamExport(data: TeamExportUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除团队数据导出 */
export function deleteTeamExport(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除团队数据导出 */
export function batchDeleteTeamExport(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出团队数据导出 */
export function exportTeamExport(params?: Partial<TeamExportListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 导入团队数据导出 */
export function importTeamExport(data: FormData) {
  return requestClient.post<{ success: number; fail: number }>(
    `${PREFIX}/import`,
    data,
  );
}

/** 下载团队数据导出导入模板 */
export function downloadImportTemplateTeamExport() {
  return requestClient.get(`${PREFIX}/import-template`, {
    responseType: 'blob',
  });
}

/** 批量编辑团队数据导出 */
export function batchUpdateTeamExport(data: { ids: string[]; exportType?: number; deployStatus?: number; status?: number; }) {
  return requestClient.put(`${PREFIX}/batch-update`, data);
}
