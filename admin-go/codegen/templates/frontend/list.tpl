<script setup lang="ts">
{{- if .HasTooltip}}
import { h, onMounted{{if .HasImport}}, ref{{end}} } from 'vue';
{{- end}}
{{- if not .HasTooltip}}
import { onMounted{{if .HasImport}}, ref{{end}} } from 'vue';
{{- end}}
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { useAccess } from '@vben/access';
import { Page, useVbenModal } from '@vben/common-ui';
import { downloadFileFromBlob } from '@vben/utils';
{{- if .HasTooltip}}
import { Button, message, Modal{{if .HasEnum}}, Tag{{end}}, Tooltip } from 'ant-design-vue';
import { QuestionCircleOutlined } from '@ant-design/icons-vue';
{{- else}}
import { Button, message, Modal{{if .HasEnum}}, Tag{{end}} } from 'ant-design-vue';
{{- end}}

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getGridSelectedIds } from '#/utils/grid-selection';
{{- if .HasTenantScope}}
import { usePlatformSuperAdmin } from '#/utils/auth-scope';
{{- end}}
{{- if .HasParentID}}
import { get{{.ModelName}}Tree, delete{{.ModelName}}, batchDelete{{.ModelName}}, export{{.ModelName}}{{if .HasImport}}, import{{.ModelName}}, downloadImportTemplate{{.ModelName}}{{end}}{{if .HasBatchEdit}}, batchUpdate{{.ModelName}}{{end}} } from '#/api/{{.AppName}}/{{.ModuleName}}';
{{- else}}
import { get{{.ModelName}}List, delete{{.ModelName}}, batchDelete{{.ModelName}}, export{{.ModelName}}{{if .HasImport}}, import{{.ModelName}}, downloadImportTemplate{{.ModelName}}{{end}}{{if .HasBatchEdit}}, batchUpdate{{.ModelName}}{{end}} } from '#/api/{{.AppName}}/{{.ModuleName}}';
{{- end}}
{{- if .HasDict}}
{{- if .AllowMissingDictModule}}
async function getDictByType(_dictType: string): Promise<Array<{ label: string; value: string | number }>> {
  return [];
}
{{- else}}
import { getDictByType } from '#/api/system/dict';
{{- end}}
{{- end}}
{{- range .SearchFields}}
{{- if and .IsForeignKey .RefTable}}
{{- if .RefIsTree}}
import { get{{.RefTableCamel}}Tree } from '#/api/{{.RefTableApp}}/{{.RefTable}}';
{{- else}}
import { get{{.RefTableCamel}}List } from '#/api/{{.RefTableApp}}/{{.RefTable}}';
{{- end}}
{{- end}}
{{- end}}
import type { {{.ModelName}}Item } from '#/api/{{.AppName}}/{{.ModuleName}}/types';
import FormModal from './modules/form.vue';
import DetailDrawer from './modules/detail-drawer.vue';
{{if .HasEnum}}
/** 标签颜色池 */
const TAG_COLORS = ['green', 'red', 'blue', 'orange', 'cyan', 'purple', 'geekblue', 'magenta'];

type EnumValue = number | string;

function getEnumLabel(map: Record<EnumValue, string>, value: EnumValue | null | undefined) {
  if (value === null || value === undefined || value === '') {
    return '-';
  }
  return map[value] ?? String(value);
}
{{- end}}

const sortableFieldMap: Record<string, string> = {
  createdAt: 'created_at',
{{- if .HasSort}}
  sort: 'sort',
{{- end}}
{{- if .HasStatus}}
  status: 'status',
{{- end}}
{{- range .Fields}}
{{- if and (not .IsHidden) (not .IsID) (not .RefFieldJSON) (or .IsMoney .IsSearchable)}}
  {{.NameLower}}: '{{.Name}}',
{{- end}}
{{- end}}
};

function resolveSortField(field?: string) {
  if (!field) {
    return '';
  }
  return sortableFieldMap[field] ?? '';
}
{{- range .Fields}}
{{- if and (not .IsHidden) (.IsEnum)}}

/** {{.Label}}选项 */
const {{.NameLower}}Options = [
{{- range .EnumValues}}
  { label: '{{.Label}}', value: {{if IsNumeric .Value}}{{.Value}}{{else}}'{{.Value}}'{{end}} },
{{- end}}
];

/** {{.Label}}映射 */
const {{.NameLower}}Map: Record<EnumValue, string> = {
{{- range .EnumValues}}
  {{if IsNumeric .Value}}{{.Value}}{{else}}'{{.Value}}'{{end}}: '{{.Label}}',
{{- end}}
};

/** {{.Label}}颜色 */
function get{{.NameCamel}}Color(val: EnumValue | null | undefined): string {
  const keys: EnumValue[] = [{{range $i, $v := .EnumValues}}{{if $i}}, {{end}}{{if IsNumeric $v.Value}}{{$v.Value}}{{else}}'{{$v.Value}}'{{end}}{{end}}];
  if (val === null || val === undefined || val === '') {
    return TAG_COLORS[0] ?? 'default';
  }
  const idx = keys.indexOf(val);
  return TAG_COLORS[idx >= 0 ? idx % TAG_COLORS.length : 0] ?? 'default';
}
{{- end}}
{{- end}}
{{- if .HasTooltip}}
/** 渲染带 Tooltip 的列标题 */
function tooltipHeader(label: string, tip: string) {
  return () => h('span', {}, [
    label + ' ',
    h(Tooltip, { title: tip }, {
      default: () => h(QuestionCircleOutlined, { style: { color: '#999', marginLeft: '4px' } }),
    }),
  ]);
}
{{- end}}

/** 表单弹窗 */
const [FormModalComp, formModalApi] = useVbenModal({
  connectedComponent: FormModal,
  destroyOnClose: true,
});

/** 详情抽屉 */
const [DetailDrawerComp, detailDrawerApi] = useVbenModal({
  connectedComponent: DetailDrawer,
  destroyOnClose: true,
});
const { hasAccessByCodes } = useAccess();
const canBatchDelete = hasAccessByCodes(['{{.AppName}}:{{.ModuleName}}:batch-delete']);
{{- if .HasTenantScope}}
const isPlatformSuperAdmin = usePlatformSuperAdmin();
{{- end}}

/** 搜索表单配置 */
const formOptions: VbenFormProps = {
  collapsed: false,
  showCollapseButton: true,
  submitOnChange: false,
  submitOnEnter: true,
  schema: [
{{- if .HasKeywordSearch}}
    {
      component: 'Input',
      componentProps: { placeholder: '请输入关键词', allowClear: true },
      fieldName: 'keyword',
      label: '关键词',
    },
{{- end}}
{{- range .SearchFields}}
{{- $isScopeField := or (eq .Name "tenant_id") (eq .Name "merchant_id")}}
{{- if and $.HasTenantScope $isScopeField}}
    ...(isPlatformSuperAdmin.value ? [
{{- end}}
{{- if .SearchRange}}
    {
      component: 'RangePicker',
      fieldName: '{{.SearchFormField}}',
      label: '{{.ShortLabel}}',
      componentProps: {
        showTime: true,
        format: 'YYYY-MM-DD HH:mm:ss',
        valueFormat: 'YYYY-MM-DD HH:mm:ss',
        class: 'w-full',
      },
    },
{{- else if eq .SearchComponent "Input"}}
    {
      component: 'Input',
      componentProps: { placeholder: '请输入{{.ShortLabel}}', allowClear: true },
      fieldName: '{{.SearchFormField}}',
      label: '{{.ShortLabel}}',
    },
{{- else if eq .SearchComponent "Select"}}
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
{{- if .IsEnum}}
        options: {{.NameLower}}Options,
{{- else}}
        options: [],
{{- end}}
        placeholder: '请选择{{.Label}}',
        class: 'w-full',
      },
      fieldName: '{{.SearchFormField}}',
      label: '{{.ShortLabel}}',
    },
{{- else if eq .SearchComponent "TreeSelect"}}
    {
      component: 'TreeSelect',
      componentProps: {
        treeData: [],
        fieldNames: { label: '{{if .RefDisplayField}}{{.RefDisplayLower}}{{else}}title{{end}}', value: 'id', children: 'children' },
        placeholder: '请选择{{.Label}}',
        allowClear: true,
        treeDefaultExpandAll: true,
        class: 'w-full',
      },
      fieldName: '{{.SearchFormField}}',
      label: '{{.ShortLabel}}',
    },
{{- end}}
{{- if and $.HasTenantScope $isScopeField}}
    ] : []),
{{- end}}
{{- end}}
    {
      component: 'RangePicker',
      fieldName: 'timeRange',
      label: '创建时间',
      componentProps: {
        showTime: true,
        format: 'YYYY-MM-DD HH:mm:ss',
        valueFormat: 'YYYY-MM-DD HH:mm:ss',
        class: 'w-full',
      },
    },
  ],
};

/** 表格列配置 */
const gridOptions: VxeGridProps<{{.ModelName}}Item> = {
  checkboxConfig: canBatchDelete ? { highlight: true } : undefined,
  columns: [
    { title: '序号', type: 'seq', width: 50 },
    ...(canBatchDelete ? [{ type: 'checkbox', width: 50 }] : []),
{{- $isTree := .HasParentID}}
{{- $firstDataCol := true}}
{{- range .Fields}}
{{- if and (not .IsHidden) (not .IsID) (not .IsParentID) (not .IsTimeField) (not .IsMultiFK) (not .IsPassword)}}
{{- $isScopeField := or (eq .Name "tenant_id") (eq .Name "merchant_id")}}
{{- if and $.HasTenantScope $isScopeField}}
    ...(isPlatformSuperAdmin.value ? [
{{- end}}
{{- $isSortable := or .IsMoney .IsSearchable (eq .Name "sort") (eq .Name "status")}}
{{- if .RefFieldJSON}}
    { field: '{{.RefFieldJSON}}', title: '{{.ShortLabel}}'{{if .TooltipText}}, slots: { header: tooltipHeader('{{.ShortLabel}}', '{{.TooltipText}}') }{{end}}{{if and $isTree $firstDataCol}}, treeNode: true{{end}} },
{{- else if .IsEnum}}
    { field: '{{.NameLower}}', title: '{{.ShortLabel}}', width: 120, slots: { default: '{{.NameLower}}_cell' }{{if and $isTree $firstDataCol}}, treeNode: true{{end}}{{if and $isSortable (not $isTree)}}, sortable: true{{end}} },
{{- else if eq .Component "ImageUpload"}}
    { field: '{{.NameLower}}', title: '{{.ShortLabel}}', width: 100, slots: { default: '{{.NameLower}}_cell' }{{if and $isTree $firstDataCol}}, treeNode: true{{end}} },
{{- else if eq .Component "InputUrl"}}
    { field: '{{.NameLower}}', title: '{{.ShortLabel}}', slots: { default: '{{.NameLower}}_cell' }{{if and $isTree $firstDataCol}}, treeNode: true{{end}} },
{{- else if eq .Component "FileUpload"}}
    { field: '{{.NameLower}}', title: '{{.ShortLabel}}', slots: { default: '{{.NameLower}}_cell' }{{if and $isTree $firstDataCol}}, treeNode: true{{end}} },
{{- else if or (eq .Component "RichText") (eq .Component "JsonEditor")}}
{{- /* 富文本和JSON字段不在列表中显示，不消耗 firstDataCol */}}
{{- else if .IsMoney}}
    { field: '{{.NameLower}}', title: '{{.ShortLabel}}'{{if .TooltipText}}, slots: { header: tooltipHeader('{{.ShortLabel}}', '{{.TooltipText}}') }{{end}}, width: 120, formatter: ({ cellValue }: any) => cellValue != null ? (cellValue / 100).toFixed(2) : '-'{{if and $isTree $firstDataCol}}, treeNode: true{{end}}{{if not $isTree}}, sortable: true{{end}} },
{{- else}}
    { field: '{{.NameLower}}', title: '{{.ShortLabel}}'{{if .TooltipText}}, slots: { header: tooltipHeader('{{.ShortLabel}}', '{{.TooltipText}}') }{{end}}{{if and $isTree $firstDataCol}}, treeNode: true{{end}}{{if and $isSortable (not $isTree)}}, sortable: true{{end}} },
{{- end}}
{{- if and $.HasTenantScope $isScopeField}}
    ] : []),
{{- end}}
{{- if not (or (eq .Component "RichText") (eq .Component "JsonEditor") (and $.HasTenantScope $isScopeField))}}
{{- $firstDataCol = false}}
{{- end}}
{{- end}}
{{- end}}
{{- range .Fields}}
{{- if and (not .IsHidden) (.IsTimeField)}}
    { field: '{{.NameLower}}', title: '{{.ShortLabel}}'{{if .TooltipText}}, slots: { header: tooltipHeader('{{.ShortLabel}}', '{{.TooltipText}}') }{{end}}, width: 180, formatter: 'formatDateTime' },
{{- end}}
{{- end}}
    { field: 'createdAt', title: '创建时间', width: 180, formatter: 'formatDateTime'{{if not $isTree}}, sortable: true{{end}} },
    { title: '操作', width: 240, fixed: 'right', slots: { default: 'action' } },
  ],
{{- if .HasParentID}}
  height: 'auto',
  pagerConfig: { enabled: false },
  treeConfig: {
    childrenField: 'children',
    expandAll: false,
  },
  proxyConfig: {
    ajax: {
      query: async (_params, formValues) => {
        const params: Record<string, any> = { ...formValues };
        if (params.timeRange && params.timeRange.length === 2) {
          params.startTime = params.timeRange[0];
          params.endTime = params.timeRange[1];
        }
        delete params.timeRange;
{{- if .HasTenantScope}}
        if (!isPlatformSuperAdmin.value) {
{{- range .SearchFields}}
{{- if or (eq .Name "tenant_id") (eq .Name "merchant_id")}}
          delete params.{{.SearchFormField}};
{{- end}}
{{- end}}
        }
{{- end}}
{{- range .SearchFields}}
{{- if .SearchRange}}
        if (params.{{.SearchFormField}} && params.{{.SearchFormField}}.length === 2) {
          params.{{.NameLower}}Start = params.{{.SearchFormField}}[0];
          params.{{.NameLower}}End = params.{{.SearchFormField}}[1];
        }
        delete params.{{.SearchFormField}};
{{- end}}
{{- end}}
        return await get{{.ModelName}}Tree(params as any) ?? [];
      },
    },
  },
{{- else}}
  height: 'auto',
  pagerConfig: {},
  proxyConfig: {
    ajax: {
      query: async ({ page, sorts }, formValues) => {
        const params: Record<string, any> = {
          pageNum: page.currentPage,
          pageSize: page.pageSize,
          ...formValues,
        };
        if (params.timeRange && params.timeRange.length === 2) {
          params.startTime = params.timeRange[0];
          params.endTime = params.timeRange[1];
        }
        delete params.timeRange;
{{- if .HasTenantScope}}
        if (!isPlatformSuperAdmin.value) {
{{- range .SearchFields}}
{{- if or (eq .Name "tenant_id") (eq .Name "merchant_id")}}
          delete params.{{.SearchFormField}};
{{- end}}
{{- end}}
        }
{{- end}}
{{- range .SearchFields}}
{{- if .SearchRange}}
        if (params.{{.SearchFormField}} && params.{{.SearchFormField}}.length === 2) {
          params.{{.NameLower}}Start = params.{{.SearchFormField}}[0];
          params.{{.NameLower}}End = params.{{.SearchFormField}}[1];
        }
        delete params.{{.SearchFormField}};
{{- end}}
{{- end}}
        if (sorts && sorts.length > 0) {
          const sort = sorts[0];
          if (sort && sort.field && sort.order) {
            params.orderBy = resolveSortField(String(sort.field));
            params.orderDir = sort.order;
          }
        }
        const res = await get{{.ModelName}}List(params as any);
        return { items: res?.list ?? [], total: res?.total ?? 0 };
      },
    },
  },
{{- end}}
{{- if not .HasParentID}}
  sortConfig: {
    remote: true,
    trigger: 'cell',
{{- if .HasSort}}
    defaultSort: { field: 'sort', order: 'asc' },
{{- else}}
    defaultSort: { field: 'createdAt', order: 'desc' },
{{- end}}
  },
{{- end}}
  toolbarConfig: {
    custom: true,
    refresh: true,
    search: true,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions,
  gridOptions,
});
{{- if .HasImport}}

const importInputRef = ref<HTMLInputElement | null>(null);
{{- end}}

async function initSearchOptions() {
{{- range .SearchFields}}
{{- if .IsParentID}}
{{- $isScopeField := or (eq .Name "tenant_id") (eq .Name "merchant_id")}}
{{- if and $.HasTenantScope $isScopeField}}
  if (isPlatformSuperAdmin.value) {
{{- end}}
  try {
    const treeRes = await get{{$.ModelName}}Tree();
    gridApi.formApi.updateSchema([
      {
        fieldName: '{{.SearchFormField}}',
        componentProps: {
          treeData: [
            { id: '0', {{if .RefDisplayLower}}{{.RefDisplayLower}}{{else}}title{{end}}: '顶级节点', children: treeRes ?? [] },
          ],
        },
      },
    ]);
  } catch {
    // ignore
  }
{{- if and $.HasTenantScope $isScopeField}}
  }
{{- end}}
{{- end}}
{{- end}}
{{- range .SearchFields}}
{{- if and .IsForeignKey .RefTable}}
{{- $isScopeField := or (eq .Name "tenant_id") (eq .Name "merchant_id")}}
{{- if and $.HasTenantScope $isScopeField}}
  if (isPlatformSuperAdmin.value) {
{{- end}}
{{- if .RefIsTree}}
  try {
    const {{.NameLower}}Tree = await get{{.RefTableCamel}}Tree();
    gridApi.formApi.updateSchema([
      {
        fieldName: '{{.SearchFormField}}',
        componentProps: { treeData: {{.NameLower}}Tree ?? [] },
      },
    ]);
  } catch {
    // ignore
  }
{{- else}}
  try {
    const {{.NameLower}}Res = await get{{.RefTableCamel}}List({ pageNum: 1, pageSize: 1000 });
    gridApi.formApi.updateSchema([
      {
        fieldName: '{{.SearchFormField}}',
        componentProps: {
          options: ({{.NameLower}}Res?.list ?? []).map((item: any) => ({
            label: item.{{.RefDisplayLower}} || item.id,
            value: item.id,
          })),
        },
      },
    ]);
  } catch {
    // ignore
  }
{{- end}}
{{- if and $.HasTenantScope $isScopeField}}
  }
{{- end}}
{{- end}}
{{- end}}
{{- range .SearchFields}}
{{- if .DictType}}
{{- $isScopeField := or (eq .Name "tenant_id") (eq .Name "merchant_id")}}
{{- if and $.HasTenantScope $isScopeField}}
  if (isPlatformSuperAdmin.value) {
{{- end}}
  try {
    const {{.NameLower}}Dict = await getDictByType('{{.DictType}}');
    gridApi.formApi.updateSchema([
      {
        fieldName: '{{.SearchFormField}}',
        componentProps: {
          options: ({{.NameLower}}Dict ?? []).map((item: any) => ({
            label: item.label,
            value: item.value,
          })),
        },
      },
    ]);
  } catch {
    // ignore
  }
{{- if and $.HasTenantScope $isScopeField}}
  }
{{- end}}
{{- end}}
{{- end}}
}

onMounted(() => {
  void initSearchOptions();
});

/** 新建 */
function handleCreate() {
  formModalApi.setData(null).open();
}

/** 查看 */
function handleView(row: {{.ModelName}}Item) {
  detailDrawerApi.setData({ id: row.id }).open();
}

/** 编辑 */
function handleEdit(row: {{.ModelName}}Item) {
  formModalApi.setData({ id: row.id }).open();
}

/** 删除 */
function handleDelete(row: {{.ModelName}}Item) {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除该{{.Comment}}吗？',
    okType: 'danger',
    async onOk() {
      await delete{{.ModelName}}(row.id);
      message.success('删除成功');
      gridApi.reload();
    },
  });
}

/** 批量删除 */
function handleBatchDelete() {
  const ids = getGridSelectedIds<{{.ModelName}}Item>(gridApi.grid as any);
  if (ids.length === 0) {
    message.warning('请先选择要删除的数据');
    return;
  }
  Modal.confirm({
    title: '确认批量删除',
    content: `确定要删除选中的 ${ids.length} 条{{.Comment}}吗？`,
    okType: 'danger',
    async onOk() {
      await batchDelete{{.ModelName}}(ids);
      message.success('批量删除成功');
      gridApi.reload();
    },
  });
}

/** 导出 */
async function handleExport() {
  try {
    const formValues = await gridApi.formApi.getValues();
    const params: Record<string, any> = { ...formValues };
{{- if not .HasParentID}}
    const sorts = gridApi.grid?.getSortColumns?.() ?? [];
{{- end}}
    if (params.timeRange && params.timeRange.length === 2) {
      params.startTime = params.timeRange[0];
      params.endTime = params.timeRange[1];
    }
    delete params.timeRange;
{{- if .HasTenantScope}}
    if (!isPlatformSuperAdmin.value) {
{{- range .SearchFields}}
{{- if or (eq .Name "tenant_id") (eq .Name "merchant_id")}}
      delete params.{{.SearchFormField}};
{{- end}}
{{- end}}
    }
{{- end}}
{{- range .SearchFields}}
{{- if .SearchRange}}
    if (params.{{.SearchFormField}} && params.{{.SearchFormField}}.length === 2) {
      params.{{.NameLower}}Start = params.{{.SearchFormField}}[0];
      params.{{.NameLower}}End = params.{{.SearchFormField}}[1];
    }
    delete params.{{.SearchFormField}};
{{- end}}
{{- end}}
{{- if not .HasParentID}}
    if (sorts.length > 0) {
      const sort = sorts[0];
      if (sort?.field && sort?.order) {
        params.orderBy = resolveSortField(String(sort.field));
        params.orderDir = sort.order;
      }
    }
{{- end}}
    const blob = await export{{.ModelName}}(params);
    downloadFileFromBlob({ fileName: '{{.Comment}}.csv', source: blob as Blob });
    message.success('导出成功');
  } catch {
    message.error('导出失败');
  }
}
{{- if .HasImport}}

function handleImportTrigger() {
  const input = importInputRef.value;
  if (!input) {
    return;
  }
  input.value = '';
  input.click();
}

/** 导入 */
async function handleImportChange(event: Event) {
  const input = event.target as HTMLInputElement | null;
  const file = input?.files?.[0];
  if (!file) {
    return;
  }
  const formData = new FormData();
  formData.append('file', file);
  try {
    const res = await import{{.ModelName}}(formData);
    message.success(`导入完成：成功 ${res?.success ?? 0} 条，失败 ${res?.fail ?? 0} 条`);
    gridApi.reload();
  } catch {
    message.error('导入失败');
  } finally {
    if (input) {
      input.value = '';
    }
  }
}

/** 下载导入模板 */
async function handleDownloadTemplate() {
  try {
    const blob = await downloadImportTemplate{{.ModelName}}();
    downloadFileFromBlob({ fileName: '{{.Comment}}导入模板.csv', source: blob as Blob });
  } catch {
    message.error('下载模板失败');
  }
}
{{- end}}
{{- if .HasBatchEdit}}

/** 批量修改状态 */
function handleBatchUpdateStatus() {
  const ids = getGridSelectedIds<{{.ModelName}}Item>(gridApi.grid as any);
  if (ids.length === 0) {
    message.warning('请先选择要修改的数据');
    return;
  }
  const rows = gridApi.grid.getCheckboxRecords() as {{.ModelName}}Item[];
  Modal.confirm({
    title: '批量修改状态',
    content: `确定要将选中的 ${ids.length} 条数据的状态切换吗？`,
    async onOk() {
      const newStatus = rows[0]?.status === 1 ? 0 : 1;
      await batchUpdate{{.ModelName}}({ ids, status: newStatus });
      message.success('批量修改成功');
      gridApi.reload();
    },
  });
}
{{- end}}
</script>

<template>
  <Page auto-content-height>
    <FormModalComp @success="() => gridApi.reload()" />
    <DetailDrawerComp />
{{- if .HasImport}}
    <input
      ref="importInputRef"
      type="file"
      accept=".csv"
      class="hidden"
      @change="handleImportChange"
    />
{{- end}}
    <Grid>
      <template #toolbar-actions>
        <Button v-access:code="'{{.AppName}}:{{.ModuleName}}:create'" type="primary" @click="handleCreate">新建</Button>
        <Button v-access:code="'{{.AppName}}:{{.ModuleName}}:batch-delete'" danger class="ml-2" @click="handleBatchDelete">批量删除</Button>
        <Button v-access:code="'{{.AppName}}:{{.ModuleName}}:export'" class="ml-2" @click="handleExport">导出</Button>
{{- if .HasImport}}
        <Button v-access:code="'{{.AppName}}:{{.ModuleName}}:import'" class="ml-2" @click="handleImportTrigger">导入</Button>
        <Button v-access:code="'{{.AppName}}:{{.ModuleName}}:import'" class="ml-2" @click="handleDownloadTemplate">模板下载</Button>
{{- end}}
{{- if .HasBatchEdit}}
        <Button v-access:code="'{{.AppName}}:{{.ModuleName}}:batch-update'" class="ml-2" @click="handleBatchUpdateStatus">批量修改状态</Button>
{{- end}}
      </template>
{{- range .Fields}}
{{- if and (not .IsHidden) (not .IsID) (not .IsParentID) (not .IsTimeField) (not .IsMultiFK) (.IsEnum)}}
      <template #{{.NameLower}}_cell="{ row }">
        <Tag :color="get{{.NameCamel}}Color(row.{{.NameLower}})">
          {{"{{"}} getEnumLabel({{.NameLower}}Map, row.{{.NameLower}}) {{"}}"}}
        </Tag>
      </template>
{{- else if and (not .IsHidden) (not .IsID) (not .IsParentID) (not .IsTimeField) (not .IsMultiFK) (eq .Component "ImageUpload")}}
      <template #{{.NameLower}}_cell="{ row }">
        <img v-if="row.{{.NameLower}} && /^https?:\/\//i.test(row.{{.NameLower}})" :src="row.{{.NameLower}}" style="width: 48px; height: 48px; object-fit: cover; border-radius: 4px;" />
        <span v-else>-</span>
      </template>
{{- else if and (not .IsHidden) (not .IsID) (not .IsParentID) (not .IsTimeField) (not .IsMultiFK) (eq .Component "InputUrl")}}
      <template #{{.NameLower}}_cell="{ row }">
        <a v-if="row.{{.NameLower}} && /^https?:\/\//i.test(row.{{.NameLower}})" :href="row.{{.NameLower}}" target="_blank" rel="noreferrer noopener" style="color: #1890ff;">{{"{{"}} row.{{.NameLower}} {{"}}"}}</a>
        <span v-else-if="row.{{.NameLower}}">{{"{{"}} row.{{.NameLower}} {{"}}"}}</span>
        <span v-else>-</span>
      </template>
{{- else if and (not .IsHidden) (not .IsID) (not .IsParentID) (not .IsTimeField) (not .IsMultiFK) (eq .Component "FileUpload")}}
      <template #{{.NameLower}}_cell="{ row }">
        <a v-if="row.{{.NameLower}} && /^https?:\/\//i.test(row.{{.NameLower}})" :href="row.{{.NameLower}}" target="_blank" rel="noreferrer noopener" style="color: #1890ff;">下载</a>
        <span v-else-if="row.{{.NameLower}}">{{"{{"}} row.{{.NameLower}} {{"}}"}}</span>
        <span v-else>-</span>
      </template>
{{- end}}
{{- end}}
      <template #action="{ row }">
        <Button v-access:code="'{{.AppName}}:{{.ModuleName}}:detail'" type="link" size="small" @click="handleView(row)">查看</Button>
        <Button v-access:code="'{{.AppName}}:{{.ModuleName}}:update'" type="link" size="small" @click="handleEdit(row)">编辑</Button>
        <Button v-access:code="'{{.AppName}}:{{.ModuleName}}:delete'" type="link" danger size="small" @click="handleDelete(row)">删除</Button>
      </template>
    </Grid>
  </Page>
</template>
