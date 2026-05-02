<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { Page } from '@vben/common-ui';
import {
  Button,
  Card,
  Form,
  FormItem,
  Input,
  message,
  Select,
  Space,
  Table,
  Tag,
} from 'ant-design-vue';
import { DownloadOutlined, ReloadOutlined } from '@ant-design/icons-vue';
import {
  getContractDownloadURL,
  getContractList,
} from '#/api/member/contract';
import type {
  ContractItem,
  ContractListParams,
} from '#/api/member/contract/types';

const loading = ref(false);
const list = ref<ContractItem[]>([]);
const total = ref(0);
const pagination = reactive({ current: 1, pageSize: 20, total: 0 });

const filter = reactive<ContractListParams>({
  userId: '',
  contractType: '',
  pageNum: 1,
  pageSize: 20,
});

const typeOptions = [
  { label: '全部', value: '' },
  { label: '注册协议', value: 'register' },
  { label: '升级协议', value: 'upgrade' },
  { label: '自定义', value: 'custom' },
];

const typeText: Record<string, string> = {
  register: '注册协议',
  upgrade: '升级协议',
  custom: '自定义',
};

function pdfTagColor(s: number) {
  return s === 2 ? 'green' : s === 3 ? 'red' : 'orange';
}

const columns = [
  { title: '合同编号', dataIndex: 'contractNo', width: 220 },
  { title: '会员', dataIndex: 'user', width: 200 },
  { title: '类型', dataIndex: 'contractType', width: 120 },
  { title: '签署时间', dataIndex: 'signedAt', width: 180 },
  { title: '签署IP', dataIndex: 'signedIp', width: 140 },
  { title: 'PDF 状态', dataIndex: 'pdfStatus', width: 120 },
  { title: '操作', dataIndex: 'op', width: 140, fixed: 'right' as const },
];

async function load() {
  loading.value = true;
  try {
    filter.pageNum = pagination.current;
    filter.pageSize = pagination.pageSize;
    const r = await getContractList({
      userId: filter.userId || undefined,
      contractType: filter.contractType || undefined,
      pageNum: filter.pageNum,
      pageSize: filter.pageSize,
    });
    list.value = r.list;
    total.value = r.total;
    pagination.total = r.total;
  } finally {
    loading.value = false;
  }
}

function reset() {
  filter.userId = '';
  filter.contractType = '';
  pagination.current = 1;
  load();
}

function search() {
  pagination.current = 1;
  load();
}

function onTableChange(p: any) {
  pagination.current = p.current;
  pagination.pageSize = p.pageSize;
  load();
}

function download(record: ContractItem) {
  if (record.pdfStatus !== 2) {
    message.warning(record.pdfStatusText || 'PDF 尚未就绪');
  }
  window.open(getContractDownloadURL(record.contractId), '_blank');
}

onMounted(load);
</script>

<template>
  <Page>
    <Card title="会员合同" :bodyStyle="{ paddingBottom: '24px' }">
      <Form layout="inline" style="margin-bottom: 16px">
        <FormItem label="会员 ID">
          <Input
            v-model:value="filter.userId"
            placeholder="按会员 ID 过滤"
            allow-clear
            style="width: 200px"
          />
        </FormItem>
        <FormItem label="合同类型">
          <Select
            v-model:value="filter.contractType"
            :options="typeOptions"
            style="width: 160px"
          />
        </FormItem>
        <FormItem>
          <Space>
            <Button type="primary" @click="search">搜索</Button>
            <Button @click="reset">
              <ReloadOutlined />
              重置
            </Button>
          </Space>
        </FormItem>
      </Form>

      <Table
        :columns="columns"
        :data-source="list"
        :loading="loading"
        :pagination="{
          current: pagination.current,
          pageSize: pagination.pageSize,
          total: pagination.total,
          showSizeChanger: true,
          showTotal: (t: number) => `共 ${t} 条`,
        }"
        row-key="contractId"
        :scroll="{ x: 1100 }"
        @change="onTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'user'">
            <div>{{ record.userNickname || record.userId }}</div>
            <div style="color: #888; font-size: 12px">
              {{ record.userPhone }}
            </div>
          </template>
          <template v-else-if="column.dataIndex === 'contractType'">
            {{ typeText[record.contractType] || record.contractType }}
          </template>
          <template v-else-if="column.dataIndex === 'pdfStatus'">
            <Tag :color="pdfTagColor(record.pdfStatus)">
              {{ record.pdfStatusText }}
            </Tag>
          </template>
          <template v-else-if="column.dataIndex === 'op'">
            <Button type="link" size="small" @click="download(record)">
              <DownloadOutlined />
              下载
            </Button>
          </template>
        </template>
      </Table>
    </Card>
  </Page>
</template>
