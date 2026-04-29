/** 体验预约类型定义 */

/** 体验预约项 */
export interface AppointmentItem {
  id: string;
  appointmentNo: string;
  customerID?: string;
  customerName?: string;
  subject: string;
  appointmentAt?: string;
  contactPhone?: string;
  address?: string;
  remark?: string;
  status?: number;
  tenantID?: string;
  tenantName?: string;
  merchantID?: string;
  merchantName?: string;
  createdAt?: string;
  updatedAt?: string;
}

/** 体验预约列表查询参数 */
export interface AppointmentListParams {
  pageNum: number;
  pageSize: number;
  orderBy?: string;
  orderDir?: string;
  startTime?: string;
  endTime?: string;
  keyword?: string;
  appointmentNo?: string;
  subject?: string;
  contactPhone?: string;
  customerID?: string;
  tenantID?: string;
  merchantID?: string;
  status?: number;
  appointmentAtStart?: string;
  appointmentAtEnd?: string;
}

/** 体验预约创建参数 */
export interface AppointmentCreateParams {
  appointmentNo: string;
  customerID?: string;
  subject: string;
  appointmentAt?: string;
  contactPhone?: string;
  address?: string;
  remark?: string;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}

/** 体验预约更新参数 */
export interface AppointmentUpdateParams {
  id: string;
  appointmentNo: string;
  customerID?: string;
  subject: string;
  appointmentAt?: string;
  contactPhone?: string;
  address?: string;
  remark?: string;
  status?: number;
  tenantID?: string;
  merchantID?: string;
}
