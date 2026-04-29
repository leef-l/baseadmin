<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { Descriptions, DescriptionsItem, Tag } from 'ant-design-vue';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { getWarehouseListingDetail } from '#/api/member/warehouse_listing';
import type { WarehouseListingItem } from '#/api/member/warehouse_listing/types';

/** 标签颜色池 */
const TAG_COLORS = ['green', 'red', 'blue', 'orange', 'cyan', 'purple', 'geekblue', 'magenta'];

type EnumValue = number | string;

function getEnumLabel(map: Record<EnumValue, string>, value: EnumValue | null | undefined) {
  if (value === null || value === undefined || value === '') {
    return '-';
  }
  return map[value] ?? String(value);
}

/** 挂卖状态映射 */
const listingStatusMap: Record<EnumValue, string> = {
  1: '挂卖中',
  2: '已售出',
  3: '已取消',
};

/** 挂卖状态颜色 */
function getListingStatusColor(val: EnumValue | null | undefined): string {
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
const detail = ref<WarehouseListingItem | null>(null);
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
      modalApi.setState({ title: '仓库挂卖记录详情' });
      try {
        const res = await getWarehouseListingDetail(data.id);
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
      <DescriptionsItem label="仓库商品">{{ detail.warehouseGoodsTitle || '-' }}</DescriptionsItem>
      <DescriptionsItem label="卖家">{{ detail.userNickname || '-' }}</DescriptionsItem>
      <DescriptionsItem label="挂卖价格">{{ detail.listingPrice != null ? (detail.listingPrice / 100).toFixed(2) : '-' }}</DescriptionsItem>
      <DescriptionsItem label="挂卖状态">
        <Tag :color="getListingStatusColor(detail.listingStatus)">{{ getEnumLabel(listingStatusMap, detail.listingStatus) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="备注">{{ displayValue(detail.remark) }}</DescriptionsItem>
      <DescriptionsItem label="状态">
        <Tag :color="getStatusColor(detail.status)">{{ getEnumLabel(statusMap, detail.status) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="租户">{{ detail.tenantName || '-' }}</DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="商户">{{ detail.merchantName || '-' }}</DescriptionsItem>
      <DescriptionsItem label="挂卖时间">{{ displayValue(detail.listedAt) }}</DescriptionsItem>
      <DescriptionsItem label="售出时间">{{ displayValue(detail.soldAt) }}</DescriptionsItem>
      <DescriptionsItem label="创建时间">{{ displayValue(detail.createdAt) }}</DescriptionsItem>
      <DescriptionsItem label="更新时间">{{ displayValue(detail.updatedAt) }}</DescriptionsItem>
    </Descriptions>
  </Modal>
</template>
