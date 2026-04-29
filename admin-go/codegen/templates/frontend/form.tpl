<script setup lang="ts">
{{- if .HasTooltip}}
import { h, ref } from 'vue';
{{- else}}
import { ref } from 'vue';
{{- end}}
import { useVbenModal } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
{{- if .HasTenantScope}}
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
{{- end}}
{{- if .HasTooltip}}
import { message, Tooltip } from 'ant-design-vue';
import { QuestionCircleOutlined } from '@ant-design/icons-vue';
{{- else}}
import { message } from 'ant-design-vue';
{{- end}}
import {
  get{{.ModelName}}Detail,
  create{{.ModelName}},
  update{{.ModelName}},{{if .HasParentID}}
  get{{.ModelName}}Tree,{{end}}
} from '#/api/{{.AppName}}/{{.ModuleName}}';
{{- if .HasDict}}
{{- if .AllowMissingDictModule}}
async function getDictByType(_dictType: string): Promise<Array<{ label: string; value: string | number }>> {
  return [];
}
{{- else}}
import { getDictByType } from '#/api/system/dict';
{{- end}}
{{- end}}
import type {
  {{.ModelName}}CreateParams,
  {{.ModelName}}UpdateParams{{if or .HasParentID .HasTreeSelect}},
  {{.ModelName}}Item{{end}}
} from '#/api/{{.AppName}}/{{.ModuleName}}/types';
{{- if or .HasParentID .HasTreeSelect}}
const treeData = ref<{{.ModelName}}Item[]>([]);
{{- end}}
{{- range .Fields}}
{{- if and .IsForeignKey (not .IsHidden) .RefTable}}
{{- if .RefIsTree}}
import { get{{.RefTableCamel}}Tree } from '#/api/{{.RefTableApp}}/{{.RefTable}}';
{{- if or (ne .RefTableApp $.AppName) (ne .RefTable $.ModuleName)}}
import type { {{.RefTableCamel}}Item } from '#/api/{{.RefTableApp}}/{{.RefTable}}/types';
{{- end}}
{{- else}}
import { get{{.RefTableCamel}}List } from '#/api/{{.RefTableApp}}/{{.RefTable}}';
{{- end}}
{{- end}}
{{- end}}
{{- range .Fields}}
{{- if and (not .IsHidden) (not .IsID) (.IsEnum) (ne .Component "Switch")}}

/** {{.Label}}选项 */
const {{.NameLower}}Options = [
{{- range .EnumValues}}
  { label: '{{.Label}}', value: {{if IsNumeric .Value}}{{.Value}}{{else}}'{{.Value}}'{{end}} },
{{- end}}
];
{{- end}}
{{- end}}
{{- range .Fields}}
{{- if and .IsForeignKey (not .IsHidden) .RefTable}}
const {{.NameLower}}Options = ref<{{if .RefIsTree}}{{if and (eq .RefTableApp $.AppName) (eq .RefTable $.ModuleName)}}{{$.ModelName}}Item{{else}}{{.RefTableCamel}}Item{{end}}[]{{else}}{ label: string; value: string | number }[]{{end}}>([]);
{{- end}}
{{- end}}
{{- range .Fields}}
{{- if and (not .IsHidden) .DictType}}
const {{.NameLower}}DictOptions = ref<{ label: string; value: string | number }[]>([]);
{{- end}}
{{- end}}
{{- if .HasTooltip}}
/** 渲染带 Tooltip 的表单 label */
function tooltipLabel(label: string, tip: string) {
  return () => h('span', {}, [
    label + ' ',
    h(Tooltip, { title: tip }, {
      default: () => h(QuestionCircleOutlined, { style: { color: '#999', marginLeft: '4px' } }),
    }),
  ]);
}
{{- end}}

const emit = defineEmits<{ success: [] }>();
{{- if .HasTenantScope}}
const isPlatformSuperAdmin = usePlatformSuperAdmin();
{{- end}}
const isEdit = ref(false);
const editId = ref('');
const openToken = ref(0);

/** 表单配置 */
const [Form, formApi] = useVbenForm({
  showDefaultActions: false,
  schema: [
{{- range .Fields}}
{{- if and (not .IsHidden) (not .IsID)}}
{{- $isScopeField := or (eq .Name "tenant_id") (eq .Name "merchant_id")}}
{{- if eq .Component "Password"}}
    {
      component: 'InputPassword',
      fieldName: '{{.NameLower}}',
      label: {{if .TooltipText}}tooltipLabel('{{.ShortLabel}}', '{{.TooltipText}}'){{else}}'{{.Label}}'{{end}},
      dependencies: {
        triggerFields: ['{{.NameLower}}'],
        rules: () => (isEdit.value ? null : 'required'),
        componentProps: () => ({
          placeholder: isEdit.value ? '不填则不修改' : '请输入{{.Label}}',
        }),
      },
    },
{{- else if eq .Component "InputNumber"}}
    {
      component: 'InputNumber',
      fieldName: '{{.NameLower}}',
      label: {{if .TooltipText}}tooltipLabel('{{.ShortLabel}}', '{{.TooltipText}}'){{else}}'{{.Label}}'{{end}},
{{- if .IsRequired}}
      rules: 'required',
{{- end}}
      componentProps: { placeholder: '请输入{{.Label}}', class: 'w-full' },
    },
{{- else if eq .Component "Textarea"}}
    {
      component: 'Textarea',
      fieldName: '{{.NameLower}}',
      label: {{if .TooltipText}}tooltipLabel('{{.ShortLabel}}', '{{.TooltipText}}'){{else}}'{{.Label}}'{{end}},
{{- if .IsRequired}}
      rules: 'required',
{{- end}}
      componentProps: { placeholder: '请输入{{.Label}}', rows: 4{{if gt .MaxLength 0}}, maxlength: {{.MaxLength}}{{end}} },
    },
{{- else if eq .Component "Switch"}}
    {
      component: 'Switch',
      fieldName: '{{.NameLower}}',
      label: {{if .TooltipText}}tooltipLabel('{{.ShortLabel}}', '{{.TooltipText}}'){{else}}'{{.Label}}'{{end}},
      componentProps: { checkedValue: 1, unCheckedValue: 0 },
      defaultValue: {{if .DefaultValue}}{{.DefaultValue}}{{else}}0{{end}},
    },
{{- else if eq .Component "Radio"}}
    {
      component: 'RadioGroup',
      fieldName: '{{.NameLower}}',
      label: {{if .TooltipText}}tooltipLabel('{{.ShortLabel}}', '{{.TooltipText}}'){{else}}'{{.Label}}'{{end}},
{{- if .IsRequired}}
      rules: 'selectRequired',
{{- end}}
      componentProps: { options: {{.NameLower}}Options },
    },
{{- else if eq .Component "Select"}}
{{- if .IsEnum}}
    {
      component: 'Select',
      fieldName: '{{.NameLower}}',
      label: {{if .TooltipText}}tooltipLabel('{{.ShortLabel}}', '{{.TooltipText}}'){{else}}'{{.Label}}'{{end}},
{{- if .IsRequired}}
      rules: 'selectRequired',
{{- end}}
      componentProps: { options: {{.NameLower}}Options, placeholder: '请选择{{.Label}}', allowClear: true, class: 'w-full' },
    },
{{- else if .DictType}}
    {
      component: 'Select',
      fieldName: '{{.NameLower}}',
      label: {{if .TooltipText}}tooltipLabel('{{.ShortLabel}}', '{{.TooltipText}}'){{else}}'{{.Label}}'{{end}},
{{- if .IsRequired}}
      rules: 'selectRequired',
{{- end}}
      componentProps: { options: [], placeholder: '请选择{{.Label}}', allowClear: true, class: 'w-full' },
    },
{{- else if .IsForeignKey}}
{{- if .RefIsTree}}
    {
      component: 'TreeSelect',
      fieldName: '{{.NameLower}}',
      label: {{if .TooltipText}}tooltipLabel('{{.ShortLabel}}', '{{.TooltipText}}'){{else}}'{{.Label}}'{{end}},
{{- if and $.HasTenantScope $isScopeField}}
      ifShow: () => isPlatformSuperAdmin.value,
{{- end}}
{{- if .IsRequired}}
      rules: 'selectRequired',
{{- end}}
      componentProps: {
        treeData: [],
        fieldNames: { label: '{{if .RefDisplayField}}{{.RefDisplayLower}}{{else}}title{{end}}', value: 'id', children: 'children' },
        placeholder: '请选择{{.Label}}',
        allowClear: true,
        treeDefaultExpandAll: true,
        class: 'w-full',
      },
    },
{{- else}}
    {
      component: 'Select',
      fieldName: '{{.NameLower}}',
      label: {{if .TooltipText}}tooltipLabel('{{.ShortLabel}}', '{{.TooltipText}}'){{else}}'{{.Label}}'{{end}},
{{- if and $.HasTenantScope $isScopeField}}
      ifShow: () => isPlatformSuperAdmin.value,
{{- end}}
{{- if .IsRequired}}
      rules: 'selectRequired',
{{- end}}
      componentProps: { options: [], placeholder: '请选择{{.Label}}', allowClear: true, class: 'w-full' },
    },
{{- end}}
{{- else}}
    {
      component: 'Input',
      fieldName: '{{.NameLower}}',
      label: {{if .TooltipText}}tooltipLabel('{{.ShortLabel}}', '{{.TooltipText}}'){{else}}'{{.Label}}'{{end}},
{{- if .IsRequired}}
      rules: 'required',
{{- end}}
      componentProps: { placeholder: '请输入{{.Label}}' },
    },
{{- end}}
{{- else if eq .Component "SelectMulti"}}
{{- if .IsEnum}}
    {
      component: 'Select',
      fieldName: '{{.NameLower}}',
      label: {{if .TooltipText}}tooltipLabel('{{.ShortLabel}}', '{{.TooltipText}}'){{else}}'{{.Label}}'{{end}},
{{- if .IsRequired}}
      rules: 'selectRequired',
{{- end}}
      componentProps: { options: {{.NameLower}}Options, placeholder: '请选择{{.Label}}', mode: 'multiple', allowClear: true, class: 'w-full' },
    },
{{- else}}
    {
      component: 'Select',
      fieldName: '{{.NameLower}}',
      label: {{if .TooltipText}}tooltipLabel('{{.ShortLabel}}', '{{.TooltipText}}'){{else}}'{{.Label}}'{{end}},
{{- if .IsRequired}}
      rules: 'selectRequired',
{{- end}}
      componentProps: { placeholder: '请输入{{.Label}}', mode: 'tags', allowClear: true, class: 'w-full' },
    },
{{- end}}
{{- else if eq .Component "TreeSelectSingle"}}
    {
      component: 'TreeSelect',
      fieldName: '{{.NameLower}}',
      label: {{if .TooltipText}}tooltipLabel('{{.ShortLabel}}', '{{.TooltipText}}'){{else}}'{{.Label}}'{{end}},
{{- if .IsRequired}}
      rules: 'selectRequired',
{{- end}}
      componentProps: {
        treeData: [],
        fieldNames: { label: '{{if .RefDisplayLower}}{{.RefDisplayLower}}{{else}}title{{end}}', value: 'id', children: 'children' },
        placeholder: '请选择{{.Label}}',
        allowClear: true,
        treeDefaultExpandAll: true,
        class: 'w-full',
      },
    },
{{- else if eq .Component "TreeSelectMulti"}}
    {
      component: 'TreeSelect',
      fieldName: '{{.NameLower}}',
      label: {{if .TooltipText}}tooltipLabel('{{.ShortLabel}}', '{{.TooltipText}}'){{else}}'{{.Label}}'{{end}},
{{- if .IsRequired}}
      rules: 'selectRequired',
{{- end}}
      componentProps: {
        treeData: [],
        fieldNames: { label: '{{if .RefDisplayLower}}{{.RefDisplayLower}}{{else}}title{{end}}', value: 'id', children: 'children' },
        placeholder: '请选择{{.Label}}',
        allowClear: true,
        treeCheckable: true,
        treeDefaultExpandAll: true,
        class: 'w-full',
      },
    },
{{- else if eq .Component "DateTimePicker"}}
    {
      component: 'DatePicker',
      fieldName: '{{.NameLower}}',
      label: {{if .TooltipText}}tooltipLabel('{{.ShortLabel}}', '{{.TooltipText}}'){{else}}'{{.Label}}'{{end}},
{{- if .IsRequired}}
      rules: 'required',
{{- end}}
      componentProps: { showTime: true, placeholder: '请选择{{.Label}}', class: 'w-full', valueFormat: 'YYYY-MM-DD HH:mm:ss' },
    },
{{- else if eq .Component "ImageUpload"}}
    {
      component: 'ImageUpload',
      fieldName: '{{.NameLower}}',
      label: '{{.Label}}',
{{- if .IsRequired}}
      rules: 'required',
{{- end}}
      componentProps: { maxCount: 1 },
    },
{{- else if eq .Component "FileUpload"}}
    {
      component: 'FileUpload',
      fieldName: '{{.NameLower}}',
      label: '{{.Label}}',
{{- if .IsRequired}}
      rules: 'required',
{{- end}}
      componentProps: { maxCount: 1 },
    },
{{- else if eq .Component "RichText"}}
    {
      component: 'RichText',
      fieldName: '{{.NameLower}}',
      label: '{{.Label}}',
{{- if .IsRequired}}
      rules: 'required',
{{- end}}
      formItemClass: 'col-span-full',
    },
{{- else if eq .Component "JsonEditor"}}
    {
      component: 'JsonEditor',
      fieldName: '{{.NameLower}}',
      label: '{{.Label}}',
{{- if .IsRequired}}
      rules: 'required',
{{- end}}
      formItemClass: 'col-span-full',
    },
{{- else if eq .Component "IconPicker"}}
    {
      component: 'IconPicker',
      fieldName: '{{.NameLower}}',
      label: {{if .TooltipText}}tooltipLabel('{{.ShortLabel}}', '{{.TooltipText}}'){{else}}'{{.Label}}'{{end}},
{{- if .IsRequired}}
      rules: 'required',
{{- end}}
      componentProps: { placeholder: '请选择图标' },
    },
{{- else if eq .Component "InputUrl"}}
    {
      component: 'Input',
      fieldName: '{{.NameLower}}',
      label: {{if .TooltipText}}tooltipLabel('{{.ShortLabel}}', '{{.TooltipText}}'){{else}}'{{.Label}}'{{end}},
{{- if .IsRequired}}
      rules: 'required',
{{- end}}
      componentProps: { placeholder: '请输入完整URL地址（含 http:// 或 https://）'{{if gt .MaxLength 0}}, maxlength: {{.MaxLength}}{{end}} },
    },
{{- else}}
    {
      component: 'Input',
      fieldName: '{{.NameLower}}',
      label: {{if .TooltipText}}tooltipLabel('{{.ShortLabel}}', '{{.TooltipText}}'){{else}}'{{.Label}}'{{end}},
{{- if .FrontendRules}}
      rules: '{{.FrontendRules}}',
{{- else if .IsRequired}}
      rules: 'required',
{{- end}}
      componentProps: { placeholder: '请输入{{.Label}}'{{if gt .MaxLength 0}}, maxlength: {{.MaxLength}}{{end}} },
    },
{{- end}}
{{- end}}
{{- end}}
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
      | {{.ModelName}}CreateParams
      | undefined;
    if (!values) return;
{{- range .Fields}}
{{- if .IsPassword}}
    if (isEdit.value && !values.{{.NameLower}}) {
      delete (values as any).{{.NameLower}};
    }
{{- end}}
{{- end}}
{{- if .HasMoney}}
{{- range .Fields}}
{{- if and .IsMoney (not .IsHidden)}}
    if (values.{{.NameLower}} != null) {
      (values as any).{{.NameLower}} = Math.round(Number(values.{{.NameLower}}) * 100);
    }
{{- end}}
{{- end}}
{{- end}}
    modalApi.lock();
    try {
      if (isEdit.value) {
        await update{{.ModelName}}({ id: editId.value, ...values } as {{.ModelName}}UpdateParams);
        message.success('更新成功');
      } else {
        await create{{.ModelName}}(values);
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
{{- if .HasParentID}}
    // 加载树形数据
    try {
      const res = await get{{.ModelName}}Tree();
      if (currentOpenToken !== openToken.value) {
        return;
      }
      treeData.value = [
        { id: '0', {{if .ParentDisplayField}}{{.ParentDisplayField}}{{else}}title{{end}}: '顶级节点', children: res ?? [] } as any,
      ];
      formApi.updateSchema([
        {
          fieldName: '{{range .Fields}}{{if .IsParentID}}{{.NameLower}}{{end}}{{end}}',
          componentProps: { treeData: treeData.value },
        },
      ]);
    } catch {
      // ignore
    }
{{- end}}
{{- range .Fields}}
{{- if and .IsForeignKey (not .IsHidden) .RefTable}}
{{- $isScopeField := and $.HasTenantScope (or (eq .Name "tenant_id") (eq .Name "merchant_id"))}}
{{- if .RefIsTree}}
{{- if $isScopeField}}
    if (isPlatformSuperAdmin.value) {
{{- end}}
    // 加载{{.Label}}树形数据
    try {
      const {{.RefTableLower}}Res = await get{{.RefTableCamel}}Tree();
      if (currentOpenToken !== openToken.value) {
        return;
      }
      {{.NameLower}}Options.value = {{.RefTableLower}}Res ?? [];
      formApi.updateSchema([
        {
          fieldName: '{{.NameLower}}',
          componentProps: { treeData: {{.NameLower}}Options.value },
        },
      ]);
    } catch {
      // ignore
    }
{{- if $isScopeField}}
    }
{{- end}}
{{- else}}
{{- if $isScopeField}}
    if (isPlatformSuperAdmin.value) {
{{- end}}
    // 加载{{.Label}}选项
    try {
      const {{.RefTableLower}}Res = await get{{.RefTableCamel}}List({ pageNum: 1, pageSize: 1000 });
      if (currentOpenToken !== openToken.value) {
        return;
      }
      {{.NameLower}}Options.value = ({{.RefTableLower}}Res?.list ?? []).map((item: any) => ({
        label: item.{{.RefDisplayLower}} || item.id,
        value: item.id,
      }));
      formApi.updateSchema([
        {
          fieldName: '{{.NameLower}}',
          componentProps: { options: {{.NameLower}}Options.value },
        },
      ]);
    } catch {
      // ignore
    }
{{- if $isScopeField}}
    }
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- range .Fields}}
{{- if and (not .IsHidden) .DictType}}
    // 加载{{.Label}}字典
    try {
      const dictRes = await getDictByType('{{.DictType}}');
      if (currentOpenToken !== openToken.value) {
        return;
      }
      {{.NameLower}}DictOptions.value = (dictRes ?? []).map((item: any) => ({
        label: item.label,
        value: item.value,
      }));
      formApi.updateSchema([
        {
          fieldName: '{{.NameLower}}',
          componentProps: { options: {{.NameLower}}DictOptions.value },
        },
      ]);
    } catch {
      // ignore
    }
{{- end}}
{{- end}}
    if (currentOpenToken !== openToken.value) {
      return;
    }
    if (data?.id) {
      isEdit.value = true;
      editId.value = data.id;
      modalApi.setState({ title: '编辑{{.Comment}}' });
      try {
        const detail = await get{{.ModelName}}Detail(data.id);
        if (currentOpenToken !== openToken.value) {
          return;
        }
        if (detail) {
{{- if .HasMoney}}
          const formData = { ...detail };
{{- range .Fields}}
{{- if and .IsMoney (not .IsHidden)}}
          if (formData.{{.NameLower}} != null) {
            formData.{{.NameLower}} = formData.{{.NameLower}} / 100;
          }
{{- end}}
{{- end}}
          formApi.setValues(formData);
{{- else}}
          formApi.setValues(detail);
{{- end}}
        }
      } catch {
        if (currentOpenToken === openToken.value) {
          message.error('获取详情失败');
        }
      }
    } else {
      isEdit.value = false;
      editId.value = '';
      modalApi.setState({ title: '新建{{.Comment}}' });
    }
  },
});
</script>

<template>
  <Modal class="w-[{{if .HasRichText}}800px{{else}}600px{{end}}]">
    <Form />
  </Modal>
</template>
