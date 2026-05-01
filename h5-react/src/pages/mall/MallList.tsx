import { InfiniteScroll, SearchBar } from 'antd-mobile';
import { useEffect, useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { mallApi } from '@/api/mall';
import { MallCategory, MallGoods } from '@/api/types';
import { formatCountdown, usePurchaseWindow } from '@/hooks/usePurchaseWindow';

export default function MallList() {
  const [categories, setCategories] = useState<MallCategory[]>([]);
  const [activeCat, setActiveCat] = useState<string>('');
  const [keyword, setKeyword] = useState('');
  const [goods, setGoods] = useState<MallGoods[]>([]);
  const [pageNum, setPageNum] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const nav = useNavigate();
  const window = usePurchaseWindow();

  useEffect(() => {
    mallApi.categories().then((r) => setCategories(r.list || []));
  }, []);

  useEffect(() => {
    setGoods([]);
    setPageNum(1);
    setHasMore(true);
  }, [activeCat, keyword]);

  const loadMore = async () => {
    const res = await mallApi.goods({
      categoryId: activeCat || undefined,
      keyword: keyword || undefined,
      pageNum,
      pageSize: 20,
    });
    setGoods((arr) => [...arr, ...res.list]);
    setPageNum((n) => n + 1);
    setHasMore(goods.length + res.list.length < res.total);
  };

  const flatCats = useMemo(() => {
    const out: MallCategory[] = [];
    const walk = (arr: MallCategory[]) =>
      arr.forEach((c) => {
        out.push(c);
        if (c.children?.length) walk(c.children);
      });
    walk(categories);
    return out;
  }, [categories]);

  return (
    <div className="bg-[#f5f5f7] min-h-screen">
      <div className="bg-white px-3 pt-3 pb-2 sticky top-0 z-30">
        <SearchBar placeholder="搜索商品" value={keyword} onChange={setKeyword} />
        {window.cfg && (
          <div
            className={`mt-2 text-xs px-3 py-1.5 rounded-lg ${
              window.isInWindow ? 'bg-green-50 text-green-700' : 'bg-orange-50 text-orange-700'
            }`}
          >
            {window.isInWindow ? (
              <>
                进货中（{window.cfg.purchaseStart}~{window.cfg.purchaseEnd}） · 距结束{' '}
                <b>{formatCountdown(window.countdownSeconds)}</b>
              </>
            ) : (
              <>
                进货时间 {window.cfg.purchaseStart}~{window.cfg.purchaseEnd}
                {window.countdownSeconds > 0 && (
                  <>
                    {' '}
                    · <b>{formatCountdown(window.countdownSeconds)}</b> 后开放
                  </>
                )}
                {!window.countdownSeconds && <> · {window.reason}</>}
              </>
            )}
          </div>
        )}
      </div>

      <div className="flex">
        <div
          className="bg-white border-r border-gray-100 overflow-y-auto"
          style={{ width: 88, height: 'calc(100vh - 60px - 60px)' }}
        >
          <div
            className={`py-3 text-center text-sm cursor-pointer ${
              activeCat === '' ? 'text-primary font-medium border-l-2 border-primary bg-orange-50' : ''
            }`}
            onClick={() => setActiveCat('')}
          >
            全部
          </div>
          {flatCats.map((c) => (
            <div
              key={c.id}
              className={`py-3 text-center text-sm cursor-pointer ${
                activeCat === c.id
                  ? 'text-primary font-medium border-l-2 border-primary bg-orange-50'
                  : ''
              }`}
              onClick={() => setActiveCat(c.id)}
            >
              {c.name}
            </div>
          ))}
        </div>

        <div className="flex-1 overflow-y-auto px-2 py-2">
          <div className="grid grid-cols-2 gap-2">
            {goods.map((g) => (
              <div
                key={g.id}
                className="bg-white rounded-xl overflow-hidden shadow-sm"
                onClick={() => nav(`/mall/detail/${g.id}`)}
              >
                <img src={g.cover} alt={g.title} className="w-full aspect-square object-cover" />
                <div className="p-2">
                  <div className="text-sm line-clamp-2 leading-tight h-10">{g.title}</div>
                  <div className="flex items-end justify-between mt-2">
                    <span className="text-primary font-bold">¥{g.price}</span>
                    <span className="text-xs text-gray-400">已售 {g.sales}</span>
                  </div>
                </div>
              </div>
            ))}
          </div>
          <InfiniteScroll loadMore={loadMore} hasMore={hasMore} />
        </div>
      </div>
    </div>
  );
}
