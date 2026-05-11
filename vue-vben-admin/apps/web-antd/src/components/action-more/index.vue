<script setup lang="ts">
import { computed } from 'vue';

import { Button, Dropdown, Menu, Space } from 'ant-design-vue';

export interface ActionMoreItem {
  danger?: boolean;
  disabled?: boolean;
  key: string;
  label: string;
  onClick: () => void;
  visible?: boolean;
}

const props = withDefaults(
  defineProps<{
    actions: ActionMoreItem[];
  }>(),
  {
    actions: () => [],
  },
);

const visibleActions = computed(() =>
  props.actions.filter((action) => action.visible !== false),
);

const directActions = computed(() => {
  if (visibleActions.value.length > 2) {
    return visibleActions.value.slice(0, 2);
  }
  return visibleActions.value;
});

const overflowActions = computed(() => {
  if (visibleActions.value.length > 2) {
    return visibleActions.value.slice(2);
  }
  return [];
});

const menuItems = computed(() =>
  overflowActions.value.map((action) => ({
    danger: action.danger,
    disabled: action.disabled,
    key: action.key,
    label: action.label,
  })),
);

function handleDirectClick(action: ActionMoreItem) {
  if (action.disabled) {
    return;
  }
  action.onClick();
}

function handleMenuClick({ key }: { key: PropertyKey }) {
  const action = overflowActions.value.find((item) => item.key === String(key));
  if (!action || action.disabled) {
    return;
  }
  action.onClick();
}
</script>

<template>
  <Space :size="0" class="action-more">
    <Button
      v-for="action in directActions"
      :key="action.key"
      :danger="action.danger"
      :disabled="action.disabled"
      size="small"
      type="link"
      @click="handleDirectClick(action)"
    >
      {{ action.label }}
    </Button>
    <Dropdown v-if="overflowActions.length" :trigger="['click']">
      <Button size="small" type="link">更多</Button>
      <template #overlay>
        <Menu :items="menuItems" @click="handleMenuClick" />
      </template>
    </Dropdown>
  </Space>
</template>
