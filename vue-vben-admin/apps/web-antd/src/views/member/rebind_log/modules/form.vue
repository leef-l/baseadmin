<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { message } from 'ant-design-vue';
import {
  getRebindLogDetail,
  createRebindLog,
  updateRebindLog,
} from '#/api/member/rebind_log';
import type {
  RebindLogCreateParams,
  RebindLogUpdateParams
} from '#/api/member/rebind_log/types';
import { getUserTree } from '#/api/member/user';
import type { UserItem } from '#/api/member/user/types';
import { getUserTree } from '#/api/member/user';
import type { UserItem } from '#/api/member/user/types';
import { getUserTree } from '#/api/member/user';
import type { UserItem } from '#/api/member/user/types';
import { getUsersList } from '#/api/system/users';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';
const userIDOptions = ref<UserItem[]>([]);
const oldParentIDOptions = ref<UserItem[]>([]);
const newParentIDOptions = ref<UserItem[]>([]);
const operatorIDOptions = ref<{ label: string; value: string | number }[]>([]);
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
      component: 'TreeSelect',
      fieldName: 'userID',
      label: '会员',
      componentProps: {
        treeData: [],
        fieldNames: { label: 'nickname', value: 'id', children: 'children' },
        placeholder: '请选择会员',
        allowClear: true,
        treeDefaultExpandAll: true,
        class: 'w-full',
      },
    },
    {
      component: 'TreeSelect',
      fieldName: 'oldParentID',
      label: '原上级',
      componentProps: {
        treeData: [],
        fieldNames: { label: 'nickname', value: 'id', children: 'children' },
        placeholder: '请选择原上级',
        allowClear: true,
        treeDefaultExpandAll: true,
        class: 'w-full',
      },
    },
    {
      component: 'TreeSelect',
      fieldName: 'newParentID',
      label: '新上级',
      componentProps: {
        treeData: [],
        fieldNames: { label: 'nickname', value: 'id', children: 'children' },
        placeholder: '请选择新上级',
        allowClear: true,
        treeDefaultExpandAll: true,
        class: 'w-full',
      },
    },
    {
      component: 'Input',
      fieldName: 'reason',
      label: '换绑原因',
      componentProps: { placeholder: '请输入换绑原因', maxlength: 500 },
    },
    {
      component: 'Select',
      fieldName: 'operatorID',
      label: '操作人',
      componentProps: { options: [], placeholder: '请选择操作人', allowClear: true, class: 'w-full' },
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
      | RebindLogCreateParams
      | undefined;
    if (!values) return;
    modalApi.lock();
    try {
      if (isEdit.value) {
        await updateRebindLog({ id: editId.value, ...values } as RebindLogUpdateParams);
        message.success('更新成功');
      } else {
        await createRebindLog(values);
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
    // 加载会员树形数据
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
    // 加载原上级树形数据
    try {
      const userRes = await getUserTree();
      if (currentOpenToken !== openToken.value) {
        return;
      }
      oldParentIDOptions.value = userRes ?? [];
      formApi.updateSchema([
        {
          fieldName: 'oldParentID',
          componentProps: { treeData: oldParentIDOptions.value },
        },
      ]);
    } catch {
      // ignore
    }
    // 加载新上级树形数据
    try {
      const userRes = await getUserTree();
      if (currentOpenToken !== openToken.value) {
        return;
      }
      newParentIDOptions.value = userRes ?? [];
      formApi.updateSchema([
        {
          fieldName: 'newParentID',
          componentProps: { treeData: newParentIDOptions.value },
        },
      ]);
    } catch {
      // ignore
    }
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
      modalApi.setState({ title: '编辑换绑上级日志' });
      try {
        const detail = await getRebindLogDetail(data.id);
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
      modalApi.setState({ title: '新建换绑上级日志' });
    }
  },
});
</script>

<template>
  <Modal class="w-[600px]">
    <Form />
  </Modal>
</template>
