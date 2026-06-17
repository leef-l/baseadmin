/** 定时任务表类型定义 */

/** 定时任务表项 */
export interface CronItem {
  id: string;
  name: string;
  expression: string;
  handler: string;
  params?: string;
  remark?: string;
  status?: number;
  lastRunAt?: string;
  lastResult?: string;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 定时任务表列表查询参数 */
export interface CronListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  keyword?: string;
  name?: string;
  tenantID?: string;
  merchantID?: string;
  status?: number;
  remark?: string;
  lastRunAtStart?: string;
  lastRunAtEnd?: string;
}

/** 定时任务表创建参数 */
export interface CronCreateParams {
  name: string;
  expression: string;
  handler: string;
  params?: string;
  remark?: string;
  status?: number;
  lastRunAt?: string;
  lastResult?: string;
  tenantID?: string;
  merchantID?: string;
}

/** 定时任务表更新参数 */
export interface CronUpdateParams {
  id: string;
  name: string;
  expression: string;
  handler: string;
  params?: string;
  remark?: string;
  status?: number;
  lastRunAt?: string;
  lastResult?: string;
  tenantID?: string;
  merchantID?: string;
}

/** 手动触发参数 */
export interface CronTriggerParams {
  name: string;
}

/** 手动触发响应 */
export interface CronTriggerResult {
  success: boolean;
  message: string;
}

/** 执行日志项 */
export interface CronLogItem {
  id: number;
  cronId: number;
  cronName: string;
  startAt: string;
  endAt?: string;
  durationMs: number;
  result: string;
  message: string;
}

/** 日志查询参数 */
export interface CronLogListParams {
  cronId: string; // Snowflake ID 超过 JS 安全整数范围，必须用 string
  pageNum: number;
  pageSize: number;
}
