import type { BasicUserInfo } from '@vben-core/typings';

/** 用户信息 */
interface UserInfo extends BasicUserInfo {
  /**
   * 用户描述
   */
  desc: string;
  /**
   * 首页地址
   */
  homePath: string;

  /**
   * 平台超级管理员标识
   */
  isAdmin?: number;

  /**
   * 商户ID
   */
  merchantId?: string;

  /**
   * accessToken
   */
  token: string;

  /**
   * 租户ID
   */
  tenantId?: string;
}

export type { UserInfo };
