<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { message } from 'ant-design-vue';
import {
  getCustomerDetail,
  createCustomer,
  updateCustomer,
} from '#/api/demo/customer';
import type {
  CustomerCreateParams,
  CustomerUpdateParams
} from '#/api/demo/customer/types';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';

/** 性别选项 */
const genderOptions = [
  { label: '未知', value: 0 },
  { label: '男', value: 1 },
  { label: '女', value: 2 },
];

/** 等级选项 */
const levelOptions = [
  { label: '普通', value: 1 },
  { label: 'VIP', value: 2 },
  { label: '付费', value: 3 },
  { label: '冻结', value: 4 },
];

/** 来源选项 */
const sourceTypeOptions = [
  { label: '官网', value: 1 },
  { label: '小程序', value: 2 },
  { label: '线下', value: 3 },
  { label: '导入', value: 4 },
];
const tenantIDOptions = ref<{ label: string; value: string | number }[]>([]);
const merchantIDOptions = ref<{ label: string; value: string | number }[]>([]);

const emit = defineEmits<{ success: [] }>();
const isPlatformSuperAdmin = usePlatformSuperAdmin();
const isEdit = ref(false);
const editId = ref('');
const openToken = ref(0);

/** 表单配置 */
const [Form, formApi] = useVbenForm({
  showDefaultActions: false,
  schema: [
    {
      component: 'ImageUpload',
      fieldName: 'avatar',
      label: '头像',
      componentProps: { maxCount: 1 },
    },
    {
      component: 'Input',
      fieldName: 'name',
      label: '客户名称',
      rules: 'required',
      componentProps: { placeholder: '请输入客户名称', maxlength: 80 },
    },
    {
      component: 'Input',
      fieldName: 'customerNo',
      label: '客户编号',
      rules: 'required',
      componentProps: { placeholder: '请输入客户编号', maxlength: 50 },
    },
    {
      component: 'Input',
      fieldName: 'phone',
      label: '联系电话',
      rules: 'phone',
      componentProps: { placeholder: '请输入联系电话', maxlength: 30 },
    },
    {
      component: 'Input',
      fieldName: 'email',
      label: '邮箱',
      rules: 'email',
      componentProps: { placeholder: '请输入邮箱', maxlength: 120 },
    },
    {
      component: 'Select',
      fieldName: 'gender',
      label: '性别',
      componentProps: { options: genderOptions, placeholder: '请选择性别', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Select',
      fieldName: 'level',
      label: '等级',
      componentProps: { options: levelOptions, placeholder: '请选择等级', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Select',
      fieldName: 'sourceType',
      label: '来源',
      componentProps: { options: sourceTypeOptions, placeholder: '请选择来源', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Switch',
      fieldName: 'isVip',
      label: '是否VIP',
      componentProps: { checkedValue: 1, unCheckedValue: 0 },
      defaultValue: 0,
    },
    {
      component: 'DatePicker',
      fieldName: 'registeredAt',
      label: '注册时间',
      componentProps: { showTime: true, placeholder: '请选择注册时间', class: 'w-full', valueFormat: 'YYYY-MM-DD HH:mm:ss' },
    },
    {
      component: 'Textarea',
      fieldName: 'remark',
      label: '备注',
      componentProps: { placeholder: '请输入备注', rows: 4, maxlength: 65535 },
    },
    {
      component: 'Switch',
      fieldName: 'status',
      label: '状态',
      componentProps: { checkedValue: 1, unCheckedValue: 0 },
      defaultValue: 1,
    },
    {
      component: 'Select',
      fieldName: 'tenantID',
      label: '租户',
      ifShow: () => isPlatformSuperAdmin.value,
      componentProps: { options: [], placeholder: '请选择租户', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Select',
      fieldName: 'merchantID',
      label: '商户',
      ifShow: () => isPlatformSuperAdmin.value,
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
      | CustomerCreateParams
      | undefined;
    if (!values) return;
    modalApi.lock();
    try {
      if (isEdit.value) {
        await updateCustomer({ id: editId.value, ...values } as CustomerUpdateParams);
        message.success('更新成功');
      } else {
        await createCustomer(values);
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
    if (isPlatformSuperAdmin.value) {
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
    }
    if (isPlatformSuperAdmin.value) {
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
    }
    if (currentOpenToken !== openToken.value) {
      return;
    }
    if (data?.id) {
      isEdit.value = true;
      editId.value = data.id;
      modalApi.setState({ title: '编辑体验客户' });
      try {
        const detail = await getCustomerDetail(data.id);
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
      modalApi.setState({ title: '新建体验客户' });
    }
  },
});
</script>

<template>
  <Modal class="w-[600px]">
    <Form />
  </Modal>
</template>
