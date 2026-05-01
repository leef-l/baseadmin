const TOKEN_KEY = 'fd_member_token';
const USER_KEY = 'fd_member_user';

export const tokenStorage = {
  get: () => localStorage.getItem(TOKEN_KEY) || '',
  set: (v: string) => localStorage.setItem(TOKEN_KEY, v),
  clear: () => localStorage.removeItem(TOKEN_KEY),
};

export const userStorage = {
  get<T = any>(): T | null {
    const raw = localStorage.getItem(USER_KEY);
    if (!raw) return null;
    try {
      return JSON.parse(raw) as T;
    } catch {
      return null;
    }
  },
  set(v: any) {
    localStorage.setItem(USER_KEY, JSON.stringify(v));
  },
  clear() {
    localStorage.removeItem(USER_KEY);
  },
};
