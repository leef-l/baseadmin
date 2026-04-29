/** 用户表类型定义 */

/** 用户表项 */
export interface UsersItem {
  id: string;
  username: string;
  nickname?: string;
  email?: string;
  avatar?: string;
  deptId?: string;
  deptTitle?: string;
  tenantId?: string;
  tenantName?: string;
  merchantId?: string;
  merchantName?: string;
  roleIds?: string[];
  roleTitles?: string[];
  status?: number;
  createdAt?: string;
  updatedAt?: string;
}

/** 用户表列表查询参数 */
export interface UsersListParams {
  pageNum: number;
  pageSize: number;
  keyword?: string;
  username?: string;
  nickname?: string;
  email?: string;
  deptId?: string;
  tenantId?: string;
  merchantId?: string;
  status?: number;
}

/** 用户表创建参数 */
export interface UsersCreateParams {
  username: string;
  password: string;
  nickname?: string;
  email?: string;
  avatar?: string;
  deptId?: string;
  tenantId?: string;
  merchantId?: string;
  roleIds?: string[];
  status?: number;
}

/** 用户表更新参数 */
export interface UsersUpdateParams {
  id: string;
  username: string;
  password?: string;
  nickname?: string;
  email?: string;
  avatar?: string;
  deptId?: string;
  tenantId?: string;
  merchantId?: string;
  roleIds?: string[];
  status?: number;
}

/** 用户表重置密码参数 */
export interface UsersResetPasswordParams {
  id: string;
  password: string;
}
