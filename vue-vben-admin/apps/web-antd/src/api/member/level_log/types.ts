/** 等级变更日志类型定义 */

/** 等级变更日志项 */
export interface LevelLogItem {
  id: string;
  userID?: string;
  userNickname?: string;
  oldLevelID?: string;
  levelName?: string;
  newLevelID?: string;
  newLevelName?: string;
  changeType?: number;
  expireAt?: string;
  remark?: string;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 等级变更日志列表查询参数 */
export interface LevelLogListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  userID?: string;
  oldLevelID?: string;
  newLevelID?: string;
  tenantID?: string;
  merchantID?: string;
  changeType?: number;
  expireAtStart?: string;
  expireAtEnd?: string;
}

/** 等级变更日志创建参数 */
export interface LevelLogCreateParams {
  userID?: string;
  oldLevelID?: string;
  newLevelID?: string;
  changeType?: number;
  expireAt?: string;
  remark?: string;
  tenantID?: string;
  merchantID?: string;
}

/** 等级变更日志更新参数 */
export interface LevelLogUpdateParams {
  id: string;
  userID?: string;
  oldLevelID?: string;
  newLevelID?: string;
  changeType?: number;
  expireAt?: string;
  remark?: string;
  tenantID?: string;
  merchantID?: string;
}
