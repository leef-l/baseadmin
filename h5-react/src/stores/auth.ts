import { create } from 'zustand';
import { tokenStorage, userStorage } from '@/utils/storage';

export interface MemberProfile {
  userID: string;
  username?: string;
  nickname?: string;
  phone?: string;
  avatar?: string;
  inviteCode?: string;
  levelName?: string;
  isQualified?: number;
}

interface AuthState {
  token: string;
  user: MemberProfile | null;
  setSession: (token: string, user: MemberProfile) => void;
  setUser: (user: MemberProfile) => void;
  clear: () => void;
  isAuthed: () => boolean;
}

export const useAuth = create<AuthState>((set, get) => ({
  token: tokenStorage.get(),
  user: userStorage.get<MemberProfile>(),
  setSession(token, user) {
    tokenStorage.set(token);
    userStorage.set(user);
    set({ token, user });
  },
  setUser(user) {
    userStorage.set(user);
    set({ user });
  },
  clear() {
    tokenStorage.clear();
    userStorage.clear();
    set({ token: '', user: null });
  },
  isAuthed: () => !!get().token,
}));
