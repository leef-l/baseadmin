<script setup lang="ts">
import { h, ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
import { message, Tooltip } from 'ant-design-vue';
import { QuestionCircleOutlined } from '@ant-design/icons-vue';
import {
  getUserDetail,
  createUser,
  updateUser,
  getUserTree,
} from '#/api/member/user';
import type {
  UserCreateParams,
  UserUpdateParams,
  UserItem
} from '#/api/member/user/types';
const treeData = ref<UserItem[]>([]);
import { getLevelList } from '#/api/member/level';
import { getTenantList } from '#/api/system/tenant';
import { getMerchantList } from '#/api/system/merchant';
const levelIDOptions = ref<{ label: string; value: string | number }[]>([]);
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
      fieldName: 'parentID',
      label: '上级会员',
      componentProps: {
        treeData: [],
        fieldNames: { label: 'username', value: 'id', children: 'children' },
        placeholder: '请选择上级会员',
        allowClear: true,
        treeDefaultExpandAll: true,
        class: 'w-full',
      },
    },
    {
      component: 'Input',
      fieldName: 'username',
      label: tooltipLabel('用户名', '登录账号'),
      rules: 'required',
      componentProps: { placeholder: '请输入用户名（登录账号）', maxlength: 50 },
    },
    {
      component: 'InputPassword',
      fieldName: 'password',
      label: tooltipLabel('密码', 'bcrypt加密'),
      dependencies: {
        triggerFields: ['password'],
        rules: () => (isEdit.value ? null : 'required'),
        componentProps: () => ({
          placeholder: isEdit.value ? '不填则不修改' : '请输入密码（bcrypt加密）',
        }),
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
      fieldName: 'phone',
      label: '手机号',
      rules: 'phone',
      componentProps: { placeholder: '请输入手机号', maxlength: 20 },
    },
    {
      component: 'ImageUpload',
      fieldName: 'avatar',
      label: '头像',
      componentProps: { maxCount: 1 },
    },
    {
      component: 'Input',
      fieldName: 'realName',
      label: '真实姓名',
      componentProps: { placeholder: '请输入真实姓名', maxlength: 50 },
    },
    {
      component: 'Select',
      fieldName: 'levelID',
      label: '当前等级',
      componentProps: { options: [], placeholder: '请选择当前等级', allowClear: true, class: 'w-full' },
    },
    {
      component: 'DatePicker',
      fieldName: 'levelExpireAt',
      label: '等级到期时间',
      componentProps: { showTime: true, placeholder: '请选择等级到期时间', class: 'w-full', valueFormat: 'YYYY-MM-DD HH:mm:ss' },
    },
    {
      component: 'Input',
      fieldName: 'teamCount',
      label: '团队总人数',
      componentProps: { placeholder: '请输入团队总人数' },
    },
    {
      component: 'Input',
      fieldName: 'directCount',
      label: '直推人数',
      componentProps: { placeholder: '请输入直推人数' },
    },
    {
      component: 'Input',
      fieldName: 'activeCount',
      label: '有效用户数',
      componentProps: { placeholder: '请输入有效用户数' },
    },
    {
      component: 'Input',
      fieldName: 'teamTurnover',
      label: tooltipLabel('团队总营业额', '分'),
      componentProps: { placeholder: '请输入团队总营业额（分）' },
    },
    {
      component: 'Switch',
      fieldName: 'isActive',
      label: '是否激活',
      componentProps: { checkedValue: 1, unCheckedValue: 0 },
      defaultValue: 1,
    },
    {
      component: 'Switch',
      fieldName: 'isQualified',
      label: '仓库资格',
      componentProps: { checkedValue: 1, unCheckedValue: 0 },
      defaultValue: 1,
    },
    {
      component: 'Input',
      fieldName: 'inviteCode',
      label: '邀请码',
      componentProps: { placeholder: '请输入邀请码', maxlength: 32 },
    },
    {
      component: 'Input',
      fieldName: 'registerIP',
      label: '注册IP',
      componentProps: { placeholder: '请输入注册IP', maxlength: 45 },
    },
    {
      component: 'DatePicker',
      fieldName: 'lastLoginAt',
      label: '最后登录时间',
      componentProps: { showTime: true, placeholder: '请选择最后登录时间', class: 'w-full', valueFormat: 'YYYY-MM-DD HH:mm:ss' },
    },
    {
      component: 'Input',
      fieldName: 'remark',
      label: '备注',
      componentProps: { placeholder: '请输入备注', maxlength: 500 },
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
      | UserCreateParams
      | undefined;
    if (!values) return;
    if (isEdit.value && !values.password) {
      delete (values as any).password;
    }
    if (values.teamTurnover != null) {
      (values as any).teamTurnover = Math.round(Number(values.teamTurnover) * 100);
    }
    modalApi.lock();
    try {
      if (isEdit.value) {
        await updateUser({ id: editId.value, ...values } as UserUpdateParams);
        message.success('更新成功');
      } else {
        await createUser(values);
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
    // 加载树形数据
    try {
      const res = await getUserTree();
      if (currentOpenToken !== openToken.value) {
        return;
      }
      treeData.value = [
        { id: '0', username: '顶级节点', children: res ?? [] } as any,
      ];
      formApi.updateSchema([
        {
          fieldName: 'parentID',
          componentProps: { treeData: treeData.value },
        },
      ]);
    } catch {
      // ignore
    }
    // 加载当前等级选项
    try {
      const levelRes = await getLevelList({ pageNum: 1, pageSize: 1000 });
      if (currentOpenToken !== openToken.value) {
        return;
      }
      levelIDOptions.value = (levelRes?.list ?? []).map((item: any) => ({
        label: item.name || item.id,
        value: item.id,
      }));
      formApi.updateSchema([
        {
          fieldName: 'levelID',
          componentProps: { options: levelIDOptions.value },
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
      modalApi.setState({ title: '编辑会员用户' });
      try {
        const detail = await getUserDetail(data.id);
        if (currentOpenToken !== openToken.value) {
          return;
        }
        if (detail) {
          const formData = { ...detail };
          if (formData.teamTurnover != null) {
            formData.teamTurnover = formData.teamTurnover / 100;
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
      modalApi.setState({ title: '新建会员用户' });
    }
  },
});
</script>

<template>
  <Modal class="w-[600px]">
    <Form />
  </Modal>
</template>
