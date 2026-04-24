import { flushPromises, mount } from '@vue/test-utils';
import { defineComponent, h } from 'vue';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import Login from '../login.vue';

const formApi = {
  getValues: vi.fn(async () => ({ username: 'alice' })),
  setFieldValue: vi.fn(),
  validate: vi.fn(async () => ({ valid: true })),
};

vi.mock('@vben-core/form-ui', () => ({
  useVbenForm: () => [
    defineComponent({
      name: 'MockForm',
      setup() {
        return () => h('div', { class: 'mock-form' });
      },
    }),
    formApi,
  ],
}));

vi.mock('@vben/locales', () => ({
  $t: (key: string) => key,
}));

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn(),
  }),
}));

vi.mock('@vben-core/shadcn-ui', () => ({
  VbenButton: defineComponent({
    name: 'VbenButton',
    emits: ['click'],
    setup(_props, { emit, slots }) {
      return () =>
        h(
          'button',
          {
            onClick: () => emit('click'),
          },
          slots.default ? slots.default() : [],
        );
    },
  }),
  VbenCheckbox: defineComponent({
    name: 'VbenCheckbox',
    props: {
      modelValue: {
        type: Boolean,
        default: false,
      },
    },
    emits: ['update:modelValue'],
    setup(props, { emit, slots }) {
      return () =>
        h('label', [
          h('input', {
            checked: props.modelValue,
            type: 'checkbox',
            onChange: (event: Event) => {
              emit(
                'update:modelValue',
                (event.target as HTMLInputElement).checked,
              );
            },
          }),
          ...(slots.default ? slots.default() : []),
        ]);
    },
  }),
}));

vi.mock('./auth-title.vue', () => ({
  default: defineComponent({
    name: 'AuthTitle',
    setup(_props, { slots }) {
      return () => h('div', slots.default ? slots.default() : []);
    },
  }),
}));

vi.mock('./third-party-login.vue', () => ({
  default: defineComponent({
    name: 'ThirdPartyLogin',
    setup() {
      return () => h('div');
    },
  }),
}));

describe('AuthenticationLogin remember me', () => {
  const rememberKey = `REMEMBER_ME_USERNAME_${location.hostname}`;

  beforeEach(() => {
    localStorage.clear();
    formApi.getValues.mockResolvedValue({ username: 'alice' });
    formApi.setFieldValue.mockReset();
    formApi.validate.mockResolvedValue({ valid: true });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  it('clears remembered username when remember me is disabled', async () => {
    localStorage.setItem(rememberKey, 'cached-user');

    const wrapper = mount(Login, {
      props: {
        formSchema: [],
        showRememberMe: false,
      },
    });

    await wrapper.find('button').trigger('click');
    await flushPromises();

    expect(localStorage.getItem(rememberKey)).toBeNull();
    expect(formApi.setFieldValue).not.toHaveBeenCalled();
    expect(wrapper.emitted('submit')?.[0]?.[0]).toEqual({ username: 'alice' });
  });
});
