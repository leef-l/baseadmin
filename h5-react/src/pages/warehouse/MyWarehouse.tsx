import { Button, Dialog, Empty, InfiniteScroll, Tag, Toast } from 'antd-mobile';
import { useEffect, useState } from 'react';
import { warehouseApi } from '@/api/warehouse';
import { MyWarehouseGoods } from '@/api/types';

export default function MyWarehouse() {
  const [list, setList] = useState<MyWarehouseGoods[]>([]);
  const [pageNum, setPageNum] = useState(1);
  const [hasMore, setHasMore] = useState(true);

  const reload = () => {
    setList([]);
    setPageNum(1);
    setHasMore(true);
  };

  const loadMore = async () => {
    const res = await warehouseApi.my({ status: 0, pageNum, pageSize: 20 });
    setList((arr) => [...arr, ...res.list]);
    setPageNum((n) => n + 1);
    setHasMore(list.length + res.list.length < res.total);
  };

  const onList = async (g: MyWarehouseGoods) => {
    const ok = await Dialog.confirm({
      content: (
        <div>
          确认挂卖 <b>{g.title}</b>？
          <br />
          挂卖价由系统自动计算：
          <span className="text-primary font-bold">¥{g.nextListingPrice}</span>
          <br />
          <span className="text-xs text-gray-500">
            （加价 {g.priceRiseRate}% / 平台抽成 {g.platformFeeRate}%）
          </span>
        </div>
      ),
    });
    if (!ok) return;
    try {
      const r = await warehouseApi.list(g.id);
      Toast.show({ icon: 'success', content: `已挂卖 ¥${r.listingPrice}` });
      reload();
    } catch {
      //
    }
  };

  const statusColor: Record<number, string> = { 1: 'success', 2: 'primary', 3: 'warning' };

  return (
    <div className="px-3 pt-3 pb-4">
      {!list.length && !hasMore && <Empty description="暂无库存商品" />}
      <div className="space-y-2">
        {list.map((g) => (
          <div key={g.id} className="bg-white rounded-xl p-3 flex gap-3 shadow-sm">
            <img src={g.cover} className="w-20 h-20 object-cover rounded" />
            <div className="flex-1">
              <div className="flex justify-between items-start">
                <div className="text-sm font-medium line-clamp-1 flex-1">{g.title}</div>
                <Tag color={statusColor[g.goodsStatus] || 'default'} className="ml-2">
                  {g.goodsStatusText}
                </Tag>
              </div>
              <div className="text-xs text-gray-500 mt-1">编号 {g.goodsNo} · 已成交 {g.tradeCount}</div>
              <div className="flex justify-between items-end mt-2">
                <div>
                  <div className="text-xs text-gray-500">当前价 ¥{g.currentPrice}</div>
                  <div className="text-primary font-bold text-sm">下次挂卖 ¥{g.nextListingPrice}</div>
                </div>
                {g.goodsStatus === 1 && (
                  <Button size="small" color="primary" onClick={() => onList(g)}>
                    一键挂卖
                  </Button>
                )}
              </div>
            </div>
          </div>
        ))}
      </div>
      <InfiniteScroll loadMore={loadMore} hasMore={hasMore} />
    </div>
  );
}
