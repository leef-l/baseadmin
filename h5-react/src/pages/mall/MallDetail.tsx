import { Button, Dialog, ImageViewer, Stepper, Swiper, Toast } from 'antd-mobile';
import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import PageHeader from '@/components/layout/PageHeader';
import { mallApi } from '@/api/mall';
import { meApi } from '@/api/me';
import { MallGoodsDetail, MeProfile } from '@/api/types';
import { formatCountdown, usePurchaseWindow } from '@/hooks/usePurchaseWindow';

export default function MallDetail() {
  const { id } = useParams();
  const [detail, setDetail] = useState<MallGoodsDetail | null>(null);
  const [profile, setProfile] = useState<MeProfile | null>(null);
  const [quantity, setQuantity] = useState(1);
  const [submitting, setSubmitting] = useState(false);
  const nav = useNavigate();
  const window = usePurchaseWindow();

  useEffect(() => {
    if (!id) return;
    mallApi.detail(id).then(setDetail);
    meApi.profile().then(setProfile).catch(() => {});
  }, [id]);

  if (!detail) {
    return (
      <div className="app-page bg-white min-h-screen">
        <PageHeader title="商品详情" />
        <div className="text-center text-gray-400 mt-20">加载中…</div>
      </div>
    );
  }

  const images = detail.images?.length ? detail.images : [detail.cover];

  const remaining = profile
    ? Math.max(0, profile.dailyPurchaseLimit - profile.todayPurchaseCount)
    : null;

  const canBuy = window.isInWindow && (remaining === null || remaining > 0);
  const blockReason = !window.isInWindow
    ? window.reason
    : remaining === 0
    ? `今日限购已用完（${profile?.todayPurchaseCount}/${profile?.dailyPurchaseLimit}）`
    : '';

  const placeOrder = async () => {
    if (!canBuy) {
      Toast.show({ icon: 'fail', content: blockReason });
      return;
    }
    const ok = await Dialog.confirm({
      content: (
        <div>
          确定购买 <b>{detail.title}</b> × {quantity}？将从优惠券钱包扣款 ¥
          {(parseFloat(detail.price) * quantity).toFixed(2)}
        </div>
      ),
    });
    if (!ok) return;
    setSubmitting(true);
    try {
      const r = await mallApi.placeOrder(detail.id, quantity);
      Toast.show({ icon: 'success', content: `下单成功 ${r.orderNo}` });
      nav('/mall/orders');
    } finally {
      setSubmitting(false);
    }
  };

  const showPics = (idx: number) => {
    ImageViewer.Multi.show({ images, defaultIndex: idx });
  };

  return (
    <div className="app-page bg-[#f5f5f7] min-h-screen pb-24">
      <PageHeader title="商品详情" />
      <Swiper autoplay loop indicatorProps={{ color: 'white' }}>
        {images.map((img, i) => (
          <Swiper.Item key={i}>
            <img
              src={img}
              alt=""
              className="w-full aspect-square object-cover"
              onClick={() => showPics(i)}
            />
          </Swiper.Item>
        ))}
      </Swiper>

      <div className="bg-white p-4">
        <div className="flex items-end gap-2">
          <span className="text-2xl font-bold text-primary">¥{detail.price}</span>
          {detail.originalPrice && parseFloat(detail.originalPrice) > 0 && (
            <span className="text-xs text-gray-400 line-through">¥{detail.originalPrice}</span>
          )}
        </div>
        <div className="text-base font-medium mt-2">{detail.title}</div>
        <div className="flex justify-between text-xs text-gray-500 mt-2">
          <span>分类：{detail.categoryName}</span>
          <span>已售 {detail.sales} · 库存 {detail.stock}</span>
        </div>
      </div>

      <div className="bg-white p-4 mt-2 text-sm">
        {window.cfg && (
          <div className="flex items-center justify-between">
            <span>进货时间</span>
            <span className={window.isInWindow ? 'text-green-600 font-medium' : 'text-gray-500'}>
              {window.cfg.purchaseStart} ~ {window.cfg.purchaseEnd}
              {window.isInWindow && (
                <span className="ml-2">距结束 {formatCountdown(window.countdownSeconds)}</span>
              )}
              {!window.isInWindow && window.countdownSeconds > 0 && (
                <span className="ml-2 text-primary">{formatCountdown(window.countdownSeconds)} 后开放</span>
              )}
            </span>
          </div>
        )}
        {profile && (
          <div className="flex items-center justify-between mt-2">
            <span>今日剩余</span>
            <span className="text-primary font-bold">
              {remaining}/{profile.dailyPurchaseLimit} 单
            </span>
          </div>
        )}
      </div>

      <div className="bg-white p-4 mt-2">
        <div className="flex items-center justify-between">
          <span className="text-sm">购买数量</span>
          <Stepper min={1} max={Math.max(1, detail.stock)} value={quantity} onChange={setQuantity} />
        </div>
      </div>

      <div className="bg-white p-4 mt-2">
        <div className="text-sm font-medium mb-2">商品详情</div>
        <div className="text-sm leading-relaxed" dangerouslySetInnerHTML={{ __html: detail.content || '' }} />
      </div>

      <div
        className="fixed bottom-0 left-0 right-0 bg-white p-3 flex flex-col gap-1"
        style={{
          borderTop: '1px solid #f0f0f0',
          paddingBottom: 'calc(12px + env(safe-area-inset-bottom))',
        }}
      >
        {!canBuy && blockReason && (
          <div className="text-center text-xs text-gray-500 pb-1">{blockReason}</div>
        )}
        <Button
          color="primary"
          fill="solid"
          loading={submitting}
          disabled={!canBuy}
          onClick={placeOrder}
          style={{ flex: 1, borderRadius: 999 }}
        >
          {canBuy ? '立即购买（优惠券支付）' : '当前不可购买'}
        </Button>
      </div>
    </div>
  );
}
