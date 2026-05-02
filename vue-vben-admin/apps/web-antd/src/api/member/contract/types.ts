export interface ContractItem {
  contractId: string;
  contractNo: string;
  userId: string;
  userNickname: string;
  userPhone: string;
  contractType: string;
  templateId: string;
  signedAt: string;
  signedIp: string;
  pdfStatus: number;
  pdfStatusText: string;
  createdAt: string;
}

export interface ContractListParams {
  userId?: string;
  contractType?: string;
  pdfStatus?: number;
  pageNum?: number;
  pageSize?: number;
}
