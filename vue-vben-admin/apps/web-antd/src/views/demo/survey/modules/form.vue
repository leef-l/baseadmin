<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { message } from 'ant-design-vue';
import {
  getSurveyDetail,
  createSurvey,
  updateSurvey,
} from '#/api/demo/survey';
import type {
  SurveyCreateParams,
  SurveyUpdateParams
} from '#/api/demo/survey/types';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';

/** 状态选项 */
const statusOptions = [
  { label: '草稿', value: 0 },
  { label: '已发布', value: 1 },
  { label: '已下架', value: 2 },
];
const tenantIDOptions = ref<{ label: string; value: string | number }[]>([]);
const merchantIDOptions = ref<{ label: string; value: string | number }[]>([]);

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
      fieldName: 'surveyNo',
      label: '问卷编号',
      rules: 'required',
      componentProps: { placeholder: '请输入问卷编号', maxlength: 50 },
    },
    {
      component: 'Input',
      fieldName: 'title',
      label: '问卷标题',
      rules: 'required',
      componentProps: { placeholder: '请输入问卷标题', maxlength: 120 },
    },
    {
      component: 'ImageUpload',
      fieldName: 'poster',
      label: '海报',
      componentProps: { maxCount: 1 },
    },
    {
      component: 'JsonEditor',
      fieldName: 'questionJSON',
      label: '问题JSON',
      formItemClass: 'col-span-full',
    },
    {
      component: 'RichText',
      fieldName: 'introContent',
      label: '问卷介绍',
      formItemClass: 'col-span-full',
    },
    {
      component: 'DatePicker',
      fieldName: 'publishAt',
      label: '发布时间',
      componentProps: { showTime: true, placeholder: '请选择发布时间', class: 'w-full', valueFormat: 'YYYY-MM-DD HH:mm:ss' },
    },
    {
      component: 'DatePicker',
      fieldName: 'expireAt',
      label: '过期时间',
      componentProps: { showTime: true, placeholder: '请选择过期时间', class: 'w-full', valueFormat: 'YYYY-MM-DD HH:mm:ss' },
    },
    {
      component: 'Switch',
      fieldName: 'isAnonymous',
      label: '是否匿名',
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
      | SurveyCreateParams
      | undefined;
    if (!values) return;
    modalApi.lock();
    try {
      if (isEdit.value) {
        await updateSurvey({ id: editId.value, ...values } as SurveyUpdateParams);
        message.success('更新成功');
      } else {
        await createSurvey(values);
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
      modalApi.setState({ title: '编辑体验问卷' });
      try {
        const detail = await getSurveyDetail(data.id);
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
      modalApi.setState({ title: '新建体验问卷' });
    }
  },
});
</script>

<template>
  <Modal class="w-[800px]">
    <Form />
  </Modal>
</template>
