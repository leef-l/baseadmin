import req from './request';
import { MeProfile, MyWallets, PageResult, TeamMember, WalletLog } from './types';

export const meApi = {
  profile: () => req.get<any, MeProfile>('/member-portal/me/profile'),

  update: (data: { nickname?: string; avatar?: string; realName?: string }) =>
    req.put<any, {}>('/member-portal/me/update', data),

  changePassword: (oldPassword: string, newPassword: string) =>
    req.post<any, {}>('/member-portal/me/change-password', { oldPassword, newPassword }),

  changePhone: (newPhone: string, smsCode: string) =>
    req.post<any, {}>('/member-portal/me/change-phone', { newPhone, smsCode }),

  wallets: () => req.get<any, MyWallets>('/member-portal/me/wallets'),

  walletLogs: (params: { walletType?: number; pageNum?: number; pageSize?: number }) =>
    req.get<any, PageResult<WalletLog>>('/member-portal/me/wallet-logs', { params }),

  team: (params: { scope?: 'direct' | 'all'; pageNum?: number; pageSize?: number }) =>
    req.get<any, PageResult<TeamMember>>('/member-portal/me/team', { params }),
};
