/** 体验审计日志类型定义 */

/** 体验审计日志项 */
export interface AuditLogItem {
  id: string;
  logNo: string;
  operatorID?: string;
  usersUsername?: string;
  action?: number;
  targetType?: number;
  targetCode?: string;
  requestJSON?: string;
  result?: number;
  clientIP?: string;
  occurredAt?: string;
  remark?: string;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 体验审计日志列表查询参数 */
export interface AuditLogListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  keyword?: string;
  logNo?: string;
  targetCode?: string;
  operatorID?: string;
  tenantID?: string;
  merchantID?: string;
  clientIP?: string;
  action?: number;
  targetType?: number;
  result?: number;
  occurredAtStart?: string;
  occurredAtEnd?: string;
}

/** 体验审计日志创建参数 */
export interface AuditLogCreateParams {
  logNo: string;
  operatorID?: string;
  action?: number;
  targetType?: number;
  targetCode?: string;
  requestJSON?: string;
  result?: number;
  clientIP?: string;
  occurredAt?: string;
  remark?: string;
  tenantID?: string;
  merchantID?: string;
}

/** 体验审计日志更新参数 */
export interface AuditLogUpdateParams {
  id: string;
  logNo: string;
  operatorID?: string;
  action?: number;
  targetType?: number;
  targetCode?: string;
  requestJSON?: string;
  result?: number;
  clientIP?: string;
  occurredAt?: string;
  remark?: string;
  tenantID?: string;
  merchantID?: string;
}
