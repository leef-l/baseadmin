<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { Descriptions, DescriptionsItem, Tag } from 'ant-design-vue';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { getTeamExportDetail } from '#/api/member/team_export';
import type { TeamExportItem } from '#/api/member/team_export/types';

/** 标签颜色池 */
const TAG_COLORS = ['green', 'red', 'blue', 'orange', 'cyan', 'purple', 'geekblue', 'magenta'];

type EnumValue = number | string;

function getEnumLabel(map: Record<EnumValue, string>, value: EnumValue | null | undefined) {
  if (value === null || value === undefined || value === '') {
    return '-';
  }
  return map[value] ?? String(value);
}

/** 导出类型映射 */
const exportTypeMap: Record<EnumValue, string> = {
  1: '手动导出',
  2: '自动升级导出',
};

/** 导出类型颜色 */
function getExportTypeColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 部署状态映射 */
const deployStatusMap: Record<EnumValue, string> = {
  0: '未部署',
  1: '部署中',
  2: '已部署',
  3: '部署失败',
};

/** 部署状态颜色 */
function getDeployStatusColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [0, 1, 2, 3];
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
const detail = ref<TeamExportItem | null>(null);
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
      modalApi.setState({ title: '团队数据导出详情' });
      try {
        const res = await getTeamExportDetail(data.id);
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
      <DescriptionsItem label="目标会员">{{ detail.userNickname || '-' }}</DescriptionsItem>
      <DescriptionsItem label="团队成员数">{{ displayValue(detail.teamMemberCount) }}</DescriptionsItem>
      <DescriptionsItem label="导出类型">
        <Tag :color="getExportTypeColor(detail.exportType)">{{ getEnumLabel(exportTypeMap, detail.exportType) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="导出文件地址">
        <a v-if="detail.fileURL && /^https?:\/\//i.test(detail.fileURL)" :href="detail.fileURL" target="_blank" rel="noreferrer noopener">{{ detail.fileURL }}</a>
        <span v-else-if="detail.fileURL">{{ detail.fileURL }}</span>
        <span v-else>-</span>
      </DescriptionsItem>
      <DescriptionsItem label="文件大小">{{ displayValue(detail.fileSize) }}</DescriptionsItem>
      <DescriptionsItem label="部署状态">
        <Tag :color="getDeployStatusColor(detail.deployStatus)">{{ getEnumLabel(deployStatusMap, detail.deployStatus) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="部署域名">{{ displayValue(detail.deployDomain) }}</DescriptionsItem>
      <DescriptionsItem label="备注">{{ displayValue(detail.remark) }}</DescriptionsItem>
      <DescriptionsItem label="状态">
        <Tag :color="getStatusColor(detail.status)">{{ getEnumLabel(statusMap, detail.status) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="租户">{{ detail.tenantName || '-' }}</DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="商户">{{ detail.merchantName || '-' }}</DescriptionsItem>
      <DescriptionsItem label="部署完成时间">{{ displayValue(detail.deployedAt) }}</DescriptionsItem>
      <DescriptionsItem label="创建时间">{{ displayValue(detail.createdAt) }}</DescriptionsItem>
      <DescriptionsItem label="更新时间">{{ displayValue(detail.updatedAt) }}</DescriptionsItem>
    </Descriptions>
  </Modal>
</template>
