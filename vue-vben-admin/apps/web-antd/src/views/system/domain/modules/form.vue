<script setup lang="ts">
import { ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { message } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import {
  createDomain,
  getDomainDetail,
  updateDomain,
} from '#/api/system/domain';
import type {
  DomainCreateParams,
  DomainUpdateParams,
} from '#/api/system/domain/types';
import { getMerchantList } from '#/api/system/merchant';
import { getTenantList } from '#/api/system/tenant';

interface SelectOption {
  label: string;
  tenantId?: string;
  value: string;
}

const emit = defineEmits<{ success: [] }>();
const isEdit = ref(false);
const editId = ref('');
const openToken = ref(0);
const tenantOptions = ref<SelectOption[]>([]);
const allMerchantOptions = ref<SelectOption[]>([]);

const ownerOptions = [
  { label: '租户', value: 1 },
  { label: '商户', value: 2 },
];

const appOptions = [{ label: '后台', value: 'admin' }];

function getMerchantOptions(tenantId?: string) {
  if (!tenantId) {
    return allMerchantOptions.value;
  }
  return allMerchantOptions.value.filter((item) => item.tenantId === tenantId);
}

const merchantDeps = {
  triggerFields: ['ownerType', 'tenantId'],
  if(values: Record<string, any>) {
    return values.ownerType === 2;
  },
  componentProps(values: Record<string, any>) {
    return {
      allowClear: true,
      class: 'w-full',
      options: getMerchantOptions(values.tenantId),
      placeholder: '请选择商户；商户账号可留空',
    };
  },
};

const [Form, formApi] = useVbenForm({
  showDefaultActions: false,
  schema: [
    {
      component: 'Input',
      componentProps: {
        maxlength: 255,
        placeholder: '请输入域名，例如 admin.example.com',
      },
      fieldName: 'domain',
      label: '域名',
      rules: 'required',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: false,
        class: 'w-full',
        options: ownerOptions,
      },
      defaultValue: 1,
      fieldName: 'ownerType',
      label: '主体类型',
      rules: 'selectRequired',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        class: 'w-full',
        options: tenantOptions.value,
        placeholder: '请选择租户；租户账号可留空',
      },
      fieldName: 'tenantId',
      label: '所属租户',
    },
    {
      component: 'Select',
      dependencies: merchantDeps,
      fieldName: 'merchantId',
      label: '所属商户',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: false,
        class: 'w-full',
        options: appOptions,
      },
      defaultValue: 'admin',
      fieldName: 'appCode',
      label: '应用',
      rules: 'selectRequired',
    },
    {
      component: 'Switch',
      componentProps: { checkedValue: 1, unCheckedValue: 0 },
      defaultValue: 1,
      fieldName: 'verifyStatus',
      label: '域名已校验',
    },
    {
      component: 'Switch',
      componentProps: { checkedValue: 1, unCheckedValue: 0 },
      defaultValue: 1,
      fieldName: 'status',
      label: '状态',
    },
    {
      component: 'Textarea',
      componentProps: { maxlength: 500, placeholder: '请输入备注', rows: 3 },
      fieldName: 'remark',
      label: '备注',
    },
  ],
});

async function loadOptions() {
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
    allMerchantOptions.value = (res?.list ?? []).map((item) => ({
      label: `${item.name}（${item.code}）`,
      tenantId: item.tenantId,
      value: item.id,
    }));
  } catch {
    allMerchantOptions.value = [];
  }
  formApi.updateSchema([
    {
      componentProps: {
        allowClear: true,
        class: 'w-full',
        options: tenantOptions.value,
        placeholder: '请选择租户；租户账号可留空',
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
      | DomainCreateParams
      | undefined;
    if (!values) return;
    modalApi.lock();
    try {
      if (isEdit.value) {
        await updateDomain({ id: editId.value, ...values } as DomainUpdateParams);
        message.success('更新成功');
      } else {
        await createDomain(values);
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
    await loadOptions();
    if (currentOpenToken !== openToken.value) {
      return;
    }

    const data = modalApi.getData<{ id?: string }>();
    if (data?.id) {
      isEdit.value = true;
      editId.value = data.id;
      modalApi.setState({ title: '编辑域名' });
      try {
        const detail = await getDomainDetail(data.id);
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
      formApi.setValues({
        appCode: 'admin',
        ownerType: 1,
        status: 1,
        verifyStatus: 1,
      });
      modalApi.setState({ title: '新建域名' });
    }
  },
});
</script>

<template>
  <Modal class="w-[760px]">
    <Form />
  </Modal>
</template>
