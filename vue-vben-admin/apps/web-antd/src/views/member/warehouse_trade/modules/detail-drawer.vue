<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { Descriptions, DescriptionsItem, Tag } from 'ant-design-vue';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { getWarehouseTradeDetail } from '#/api/member/warehouse_trade';
import type { WarehouseTradeItem } from '#/api/member/warehouse_trade/types';

/** 标签颜色池 */
const TAG_COLORS = ['green', 'red', 'blue', 'orange', 'cyan', 'purple', 'geekblue', 'magenta'];

type EnumValue = number | string;

function getEnumLabel(map: Record<EnumValue, string>, value: EnumValue | null | undefined) {
  if (value === null || value === undefined || value === '') {
    return '-';
  }
  return map[value] ?? String(value);
}

/** 交易状态映射 */
const tradeStatusMap: Record<EnumValue, string> = {
  1: '待卖家确认',
  2: '已确认完成',
  3: '已取消',
};

/** 交易状态颜色 */
function getTradeStatusColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2, 3];
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
const detail = ref<WarehouseTradeItem | null>(null);
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
      modalApi.setState({ title: '仓库交易记录详情' });
      try {
        const res = await getWarehouseTradeDetail(data.id);
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
      <DescriptionsItem label="交易编号">{{ displayValue(detail.tradeNo) }}</DescriptionsItem>
      <DescriptionsItem label="仓库商品">{{ detail.warehouseGoodsTitle || '-' }}</DescriptionsItem>
      <DescriptionsItem label="挂卖记录">{{ detail.warehouseListingID || '-' }}</DescriptionsItem>
      <DescriptionsItem label="卖家">{{ detail.userNickname || '-' }}</DescriptionsItem>
      <DescriptionsItem label="买家">{{ detail.buyerNickname || '-' }}</DescriptionsItem>
      <DescriptionsItem label="成交价格">{{ detail.tradePrice != null ? (detail.tradePrice / 100).toFixed(2) : '-' }}</DescriptionsItem>
      <DescriptionsItem label="平台扣除费用">{{ detail.platformFee != null ? (detail.platformFee / 100).toFixed(2) : '-' }}</DescriptionsItem>
      <DescriptionsItem label="卖家实收">{{ detail.sellerIncome != null ? (detail.sellerIncome / 100).toFixed(2) : '-' }}</DescriptionsItem>
      <DescriptionsItem label="交易状态">
        <Tag :color="getTradeStatusColor(detail.tradeStatus)">{{ getEnumLabel(tradeStatusMap, detail.tradeStatus) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="备注">{{ displayValue(detail.remark) }}</DescriptionsItem>
      <DescriptionsItem label="状态">
        <Tag :color="getStatusColor(detail.status)">{{ getEnumLabel(statusMap, detail.status) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="租户">{{ detail.tenantName || '-' }}</DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="商户">{{ detail.merchantName || '-' }}</DescriptionsItem>
      <DescriptionsItem label="确认时间">{{ displayValue(detail.confirmedAt) }}</DescriptionsItem>
      <DescriptionsItem label="创建时间">{{ displayValue(detail.createdAt) }}</DescriptionsItem>
      <DescriptionsItem label="更新时间">{{ displayValue(detail.updatedAt) }}</DescriptionsItem>
    </Descriptions>
  </Modal>
</template>
