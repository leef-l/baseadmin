/** 换绑上级日志类型定义 */

/** 换绑上级日志项 */
export interface RebindLogItem {
  id: string;
  userID?: string;
  userNickname?: string;
  oldParentID?: string;
  oldParentNickname?: string;
  newParentID?: string;
  newParentNickname?: string;
  reason?: string;
  operatorID?: string;
  usersUsername?: string;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 换绑上级日志列表查询参数 */
export interface RebindLogListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  userID?: string;
  oldParentID?: string;
  newParentID?: string;
  operatorID?: string;
  tenantID?: string;
  merchantID?: string;
}

/** 换绑上级日志创建参数 */
export interface RebindLogCreateParams {
  userID?: string;
  oldParentID?: string;
  newParentID?: string;
  reason?: string;
  operatorID?: string;
  tenantID?: string;
  merchantID?: string;
}

/** 换绑上级日志更新参数 */
export interface RebindLogUpdateParams {
  id: string;
  userID?: string;
  oldParentID?: string;
  newParentID?: string;
  reason?: string;
  operatorID?: string;
  tenantID?: string;
  merchantID?: string;
}
