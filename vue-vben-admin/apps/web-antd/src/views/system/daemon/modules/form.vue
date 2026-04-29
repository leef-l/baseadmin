<script setup lang="ts">
import { computed, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { message } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import {
  createDaemon,
  getDaemonDetail,
  updateDaemon,
} from '#/api/system/daemon';
import type {
  DaemonCreateParams,
  DaemonUpdateParams,
} from '#/api/system/daemon/types';

const emit = defineEmits<{ success: [] }>();
const isEdit = ref(false);
const editId = ref('');
const openToken = ref(0);

const programDisabled = computed(() => isEdit.value);

const [Form, formApi] = useVbenForm({
  showDefaultActions: false,
  schema: [
    {
      component: 'Input',
      componentProps: { maxlength: 80, placeholder: '请输入显示名称' },
      fieldName: 'name',
      label: '显示名称',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: () => ({
        disabled: programDisabled.value,
        maxlength: 80,
        placeholder: '仅支持字母、数字、中横线、下划线',
      }),
      fieldName: 'program',
      label: '进程名',
      rules: 'required',
    },
    {
      component: 'Textarea',
      componentProps: { maxlength: 1000, placeholder: '请输入启动命令', rows: 3 },
      fieldName: 'command',
      label: '启动命令',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: { maxlength: 500, placeholder: '请输入运行目录' },
      fieldName: 'directory',
      label: '运行目录',
      rules: 'required',
    },
    {
      component: 'Input',
      componentProps: { maxlength: 80, placeholder: '默认 root' },
      defaultValue: 'root',
      fieldName: 'runUser',
      label: '运行用户',
    },
    {
      component: 'InputNumber',
      componentProps: { class: 'w-full', max: 64, min: 1 },
      defaultValue: 1,
      fieldName: 'numprocs',
      label: '进程数量',
    },
    {
      component: 'InputNumber',
      componentProps: { class: 'w-full', max: 9999, min: 1 },
      defaultValue: 999,
      fieldName: 'priority',
      label: '启动优先级',
    },
    {
      component: 'InputNumber',
      componentProps: { class: 'w-full', max: 3600, min: 0 },
      defaultValue: 3,
      fieldName: 'startsecs',
      label: '稳定秒数',
    },
    {
      component: 'InputNumber',
      componentProps: { class: 'w-full', max: 100, min: 0 },
      defaultValue: 3,
      fieldName: 'startretries',
      label: '重试次数',
    },
    {
      component: 'Select',
      componentProps: {
        class: 'w-full',
        options: ['TERM', 'HUP', 'INT', 'QUIT', 'KILL', 'USR1', 'USR2'].map(
          (item) => ({ label: item, value: item }),
        ),
      },
      defaultValue: 'QUIT',
      fieldName: 'stopSignal',
      label: '停止信号',
    },
    {
      component: 'Switch',
      componentProps: { checkedValue: 1, unCheckedValue: 0 },
      defaultValue: 1,
      fieldName: 'autostart',
      label: '随服务启动',
    },
    {
      component: 'Switch',
      componentProps: { checkedValue: 1, unCheckedValue: 0 },
      defaultValue: 1,
      fieldName: 'autorestart',
      label: '自动重启',
    },
    {
      component: 'Textarea',
      componentProps: {
        maxlength: 1000,
        placeholder: '例如 GOGC="100",GOMEMLIMIT="512MiB"',
        rows: 2,
      },
      fieldName: 'environment',
      label: '环境变量',
    },
    {
      component: 'Textarea',
      componentProps: { maxlength: 500, placeholder: '请输入备注', rows: 3 },
      fieldName: 'remark',
      label: '备注',
    },
  ],
});

const [Modal, modalApi] = useVbenModal({
  fullscreenButton: false,
  onCancel() {
    modalApi.close();
  },
  onConfirm: async () => {
    const values = (await formApi.validateAndSubmitForm()) as
      | DaemonCreateParams
      | undefined;
    if (!values) return;
    modalApi.lock();
    try {
      if (isEdit.value) {
        await updateDaemon({ id: editId.value, ...values } as DaemonUpdateParams);
        message.success('更新成功');
      } else {
        await createDaemon(values);
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
    if (data?.id) {
      isEdit.value = true;
      editId.value = data.id;
      modalApi.setState({ title: '编辑守护进程' });
      try {
        const detail = await getDaemonDetail(data.id);
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
        autostart: 1,
        autorestart: 1,
        numprocs: 1,
        priority: 999,
        runUser: 'root',
        startretries: 3,
        startsecs: 3,
        stopSignal: 'QUIT',
      });
      modalApi.setState({ title: '新建守护进程' });
    }
  },
});
</script>

<template>
  <Modal class="w-[820px]">
    <Form />
  </Modal>
</template>
