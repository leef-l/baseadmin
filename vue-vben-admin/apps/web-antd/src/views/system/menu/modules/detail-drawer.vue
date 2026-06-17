<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { Descriptions, DescriptionsItem, Tag } from 'ant-design-vue';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { getMenuDetail } from '#/api/system/menu';
import type { MenuItem } from '#/api/system/menu/types';

/** 标签颜色池 */
const TAG_COLORS = ['green', 'red', 'blue', 'orange', 'cyan', 'purple', 'geekblue', 'magenta'];

type EnumValue = number | string;

function getEnumLabel(map: Record<EnumValue, string>, value: EnumValue | null | undefined) {
  if (value === null || value === undefined || value === '') {
    return '-';
  }
  return map[value] ?? String(value);
}

/** 类型映射 */
const typeMap: Record<EnumValue, string> = {
  1: '目录',
  2: '菜单',
  3: '按钮',
  4: '外链',
  5: '内链',
};

/** 类型颜色 */
function getTypeColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2, 3, 4, 5];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 是否显示映射 */
const isShowMap: Record<EnumValue, string> = {
  0: '隐藏',
  1: '显示',
};

/** 是否显示颜色 */
function getIsShowColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [0, 1];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 是否缓存映射 */
const isCacheMap: Record<EnumValue, string> = {
  0: '不缓存',
  1: '缓存',
};

/** 是否缓存颜色 */
function getIsCacheColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [0, 1];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 状态映射 */
const statusMap: Record<EnumValue, string> = {
  0: '关闭',
  1: '开启',
};

/** 状态颜色 */
function getStatusColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [0, 1];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

const isPlatformSuperAdmin = usePlatformSuperAdmin();
const detail = ref<MenuItem | null>(null);
const openToken = ref(0);

function displayValue(value: null | number | string | undefined) {
  if (value === null || value === undefined || value === '') {
    return '-';
  }
  return value;
}

const [Modal, modalApi] = useVbenModal({
  fullscreenButton: false,
  footer: false,
  async onOpenChange(isOpen: boolean) {
    if (!isOpen) {
      openToken.value += 1;
      detail.value = null;
      return;
    }

    const currentOpenToken = ++openToken.value;
    const data = modalApi.getData<{ id: string }>();
    if (data?.id) {
      modalApi.setState({ title: '菜单表详情' });
      try {
        const res = await getMenuDetail(data.id);
        if (currentOpenToken !== openToken.value) {
          return;
        }
        detail.value = res;
      } catch {
        if (currentOpenToken === openToken.value) {
          detail.value = null;
        }
      }
    }
  },
});
</script>

<template>
  <Modal class="w-[600px]">
    <Descriptions v-if="detail" bordered :column="1" size="small">
      <DescriptionsItem label="ID">{{ detail.id }}</DescriptionsItem>
      <DescriptionsItem label="上级菜单ID，0 表示顶级菜单">{{ detail.menuTitle || '-' }}</DescriptionsItem>
      <DescriptionsItem label="菜单名称">{{ displayValue(detail.title) }}</DescriptionsItem>
      <DescriptionsItem label="类型">
        <Tag :color="getTypeColor(detail.type)">{{ getEnumLabel(typeMap, detail.type) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="前端路由路径">{{ displayValue(detail.path) }}</DescriptionsItem>
      <DescriptionsItem label="前端组件路径">{{ displayValue(detail.component) }}</DescriptionsItem>
      <DescriptionsItem label="权限标识（如 system">{{ displayValue(detail.permission) }}</DescriptionsItem>
      <DescriptionsItem label="菜单图标">{{ displayValue(detail.icon) }}</DescriptionsItem>
      <DescriptionsItem label="排序">{{ displayValue(detail.sort) }}</DescriptionsItem>
      <DescriptionsItem label="是否显示">
        <Tag :color="getIsShowColor(detail.isShow)">{{ getEnumLabel(isShowMap, detail.isShow) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="是否缓存">
        <Tag :color="getIsCacheColor(detail.isCache)">{{ getEnumLabel(isCacheMap, detail.isCache) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="外链/内链地址">
        <a v-if="detail.linkURL && /^https?:\/\//i.test(detail.linkURL)" :href="detail.linkURL" target="_blank" rel="noreferrer noopener">{{ detail.linkURL }}</a>
        <span v-else-if="detail.linkURL">{{ detail.linkURL }}</span>
        <span v-else>-</span>
      </DescriptionsItem>
      <DescriptionsItem label="状态">
        <Tag :color="getStatusColor(detail.status)">{{ getEnumLabel(statusMap, detail.status) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="租户">{{ detail.tenantName || '-' }}</DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="商户">{{ detail.merchantName || '-' }}</DescriptionsItem>
      <DescriptionsItem label="创建时间">{{ displayValue(detail.createdAt) }}</DescriptionsItem>
      <DescriptionsItem label="更新时间">{{ displayValue(detail.updatedAt) }}</DescriptionsItem>
    </Descriptions>
  </Modal>
</template>
