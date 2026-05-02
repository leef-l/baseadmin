/** 进货时间窗 + 工作日 */
export interface BizConfigPurchase {
  startTime: string;
  endTime: string;
  allowedWeekdays: number[];
}

/** 寄售时间窗（endTime 为 null 表示无截止） */
export interface BizConfigConsign {
  startTime: string;
  endTime: string | null;
}

/** 自购阶梯返佣档位 */
export interface BizConfigRebate {
  nthOrder: number;
  rewardYuan: number;
}

/** 业务配置（与后端 logic.bizconfig.Config 同构） */
export interface BizConfigData {
  purchase: BizConfigPurchase;
  consign: BizConfigConsign;
  selfRebateTiers: BizConfigRebate[];
  selfTurnoverRewardRate: number;
  directPromoteRate: number;
}
