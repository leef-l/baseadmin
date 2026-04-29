<script setup lang="ts">
import { ref } from 'vue';
import { isPlatformSuperAdminUser } from '@/utils/auth-scope';
import { useVbenModal } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { message } from 'ant-design-vue';
import {
  getAppointmentDetail,
  createAppointment,
  updateAppointment,
} from '#/api/demo/appointment';
import type {
  AppointmentCreateParams,
  AppointmentUpdateParams
} from '#/api/demo/appointment/types';
import { getCustomerList } from '#/api/demo/customer';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';

/** 状态选项 */
const statusOptions = [
  { label: '待确认', value: 0 },
  { label: '已确认', value: 1 },
  { label: '已完成', value: 2 },
  { label: '已取消', value: 3 },
];
const customerIDOptions = ref<{ label: string; value: string | number }[]>([]);
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
      fieldName: 'appointmentNo',
      label: '预约编号',
      rules: 'required',
      componentProps: { placeholder: '请输入预约编号', maxlength: 50 },
    },
    {
      component: 'Select',
      fieldName: 'customerID',
      label: '客户',
      componentProps: { options: [], placeholder: '请选择客户', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Input',
      fieldName: 'subject',
      label: '预约主题',
      rules: 'required',
      componentProps: { placeholder: '请输入预约主题', maxlength: 120 },
    },
    {
      component: 'DatePicker',
      fieldName: 'appointmentAt',
      label: '预约时间',
      componentProps: { showTime: true, placeholder: '请选择预约时间', class: 'w-full', valueFormat: 'YYYY-MM-DD HH:mm:ss' },
    },
    {
      component: 'Input',
      fieldName: 'contactPhone',
      label: '联系电话',
      componentProps: { placeholder: '请输入联系电话', maxlength: 30 },
    },
    {
      component: 'Input',
      fieldName: 'address',
      label: '预约地址',
      componentProps: { placeholder: '请输入预约地址', maxlength: 255 },
    },
    {
      component: 'Textarea',
      fieldName: 'remark',
      label: '备注',
      componentProps: { placeholder: '请输入备注', rows: 4, maxlength: 65535 },
    },
    {
      component: 'RadioGroup',
      fieldName: 'status',
      label: '状态',
      componentProps: { options: statusOptions },
    },
    {
      component: 'Select',
      ifShow: () => isPlatformSuperAdminUser(),
      fieldName: 'tenantID',
      label: '租户',
      componentProps: { options: [], placeholder: '请选择租户', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Select',
      ifShow: () => isPlatformSuperAdminUser(),
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
      | AppointmentCreateParams
      | undefined;
    if (!values) return;
    modalApi.lock();
    try {
      if (isEdit.value) {
        await updateAppointment({ id: editId.value, ...values } as AppointmentUpdateParams);
        message.success('更新成功');
      } else {
        await createAppointment(values);
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
    // 加载客户选项
    try {
      const customerRes = await getCustomerList({ pageNum: 1, pageSize: 1000 });
      if (currentOpenToken !== openToken.value) {
        return;
      }
      customerIDOptions.value = (customerRes?.list ?? []).map((item: any) => ({
        label: item.name || item.id,
        value: item.id,
      }));
      formApi.updateSchema([
        {
          fieldName: 'customerID',
          componentProps: { options: customerIDOptions.value },
        },
      ]);
    } catch {
      // ignore
    }
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
      modalApi.setState({ title: '编辑体验预约' });
      try {
        const detail = await getAppointmentDetail(data.id);
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
      modalApi.setState({ title: '新建体验预约' });
    }
  },
});
</script>

<template>
  <Modal class="w-[600px]">
    <Form />
  </Modal>
</template>
