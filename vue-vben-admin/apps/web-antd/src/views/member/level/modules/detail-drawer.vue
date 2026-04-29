<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { Descriptions, DescriptionsItem, Tag } from 'ant-design-vue';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { getLevelDetail } from '#/api/member/level';
import type { LevelItem } from '#/api/member/level/types';

/** 标签颜色池 */
const TAG_COLORS = ['green', 'red', 'blue', 'orange', 'cyan', 'purple', 'geekblue', 'magenta'];

type EnumValue = number | string;

function getEnumLabel(map: Record<EnumValue, string>, value: EnumValue | null | undefined) {
  if (value === null || value === undefined || value === '') {
    return '-';
  }
  return map[value] ?? String(value);
}

/** 是否最高等级映射 */
const isTopMap: Record<EnumValue, string> = {
  0: '否',
  1: '是',
};

/** 是否最高等级颜色 */
function getIsTopColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [0, 1];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 到达后自动部署站点映射 */
const autoDeployMap: Record<EnumValue, string> = {
  0: '否',
  1: '是',
};

/** 到达后自动部署站点颜色 */
function getAutoDeployColor(val: EnumValue | null | undefined): string {
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
const detail = ref<LevelItem | null>(null);
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
      modalApi.setState({ title: '会员等级配置详情' });
      try {
        const res = await getLevelDetail(data.id);
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
      <DescriptionsItem label="等级名称">{{ displayValue(detail.name) }}</DescriptionsItem>
      <DescriptionsItem label="等级编号">{{ displayValue(detail.levelNo) }}</DescriptionsItem>
      <DescriptionsItem label="等级图标">{{ displayValue(detail.icon) }}</DescriptionsItem>
      <DescriptionsItem label="有效天数">{{ displayValue(detail.durationDays) }}</DescriptionsItem>
      <DescriptionsItem label="升级所需有效用户数">{{ displayValue(detail.needActiveCount) }}</DescriptionsItem>
      <DescriptionsItem label="升级所需团队营业额">{{ detail.needTeamTurnover != null ? (detail.needTeamTurnover / 100).toFixed(2) : '-' }}</DescriptionsItem>
      <DescriptionsItem label="是否最高等级">
        <Tag :color="getIsTopColor(detail.isTop)">{{ getEnumLabel(isTopMap, detail.isTop) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="到达后自动部署站点">
        <Tag :color="getAutoDeployColor(detail.autoDeploy)">{{ getEnumLabel(autoDeployMap, detail.autoDeploy) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="等级说明">{{ displayValue(detail.remark) }}</DescriptionsItem>
      <DescriptionsItem label="排序">{{ displayValue(detail.sort) }}</DescriptionsItem>
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
