import req from './request';
import { InvitePreview, LoginResult } from './types';

export type SmsScene = 'register' | 'forget_password' | 'change_phone';

export const authApi = {
  sendSmsCode: (phone: string, scene: SmsScene) =>
    req.post<any, { expiresIn: number }>('/member-portal/sms/code', { phone, scene }),

  register: (data: {
    phone: string;
    smsCode: string;
    password: string;
    inviteCode: string;
    nickname?: string;
  }) => req.post<any, LoginResult>('/member-portal/auth/register', data),

  login: (account: string, password: string) =>
    req.post<any, LoginResult>('/member-portal/auth/login', { account, password }),

  forgetPassword: (data: { phone: string; smsCode: string; newPassword: string }) =>
    req.post<any, {}>('/member-portal/auth/forget-password', data),

  invitePreview: (inviteCode: string) =>
    req.get<any, InvitePreview>('/member-portal/auth/invite-preview', { params: { inviteCode } }),
};
