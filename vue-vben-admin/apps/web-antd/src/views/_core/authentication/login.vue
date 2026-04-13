<script lang="ts" setup>
import type { VbenFormSchema } from '@vben/common-ui';
import type { Recordable } from '@vben/types';

import { computed, onMounted, markRaw } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { AuthenticationLogin, SliderCaptcha, z } from '@vben/common-ui';
import { $t } from '@vben/locales';
import { preferences } from '@vben/preferences';

import { useAuthStore } from '#/store';

defineOptions({ name: 'Login' });

const authStore = useAuthStore();
const route = useRoute();
const router = useRouter();

const formSchema = computed((): VbenFormSchema[] => {
  return [
    {
      component: 'VbenInput',
      componentProps: {
        placeholder: $t('authentication.usernameTip'),
      },
      fieldName: 'username',
      label: $t('authentication.username'),
      rules: z.string().min(1, { message: $t('authentication.usernameTip') }),
    },
    {
      component: 'VbenInputPassword',
      componentProps: {
        placeholder: $t('authentication.password'),
      },
      fieldName: 'password',
      label: $t('authentication.password'),
      rules: z.string().min(1, { message: $t('authentication.passwordTip') }),
    },
    {
      component: markRaw(SliderCaptcha),
      fieldName: 'captcha',
      rules: z.boolean().refine((value) => value, {
        message: $t('authentication.verifyRequiredTip'),
      }),
    },
  ];
});

function resolveRedirectPath() {
  const redirect = route.query.redirect;
  if (typeof redirect !== 'string' || !redirect.trim()) {
    return preferences.app.defaultHomePath;
  }
  const next = redirect.trim();
  if (!next.startsWith('/')) {
    return preferences.app.defaultHomePath;
  }
  try {
    const decoded = decodeURIComponent(next);
    return decoded.startsWith('/') ? decoded : preferences.app.defaultHomePath;
  } catch {
    return next;
  }
}

async function handleLoginSuccess() {
  await router.replace(resolveRedirectPath());
}

async function handleSubmit(values: Recordable<any>) {
  await authStore.authLogin(values, handleLoginSuccess);
}

async function tryTicketLogin() {
  const ticket = typeof route.query.ticket === 'string' ? route.query.ticket.trim() : '';
  if (!ticket) {
    return;
  }
  await authStore.authLoginByTicket(ticket, handleLoginSuccess);
}

onMounted(() => {
  void tryTicketLogin();
});
</script>

<template>
  <AuthenticationLogin
    :form-schema="formSchema"
    :loading="authStore.loginLoading"
    :show-code-login="false"
    :show-forget-password="false"
    :show-qrcode-login="false"
    :show-register="false"
    :show-third-party-login="false"
    @submit="handleSubmit"
  />
</template>
