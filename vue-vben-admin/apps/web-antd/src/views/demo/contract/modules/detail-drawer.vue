<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { Descriptions, DescriptionsItem, Tag } from 'ant-design-vue';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { getContractDetail } from '#/api/demo/contract';
import type { ContractItem } from '#/api/demo/contract/types';

/** 标签颜色池 */
const TAG_COLORS = ['green', 'red', 'blue', 'orange', 'cyan', 'purple', 'geekblue', 'magenta'];

type EnumValue = number | string;

function getEnumLabel(map: Record<EnumValue, string>, value: EnumValue | null | undefined) {
  if (value === null || value === undefined || value === '') {
    return '-';
  }
  return map[value] ?? String(value);
}

/** 状态映射 */
const statusMap: Record<EnumValue, string> = {
  0: '待审核',
  1: '已通过',
  2: '已拒绝',
  3: '已取消',
};

/** 状态颜色 */
function getStatusColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [0, 1, 2, 3];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

const isPlatformSuperAdmin = usePlatformSuperAdmin();
const detail = ref<ContractItem | null>(null);
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
      modalApi.setState({ title: '体验合同详情' });
      try {
        const res = await getContractDetail(data.id);
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
      <DescriptionsItem label="合同编号">{{ displayValue(detail.contractNo) }}</DescriptionsItem>
      <DescriptionsItem label="客户">{{ detail.customerName || '-' }}</DescriptionsItem>
      <DescriptionsItem label="订单">{{ detail.orderOrderNo || '-' }}</DescriptionsItem>
      <DescriptionsItem label="合同标题">{{ displayValue(detail.title) }}</DescriptionsItem>
      <DescriptionsItem label="合同文件">
        <a v-if="detail.contractFile && /^https?:\/\//i.test(detail.contractFile)" :href="detail.contractFile" target="_blank" rel="noreferrer noopener">查看文件</a>
        <span v-else-if="detail.contractFile">{{ detail.contractFile }}</span>
        <span v-else>-</span>
      </DescriptionsItem>
      <DescriptionsItem label="签章图片">
        <img v-if="detail.signImage && /^https?:\/\//i.test(detail.signImage)" :src="detail.signImage" style="max-width: 200px; max-height: 200px; object-fit: contain;" />
        <span v-else>-</span>
      </DescriptionsItem>
      <DescriptionsItem label="合同金额">{{ detail.contractAmount != null ? (detail.contractAmount / 100).toFixed(2) : '-' }}</DescriptionsItem>
      <DescriptionsItem label="状态">
        <Tag :color="getStatusColor(detail.status)">{{ getEnumLabel(statusMap, detail.status) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="租户">{{ detail.tenantName || '-' }}</DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="商户">{{ detail.merchantName || '-' }}</DescriptionsItem>
      <DescriptionsItem label="签署时间">{{ displayValue(detail.signedAt) }}</DescriptionsItem>
      <DescriptionsItem label="到期时间">{{ displayValue(detail.expiresAt) }}</DescriptionsItem>
      <DescriptionsItem label="创建时间">{{ displayValue(detail.createdAt) }}</DescriptionsItem>
      <DescriptionsItem label="更新时间">{{ displayValue(detail.updatedAt) }}</DescriptionsItem>
    </Descriptions>
  </Modal>
</template>
