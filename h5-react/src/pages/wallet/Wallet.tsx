import { Card, InfiniteScroll, List, Tabs } from 'antd-mobile';
import { useEffect, useState } from 'react';
import PageHeader from '@/components/layout/PageHeader';
import { meApi } from '@/api/me';
import { MyWallets, WalletInfo, WalletLog } from '@/api/types';

const TYPE_OPTIONS = [
  { type: 1, name: '优惠券', color: '#ff6a00' },
  { type: 2, name: '奖金', color: '#52c41a' },
  { type: 3, name: '推广奖', color: '#1677ff' },
];

export default function Wallet() {
  const [wallets, setWallets] = useState<MyWallets | null>(null);
  const [activeType, setActiveType] = useState<number>(1);
  const [logs, setLogs] = useState<WalletLog[]>([]);
  const [pageNum, setPageNum] = useState(1);
  const [hasMore, setHasMore] = useState(true);

  useEffect(() => {
    meApi.wallets().then(setWallets);
  }, []);

  useEffect(() => {
    setLogs([]);
    setPageNum(1);
    setHasMore(true);
  }, [activeType]);

  const loadMore = async () => {
    const res = await meApi.walletLogs({ walletType: activeType, pageNum, pageSize: 20 });
    setLogs((arr) => [...arr, ...res.list]);
    setPageNum((n) => n + 1);
    setHasMore(logs.length + res.list.length < res.total);
  };

  const currentWallet: WalletInfo | undefined =
    wallets && (activeType === 1 ? wallets.coupon : activeType === 2 ? wallets.reward : wallets.promote);
  const themeColor = TYPE_OPTIONS.find((t) => t.type === activeType)?.color || '#ff6a00';

  return (
    <div className="app-page bg-[#f5f5f7] min-h-screen">
      <PageHeader title="我的钱包" />

      <Card
        className="m-3 rounded-xl text-white"
        style={{ background: `linear-gradient(135deg, ${themeColor}cc, ${themeColor})` }}
      >
        <div className="text-xs opacity-80">
          {TYPE_OPTIONS.find((t) => t.type === activeType)?.name}余额（元）
        </div>
        <div className="text-3xl font-bold mt-2">¥{currentWallet?.balance ?? '0.00'}</div>
        <div className="grid grid-cols-3 gap-2 mt-4 text-center text-xs">
          <Cell label="累计收入" value={currentWallet?.totalIncome} />
          <Cell label="累计支出" value={currentWallet?.totalExpense} />
          <Cell label="冻结" value={currentWallet?.frozenAmount} />
        </div>
      </Card>

      <Tabs
        activeKey={String(activeType)}
        onChange={(k) => setActiveType(Number(k))}
        style={{ '--active-line-color': themeColor, '--active-title-color': themeColor } as any}
      >
        {TYPE_OPTIONS.map((t) => (
          <Tabs.Tab key={t.type} title={t.name} />
        ))}
      </Tabs>

      <List header="流水明细" style={{ '--padding-left': '12px' } as any}>
        {logs.map((log) => (
          <List.Item
            key={log.id}
            description={
              <div className="text-xs text-gray-500">
                {log.changeTypeText} · {log.createdAt}
                {log.relatedOrderNo && <span className="ml-2">订单 {log.relatedOrderNo}</span>}
              </div>
            }
            extra={
              <span
                style={{
                  color: log.changeAmount.startsWith('-') ? '#888' : themeColor,
                  fontWeight: 600,
                }}
              >
                {log.changeAmount}
              </span>
            }
          >
            {log.remark || log.changeTypeText}
          </List.Item>
        ))}
      </List>
      <InfiniteScroll loadMore={loadMore} hasMore={hasMore} />
    </div>
  );
}

function Cell({ label, value }: { label: string; value: string | undefined }) {
  return (
    <div>
      <div className="font-bold">¥{value ?? '0.00'}</div>
      <div className="opacity-80 mt-0.5">{label}</div>
    </div>
  );
}
