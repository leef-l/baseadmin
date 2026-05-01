import { Empty, InfiniteScroll, List } from 'antd-mobile';
import { useState } from 'react';
import PageHeader from '@/components/layout/PageHeader';
import { mallApi } from '@/api/mall';
import { MallOrder } from '@/api/types';

export default function MallOrders() {
  const [list, setList] = useState<MallOrder[]>([]);
  const [pageNum, setPageNum] = useState(1);
  const [hasMore, setHasMore] = useState(true);

  const loadMore = async () => {
    const res = await mallApi.myOrders({ pageNum, pageSize: 20 });
    setList((arr) => [...arr, ...res.list]);
    setPageNum((n) => n + 1);
    setHasMore(list.length + res.list.length < res.total);
  };

  return (
    <div className="app-page bg-[#f5f5f7] min-h-screen">
      <PageHeader title="我的订单" />
      {!list.length && !hasMore && <Empty description="暂无订单" />}
      <List>
        {list.map((o) => (
          <List.Item
            key={o.orderId}
            prefix={<img src={o.goodsCover} className="w-12 h-12 rounded object-cover" />}
            description={
              <div className="text-xs text-gray-500">
                {o.orderNo} · {o.createdAt}
              </div>
            }
            extra={
              <div className="text-right">
                <div className="text-primary font-bold">¥{o.totalPrice}</div>
                <div className="text-xs text-gray-500 mt-1">{o.statusText}</div>
              </div>
            }
          >
            {o.goodsTitle} × {o.quantity}
          </List.Item>
        ))}
      </List>
      <InfiniteScroll loadMore={loadMore} hasMore={hasMore} />
    </div>
  );
}
