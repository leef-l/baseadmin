import { Button, Dialog, ImageViewer, Stepper, Swiper, Toast } from 'antd-mobile';
import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import PageHeader from '@/components/layout/PageHeader';
import { mallApi } from '@/api/mall';
import { MallGoodsDetail } from '@/api/types';

export default function MallDetail() {
  const { id } = useParams();
  const [detail, setDetail] = useState<MallGoodsDetail | null>(null);
  const [quantity, setQuantity] = useState(1);
  const [submitting, setSubmitting] = useState(false);
  const nav = useNavigate();

  useEffect(() => {
    if (!id) return;
    mallApi.detail(id).then(setDetail);
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

  const placeOrder = async () => {
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
        className="fixed bottom-0 left-0 right-0 bg-white p-3 flex gap-2"
        style={{
          borderTop: '1px solid #f0f0f0',
          paddingBottom: 'calc(12px + env(safe-area-inset-bottom))',
        }}
      >
        <Button
          color="primary"
          fill="solid"
          loading={submitting}
          onClick={placeOrder}
          style={{ flex: 1, borderRadius: 999 }}
        >
          立即购买（优惠券支付）
        </Button>
      </div>
    </div>
  );
}
