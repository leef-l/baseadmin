<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { Descriptions, DescriptionsItem, Tag } from 'ant-design-vue';
import { getWorkOrderDetail } from '#/api/demo/work_order';
import type { WorkOrderItem } from '#/api/demo/work_order/types';

/** 标签颜色池 */
const TAG_COLORS = ['green', 'red', 'blue', 'orange', 'cyan', 'purple', 'geekblue', 'magenta'];

type EnumValue = number | string;

function getEnumLabel(map: Record<EnumValue, string>, value: EnumValue | null | undefined) {
  if (value === null || value === undefined || value === '') {
    return '-';
  }
  return map[value] ?? String(value);
}

/** 优先级映射 */
const priorityMap: Record<EnumValue, string> = {
  1: '低',
  2: '普通',
  3: '高',
  4: '紧急',
};

/** 优先级颜色 */
function getPriorityColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2, 3, 4];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 来源映射 */
const sourceTypeMap: Record<EnumValue, string> = {
  1: '官网',
  2: '电话',
  3: '微信',
  4: '后台',
};

/** 来源颜色 */
function getSourceTypeColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2, 3, 4];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 状态映射 */
const statusMap: Record<EnumValue, string> = {
  0: '待处理',
  1: '进行中',
  2: '已完成',
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

const detail = ref<WorkOrderItem | null>(null);
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
      modalApi.setState({ title: '体验工单详情' });
      try {
        const res = await getWorkOrderDetail(data.id);
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
      <DescriptionsItem label="工单号">{{ displayValue(detail.ticketNo) }}</DescriptionsItem>
      <DescriptionsItem label="客户">{{ detail.customerName || '-' }}</DescriptionsItem>
      <DescriptionsItem label="商品">{{ detail.productSkuNo || '-' }}</DescriptionsItem>
      <DescriptionsItem label="订单">{{ detail.orderOrderNo || '-' }}</DescriptionsItem>
      <DescriptionsItem label="工单标题">{{ displayValue(detail.title) }}</DescriptionsItem>
      <DescriptionsItem label="优先级">
        <Tag :color="getPriorityColor(detail.priority)">{{ getEnumLabel(priorityMap, detail.priority) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="来源">
        <Tag :color="getSourceTypeColor(detail.sourceType)">{{ getEnumLabel(sourceTypeMap, detail.sourceType) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="问题描述">{{ displayValue(detail.description) }}</DescriptionsItem>
      <DescriptionsItem label="附件">
        <a v-if="detail.attachmentFile" :href="detail.attachmentFile" target="_blank" rel="noreferrer noopener">查看文件</a>
        <span v-else>-</span>
      </DescriptionsItem>
      <DescriptionsItem label="状态">
        <Tag :color="getStatusColor(detail.status)">{{ getEnumLabel(statusMap, detail.status) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="租户">{{ detail.tenantName || '-' }}</DescriptionsItem>
      <DescriptionsItem label="商户">{{ detail.merchantName || '-' }}</DescriptionsItem>
      <DescriptionsItem label="截止时间">{{ displayValue(detail.dueAt) }}</DescriptionsItem>
      <DescriptionsItem label="创建时间">{{ displayValue(detail.createdAt) }}</DescriptionsItem>
      <DescriptionsItem label="更新时间">{{ displayValue(detail.updatedAt) }}</DescriptionsItem>
    </Descriptions>
  </Modal>
</template>
