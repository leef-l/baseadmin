import { requestClient } from '#/api/request';

import type { ContractItem, ContractListParams } from './types';

const PREFIX = '/member/contract';

/** 获取合同列表 */
export function getContractList(params: ContractListParams) {
  return requestClient.get<{ list: ContractItem[]; total: number }>(
    `${PREFIX}/list`,
    { params },
  );
}

/** 下载合同 URL（浏览器跳转，不走 axios） */
export function getContractDownloadURL(contractId: string) {
  return `/api/member/contract/download?contractId=${encodeURIComponent(contractId)}`;
}
