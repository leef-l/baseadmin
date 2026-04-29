/** 体验活动类型定义 */

/** 体验活动项 */
export interface CampaignItem {
  id: string;
  campaignNo: string;
  title: string;
  banner?: string;
  type?: number;
  channel?: number;
  budgetAmount?: number;
  landingURL?: string;
  ruleJSON?: string;
  introContent?: string;
  startAt?: string;
  endAt?: string;
  isPublic?: number;
  status?: number;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 体验活动列表查询参数 */
export interface CampaignListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  campaignNo?: string;
  title?: string;
  tenantID?: string;
  merchantID?: string;
  type?: number;
  channel?: number;
  isPublic?: number;
  status?: number;
  startAtStart?: string;
  startAtEnd?: string;
  endAtStart?: string;
  endAtEnd?: string;
}

/** 体验活动创建参数 */
export interface CampaignCreateParams {
  campaignNo: string;
  title: string;
  banner?: string;
  type?: number;
  channel?: number;
  budgetAmount?: number;
  landingURL?: string;
  ruleJSON?: string;
  introContent?: string;
  startAt?: string;
  endAt?: string;
  isPublic?: number;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}

/** 体验活动更新参数 */
export interface CampaignUpdateParams {
  id: string;
  campaignNo: string;
  title: string;
  banner?: string;
  type?: number;
  channel?: number;
  budgetAmount?: number;
  landingURL?: string;
  ruleJSON?: string;
  introContent?: string;
  startAt?: string;
  endAt?: string;
  isPublic?: number;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}
