<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { Descriptions, DescriptionsItem, Tag } from 'ant-design-vue';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { getOrderDetail } from '#/api/demo/order';
import type { OrderItem } from '#/api/demo/order/types';

/** 标签颜色池 */
const TAG_COLORS = ['green', 'red', 'blue', 'orange', 'cyan', 'purple', 'geekblue', 'magenta'];

type EnumValue = number | string;

function getEnumLabel(map: Record<EnumValue, string>, value: EnumValue | null | undefined) {
  if (value === null || value === undefined || value === '') {
    return '-';
  }
  return map[value] ?? String(value);
}

/** 支付状态映射 */
const payStatusMap: Record<EnumValue, string> = {
  0: '待支付',
  1: '已支付',
  2: '已退款',
};

/** 支付状态颜色 */
function getPayStatusColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [0, 1, 2];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 发货状态映射 */
const deliverStatusMap: Record<EnumValue, string> = {
  0: '待发货',
  1: '已发货',
  2: '已签收',
};

/** 发货状态颜色 */
function getDeliverStatusColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [0, 1, 2];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 状态映射 */
const statusMap: Record<EnumValue, string> = {
  0: '待确认',
  1: '已确认',
  2: '已取消',
};

/** 状态颜色 */
function getStatusColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [0, 1, 2];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

const isPlatformSuperAdmin = usePlatformSuperAdmin();
const detail = ref<OrderItem | null>(null);
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
      modalApi.setState({ title: '体验订单详情' });
      try {
        const res = await getOrderDetail(data.id);
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
      <DescriptionsItem label="订单号">{{ displayValue(detail.orderNo) }}</DescriptionsItem>
      <DescriptionsItem label="客户">{{ detail.customerName || '-' }}</DescriptionsItem>
      <DescriptionsItem label="商品">{{ detail.productSkuNo || '-' }}</DescriptionsItem>
      <DescriptionsItem label="购买数量">{{ displayValue(detail.quantity) }}</DescriptionsItem>
      <DescriptionsItem label="订单金额">{{ detail.amount != null ? (detail.amount / 100).toFixed(2) : '-' }}</DescriptionsItem>
      <DescriptionsItem label="支付状态">
        <Tag :color="getPayStatusColor(detail.payStatus)">{{ getEnumLabel(payStatusMap, detail.payStatus) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="发货状态">
        <Tag :color="getDeliverStatusColor(detail.deliverStatus)">{{ getEnumLabel(deliverStatusMap, detail.deliverStatus) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="收货电话">{{ displayValue(detail.receiverPhone) }}</DescriptionsItem>
      <DescriptionsItem label="收货地址">{{ displayValue(detail.address) }}</DescriptionsItem>
      <DescriptionsItem label="备注">{{ displayValue(detail.remark) }}</DescriptionsItem>
      <DescriptionsItem label="状态">
        <Tag :color="getStatusColor(detail.status)">{{ getEnumLabel(statusMap, detail.status) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="租户">{{ detail.tenantName || '-' }}</DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="商户">{{ detail.merchantName || '-' }}</DescriptionsItem>
      <DescriptionsItem label="支付时间">{{ displayValue(detail.paidAt) }}</DescriptionsItem>
      <DescriptionsItem label="发货时间">{{ displayValue(detail.deliverAt) }}</DescriptionsItem>
      <DescriptionsItem label="创建时间">{{ displayValue(detail.createdAt) }}</DescriptionsItem>
      <DescriptionsItem label="更新时间">{{ displayValue(detail.updatedAt) }}</DescriptionsItem>
    </Descriptions>
  </Modal>
</template>
