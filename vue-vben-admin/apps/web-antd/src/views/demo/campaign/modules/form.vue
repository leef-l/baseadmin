<script setup lang="ts">
import { h, ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { message, Tooltip } from 'ant-design-vue';
import { QuestionCircleOutlined } from '@ant-design/icons-vue';
import {
  getCampaignDetail,
  createCampaign,
  updateCampaign,
} from '#/api/demo/campaign';
import type {
  CampaignCreateParams,
  CampaignUpdateParams
} from '#/api/demo/campaign/types';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';

/** 活动类型选项 */
const typeOptions = [
  { label: '免费', value: 1 },
  { label: '付费', value: 2 },
  { label: '公开', value: 3 },
  { label: '私密', value: 4 },
];

/** 投放渠道选项 */
const channelOptions = [
  { label: '官网', value: 1 },
  { label: '小程序', value: 2 },
  { label: '短信', value: 3 },
  { label: '线下', value: 4 },
];

/** 状态选项 */
const statusOptions = [
  { label: '草稿', value: 0 },
  { label: '已发布', value: 1 },
  { label: '已下架', value: 2 },
];
const tenantIDOptions = ref<{ label: string; value: string | number }[]>([]);
const merchantIDOptions = ref<{ label: string; value: string | number }[]>([]);
/** 渲染带 Tooltip 的表单 label */
function tooltipLabel(label: string, tip: string) {
  return () => h('span', {}, [
    label + ' ',
    h(Tooltip, { title: tip }, {
      default: () => h(QuestionCircleOutlined, { style: { color: '#999', marginLeft: '4px' } }),
    }),
  ]);
}

const emit = defineEmits<{ success: [] }>();
const isEdit = ref(false);
const editId = ref('');
const openToken = ref(0);

/** 表单配置 */
const [Form, formApi] = useVbenForm({
  showDefaultActions: false,
  schema: [
    {
      component: 'Input',
      fieldName: 'campaignNo',
      label: '活动编号',
      rules: 'required',
      componentProps: { placeholder: '请输入活动编号', maxlength: 50 },
    },
    {
      component: 'Input',
      fieldName: 'title',
      label: '活动标题',
      rules: 'required',
      componentProps: { placeholder: '请输入活动标题', maxlength: 120 },
    },
    {
      component: 'ImageUpload',
      fieldName: 'banner',
      label: '横幅图',
      componentProps: { maxCount: 1 },
    },
    {
      component: 'Select',
      fieldName: 'type',
      label: '活动类型',
      componentProps: { options: typeOptions, placeholder: '请选择活动类型', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Select',
      fieldName: 'channel',
      label: '投放渠道',
      componentProps: { options: channelOptions, placeholder: '请选择投放渠道', allowClear: true, class: 'w-full' },
    },
    {
      component: 'InputNumber',
      fieldName: 'budgetAmount',
      label: tooltipLabel('预算金额', '分'),
      componentProps: { placeholder: '请输入预算金额（分）', class: 'w-full' },
    },
    {
      component: 'Input',
      fieldName: 'landingURL',
      label: '落地页URL',
      componentProps: { placeholder: '请输入URL地址', maxlength: 500, addonBefore: 'https://' },
    },
    {
      component: 'JsonEditor',
      fieldName: 'ruleJSON',
      label: '规则JSON',
      formItemClass: 'col-span-full',
    },
    {
      component: 'RichText',
      fieldName: 'introContent',
      label: '活动介绍',
      formItemClass: 'col-span-full',
    },
    {
      component: 'DatePicker',
      fieldName: 'startAt',
      label: '开始时间',
      componentProps: { showTime: true, placeholder: '请选择开始时间', class: 'w-full', valueFormat: 'YYYY-MM-DD HH:mm:ss' },
    },
    {
      component: 'DatePicker',
      fieldName: 'endAt',
      label: '结束时间',
      componentProps: { showTime: true, placeholder: '请选择结束时间', class: 'w-full', valueFormat: 'YYYY-MM-DD HH:mm:ss' },
    },
    {
      component: 'Switch',
      fieldName: 'isPublic',
      label: '是否公开',
      componentProps: { checkedValue: 1, unCheckedValue: 0 },
      defaultValue: 1,
    },
    {
      component: 'RadioGroup',
      fieldName: 'status',
      label: '状态',
      componentProps: { options: statusOptions },
    },
    {
      component: 'Select',
      fieldName: 'tenantID',
      label: '租户',
      componentProps: { options: [], placeholder: '请选择租户', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Select',
      fieldName: 'merchantID',
      label: '商户',
      componentProps: { options: [], placeholder: '请选择商户', allowClear: true, class: 'w-full' },
    },
  ],
});

/** Modal 配置 */
const [Modal, modalApi] = useVbenModal({
  fullscreenButton: false,
  onCancel() {
    modalApi.close();
  },
  onConfirm: async () => {
    const values = await formApi.validateAndSubmitForm() as
      | CampaignCreateParams
      | undefined;
    if (!values) return;
    modalApi.lock();
    try {
      if (isEdit.value) {
        await updateCampaign({ id: editId.value, ...values } as CampaignUpdateParams);
        message.success('更新成功');
      } else {
        await createCampaign(values);
        message.success('创建成功');
      }
      emit('success');
      modalApi.close();
    } finally {
      modalApi.lock(false);
    }
  },
  async onOpenChange(isOpen: boolean) {
    if (!isOpen) {
      openToken.value += 1;
      return;
    }

    const currentOpenToken = ++openToken.value;
    formApi.resetForm();
    const data = modalApi.getData<{ id?: string }>();
    // 加载租户选项
    try {
      const tenantRes = await getTenantList({ pageNum: 1, pageSize: 1000 });
      if (currentOpenToken !== openToken.value) {
        return;
      }
      tenantIDOptions.value = (tenantRes?.list ?? []).map((item: any) => ({
        label: item.name || item.id,
        value: item.id,
      }));
      formApi.updateSchema([
        {
          fieldName: 'tenantID',
          componentProps: { options: tenantIDOptions.value },
        },
      ]);
    } catch {
      // ignore
    }
    // 加载商户选项
    try {
      const merchantRes = await getMerchantList({ pageNum: 1, pageSize: 1000 });
      if (currentOpenToken !== openToken.value) {
        return;
      }
      merchantIDOptions.value = (merchantRes?.list ?? []).map((item: any) => ({
        label: item.name || item.id,
        value: item.id,
      }));
      formApi.updateSchema([
        {
          fieldName: 'merchantID',
          componentProps: { options: merchantIDOptions.value },
        },
      ]);
    } catch {
      // ignore
    }
    if (currentOpenToken !== openToken.value) {
      return;
    }
    if (data?.id) {
      isEdit.value = true;
      editId.value = data.id;
      modalApi.setState({ title: '编辑体验活动' });
      try {
        const detail = await getCampaignDetail(data.id);
        if (currentOpenToken !== openToken.value) {
          return;
        }
        if (detail) {
          formApi.setValues(detail);
        }
      } catch {
        if (currentOpenToken === openToken.value) {
          message.error('获取详情失败');
        }
      }
    } else {
      isEdit.value = false;
      editId.value = '';
      modalApi.setState({ title: '新建体验活动' });
    }
  },
});
</script>

<template>
  <Modal class="w-[800px]">
    <Form />
  </Modal>
</template>
