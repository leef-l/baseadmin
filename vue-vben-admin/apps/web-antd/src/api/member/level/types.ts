/** 会员等级配置类型定义 */

/** 会员等级配置项 */
export interface LevelItem {
  id: string;
  name: string;
  levelNo?: number;
  icon?: string;
  durationDays?: number;
  needActiveCount?: number;
  needTeamTurnover?: string;
  isTop?: number;
  autoDeploy?: number;
  remark?: string;
  sort?: number;
  status?: number;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 会员等级配置列表查询参数 */
export interface LevelListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  name?: string;
  levelNo?: string;
  tenantID?: string;
  merchantID?: string;
  isTop?: number;
  autoDeploy?: number;
  status?: number;
}

/** 会员等级配置创建参数 */
export interface LevelCreateParams {
  name: string;
  levelNo?: number;
  icon?: string;
  durationDays?: number;
  needActiveCount?: number;
  needTeamTurnover?: string;
  isTop?: number;
  autoDeploy?: number;
  remark?: string;
  sort?: number;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}

/** 会员等级配置更新参数 */
export interface LevelUpdateParams {
  id: string;
  name: string;
  levelNo?: number;
  icon?: string;
  durationDays?: number;
  needActiveCount?: number;
  needTeamTurnover?: string;
  isTop?: number;
  autoDeploy?: number;
  remark?: string;
  sort?: number;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}
