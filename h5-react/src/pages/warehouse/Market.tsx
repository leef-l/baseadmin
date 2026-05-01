import { Button, Dialog, Dropdown, Empty, InfiniteScroll, SearchBar, Toast } from 'antd-mobile';
import { useEffect, useState } from 'react';
import { warehouseApi } from '@/api/warehouse';
import { WarehouseMarketListing } from '@/api/types';

type OrderBy = 'latest' | 'price_asc' | 'price_desc';

export default function Market() {
  const [keyword, setKeyword] = useState('');
  const [orderBy, setOrderBy] = useState<OrderBy>('latest');
  const [list, setList] = useState<WarehouseMarketListing[]>([]);
  const [pageNum, setPageNum] = useState(1);
  const [hasMore, setHasMore] = useState(true);

  useEffect(() => {
    setList([]);
    setPageNum(1);
    setHasMore(true);
  }, [keyword, orderBy]);

  const loadMore = async () => {
    const res = await warehouseApi.market({
      keyword: keyword || undefined,
      orderBy,
      pageNum,
      pageSize: 20,
    });
    setList((arr) => [...arr, ...res.list]);
    setPageNum((n) => n + 1);
    setHasMore(list.length + res.list.length < res.total);
  };

  const buy = async (l: WarehouseMarketListing) => {
    const ok = await Dialog.confirm({
      content: (
        <div>
          确认购买 <b>{l.title}</b>？
          <br />
          成交价 <span className="text-primary font-bold">¥{l.listingPrice}</span>
          <br />
          <span className="text-xs text-gray-500">
            仓库交易不扣款，由卖家确认后转移所有权
          </span>
        </div>
      ),
    });
    if (!ok) return;
    try {
      const r = await warehouseApi.placeTrade(l.listingId);
      Toast.show({ icon: 'success', content: `已下单 ${r.tradeNo}` });
      setList([]);
      setPageNum(1);
      setHasMore(true);
    } catch {
      // 全局已 toast
    }
  };

  const orderText = { latest: '最新', price_asc: '价格升序', price_desc: '价格降序' }[orderBy];

  return (
    <div className="px-3 pt-3 pb-4">
      <div className="bg-white rounded-xl p-2 flex items-center gap-2 mb-3">
        <SearchBar placeholder="搜索仓库商品" value={keyword} onChange={setKeyword} style={{ flex: 1 }} />
        <Dropdown>
          <Dropdown.Item key="sort" title={orderText}>
            <div className="p-3">
              {(['latest', 'price_asc', 'price_desc'] as OrderBy[]).map((o) => (
                <div
                  key={o}
                  className={`py-2 ${orderBy === o ? 'text-primary font-bold' : ''}`}
                  onClick={() => setOrderBy(o)}
                >
                  {o === 'latest' ? '最新' : o === 'price_asc' ? '价格升序' : '价格降序'}
                </div>
              ))}
            </div>
          </Dropdown.Item>
        </Dropdown>
      </div>

      {!list.length && !hasMore && <Empty description="暂无挂卖" />}

      <div className="grid grid-cols-2 gap-2">
        {list.map((l) => (
          <div key={l.listingId} className="bg-white rounded-xl overflow-hidden shadow-sm">
            <img src={l.cover} className="w-full aspect-square object-cover" />
            <div className="p-2">
              <div className="text-sm line-clamp-2 leading-tight h-10">{l.title}</div>
              <div className="text-xs text-gray-400 mt-1">编号 {l.goodsNo}</div>
              <div className="flex items-end justify-between mt-2">
                <span className="text-primary font-bold">¥{l.listingPrice}</span>
                <Button size="mini" color="primary" onClick={() => buy(l)}>
                  购买
                </Button>
              </div>
            </div>
          </div>
        ))}
      </div>
      <InfiniteScroll loadMore={loadMore} hasMore={hasMore} />
    </div>
  );
}
