import req from './request';
import { PageResult } from './types';

export interface ContractTemplate {
  templateId: string;
  templateName: string;
  content: string;
}

export interface ContractItem {
  contractId: string;
  contractNo: string;
  contractType: string;
  contractTypeText: string;
  signedAt: string;
  pdfStatus: number;
  pdfStatusText: string;
}

export interface SignResult {
  contractId: string;
  contractNo: string;
}

export type ContractType = 'register' | 'upgrade' | 'custom';

export const contractApi = {
  template: (contractType: ContractType = 'register') =>
    req.get<any, ContractTemplate>('/member-portal/contract/template', { params: { contractType } }),

  status: (contractType: ContractType = 'register') =>
    req.get<any, { hasSign: boolean }>('/member-portal/contract/status', { params: { contractType } }),

  sign: (data: {
    contractType?: ContractType;
    templateId?: string;
    signatureImage: string;
    relatedId?: string;
  }) => req.post<any, SignResult>('/member-portal/contract/sign', data),

  list: (params: { pageNum?: number; pageSize?: number }) =>
    req.get<any, PageResult<ContractItem>>('/member-portal/contract/list', { params }),

  // 下载是浏览器跳转，不走 axios
  downloadURL: (contractId: string) =>
    `/api/member-portal/contract/download?contractId=${encodeURIComponent(contractId)}`,
};
