import { useEffect, useMemo, useState } from 'react';
import { bizConfigApi } from '@/api/bizConfig';
import { PortalBizConfig } from '@/api/types';

export interface PurchaseWindowState {
  cfg: PortalBizConfig | null;
  now: Date;
  isToday: boolean;          // 今日是否在允许工作日内
  isInWindow: boolean;       // 当前是否在购买时间窗内
  countdownSeconds: number;  // 距开始或距结束剩余秒数
  reason: string;            // 不可购买时的原因
}

const HHMMtoToday = (now: Date, hhmm: string): Date => {
  const [h, m] = hhmm.split(':').map(Number);
  const d = new Date(now);
  d.setHours(h, m, 0, 0);
  return d;
};

const calc = (cfg: PortalBizConfig, now: Date): PurchaseWindowState => {
  const wd = ((now.getDay() + 6) % 7) + 1; // 1=Mon...7=Sun
  const isToday = (cfg.purchaseDays || []).includes(wd);
  if (!isToday) {
    return {
      cfg,
      now,
      isToday: false,
      isInWindow: false,
      countdownSeconds: 0,
      reason: '今天不在工作日内（周末禁购）',
    };
  }
  const start = HHMMtoToday(now, cfg.purchaseStart);
  const end = HHMMtoToday(now, cfg.purchaseEnd);
  if (now < start) {
    return {
      cfg,
      now,
      isToday: true,
      isInWindow: false,
      countdownSeconds: Math.floor((start.getTime() - now.getTime()) / 1000),
      reason: `距进货开始 ${cfg.purchaseStart}`,
    };
  }
  if (now >= end) {
    return {
      cfg,
      now,
      isToday: true,
      isInWindow: false,
      countdownSeconds: 0,
      reason: `今日进货已结束（${cfg.purchaseStart}~${cfg.purchaseEnd}）`,
    };
  }
  return {
    cfg,
    now,
    isToday: true,
    isInWindow: true,
    countdownSeconds: Math.floor((end.getTime() - now.getTime()) / 1000),
    reason: '',
  };
};

export function formatCountdown(seconds: number): string {
  if (seconds <= 0) return '00:00';
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = seconds % 60;
  const pad = (n: number) => String(n).padStart(2, '0');
  return h > 0 ? `${pad(h)}:${pad(m)}:${pad(s)}` : `${pad(m)}:${pad(s)}`;
}

/**
 * usePurchaseWindow 监听服务器时间窗 + 倒计时（每秒更新）。
 * 校时：以 serverTimestamp 为基准，与本地时间偏移用于 now 计算。
 */
export function usePurchaseWindow(): PurchaseWindowState {
  const [cfg, setCfg] = useState<PortalBizConfig | null>(null);
  const [tick, setTick] = useState(0);
  const [offsetMs, setOffsetMs] = useState(0);

  useEffect(() => {
    bizConfigApi.get().then((c) => {
      setCfg(c);
      setOffsetMs(c.serverTimestamp * 1000 - Date.now());
    });
    const id = setInterval(() => setTick((n) => n + 1), 1000);
    return () => clearInterval(id);
  }, []);

  return useMemo(() => {
    if (!cfg) {
      return {
        cfg: null,
        now: new Date(),
        isToday: false,
        isInWindow: false,
        countdownSeconds: 0,
        reason: '加载中...',
      };
    }
    return calc(cfg, new Date(Date.now() + offsetMs));
  }, [cfg, offsetMs, tick]);
}

/**
 * useConsignWindow 寄售时间窗（14:30 起，可选结束）。
 */
export function useConsignWindow(): PurchaseWindowState {
  const [cfg, setCfg] = useState<PortalBizConfig | null>(null);
  const [tick, setTick] = useState(0);
  const [offsetMs, setOffsetMs] = useState(0);

  useEffect(() => {
    bizConfigApi.get().then((c) => {
      setCfg(c);
      setOffsetMs(c.serverTimestamp * 1000 - Date.now());
    });
    const id = setInterval(() => setTick((n) => n + 1), 1000);
    return () => clearInterval(id);
  }, []);

  return useMemo(() => {
    if (!cfg) {
      return {
        cfg: null,
        now: new Date(),
        isToday: true,
        isInWindow: false,
        countdownSeconds: 0,
        reason: '加载中...',
      };
    }
    const now = new Date(Date.now() + offsetMs);
    const start = HHMMtoToday(now, cfg.consignStart);
    if (now < start) {
      return {
        cfg,
        now,
        isToday: true,
        isInWindow: false,
        countdownSeconds: Math.floor((start.getTime() - now.getTime()) / 1000),
        reason: `距寄售开放 ${cfg.consignStart}`,
      };
    }
    if (cfg.consignEnd) {
      const end = HHMMtoToday(now, cfg.consignEnd);
      if (now >= end) {
        return {
          cfg,
          now,
          isToday: true,
          isInWindow: false,
          countdownSeconds: 0,
          reason: `今日寄售已结束（${cfg.consignStart}~${cfg.consignEnd}）`,
        };
      }
    }
    return {
      cfg,
      now,
      isToday: true,
      isInWindow: true,
      countdownSeconds: 0,
      reason: '',
    };
  }, [cfg, offsetMs, tick]);
}
