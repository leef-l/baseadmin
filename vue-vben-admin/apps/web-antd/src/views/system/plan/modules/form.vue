<script setup lang="ts">
import { h, ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { message, Tooltip } from 'ant-design-vue';
import { QuestionCircleOutlined } from '@ant-design/icons-vue';
import {
  getPlanDetail,
  createPlan,
  updatePlan,
} from '#/api/system/plan';
import type {
  PlanCreateParams,
  PlanUpdateParams
} from '#/api/system/plan/types';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';
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
const isPlatformSuperAdmin = usePlatformSuperAdmin();
const isEdit = ref(false);
const editId = ref('');
const openToken = ref(0);

/** 表单配置 */
const [Form, formApi] = useVbenForm({
  showDefaultActions: false,
  schema: [
    {
      component: 'Input',
      fieldName: 'name',
      label: '套餐名称',
      rules: 'required',
      componentProps: { placeholder: '请输入套餐名称', maxlength: 80 },
    },
    {
      component: 'Input',
      fieldName: 'code',
      label: '套餐编码',
      rules: 'required',
      componentProps: { placeholder: '请输入套餐编码', maxlength: 50 },
    },
    {
      component: 'Input',
      fieldName: 'description',
      label: '套餐描述',
      componentProps: { placeholder: '请输入套餐描述', maxlength: 500 },
    },
    {
      component: 'Input',
      fieldName: 'userLimit',
      label: tooltipLabel('用户数量上限', '-1=无限'),
      componentProps: { placeholder: '请输入用户数量上限（-1=无限）' },
    },
    {
      component: 'Input',
      fieldName: 'storageMb',
      label: tooltipLabel('存储空间上限MB', '-1=无限'),
      componentProps: { placeholder: '请输入存储空间上限MB（-1=无限）' },
    },
    {
      component: 'Textarea',
      fieldName: 'features',
      label: tooltipLabel('功能开关列表', 'JSON数组'),
      componentProps: { placeholder: '请输入功能开关列表（JSON数组）', rows: 4, maxlength: 65535 },
    },
    {
      component: 'InputNumber',
      fieldName: 'price',
      label: tooltipLabel('价格', '分'),
      componentProps: { placeholder: '请输入价格（分）', class: 'w-full' },
    },
    {
      component: 'InputNumber',
      fieldName: 'sort',
      label: tooltipLabel('排序', '升序'),
      componentProps: { placeholder: '请输入排序（升序）', class: 'w-full' },
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
      | PlanCreateParams
      | undefined;
    if (!values) return;
    if (values.price != null) {
      (values as any).price = Math.round(Number(values.price) * 100);
    }
    modalApi.lock();
    try {
      if (isEdit.value) {
        await updatePlan({ id: editId.value, ...values } as PlanUpdateParams);
        message.success('更新成功');
      } else {
        await createPlan(values);
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
      modalApi.setState({ title: '编辑套餐表' });
      try {
        const detail = await getPlanDetail(data.id);
        if (currentOpenToken !== openToken.value) {
          return;
        }
        if (detail) {
          const formData = { ...detail };
          if (formData.price != null) {
            formData.price = formData.price / 100;
          }
          formApi.setValues(formData);
        }
      } catch {
        if (currentOpenToken === openToken.value) {
          message.error('获取详情失败');
        }
      }
    } else {
      isEdit.value = false;
      editId.value = '';
      modalApi.setState({ title: '新建套餐表' });
    }
  },
});
</script>

<template>
  <Modal class="w-[600px]">
    <Form />
  </Modal>
</template>
