import req from './request';
import { PortalBizConfig } from './types';

export const bizConfigApi = {
  get: () => req.get<any, PortalBizConfig>('/member-portal/biz-config'),
};
