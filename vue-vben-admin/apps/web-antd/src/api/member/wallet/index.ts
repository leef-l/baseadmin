import { requestClient } from '#/api/request';

import type {
  WalletItem,
  WalletListParams,
  WalletCreateParams,
  WalletUpdateParams,
} from './types';

/** API 前缀 */
const PREFIX = '/member/wallet';

/** 获取会员钱包列表 */
export function getWalletList(params: WalletListParams) {
  return requestClient.get<{ list: WalletItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 获取会员钱包详情 */
export function getWalletDetail(id: string) {
  return requestClient.get<WalletItem>(`${PREFIX}/detail`, {
    params: { id },
  });
}

/** 创建会员钱包 */
export function createWallet(data: WalletCreateParams) {
  return requestClient.post(`${PREFIX}/create`, data);
}

/** 更新会员钱包 */
export function updateWallet(data: WalletUpdateParams) {
  return requestClient.put(`${PREFIX}/update`, data);
}

/** 删除会员钱包 */
export function deleteWallet(id: string) {
  return requestClient.delete(`${PREFIX}/delete`, { data: { id } });
}

/** 批量删除会员钱包 */
export function batchDeleteWallet(ids: string[]) {
  return requestClient.delete(`${PREFIX}/batch-delete`, { data: { ids } });
}

/** 导出会员钱包 */
export function exportWallet(params?: Partial<WalletListParams>) {
  return requestClient.get(`${PREFIX}/export`, {
    params,
    responseType: 'blob',
  });
}

/** 导入会员钱包 */
export function importWallet(data: FormData) {
  return requestClient.post<{ success: number; fail: number }>(
    `${PREFIX}/import`,
    data,
  );
}

/** 下载会员钱包导入模板 */
export function downloadImportTemplateWallet() {
  return requestClient.get(`${PREFIX}/import-template`, {
    responseType: 'blob',
  });
}

/** 批量编辑会员钱包 */
export function batchUpdateWallet(data: { ids: string[]; walletType?: number; status?: number; }) {
  return requestClient.put(`${PREFIX}/batch-update`, data);
}
