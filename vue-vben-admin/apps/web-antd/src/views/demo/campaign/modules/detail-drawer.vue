<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { Descriptions, DescriptionsItem, Tag } from 'ant-design-vue';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { getCampaignDetail } from '#/api/demo/campaign';
import type { CampaignItem } from '#/api/demo/campaign/types';
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

/** 活动类型映射 */
const typeMap: Record<EnumValue, string> = {
  1: '免费',
  2: '付费',
  3: '公开',
  4: '私密',
};

/** 活动类型颜色 */
function getTypeColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2, 3, 4];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 投放渠道映射 */
const channelMap: Record<EnumValue, string> = {
  1: '官网',
  2: '小程序',
  3: '短信',
  4: '线下',
};

/** 投放渠道颜色 */
function getChannelColor(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [1, 2, 3, 4];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}

/** 是否公开映射 */
const isPublicMap: Record<EnumValue, string> = {
  0: '否',
  1: '是',
};

/** 是否公开颜色 */
function getIsPublicColor(val: EnumValue | null | undefined): string {
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
  1: '已发布',
  2: '已下架',
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
const detail = ref<CampaignItem | null>(null);
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
      modalApi.setState({ title: '体验活动详情' });
      try {
        const res = await getCampaignDetail(data.id);
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
      <DescriptionsItem label="活动编号">{{ displayValue(detail.campaignNo) }}</DescriptionsItem>
      <DescriptionsItem label="活动标题">{{ displayValue(detail.title) }}</DescriptionsItem>
      <DescriptionsItem label="横幅图">
        <img v-if="detail.banner && /^https?:\/\//i.test(detail.banner)" :src="detail.banner" style="max-width: 200px; max-height: 200px; object-fit: contain;" />
        <span v-else>-</span>
      </DescriptionsItem>
      <DescriptionsItem label="活动类型">
        <Tag :color="getTypeColor(detail.type)">{{ getEnumLabel(typeMap, detail.type) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="投放渠道">
        <Tag :color="getChannelColor(detail.channel)">{{ getEnumLabel(channelMap, detail.channel) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="预算金额">{{ detail.budgetAmount != null ? (detail.budgetAmount / 100).toFixed(2) : '-' }}</DescriptionsItem>
      <DescriptionsItem label="落地页URL">
        <a v-if="detail.landingURL && /^https?:\/\//i.test(detail.landingURL)" :href="detail.landingURL" target="_blank" rel="noreferrer noopener">{{ detail.landingURL }}</a>
        <span v-else-if="detail.landingURL">{{ detail.landingURL }}</span>
        <span v-else>-</span>
      </DescriptionsItem>
      <DescriptionsItem label="规则JSON">
        <pre style="max-height: 300px; overflow: auto; white-space: pre-wrap; word-break: break-all; margin: 0; font-size: 12px;">{{ (() => { const value = detail.ruleJSON; if (!value) return '-'; try { return JSON.stringify(JSON.parse(value), null, 2) } catch { return value } })() }}</pre>
      </DescriptionsItem>
      <DescriptionsItem label="活动介绍">
        <RichText v-if="detail.introContent" :value="detail.introContent" disabled :height="260" />
        <span v-else>-</span>
      </DescriptionsItem>
      <DescriptionsItem label="是否公开">
        <Tag :color="getIsPublicColor(detail.isPublic)">{{ getEnumLabel(isPublicMap, detail.isPublic) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem label="状态">
        <Tag :color="getStatusColor(detail.status)">{{ getEnumLabel(statusMap, detail.status) }}</Tag>
      </DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="租户">{{ detail.tenantName || '-' }}</DescriptionsItem>
      <DescriptionsItem v-if="isPlatformSuperAdmin" label="商户">{{ detail.merchantName || '-' }}</DescriptionsItem>
      <DescriptionsItem label="开始时间">{{ displayValue(detail.startAt) }}</DescriptionsItem>
      <DescriptionsItem label="结束时间">{{ displayValue(detail.endAt) }}</DescriptionsItem>
      <DescriptionsItem label="创建时间">{{ displayValue(detail.createdAt) }}</DescriptionsItem>
      <DescriptionsItem label="更新时间">{{ displayValue(detail.updatedAt) }}</DescriptionsItem>
    </Descriptions>
  </Modal>
</template>
