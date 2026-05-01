import req from './request';
import { HomeAggregate } from './types';

export const homeApi = {
  get: () => req.get<any, HomeAggregate>('/member-portal/home'),
};
