/** 体验问卷类型定义 */

/** 体验问卷项 */
export interface SurveyItem {
  id: string;
  surveyNo: string;
  title: string;
  poster?: string;
  questionJSON?: string;
  introContent?: string;
  publishAt?: string;
  expireAt?: string;
  isAnonymous?: number;
  status?: number;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 体验问卷列表查询参数 */
export interface SurveyListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  surveyNo?: string;
  title?: string;
  tenantID?: string;
  merchantID?: string;
  isAnonymous?: number;
  status?: number;
  publishAtStart?: string;
  publishAtEnd?: string;
  expireAtStart?: string;
  expireAtEnd?: string;
}

/** 体验问卷创建参数 */
export interface SurveyCreateParams {
  surveyNo: string;
  title: string;
  poster?: string;
  questionJSON?: string;
  introContent?: string;
  publishAt?: string;
  expireAt?: string;
  isAnonymous?: number;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}

/** 体验问卷更新参数 */
export interface SurveyUpdateParams {
  id: string;
  surveyNo: string;
  title: string;
  poster?: string;
  questionJSON?: string;
  introContent?: string;
  publishAt?: string;
  expireAt?: string;
  isAnonymous?: number;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}
