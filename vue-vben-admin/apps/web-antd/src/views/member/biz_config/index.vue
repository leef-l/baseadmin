<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { Page } from '@vben/common-ui';
import {
  Button,
  Card,
  Checkbox,
  CheckboxGroup,
  Form,
  FormItem,
  Input,
  InputNumber,
  message,
  Space,
  Spin,
  Table,
  TimePicker,
} from 'ant-design-vue';
import { DeleteOutlined, PlusOutlined } from '@ant-design/icons-vue';
import dayjs, { type Dayjs } from 'dayjs';
import { getBizConfig, saveBizConfig } from '#/api/member/biz_config';
import type {
  BizConfigData,
  BizConfigRebate,
} from '#/api/member/biz_config/types';

const loading = ref(false);
const submitting = ref(false);

const weekdayOptions = [
  { label: '周一', value: 1 },
  { label: '周二', value: 2 },
  { label: '周三', value: 3 },
  { label: '周四', value: 4 },
  { label: '周五', value: 5 },
  { label: '周六', value: 6 },
  { label: '周日', value: 7 },
];

const state = reactive<{
  purchaseStart: Dayjs | null;
  purchaseEnd: Dayjs | null;
  allowedWeekdays: number[];
  consignStart: Dayjs | null;
  consignEnd: Dayjs | null;
  selfTurnoverRewardRate: number;
  directPromoteRate: number;
  tiers: BizConfigRebate[];
}>({
  purchaseStart: null,
  purchaseEnd: null,
  allowedWeekdays: [1, 2, 3, 4, 5],
  consignStart: null,
  consignEnd: null,
  selfTurnoverRewardRate: 1,
  directPromoteRate: 0.4,
  tiers: [],
});

const tierColumns = [
  { title: '第几单', dataIndex: 'nthOrder', width: 120 },
  { title: '奖励金额（元）', dataIndex: 'rewardYuan', width: 180 },
  { title: '操作', dataIndex: 'op', width: 100 },
];

function parseHHMM(s: string): Dayjs | null {
  if (!s) return null;
  return dayjs(s, 'HH:mm');
}
function fmtHHMM(d: Dayjs | null): string {
  return d ? d.format('HH:mm') : '';
}

async function load() {
  loading.value = true;
  try {
    const r = await getBizConfig();
    state.purchaseStart = parseHHMM(r.purchase.startTime);
    state.purchaseEnd = parseHHMM(r.purchase.endTime);
    state.allowedWeekdays = r.purchase.allowedWeekdays || [];
    state.consignStart = parseHHMM(r.consign.startTime);
    state.consignEnd = r.consign.endTime ? parseHHMM(r.consign.endTime) : null;
    state.selfTurnoverRewardRate = r.selfTurnoverRewardRate || 0;
    state.directPromoteRate = r.directPromoteRate || 0;
    state.tiers = (r.selfRebateTiers || []).map((t) => ({ ...t }));
  } finally {
    loading.value = false;
  }
}

function addTier() {
  const next =
    state.tiers.length === 0
      ? 2
      : Math.max(...state.tiers.map((t) => t.nthOrder)) + 1;
  state.tiers.push({ nthOrder: next, rewardYuan: 0 });
}
function removeTier(i: number) {
  state.tiers.splice(i, 1);
}

async function submit() {
  if (!state.purchaseStart || !state.purchaseEnd) {
    message.error('请填写进货时间窗');
    return;
  }
  if (!state.consignStart) {
    message.error('请填写寄售开始时间');
    return;
  }
  if (state.allowedWeekdays.length === 0) {
    message.error('至少选择一个允许进货的工作日');
    return;
  }
  // 校验阶梯档位 nthOrder 唯一
  const seen = new Set<number>();
  for (const t of state.tiers) {
    if (!t.nthOrder || t.nthOrder < 1) {
      message.error('阶梯档位"第几单"必须 ≥ 1');
      return;
    }
    if (seen.has(t.nthOrder)) {
      message.error(`阶梯档位"第 ${t.nthOrder} 单"重复`);
      return;
    }
    seen.add(t.nthOrder);
  }

  const payload: BizConfigData = {
    purchase: {
      startTime: fmtHHMM(state.purchaseStart),
      endTime: fmtHHMM(state.purchaseEnd),
      allowedWeekdays: [...state.allowedWeekdays].sort((a, b) => a - b),
    },
    consign: {
      startTime: fmtHHMM(state.consignStart),
      endTime: state.consignEnd ? fmtHHMM(state.consignEnd) : null,
    },
    selfRebateTiers: state.tiers
      .slice()
      .sort((a, b) => a.nthOrder - b.nthOrder),
    selfTurnoverRewardRate: Number(state.selfTurnoverRewardRate),
    directPromoteRate: Number(state.directPromoteRate),
  };

  submitting.value = true;
  try {
    await saveBizConfig(payload);
    message.success('保存成功');
    await load();
  } finally {
    submitting.value = false;
  }
}

onMounted(load);
</script>

<template>
  <Page>
    <Spin :spinning="loading">
      <Card title="会员业务配置" :bodyStyle="{ paddingBottom: '24px' }">
        <Form layout="vertical">
          <Card size="small" title="进货时间窗" style="margin-bottom: 16px">
            <FormItem label="开始时间" required>
              <TimePicker
                v-model:value="state.purchaseStart"
                format="HH:mm"
                :minute-step="5"
                placeholder="如 10:00"
              />
            </FormItem>
            <FormItem label="结束时间" required>
              <TimePicker
                v-model:value="state.purchaseEnd"
                format="HH:mm"
                :minute-step="5"
                placeholder="如 10:30"
              />
            </FormItem>
            <FormItem label="允许进货的工作日" required>
              <CheckboxGroup
                v-model:value="state.allowedWeekdays"
                :options="weekdayOptions"
              />
            </FormItem>
          </Card>

          <Card size="small" title="寄售时间窗" style="margin-bottom: 16px">
            <FormItem label="开始时间" required>
              <TimePicker
                v-model:value="state.consignStart"
                format="HH:mm"
                :minute-step="5"
                placeholder="如 14:30"
              />
            </FormItem>
            <FormItem label="结束时间（留空表示当天无截止）">
              <TimePicker
                v-model:value="state.consignEnd"
                format="HH:mm"
                :minute-step="5"
                placeholder="留空 = 不截止"
                allow-clear
              />
            </FormItem>
          </Card>

          <Card size="small" title="自购阶梯返佣（一生维度，每档位仅奖一次）" style="margin-bottom: 16px">
            <Table
              :columns="tierColumns"
              :data-source="state.tiers"
              :pagination="false"
              size="small"
              row-key="nthOrder"
            >
              <template #bodyCell="{ column, index, record }">
                <template v-if="column.dataIndex === 'nthOrder'">
                  <InputNumber
                    v-model:value="record.nthOrder"
                    :min="1"
                    :max="999"
                    style="width: 100%"
                  />
                </template>
                <template v-else-if="column.dataIndex === 'rewardYuan'">
                  <InputNumber
                    v-model:value="record.rewardYuan"
                    :min="0"
                    :step="1"
                    :precision="2"
                    style="width: 100%"
                  />
                </template>
                <template v-else-if="column.dataIndex === 'op'">
                  <Button danger size="small" @click="removeTier(index)">
                    <DeleteOutlined />
                  </Button>
                </template>
              </template>
            </Table>
            <Button
              type="dashed"
              block
              style="margin-top: 12px"
              @click="addTier"
            >
              <PlusOutlined /> 新增阶梯
            </Button>
          </Card>

          <Card size="small" title="比例返佣" style="margin-bottom: 16px">
            <FormItem label="自购返奖比例（%）">
              <InputNumber
                v-model:value="state.selfTurnoverRewardRate"
                :min="0"
                :max="100"
                :precision="2"
                :step="0.1"
                style="width: 200px"
              />
              <span style="color: #888; margin-left: 8px">
                每笔自购订单按该比例入"奖励钱包"
              </span>
            </FormItem>
            <FormItem label="直推返奖比例（%）">
              <InputNumber
                v-model:value="state.directPromoteRate"
                :min="0"
                :max="100"
                :precision="2"
                :step="0.1"
                style="width: 200px"
              />
              <span style="color: #888; margin-left: 8px">
                直推下级每笔进货按该比例入上级"推广钱包"（仅 1 层）
              </span>
            </FormItem>
          </Card>

          <Space>
            <Button type="primary" :loading="submitting" @click="submit">
              保存
            </Button>
            <Button @click="load">重置</Button>
          </Space>
        </Form>
      </Card>
    </Spin>
  </Page>
</template>
