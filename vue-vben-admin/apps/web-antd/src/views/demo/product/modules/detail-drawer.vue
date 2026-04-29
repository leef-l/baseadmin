<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { Descriptions, DescriptionsItem, Tag } from 'ant-design-vue';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { getProductDetail } from '#/api/demo/product';
import type { ProductItem } from '#/api/demo/product/types';
import RichText from '#/components/tinymce/index.vue';

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
  1: '普通',
  2: '置顶',
  3: '推荐',
  4: '热门',
};

/** 类型颜色 */
function getTypeColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2, 3, 4];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 是否推荐映射 */
const isRecommendMap: Record<EnumValue, string> = {
  0: '否',
  1: '是',
};

/** 是否推荐颜色 */
function getIsRecommendColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [0, 1];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 状态映射 */
const statusMap: Record<EnumValue, string> = {
  0: '草稿',
  1: '上架',
  2: '下架',
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
const detail = ref<ProductItem | null>(null);
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
      modalApi.setState({ title: '体验商品详情' });
      try {
        const res = await getProductDetail(data.id);
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
  <Modal class="w-[800px]">
    <Descriptions v-if="detail" bordered :column="1" size="small">
      <DescriptionsItem label="ID">{{ detail.id }}</DescriptionsItem>
      <DescriptionsItem label="商品分类">{{ detail.categoryName || '-' }}</DescriptionsItem>
      <DescriptionsItem label="SKU编号">{{ displayValue(detail.skuNo) }}</DescriptionsItem>
      <DescriptionsItem label="商品名称">{{ displayValue(detail.name) }}</DescriptionsItem>
      <DescriptionsItem label="封面">
        <img v-if="detail.cover && /^https?:\/\//i.test(detail.cover)" :src="detail.cover" style="max-width: 200px; max-height: 200px; object-fit: contain;" />
        <span v-else>-</span>
      </DescriptionsItem>
      <DescriptionsItem label="说明书文件">
        <a v-if="detail.manualFile && /^https?:\/\//i.test(detail.manualFile)" :href="detail.manualFile" target="_blank" rel="noreferrer noopener">查看文件</a>
        <span v-else-if="detail.manualFile">{{ detail.manualFile }}</span>
        <span v-else>-</span>
      </DescriptionsItem>
      <DescriptionsItem label="详情内容">
        <RichText v-if="detail.detailContent" :value="detail.detailContent" disabled :height="260" />
        <span v-else>-</span>
      </DescriptionsItem>
      <DescriptionsItem label="规格JSON">
        <pre style="max-height: 300px; overflow: auto; white-space: pre-wrap; word-break: break-all; margin: 0; font-size: 12px;">{{ (() => { const value = detail.specJSON; if (!value) return '-'; try { return JSON.stringify(JSON.parse(value), null, 2) } catch { return value } })() }}</pre>
      </DescriptionsItem>
      <DescriptionsItem label="官网URL">
        <a v-if="detail.websiteURL && /^https?:\/\//i.test(detail.websiteURL)" :href="detail.websiteURL" target="_blank" rel="noreferrer noopener">{{ detail.websiteURL }}</a>
        <span v-else-if="detail.websiteURL">{{ detail.websiteURL }}</span>
        <span v-else>-</span>
      </DescriptionsItem>
      <DescriptionsItem label="类型">
        <Tag :color="getTypeColor(detail.type)">{{ getEnumLabel(typeMap, detail.type) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="是否推荐">
        <Tag :color="getIsRecommendColor(detail.isRecommend)">{{ getEnumLabel(isRecommendMap, detail.isRecommend) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="销售价">{{ detail.salePrice != null ? (detail.salePrice / 100).toFixed(2) : '-' }}</DescriptionsItem>
      <DescriptionsItem label="库存数量">{{ displayValue(detail.stockNum) }}</DescriptionsItem>
      <DescriptionsItem label="重量">{{ displayValue(detail.weightNum) }}</DescriptionsItem>
      <DescriptionsItem label="排序">{{ displayValue(detail.sort) }}</DescriptionsItem>
      <DescriptionsItem label="图标">{{ displayValue(detail.icon) }}</DescriptionsItem>
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
