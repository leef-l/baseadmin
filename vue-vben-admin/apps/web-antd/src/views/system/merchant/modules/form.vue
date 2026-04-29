<script setup lang="ts">
import { ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { message } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import {
  createMerchant,
  getMerchantDetail,
  updateMerchant,
} from '#/api/system/merchant';
import type {
  MerchantCreateParams,
  MerchantItem,
  MerchantUpdateParams,
} from '#/api/system/merchant/types';
import { getTenantList } from '#/api/system/tenant';

interface SelectOption {
  label: string;
  value: string;
}

const emit = defineEmits<{ success: [] }>();
const isEdit = ref(false);
const editId = ref('');
const openToken = ref(0);
const tenantOptions = ref<SelectOption[]>([]);

const createAdminDeps = {
  triggerFields: ['_mode'],
  if: () => !isEdit.value,
};

const adminFieldDeps = {
  triggerFields: ['createAdmin', '_mode'],
  if(values: Record<string, any>) {
    return !isEdit.value && values.createAdmin === 1;
  },
};

const [Form, formApi] = useVbenForm({
  showDefaultActions: false,
  schema: [
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        class: 'w-full',
        options: tenantOptions.value,
        placeholder: '请选择租户',
      },
      fieldName: 'tenantId',
      label: '所属租户',
    },
    {
      component: 'Input',
      componentProps: { maxlength: 80, placeholder: '请输入商户名称' },
      fieldName: 'name',
      label: '商户名称',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: { maxlength: 50, placeholder: '请输入商户编码' },
      fieldName: 'code',
      label: '商户编码',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: { maxlength: 50, placeholder: '请输入联系人' },
      fieldName: 'contactName',
      label: '联系人',
    },
    {
      component: 'Input',
      componentProps: { maxlength: 30, placeholder: '请输入联系电话' },
      fieldName: 'contactPhone',
      label: '联系电话',
    },
    {
      component: 'Input',
      componentProps: { maxlength: 255, placeholder: '请输入商户地址' },
      fieldName: 'address',
      label: '商户地址',
    },
    {
      component: 'Switch',
      componentProps: { checkedValue: 1, unCheckedValue: 0 },
      defaultValue: 1,
      fieldName: 'status',
      label: '状态',
    },
    {
      component: 'Switch',
      componentProps: { checkedValue: 1, unCheckedValue: 0 },
      defaultValue: 1,
      dependencies: createAdminDeps,
      fieldName: 'createAdmin',
      label: '同步创建管理员',
    },
    {
      component: 'Input',
      componentProps: {
        maxlength: 50,
        placeholder: '留空默认商户编码_admin',
      },
      dependencies: adminFieldDeps,
      fieldName: 'adminUsername',
      label: '管理员用户名',
    },
    {
      component: 'InputPassword',
      componentProps: { maxlength: 64, placeholder: '请输入管理员密码' },
      dependencies: adminFieldDeps,
      fieldName: 'adminPassword',
      label: '管理员密码',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: { maxlength: 50, placeholder: '请输入管理员昵称' },
      dependencies: adminFieldDeps,
      fieldName: 'adminNickname',
      label: '管理员昵称',
    },
    {
      component: 'Input',
      componentProps: { maxlength: 100, placeholder: '请输入管理员邮箱' },
      dependencies: adminFieldDeps,
      fieldName: 'adminEmail',
      label: '管理员邮箱',
    },
    {
      component: 'Textarea',
      componentProps: { maxlength: 255, placeholder: '请输入备注', rows: 3 },
      fieldName: 'remark',
      label: '备注',
    },
  ],
});

function mergeCurrentTenantOption(detail?: MerchantItem) {
  if (!detail?.tenantId || !detail.tenantName) {
    return;
  }
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

async function loadTenantOptions() {
  try {
    const res = await getTenantList({ pageNum: 1, pageSize: 500 });
    tenantOptions.value = (res?.list ?? []).map((item) => ({
      label: `${item.name}（${item.code}）`,
      value: item.id,
    }));
  } catch {
    tenantOptions.value = [];
  }
  formApi.updateSchema([
    {
      componentProps: {
        allowClear: true,
        class: 'w-full',
        options: tenantOptions.value,
        placeholder: '请选择租户',
      },
      fieldName: 'tenantId',
    },
  ]);
}

const [Modal, modalApi] = useVbenModal({
  fullscreenButton: false,
  onCancel() {
    modalApi.close();
  },
  onConfirm: async () => {
    const values = (await formApi.validateAndSubmitForm()) as
      | MerchantCreateParams
      | undefined;
    if (!values) return;
    modalApi.lock();
    try {
      if (isEdit.value) {
        await updateMerchant({
          id: editId.value,
          ...values,
        } as MerchantUpdateParams);
        message.success('更新成功');
      } else {
        await createMerchant(values);
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
    await loadTenantOptions();
    if (currentOpenToken !== openToken.value) {
      return;
    }
    if (data?.id) {
      isEdit.value = true;
      editId.value = data.id;
      formApi.setValues({ _mode: 'edit', createAdmin: 0 });
      modalApi.setState({ title: '编辑商户' });
      try {
        const detail = await getMerchantDetail(data.id);
        if (currentOpenToken !== openToken.value) {
          return;
        }
        if (detail) {
          mergeCurrentTenantOption(detail);
          formApi.updateSchema([
            {
              componentProps: {
                allowClear: true,
                class: 'w-full',
                options: tenantOptions.value,
                placeholder: '请选择租户',
              },
              fieldName: 'tenantId',
            },
          ]);
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
      formApi.setValues({ _mode: 'create', createAdmin: 1 });
      modalApi.setState({ title: '新建商户' });
    }
  },
});
</script>

<template>
  <Modal class="w-[760px]">
    <Form />
  </Modal>
</template>
