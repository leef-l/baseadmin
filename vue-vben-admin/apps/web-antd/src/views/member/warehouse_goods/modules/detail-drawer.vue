<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { Descriptions, DescriptionsItem, Tag } from 'ant-design-vue';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { getWarehouseGoodsDetail } from '#/api/member/warehouse_goods';
import type { WarehouseGoodsItem } from '#/api/member/warehouse_goods/types';

/** 标签颜色池 */
const TAG_COLORS = ['green', 'red', 'blue', 'orange', 'cyan', 'purple', 'geekblue', 'magenta'];

type EnumValue = number | string;

function getEnumLabel(map: Record<EnumValue, string>, value: EnumValue | null | undefined) {
  if (value === null || value === undefined || value === '') {
    return '-';
  }
  return map[value] ?? String(value);
}

/** 商品状态映射 */
const goodsStatusMap: Record<EnumValue, string> = {
  1: '持有中',
  2: '挂卖中',
  3: '交易中',
};

/** 商品状态颜色 */
function getGoodsStatusColor(val: EnumValue | null | undefined): string {
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
const detail = ref<WarehouseGoodsItem | null>(null);
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
      modalApi.setState({ title: '仓库商品详情' });
      try {
        const res = await getWarehouseGoodsDetail(data.id);
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
      <DescriptionsItem label="商品编号">{{ displayValue(detail.goodsNo) }}</DescriptionsItem>
      <DescriptionsItem label="商品名称">{{ displayValue(detail.title) }}</DescriptionsItem>
      <DescriptionsItem label="商品封面">
        <img v-if="detail.cover && /^https?:\/\//i.test(detail.cover)" :src="detail.cover" style="max-width: 200px; max-height: 200px; object-fit: contain;" />
        <span v-else>-</span>
      </DescriptionsItem>
      <DescriptionsItem label="初始价格">{{ detail.initPrice != null ? (detail.initPrice / 100).toFixed(2) : '-' }}</DescriptionsItem>
      <DescriptionsItem label="当前价格">{{ detail.currentPrice != null ? (detail.currentPrice / 100).toFixed(2) : '-' }}</DescriptionsItem>
      <DescriptionsItem label="每次加价比例">{{ detail.priceRiseRate != null ? `${detail.priceRiseRate}%` : '-' }}</DescriptionsItem>
      <DescriptionsItem label="平台扣除比例">{{ detail.platformFeeRate != null ? `${detail.platformFeeRate}%` : '-' }}</DescriptionsItem>
      <DescriptionsItem label="当前持有人">{{ detail.userNickname || '-' }}</DescriptionsItem>
      <DescriptionsItem label="流转次数">{{ displayValue(detail.tradeCount) }}</DescriptionsItem>
      <DescriptionsItem label="商品状态">
        <Tag :color="getGoodsStatusColor(detail.goodsStatus)">{{ getEnumLabel(goodsStatusMap, detail.goodsStatus) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="备注">{{ displayValue(detail.remark) }}</DescriptionsItem>
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
