import { baseRequestClient, requestClient } from '#/api/request';
import { extractAccessCodes, getAuthInfoApi } from './auth-info';

export namespace AuthApi {
  /** 登录接口参数 */
  export interface LoginParams {
    password?: string;
    username?: string;
  }

  /** 票据登录接口参数 */
  export interface TicketLoginParams {
    ticket: string;
  }

  /** 生成应用票据接口参数 */
  export interface IssueTicketParams {
    targetApp: string;
  }

  /** 登录接口返回值 */
  export interface LoginResult {
    token: string;
    userId: string;
    username: string;
    nickname: string;
    avatar: string;
  }

  /** 生成应用票据接口返回值 */
  export interface IssueTicketResult {
    ticket: string;
    sourceApp: string;
    targetApp: string;
    expiresIn: number;
  }
}

/**
 * 登录
 */
export async function loginApi(data: AuthApi.LoginParams) {
  return requestClient.post<AuthApi.LoginResult>(
    '/system/auth/login',
    data,
    {
      // 登录接口不需要 token
    },
  );
}

/**
 * 票据登录
 */
export async function ticketLoginApi(data: AuthApi.TicketLoginParams) {
  return requestClient.post<AuthApi.LoginResult>(
    '/system/auth/ticket-login',
    data,
  );
}

/**
 * 生成应用间票据
 */
export async function issueTicketApi(data: AuthApi.IssueTicketParams) {
  return requestClient.post<AuthApi.IssueTicketResult>(
    '/system/auth/ticket',
    data,
  );
}

/**
 * 刷新accessToken（暂不支持，直接返回空）
 */
export async function refreshTokenApi() {
  return baseRequestClient.post<{ data: string; status: number }>(
    '/system/auth/refresh',
    {},
    { withCredentials: true },
  );
}

/**
 * 退出登录
 */
export async function logoutApi() {
  try {
    await requestClient.post('/system/auth/logout');
  } catch {
    // 即使后端调用失败，前端仍需清除本地 token
  }
}

/**
 * 获取用户权限码
 */
export async function getAccessCodesApi() {
  const res = await getAuthInfoApi();
  return extractAccessCodes(res);
}
