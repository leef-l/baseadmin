<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { Descriptions, DescriptionsItem, Tag } from 'ant-design-vue';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { getDomainDetail } from '#/api/system/domain';
import type { DomainItem } from '#/api/system/domain/types';

/** 标签颜色池 */
const TAG_COLORS = ['green', 'red', 'blue', 'orange', 'cyan', 'purple', 'geekblue', 'magenta'];

type EnumValue = number | string;

function getEnumLabel(map: Record<EnumValue, string>, value: EnumValue | null | undefined) {
  if (value === null || value === undefined || value === '') {
    return '-';
  }
  return map[value] ?? String(value);
}

/** 主体类型映射 */
const ownerTypeMap: Record<EnumValue, string> = {
  1: '租户',
  2: '商户',
};

/** 主体类型颜色 */
function getOwnerTypeColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 校验状态映射 */
const verifyStatusMap: Record<EnumValue, string> = {
  0: '未校验',
  1: '已校验',
};

/** 校验状态颜色 */
function getVerifyStatusColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [0, 1];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** SSL状态映射 */
const sslStatusMap: Record<EnumValue, string> = {
  0: '未配置',
  1: '已配置',
};

/** SSL状态颜色 */
function getSslStatusColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [0, 1];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** Nginx配置状态映射 */
const nginxStatusMap: Record<EnumValue, string> = {
  0: '未应用',
  1: '已应用',
};

/** Nginx配置状态颜色 */
function getNginxStatusColor(val: EnumValue | null | undefined): string {
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
const detail = ref<DomainItem | null>(null);
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
      modalApi.setState({ title: '域名绑定表详情' });
      try {
        const res = await getDomainDetail(data.id);
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
      <DescriptionsItem label="绑定域名">{{ displayValue(detail.domain) }}</DescriptionsItem>
      <DescriptionsItem label="主体类型">
        <Tag :color="getOwnerTypeColor(detail.ownerType)">{{ getEnumLabel(ownerTypeMap, detail.ownerType) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="租户">{{ detail.tenantName || '-' }}</DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="商户">{{ detail.merchantName || '-' }}</DescriptionsItem>
      <DescriptionsItem label="应用编码">{{ displayValue(detail.appCode) }}</DescriptionsItem>
      <DescriptionsItem label="域名校验令牌">{{ displayValue(detail.verifyToken) }}</DescriptionsItem>
      <DescriptionsItem label="校验状态">
        <Tag :color="getVerifyStatusColor(detail.verifyStatus)">{{ getEnumLabel(verifyStatusMap, detail.verifyStatus) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="SSL状态">
        <Tag :color="getSslStatusColor(detail.sslStatus)">{{ getEnumLabel(sslStatusMap, detail.sslStatus) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="Nginx配置状态">
        <Tag :color="getNginxStatusColor(detail.nginxStatus)">{{ getEnumLabel(nginxStatusMap, detail.nginxStatus) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="状态">
        <Tag :color="getStatusColor(detail.status)">{{ getEnumLabel(statusMap, detail.status) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="备注">{{ displayValue(detail.remark) }}</DescriptionsItem>
      <DescriptionsItem label="创建时间">{{ displayValue(detail.createdAt) }}</DescriptionsItem>
      <DescriptionsItem label="更新时间">{{ displayValue(detail.updatedAt) }}</DescriptionsItem>
    </Descriptions>
  </Modal>
</template>
