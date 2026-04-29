import { computed } from 'vue';

import { useUserStore } from '@vben/stores';

function normalizeScopeId(value: unknown) {
  if (value === null || value === undefined || value === '') {
    return null;
  }
  const numeric = Number(value);
  return Number.isFinite(numeric) ? numeric : null;
}

export function isPlatformSuperAdminUser(
  userInfo: null | Record<string, unknown> | undefined,
) {
  return (
    Number(userInfo?.isAdmin ?? 0) === 1 &&
    normalizeScopeId(userInfo?.tenantId) === 0 &&
    normalizeScopeId(userInfo?.merchantId) === 0
  );
}

export function usePlatformSuperAdmin() {
  const userStore = useUserStore();
  return computed(() =>
    isPlatformSuperAdminUser(
      userStore.userInfo as null | Record<string, unknown> | undefined,
    ),
  );
}
