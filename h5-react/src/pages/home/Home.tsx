import { Avatar, Card, ProgressBar, Skeleton, Swiper } from 'antd-mobile';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { homeApi } from '@/api/home';
import { HomeAggregate } from '@/api/types';
import { useAuth } from '@/stores/auth';

export default function Home() {
  const [data, setData] = useState<HomeAggregate | null>(null);
  const [loading, setLoading] = useState(true);
  const user = useAuth((s) => s.user);
  const nav = useNavigate();

  useEffect(() => {
    homeApi
      .get()
      .then(setData)
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="bg-[#f5f5f7] min-h-screen">
      <div className="gradient-primary text-white pt-3 pb-20 px-4">
        <div className="flex items-center gap-3">
          <Avatar
            src={user?.avatar || ''}
            style={{ '--size': '48px', '--border-radius': '24px' } as any}
          />
          <div className="flex-1">
            <div className="text-base font-medium">{user?.nickname || '欢迎'}</div>
            <div className="text-xs opacity-80 mt-1">
              邀请码 {user?.inviteCode || '-'}
              <span className="ml-2 px-2 py-0.5 bg-white/20 rounded">
                {user?.isQualified === 1 ? '资格有效' : '未激活'}
              </span>
            </div>
          </div>
        </div>
      </div>

      <div className="-mt-16 px-3 space-y-3">
        <Card className="rounded-xl shadow-sm">
          {loading || !data ? (
            <Skeleton.Paragraph lineCount={3} animated />
          ) : (
            <div onClick={() => nav('/wallet')}>
              <div className="text-sm text-gray-500 mb-2">我的资产</div>
              <div className="grid grid-cols-3 gap-2 text-center">
                <WalletBlock label="优惠券" value={data.walletBriefs.coupon} color="#ff6a00" />
                <WalletBlock label="奖金" value={data.walletBriefs.reward} color="#52c41a" />
                <WalletBlock label="推广奖" value={data.walletBriefs.promote} color="#1677ff" />
              </div>
            </div>
          )}
        </Card>

        <Card className="rounded-xl shadow-sm">
          {loading || !data ? (
            <Skeleton.Paragraph lineCount={2} animated />
          ) : (
            <div>
              <div className="flex justify-between text-sm mb-2">
                <span className="font-medium">
                  当前等级：
                  <span className="text-primary">{data.levelProgress.currentLevelName || '普通会员'}</span>
                </span>
                {!data.levelProgress.isTopLevel && (
                  <span className="text-gray-500">
                    下一级 {data.levelProgress.nextLevelName}
                  </span>
                )}
              </div>
              {data.levelProgress.isTopLevel ? (
                <div className="text-xs text-gray-500">已达最高等级</div>
              ) : (
                <>
                  <ProgressBar percent={50} style={{ '--fill-color': '#ff6a00' } as any} />
                  <div className="text-xs text-gray-500 mt-1">
                    距下一级：{data.levelProgress.needActiveCount} 活跃 / {data.levelProgress.needTurnover} 元业绩
                  </div>
                </>
              )}
            </div>
          )}
        </Card>

        {data?.banners && data.banners.length > 0 && (
          <Swiper autoplay loop indicatorProps={{ color: 'primary' }}>
            {data.banners.map((b, i) => (
              <Swiper.Item key={i}>
                <div
                  className="rounded-xl overflow-hidden"
                  style={{
                    height: 140,
                    backgroundImage: `url(${b.image})`,
                    backgroundSize: 'cover',
                    backgroundPosition: 'center',
                  }}
                  onClick={() => b.link && (window.location.href = b.link)}
                />
              </Swiper.Item>
            ))}
          </Swiper>
        )}

        <Card title={<div className="font-medium">推荐商品</div>} className="rounded-xl shadow-sm">
          <div className="flex gap-3 overflow-x-auto pb-1">
            {(data?.recommendedGoods || []).map((g) => (
              <div
                key={g.id}
                className="flex-shrink-0 w-32"
                onClick={() => nav(`/mall/detail/${g.id}`)}
              >
                <img src={g.cover} alt={g.title} className="w-full h-32 object-cover rounded-lg" />
                <div className="text-sm mt-1 line-clamp-2 leading-tight">{g.title}</div>
                <div className="text-primary text-sm font-bold mt-1">¥{g.price}</div>
              </div>
            ))}
          </div>
        </Card>

        <Card title={<div className="font-medium">仓库市场</div>} className="rounded-xl shadow-sm">
          {(data?.warehouseListings || []).map((l) => (
            <div
              key={l.listingId}
              className="flex items-center gap-3 py-2 border-b last:border-b-0 border-gray-100"
              onClick={() => nav('/warehouse')}
            >
              <img src={l.cover} alt={l.title} className="w-12 h-12 object-cover rounded" />
              <div className="flex-1">
                <div className="text-sm font-medium">{l.title}</div>
                <div className="text-xs text-gray-500">卖家 {l.sellerName}</div>
              </div>
              <div className="text-primary font-bold">¥{l.listingPrice}</div>
            </div>
          ))}
        </Card>
      </div>
    </div>
  );
}

function WalletBlock({ label, value, color }: { label: string; value: string; color: string }) {
  return (
    <div>
      <div className="text-base font-bold" style={{ color }}>
        ¥{value}
      </div>
      <div className="text-xs text-gray-500 mt-0.5">{label}</div>
    </div>
  );
}

