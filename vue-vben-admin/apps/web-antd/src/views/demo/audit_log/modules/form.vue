<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { message } from 'ant-design-vue';
import {
  getAuditLogDetail,
  createAuditLog,
  updateAuditLog,
} from '#/api/demo/audit_log';
import type {
  AuditLogCreateParams,
  AuditLogUpdateParams
} from '#/api/demo/audit_log/types';
import { getUsersList } from '#/api/system/users';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';

/** 动作选项 */
const actionOptions = [
  { label: '创建', value: 1 },
  { label: '修改', value: 2 },
  { label: '删除', value: 3 },
  { label: '导出', value: 4 },
  { label: '导入', value: 5 },
];

/** 对象类型选项 */
const targetTypeOptions = [
  { label: '客户', value: 1 },
  { label: '商品', value: 2 },
  { label: '订单', value: 3 },
  { label: '工单', value: 4 },
];

/** 结果选项 */
const resultOptions = [
  { label: '失败', value: 0 },
  { label: '成功', value: 1 },
];
const operatorIDOptions = ref<{ label: string; value: string | number }[]>([]);
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
      fieldName: 'logNo',
      label: '日志编号',
      rules: 'required',
      componentProps: { placeholder: '请输入日志编号', maxlength: 50 },
    },
    {
      component: 'Select',
      fieldName: 'operatorID',
      label: '操作人',
      componentProps: { options: [], placeholder: '请选择操作人', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Select',
      fieldName: 'action',
      label: '动作',
      componentProps: { options: actionOptions, placeholder: '请选择动作', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Select',
      fieldName: 'targetType',
      label: '对象类型',
      componentProps: { options: targetTypeOptions, placeholder: '请选择对象类型', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Input',
      fieldName: 'targetCode',
      label: '对象编号',
      componentProps: { placeholder: '请输入对象编号', maxlength: 80 },
    },
    {
      component: 'JsonEditor',
      fieldName: 'requestJSON',
      label: '请求JSON',
      formItemClass: 'col-span-full',
    },
    {
      component: 'Select',
      fieldName: 'result',
      label: '结果',
      componentProps: { options: resultOptions, placeholder: '请选择结果', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Input',
      fieldName: 'clientIP',
      label: '客户端IP',
      componentProps: { placeholder: '请输入客户端IP', maxlength: 50 },
    },
    {
      component: 'DatePicker',
      fieldName: 'occurredAt',
      label: '发生时间',
      componentProps: { showTime: true, placeholder: '请选择发生时间', class: 'w-full', valueFormat: 'YYYY-MM-DD HH:mm:ss' },
    },
    {
      component: 'Textarea',
      fieldName: 'remark',
      label: '备注',
      componentProps: { placeholder: '请输入备注', rows: 4, maxlength: 65535 },
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
      | AuditLogCreateParams
      | undefined;
    if (!values) return;
    modalApi.lock();
    try {
      if (isEdit.value) {
        await updateAuditLog({ id: editId.value, ...values } as AuditLogUpdateParams);
        message.success('更新成功');
      } else {
        await createAuditLog(values);
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
    // 加载操作人选项
    try {
      const usersRes = await getUsersList({ pageNum: 1, pageSize: 1000 });
      if (currentOpenToken !== openToken.value) {
        return;
      }
      operatorIDOptions.value = (usersRes?.list ?? []).map((item: any) => ({
        label: item.username || item.id,
        value: item.id,
      }));
      formApi.updateSchema([
        {
          fieldName: 'operatorID',
          componentProps: { options: operatorIDOptions.value },
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
      modalApi.setState({ title: '编辑体验审计日志' });
      try {
        const detail = await getAuditLogDetail(data.id);
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
      modalApi.setState({ title: '新建体验审计日志' });
    }
  },
});
</script>

<template>
  <Modal class="w-[800px]">
    <Form />
  </Modal>
</template>
