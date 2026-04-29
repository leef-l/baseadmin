<script setup lang="ts">
import { h, ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { message, Tooltip } from 'ant-design-vue';
import { QuestionCircleOutlined } from '@ant-design/icons-vue';
import {
  getLevelDetail,
  createLevel,
  updateLevel,
} from '#/api/member/level';
import type {
  LevelCreateParams,
  LevelUpdateParams
} from '#/api/member/level/types';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';

/** 到达后自动部署站点选项 */
const autoDeployOptions = [
  { label: '否', value: 0 },
  { label: '是', value: 1 },
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
      label: '等级名称',
      rules: 'required',
      componentProps: { placeholder: '请输入等级名称', maxlength: 50 },
    },
    {
      component: 'Input',
      fieldName: 'levelNo',
      label: tooltipLabel('等级编号', '越大越高'),
      componentProps: { placeholder: '请输入等级编号（越大越高）' },
    },
    {
      component: 'IconPicker',
      fieldName: 'icon',
      label: '等级图标',
      componentProps: { placeholder: '请选择图标' },
    },
    {
      component: 'Input',
      fieldName: 'durationDays',
      label: tooltipLabel('有效天数', '0=永久'),
      componentProps: { placeholder: '请输入有效天数（0=永久）' },
    },
    {
      component: 'Input',
      fieldName: 'needActiveCount',
      label: '升级所需有效用户数',
      componentProps: { placeholder: '请输入升级所需有效用户数' },
    },
    {
      component: 'Input',
      fieldName: 'needTeamTurnover',
      label: tooltipLabel('升级所需团队营业额', '分'),
      componentProps: { placeholder: '请输入升级所需团队营业额（分）' },
    },
    {
      component: 'Switch',
      fieldName: 'isTop',
      label: '是否最高等级',
      componentProps: { checkedValue: 1, unCheckedValue: 0 },
      defaultValue: 0,
    },
    {
      component: 'Select',
      fieldName: 'autoDeploy',
      label: '到达后自动部署站点',
      componentProps: { options: autoDeployOptions, placeholder: '请选择到达后自动部署站点', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Input',
      fieldName: 'remark',
      label: '等级说明',
      componentProps: { placeholder: '请输入等级说明', maxlength: 500 },
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
      | LevelCreateParams
      | undefined;
    if (!values) return;
    if (values.needTeamTurnover != null) {
      (values as any).needTeamTurnover = Math.round(Number(values.needTeamTurnover) * 100);
    }
    modalApi.lock();
    try {
      if (isEdit.value) {
        await updateLevel({ id: editId.value, ...values } as LevelUpdateParams);
        message.success('更新成功');
      } else {
        await createLevel(values);
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
      modalApi.setState({ title: '编辑会员等级配置' });
      try {
        const detail = await getLevelDetail(data.id);
        if (currentOpenToken !== openToken.value) {
          return;
        }
        if (detail) {
          const formData = { ...detail };
          if (formData.needTeamTurnover != null) {
            formData.needTeamTurnover = formData.needTeamTurnover / 100;
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
      modalApi.setState({ title: '新建会员等级配置' });
    }
  },
});
</script>

<template>
  <Modal class="w-[600px]">
    <Form />
  </Modal>
</template>
