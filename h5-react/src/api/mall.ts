import req from './request';
import {
  MallCategory,
  MallGoods,
  MallGoodsDetail,
  MallOrder,
  PageResult,
  PlaceMallOrderResult,
} from './types';

export const mallApi = {
  categories: () =>
    req.get<any, { list: MallCategory[] }>('/member-portal/mall/categories'),

  goods: (params: {
    categoryId?: string;
    keyword?: string;
    isRecommend?: number;
    pageNum?: number;
    pageSize?: number;
  }) => req.get<any, PageResult<MallGoods>>('/member-portal/mall/goods', { params }),

  detail: (id: string) =>
    req.get<any, MallGoodsDetail>('/member-portal/mall/goods/detail', { params: { id } }),

  placeOrder: (goodsId: string, quantity: number, remark?: string) =>
    req.post<any, PlaceMallOrderResult>('/member-portal/mall/order/place', {
      goodsId,
      quantity,
      remark,
    }),

  myOrders: (params: { pageNum?: number; pageSize?: number }) =>
    req.get<any, PageResult<MallOrder>>('/member-portal/mall/orders', { params }),
};
