import { Button, Dialog, Empty, InfiniteScroll, Tabs, Tag, Toast } from 'antd-mobile';
import { useEffect, useState } from 'react';
import { warehouseApi } from '@/api/warehouse';
import { TradeRecord } from '@/api/types';

type Role = 'buyer' | 'seller';

export default function MyTrades() {
  const [role, setRole] = useState<Role>('buyer');
  const [list, setList] = useState<TradeRecord[]>([]);
  const [pageNum, setPageNum] = useState(1);
  const [hasMore, setHasMore] = useState(true);

  useEffect(() => {
    setList([]);
    setPageNum(1);
    setHasMore(true);
  }, [role]);

  const loadMore = async () => {
    const res = await warehouseApi.myTrades({ role, status: 0, pageNum, pageSize: 20 });
    setList((arr) => [...arr, ...res.list]);
    setPageNum((n) => n + 1);
    setHasMore(list.length + res.list.length < res.total);
  };

  const confirmTrade = async (t: TradeRecord) => {
    const ok = await Dialog.confirm({
      content: (
        <div>
          确认交易 <b>{t.goodsTitle}</b>？
          <br />
          成交价 ¥{t.tradePrice}（确认后增值差额扣除平台抽成入账到你的奖金钱包）
        </div>
      ),
    });
    if (!ok) return;
    try {
      const r = await warehouseApi.confirmTrade(t.tradeId);
      Toast.show({ content: `已确认，奖金 +¥${r.sellerIncome}`, icon: 'success' });
      setList([]);
      setPageNum(1);
      setHasMore(true);
    } catch {
      //
    }
  };

  return (
    <>
      <Tabs activeKey={role} onChange={(k) => setRole(k as Role)}>
        <Tabs.Tab title="买入" key="buyer" />
        <Tabs.Tab title="卖出" key="seller" />
      </Tabs>

      <div className="px-3 pt-2 pb-4 space-y-2">
        {!list.length && !hasMore && <Empty description="暂无交易记录" />}
        {list.map((t) => (
          <div key={t.tradeId} className="bg-white rounded-xl p-3 shadow-sm">
            <div className="flex gap-3">
              <img src={t.goodsCover} className="w-16 h-16 rounded object-cover" />
              <div className="flex-1">
                <div className="flex justify-between">
                  <span className="text-sm font-medium line-clamp-1">{t.goodsTitle}</span>
                  <Tag
                    color={
                      t.tradeStatus === 1 ? 'warning' : t.tradeStatus === 2 ? 'success' : 'default'
                    }
                  >
                    {t.tradeStatusText}
                  </Tag>
                </div>
                <div className="text-xs text-gray-500 mt-1">单号 {t.tradeNo}</div>
                <div className="text-xs text-gray-500">{t.createdAt}</div>
                <div className="text-xs text-gray-500">对方：{t.counterparty}</div>
              </div>
            </div>
            <div className="flex justify-between items-end mt-2 pt-2 border-t border-gray-100">
              <div>
                <div className="text-xs text-gray-500">成交价</div>
                <div className="text-primary font-bold">¥{t.tradePrice}</div>
              </div>
              {role === 'seller' && (
                <div className="text-xs text-gray-500 text-right">
                  <div>抽成 ¥{t.platformFee}</div>
                  <div>到账 ¥{t.sellerIncome}</div>
                </div>
              )}
              {role === 'seller' && t.tradeStatus === 1 && (
                <Button size="small" color="primary" onClick={() => confirmTrade(t)}>
                  确认成交
                </Button>
              )}
            </div>
          </div>
        ))}
        <InfiniteScroll loadMore={loadMore} hasMore={hasMore} />
      </div>
    </>
  );
}
