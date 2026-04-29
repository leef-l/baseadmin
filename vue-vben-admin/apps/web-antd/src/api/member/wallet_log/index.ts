import { requestClient } from '#/api/request';

import type {
  WalletLogItem,
  WalletLogListParams,
  WalletLogCreateParams,
  WalletLogUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/member/wallet_log';

/** 获取钱包流水记录列表 */
export function getWalletLogList(params: WalletLogListParams) {
  return requestClient.get<{ list: WalletLogItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取钱包流水记录详情 */
export function getWalletLogDetail(id: string) {
  return requestClient.get<WalletLogItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建钱包流水记录 */
export function createWalletLog(data: WalletLogCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新钱包流水记录 */
export function updateWalletLog(data: WalletLogUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除钱包流水记录 */
export function deleteWalletLog(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除钱包流水记录 */
export function batchDeleteWalletLog(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出钱包流水记录 */
export function exportWalletLog(params?: Partial<WalletLogListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 导入钱包流水记录 */
export function importWalletLog(data: FormData) {
  return requestClient.post<{ success: number; fail: number }>(
    `${PREFIX}/import`,
    data,
  );
}

/** 下载钱包流水记录导入模板 */
export function downloadImportTemplateWalletLog() {
  return requestClient.get(`${PREFIX}/import-template`, {
    responseType: 'blob',
  });
}
