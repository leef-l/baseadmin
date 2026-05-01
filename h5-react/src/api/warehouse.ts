import req from './request';
import {
  ConfirmTradeResult,
  ListGoodsResult,
  MyWarehouseGoods,
  PageResult,
  PlaceTradeResult,
  TradeRecord,
  WarehouseMarketListing,
} from './types';

export const warehouseApi = {
  my: (params: { status?: number; pageNum?: number; pageSize?: number }) =>
    req.get<any, PageResult<MyWarehouseGoods>>('/member-portal/warehouse/my', { params }),

  market: (params: {
    keyword?: string;
    orderBy?: 'price_asc' | 'price_desc' | 'latest';
    pageNum?: number;
    pageSize?: number;
  }) => req.get<any, PageResult<WarehouseMarketListing>>('/member-portal/warehouse/market', { params }),

  list: (goodsId: string) =>
    req.post<any, ListGoodsResult>('/member-portal/warehouse/list', { goodsId }),

  placeTrade: (listingId: string) =>
    req.post<any, PlaceTradeResult>('/member-portal/warehouse/trade/place', { listingId }),

  confirmTrade: (tradeId: string) =>
    req.post<any, ConfirmTradeResult>('/member-portal/warehouse/trade/confirm', { tradeId }),

  myTrades: (params: {
    role?: 'buyer' | 'seller';
    status?: number;
    pageNum?: number;
    pageSize?: number;
  }) => req.get<any, PageResult<TradeRecord>>('/member-portal/warehouse/my-trades', { params }),
};
