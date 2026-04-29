/** 团队数据导出类型定义 */

/** 团队数据导出项 */
export interface TeamExportItem {
  id: string;
  userID?: string;
  userNickname?: string;
  teamMemberCount?: number;
  exportType?: number;
  fileURL?: string;
  fileSize?: string;
  deployStatus?: number;
  deployDomain?: string;
  deployedAt?: string;
  remark?: string;
  status?: number;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 团队数据导出列表查询参数 */
export interface TeamExportListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  userID?: string;
  tenantID?: string;
  merchantID?: string;
  exportType?: number;
  deployStatus?: number;
  status?: number;
  deployDomain?: string;
  deployedAtStart?: string;
  deployedAtEnd?: string;
}

/** 团队数据导出创建参数 */
export interface TeamExportCreateParams {
  userID?: string;
  teamMemberCount?: number;
  exportType?: number;
  fileURL?: string;
  fileSize?: string;
  deployStatus?: number;
  deployDomain?: string;
  deployedAt?: string;
  remark?: string;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}

/** 团队数据导出更新参数 */
export interface TeamExportUpdateParams {
  id: string;
  userID?: string;
  teamMemberCount?: number;
  exportType?: number;
  fileURL?: string;
  fileSize?: string;
  deployStatus?: number;
  deployDomain?: string;
  deployedAt?: string;
  remark?: string;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}
