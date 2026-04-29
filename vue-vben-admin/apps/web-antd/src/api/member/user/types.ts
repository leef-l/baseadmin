/** 会员用户类型定义 */

/** 会员用户项 */
export interface UserItem {
  id: string;
  parentID?: string;
  userUsername?: string;
  username: string;
  nickname?: string;
  phone?: string;
  avatar?: string;
  realName?: string;
  levelID?: string;
  levelName?: string;
  levelExpireAt?: string;
  teamCount?: number;
  directCount?: number;
  activeCount?: number;
  teamTurnover?: string;
  isActive?: number;
  isQualified?: number;
  inviteCode?: string;
  registerIP?: string;
  lastLoginAt?: string;
  remark?: string;
  sort?: number;
  status?: number;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
  children?: UserItem[];
}

/** 会员用户列表查询参数 */
export interface UserListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  keyword?: string;
  username?: string;
  inviteCode?: string;
  nickname?: string;
  realName?: string;
  phone?: string;
  parentID?: string;
  levelID?: string;
  tenantID?: string;
  merchantID?: string;
  isActive?: number;
  isQualified?: number;
  status?: number;
  levelExpireAtStart?: string;
  levelExpireAtEnd?: string;
  lastLoginAtStart?: string;
  lastLoginAtEnd?: string;
}
/** 会员用户树形查询参数 */
export interface UserTreeParams {
  startTime?: string;
  endTime?: string;
  keyword?: string;
  username?: string;
  inviteCode?: string;
  nickname?: string;
  realName?: string;
  phone?: string;
  parentID?: string;
  levelID?: string;
  tenantID?: string;
  merchantID?: string;
  isActive?: number;
  isQualified?: number;
  status?: number;
  levelExpireAtStart?: string;
  levelExpireAtEnd?: string;
  lastLoginAtStart?: string;
  lastLoginAtEnd?: string;
}

/** 会员用户创建参数 */
export interface UserCreateParams {
  parentID?: string;
  username: string;
  password?: string;
  nickname?: string;
  phone?: string;
  avatar?: string;
  realName?: string;
  levelID?: string;
  levelExpireAt?: string;
  teamCount?: number;
  directCount?: number;
  activeCount?: number;
  teamTurnover?: string;
  isActive?: number;
  isQualified?: number;
  inviteCode?: string;
  registerIP?: string;
  lastLoginAt?: string;
  remark?: string;
  sort?: number;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}

/** 会员用户更新参数 */
export interface UserUpdateParams {
  id: string;
  parentID?: string;
  username: string;
  password?: string;
  nickname?: string;
  phone?: string;
  avatar?: string;
  realName?: string;
  levelID?: string;
  levelExpireAt?: string;
  teamCount?: number;
  directCount?: number;
  activeCount?: number;
  teamTurnover?: string;
  isActive?: number;
  isQualified?: number;
  inviteCode?: string;
  registerIP?: string;
  lastLoginAt?: string;
  remark?: string;
  sort?: number;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}
