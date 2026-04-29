import { requestClient } from '#/api/request';

import type {
  AppointmentItem,
  AppointmentListParams,
  AppointmentCreateParams,
  AppointmentUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/demo/appointment';

/** 获取体验预约列表 */
export function getAppointmentList(params: AppointmentListParams) {
  return requestClient.get<{ list: AppointmentItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取体验预约详情 */
export function getAppointmentDetail(id: string) {
  return requestClient.get<AppointmentItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建体验预约 */
export function createAppointment(data: AppointmentCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新体验预约 */
export function updateAppointment(data: AppointmentUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除体验预约 */
export function deleteAppointment(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除体验预约 */
export function batchDeleteAppointment(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出体验预约 */
export function exportAppointment(params?: Partial<AppointmentListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 导入体验预约 */
export function importAppointment(data: FormData) {
  return requestClient.post<{ success: number; fail: number }>(
    `${PREFIX}/import`,
    data,
  );
}

/** 下载体验预约导入模板 */
export function downloadImportTemplateAppointment() {
  return requestClient.get(`${PREFIX}/import-template`, {
    responseType: 'blob',
  });
}

/** 批量编辑体验预约 */
export function batchUpdateAppointment(data: { ids: string[]; status?: number; }) {
  return requestClient.put(`${PREFIX}/batch-update`, data);
}
