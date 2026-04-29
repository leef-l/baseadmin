<script setup lang="ts">
import { ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { message } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import { getDeptTree } from '#/api/system/dept';
import type { DeptItem } from '#/api/system/dept/types';
import { getMerchantList } from '#/api/system/merchant';
import { getRoleTree } from '#/api/system/role';
import { getTenantList } from '#/api/system/tenant';
import {
  createUsers,
  getUsersDetail,
  updateUsers,
} from '#/api/system/users';
import type {
  UsersCreateParams,
  UsersItem,
  UsersUpdateParams,
} from '#/api/system/users/types';

import { normalizeRoleIds } from './role-ids';

interface SelectOption {
  label: string;
  value: string;
}

const emit = defineEmits<{ success: [] }>();
const isEdit = ref(false);
const editId = ref('');
const deptTreeData = ref<DeptItem[]>([]);
const tenantOptions = ref<SelectOption[]>([]);
const merchantOptions = ref<SelectOption[]>([]);
const openToken = ref(0);

/** 表单配置 */
const [Form, formApi] = useVbenForm({
  showDefaultActions: false,
  schema: [
    {
      component: 'Input',
      fieldName: 'username',
      label: '登录用户名',
      rules: 'required',
      componentProps: { placeholder: '请输入登录用户名', maxlength: 50 },
    },
    {
      component: 'InputPassword',
      fieldName: 'password',
      label: '密码',
      rules: 'required',
      componentProps: { placeholder: '请输入密码' },
      dependencies: {
        triggerFields: ['_mode'],
        if: () => !isEdit.value,
      },
    },
    {
      component: 'Input',
      fieldName: 'nickname',
      label: '昵称',
      componentProps: { placeholder: '请输入昵称', maxlength: 50 },
    },
    {
      component: 'Input',
      fieldName: 'email',
      label: '邮箱',
      componentProps: { placeholder: '请输入邮箱', maxlength: 100 },
    },
    {
      component: 'Select',
      fieldName: 'tenantId',
      label: '所属租户',
      componentProps: {
        options: tenantOptions.value,
        placeholder: '请选择所属租户',
        allowClear: true,
        class: 'w-full',
      },
    },
    {
      component: 'Select',
      fieldName: 'merchantId',
      label: '所属商户',
      componentProps: {
        options: merchantOptions.value,
        placeholder: '请选择所属商户',
        allowClear: true,
        class: 'w-full',
      },
    },
    {
      component: 'TreeSelect',
      fieldName: 'deptId',
      label: '所属部门',
      componentProps: {
        treeData: [],
        fieldNames: { label: 'title', value: 'id', children: 'children' },
        placeholder: '请选择所属部门',
        allowClear: true,
        treeDefaultExpandAll: true,
        class: 'w-full',
      },
    },
    {
      component: 'TreeSelect',
      fieldName: 'roleIds',
      label: '角色',
      componentProps: {
        treeData: [],
        fieldNames: { label: 'title', value: 'id', children: 'children' },
        placeholder: '请选择角色',
        multiple: true,
        allowClear: true,
        treeCheckable: true,
        treeDefaultExpandAll: true,
        class: 'w-full',
      },
    },
    {
      component: 'Switch',
      fieldName: 'status',
      label: '状态',
      componentProps: { checkedValue: 1, unCheckedValue: 0 },
      defaultValue: 1,
    },
  ],
});

function mergeCurrentOwnershipOptions(detail?: UsersItem) {
  if (detail?.tenantId && detail.tenantName) {
    const exists = tenantOptions.value.some(
      (item) => item.value === detail.tenantId,
    );
    if (!exists) {
      tenantOptions.value = [
        { label: detail.tenantName, value: detail.tenantId },
        ...tenantOptions.value,
      ];
    }
  }
  if (detail?.merchantId && detail.merchantName) {
    const exists = merchantOptions.value.some(
      (item) => item.value === detail.merchantId,
    );
    if (!exists) {
      merchantOptions.value = [
        { label: detail.merchantName, value: detail.merchantId },
        ...merchantOptions.value,
      ];
    }
  }
}

function updateOwnershipSchemas() {
  formApi.updateSchema([
    {
      fieldName: 'tenantId',
      componentProps: {
        allowClear: true,
        class: 'w-full',
        options: tenantOptions.value,
        placeholder: '请选择所属租户',
      },
    },
    {
      fieldName: 'merchantId',
      componentProps: {
        allowClear: true,
        class: 'w-full',
        options: merchantOptions.value,
        placeholder: '请选择所属商户',
      },
    },
  ]);
}

async function loadOwnershipOptions() {
  try {
    const res = await getTenantList({ pageNum: 1, pageSize: 500 });
    tenantOptions.value = (res?.list ?? []).map((item) => ({
      label: `${item.name}（${item.code}）`,
      value: item.id,
    }));
  } catch {
    tenantOptions.value = [];
  }

  try {
    const res = await getMerchantList({ pageNum: 1, pageSize: 500 });
    merchantOptions.value = (res?.list ?? []).map((item) => ({
      label: item.tenantName
        ? `${item.name}（${item.tenantName}）`
        : `${item.name}（${item.code}）`,
      value: item.id,
    }));
  } catch {
    merchantOptions.value = [];
  }

  updateOwnershipSchemas();
}

/** Modal 配置 */
const [Modal, modalApi] = useVbenModal({
  fullscreenButton: false,
  onCancel() {
    modalApi.close();
  },
  onConfirm: async () => {
    const values = (await formApi.validateAndSubmitForm()) as
      | UsersCreateParams
      | undefined;
    if (!values) return;
    const submitValues: UsersCreateParams = {
      ...values,
      roleIds: normalizeRoleIds(values.roleIds),
    };
    modalApi.lock();
    try {
      if (isEdit.value) {
        await updateUsers({
          id: editId.value,
          ...submitValues,
        } as UsersUpdateParams);
        message.success('更新成功');
      } else {
        await createUsers(submitValues);
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

    await loadOwnershipOptions();
    if (currentOpenToken !== openToken.value) {
      return;
    }

    // 加载部门树
    try {
      const res = await getDeptTree();
      if (currentOpenToken !== openToken.value) {
        return;
      }
      deptTreeData.value = [
        { id: '0', title: '顶级部门', children: res ?? [] } as any,
      ];
      formApi.updateSchema([
        {
          fieldName: 'deptId',
          componentProps: { treeData: deptTreeData.value },
        },
      ]);
    } catch { /* ignore */ }

    // 加载角色树
    try {
      const res = await getRoleTree({ assignableOnly: true });
      if (currentOpenToken !== openToken.value) {
        return;
      }
      formApi.updateSchema([
        {
          fieldName: 'roleIds',
          componentProps: { treeData: res ?? [] },
        },
      ]);
    } catch { /* ignore */ }

    if (currentOpenToken !== openToken.value) {
      return;
    }
    if (data?.id) {
      isEdit.value = true;
      editId.value = data.id;
      modalApi.setState({ title: '编辑用户' });
      try {
        const detail = await getUsersDetail(data.id);
        if (currentOpenToken !== openToken.value) {
          return;
        }
        if (detail) {
          mergeCurrentOwnershipOptions(detail);
          updateOwnershipSchemas();
          formApi.setValues({
            ...detail,
            roleIds: normalizeRoleIds(detail.roleIds),
          });
        }
      } catch {
        if (currentOpenToken === openToken.value) {
          message.error('获取详情失败');
        }
      }
    } else {
      isEdit.value = false;
      editId.value = '';
      modalApi.setState({ title: '新建用户' });
    }
  },
});
</script>

<template>
  <Modal class="w-[600px]">
    <Form />
  </Modal>
</template>
