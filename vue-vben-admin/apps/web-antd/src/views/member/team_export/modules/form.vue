<script setup lang="ts">
import { h, ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { message, Tooltip } from 'ant-design-vue';
import { QuestionCircleOutlined } from '@ant-design/icons-vue';
import {
  getTeamExportDetail,
  createTeamExport,
  updateTeamExport,
} from '#/api/member/team_export';
import type {
  TeamExportCreateParams,
  TeamExportUpdateParams
} from '#/api/member/team_export/types';
import { getUserTree } from '#/api/member/user';
import type { UserItem } from '#/api/member/user/types';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';

/** 导出类型选项 */
const exportTypeOptions = [
  { label: '手动导出', value: 1 },
  { label: '自动升级导出', value: 2 },
];

/** 部署状态选项 */
const deployStatusOptions = [
  { label: '未部署', value: 0 },
  { label: '部署中', value: 1 },
  { label: '已部署', value: 2 },
  { label: '部署失败', value: 3 },
];
const userIDOptions = ref<UserItem[]>([]);
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
      component: 'TreeSelect',
      fieldName: 'userID',
      label: '目标会员',
      componentProps: {
        treeData: [],
        fieldNames: { label: 'nickname', value: 'id', children: 'children' },
        placeholder: '请选择目标会员',
        allowClear: true,
        treeDefaultExpandAll: true,
        class: 'w-full',
      },
    },
    {
      component: 'Input',
      fieldName: 'teamMemberCount',
      label: '团队成员数',
      componentProps: { placeholder: '请输入团队成员数' },
    },
    {
      component: 'Select',
      fieldName: 'exportType',
      label: '导出类型',
      componentProps: { options: exportTypeOptions, placeholder: '请选择导出类型', allowClear: true, class: 'w-full' },
    },
    {
      component: 'InputUrl',
      fieldName: 'fileURL',
      label: '导出文件地址',
      componentProps: { placeholder: '请输入完整URL地址（含 http:// 或 https://）', maxlength: 500 },
    },
    {
      component: 'Input',
      fieldName: 'fileSize',
      label: tooltipLabel('文件大小', '字节'),
      componentProps: { placeholder: '请输入文件大小（字节）' },
    },
    {
      component: 'Select',
      fieldName: 'deployStatus',
      label: '部署状态',
      componentProps: { options: deployStatusOptions, placeholder: '请选择部署状态', allowClear: true, class: 'w-full' },
    },
    {
      component: 'Input',
      fieldName: 'deployDomain',
      label: '部署域名',
      componentProps: { placeholder: '请输入部署域名', maxlength: 200 },
    },
    {
      component: 'DatePicker',
      fieldName: 'deployedAt',
      label: '部署完成时间',
      componentProps: { showTime: true, placeholder: '请选择部署完成时间', class: 'w-full', valueFormat: 'YYYY-MM-DD HH:mm:ss' },
    },
    {
      component: 'Input',
      fieldName: 'remark',
      label: '备注',
      componentProps: { placeholder: '请输入备注', maxlength: 500 },
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
      | TeamExportCreateParams
      | undefined;
    if (!values) return;
    modalApi.lock();
    try {
      if (isEdit.value) {
        await updateTeamExport({ id: editId.value, ...values } as TeamExportUpdateParams);
        message.success('更新成功');
      } else {
        await createTeamExport(values);
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
    // 加载目标会员树形数据
    try {
      const userRes = await getUserTree();
      if (currentOpenToken !== openToken.value) {
        return;
      }
      userIDOptions.value = userRes ?? [];
      formApi.updateSchema([
        {
          fieldName: 'userID',
          componentProps: { treeData: userIDOptions.value },
        },
      ]);
    } catch {
      // ignore
    }
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
      modalApi.setState({ title: '编辑团队数据导出' });
      try {
        const detail = await getTeamExportDetail(data.id);
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
      modalApi.setState({ title: '新建团队数据导出' });
    }
  },
});
</script>

<template>
  <Modal class="w-[600px]">
    <Form />
  </Modal>
</template>
