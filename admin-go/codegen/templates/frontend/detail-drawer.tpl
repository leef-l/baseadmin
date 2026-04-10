<script setup lang="ts">
import { ref } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { Descriptions, DescriptionsItem{{if .HasEnum}}, Tag{{end}} } from 'ant-design-vue';
import { get{{.ModelName}}Detail } from '#/api/{{.AppName}}/{{.ModuleName}}';
import type { {{.ModelName}}Item } from '#/api/{{.AppName}}/{{.ModuleName}}/types';
{{- if .HasEnum}}

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
{{range .Fields}}
{{- if and (not .IsHidden) (not .IsID) (.IsEnum)}}
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
{{end}}
{{- end}}
const detail = ref<{{.ModelName}}Item | null>(null);
const openToken = ref(0);

function displayValue(value: null | number | string | undefined) {
  if (value === null || value === undefined || value === '') {
    return '-';
  }
  return value;
}

const [Modal, modalApi] = useVbenModal({
  fullscreenButton: false,
  footer: false,
  async onOpenChange(isOpen: boolean) {
    if (!isOpen) {
      openToken.value += 1;
      detail.value = null;
      return;
    }

    const currentOpenToken = ++openToken.value;
    const data = modalApi.getData<{ id: string }>();
    if (data?.id) {
      modalApi.setState({ title: '{{.Comment}}详情' });
      try {
        const res = await get{{.ModelName}}Detail(data.id);
        if (currentOpenToken !== openToken.value) {
          return;
        }
        detail.value = res;
      } catch {
        if (currentOpenToken === openToken.value) {
          detail.value = null;
        }
      }
    }
  },
});
</script>

<template>
  <Modal class="w-[600px]">
    <Descriptions v-if="detail" bordered :column="1" size="small">
      <DescriptionsItem label="ID">{{"{{"}} detail.id {{"}}"}}</DescriptionsItem>
{{- range .Fields}}
{{- if and (not .IsHidden) (not .IsID) (not .IsPassword) (not .IsTimeField)}}
{{- if .RefFieldJSON}}
      <DescriptionsItem label="{{.ShortLabel}}">{{"{{"}} detail.{{.RefFieldJSON}} || '-' {{"}}"}}</DescriptionsItem>
{{- else if .IsEnum}}
      <DescriptionsItem label="{{.ShortLabel}}">
        <Tag :color="get{{.NameCamel}}Color(detail.{{.NameLower}})">{{"{{"}} getEnumLabel({{.NameLower}}Map, detail.{{.NameLower}}) {{"}}"}}</Tag>
      </DescriptionsItem>
{{- else if .IsMoney}}
      <DescriptionsItem label="{{.ShortLabel}}">{{"{{"}} detail.{{.NameLower}} != null ? (detail.{{.NameLower}} / 100).toFixed(2) : '-' {{"}}"}}</DescriptionsItem>
{{- else if eq .Component "ImageUpload"}}
      <DescriptionsItem label="{{.ShortLabel}}">
        <img v-if="detail.{{.NameLower}}" :src="detail.{{.NameLower}}" style="max-width: 200px; max-height: 200px; object-fit: contain;" />
        <span v-else>-</span>
      </DescriptionsItem>
{{- else if eq .Component "FileUpload"}}
      <DescriptionsItem label="{{.ShortLabel}}">
        <a v-if="detail.{{.NameLower}}" :href="detail.{{.NameLower}}" target="_blank" rel="noreferrer noopener">查看文件</a>
        <span v-else>-</span>
      </DescriptionsItem>
{{- else if eq .Component "InputUrl"}}
      <DescriptionsItem label="{{.ShortLabel}}">
        <a v-if="detail.{{.NameLower}}" :href="detail.{{.NameLower}}" target="_blank" rel="noreferrer noopener">{{"{{"}} detail.{{.NameLower}} {{"}}"}}</a>
        <span v-else>-</span>
      </DescriptionsItem>
{{- else if eq .Component "RichText"}}
      <DescriptionsItem label="{{.ShortLabel}}">
        <div v-html="detail.{{.NameLower}}" style="max-height: 300px; overflow: auto;" />
      </DescriptionsItem>
{{- else if eq .Component "JsonEditor"}}
      <DescriptionsItem label="{{.ShortLabel}}">
        <pre style="max-height: 300px; overflow: auto; white-space: pre-wrap; word-break: break-all; margin: 0; font-size: 12px;">{{"{{"}} (() => { const value = detail.{{.NameLower}}; if (!value) return '-'; try { return JSON.stringify(JSON.parse(value), null, 2) } catch { return value } })() {{"}}"}}</pre>
      </DescriptionsItem>
{{- else}}
      <DescriptionsItem label="{{.ShortLabel}}">{{"{{"}} displayValue(detail.{{.NameLower}}) {{"}}"}}</DescriptionsItem>
{{- end}}
{{- end}}
{{- end}}
{{- range .Fields}}
{{- if and (not .IsHidden) (not .IsID) (not .IsPassword) (.IsTimeField)}}
      <DescriptionsItem label="{{.ShortLabel}}">{{"{{"}} displayValue(detail.{{.NameLower}}) {{"}}"}}</DescriptionsItem>
{{- end}}
{{- end}}
      <DescriptionsItem label="创建时间">{{"{{"}} displayValue(detail.createdAt) {{"}}"}}</DescriptionsItem>
      <DescriptionsItem label="更新时间">{{"{{"}} displayValue(detail.updatedAt) {{"}}"}}</DescriptionsItem>
    </Descriptions>
  </Modal>
</template>
