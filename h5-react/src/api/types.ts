// 与 admin-go/app/member/api/portal/v1 对齐的前端类型。

export interface LoginResult {
  token: string;
  memberId: string;
  phone: string;
  nickname: string;
  avatar: string;
  inviteCode: string;
  levelId: string;
  isQualified: number;
}

export interface InvitePreview {
  found: boolean;
  nickname?: string;
  avatar?: string;
}

export interface MeProfile {
  memberId: string;
  phone: string;
  username: string;
  nickname: string;
  avatar: string;
  realName: string;
  inviteCode: string;
  parentId: string;
  levelId: string;
  levelName: string;
  levelExpireAt: string;
  isActive: number;
  isQualified: number;
  teamCount: number;
  directCount: number;
  activeCount: number;
  teamTurnover: number;
  inviteUrl: string;
}

export interface WalletInfo {
  balance: string;
  balanceCent: number;
  totalIncome: string;
  totalExpense: string;
  frozenAmount: string;
}

export interface MyWallets {
  coupon: WalletInfo;
  reward: WalletInfo;
  promote: WalletInfo;
}

export interface WalletLog {
  id: string;
  walletType: number;
  walletTypeText: string;
  changeType: number;
  changeTypeText: string;
  changeAmount: string;
  beforeBalance: string;
  afterBalance: string;
  relatedOrderNo: string;
  remark: string;
  createdAt: string;
}

export interface PageResult<T> {
  total: number;
  list: T[];
}

export interface TeamMember {
  memberId: string;
  nickname: string;
  avatar: string;
  phone: string;
  levelName: string;
  isQualified: number;
  joinedAt: string;
}

export interface MallCategory {
  id: string;
  parentId: string;
  name: string;
  icon: string;
  sort: number;
  children: MallCategory[];
}

export interface MallGoods {
  id: string;
  title: string;
  cover: string;
  price: string;
  originalPrice: string;
  sales: number;
  stock: number;
  isRecommend: number;
}

export interface MallGoodsDetail {
  id: string;
  categoryId: string;
  categoryName: string;
  title: string;
  cover: string;
  images: string[];
  price: string;
  originalPrice: string;
  stock: number;
  sales: number;
  content: string;
  isRecommend: number;
  status: number;
}

export interface MallOrder {
  orderId: string;
  orderNo: string;
  goodsId: string;
  goodsTitle: string;
  goodsCover: string;
  quantity: number;
  totalPrice: string;
  status: number;
  statusText: string;
  remark: string;
  createdAt: string;
}

export interface PlaceMallOrderResult {
  orderId: string;
  orderNo: string;
  totalPrice: string;
}

export interface MyWarehouseGoods {
  id: string;
  goodsNo: string;
  title: string;
  cover: string;
  initPrice: string;
  currentPrice: string;
  nextListingPrice: string;
  priceRiseRate: number;
  platformFeeRate: number;
  tradeCount: number;
  goodsStatus: number;
  goodsStatusText: string;
  activeListingId?: string;
}

export interface WarehouseMarketListing {
  listingId: string;
  goodsId: string;
  goodsNo: string;
  title: string;
  cover: string;
  listingPrice: string;
  sellerName: string;
  tradeCount: number;
  listedAt: string;
}

export interface ListGoodsResult {
  listingId: string;
  listingPrice: string;
}

export interface PlaceTradeResult {
  tradeId: string;
  tradeNo: string;
}

export interface ConfirmTradeResult {
  tradeId: string;
  tradePrice: string;
  platformFee: string;
  sellerIncome: string;
}

export interface TradeRecord {
  tradeId: string;
  tradeNo: string;
  goodsId: string;
  goodsNo: string;
  goodsTitle: string;
  goodsCover: string;
  tradePrice: string;
  platformFee: string;
  sellerIncome: string;
  tradeStatus: number;
  tradeStatusText: string;
  counterparty: string;
  createdAt: string;
  confirmedAt: string;
}

export interface HomeBanner {
  image: string;
  link: string;
  title: string;
}
export interface HomeLevelProgress {
  currentLevelName: string;
  nextLevelName: string;
  needActiveCount: number;
  needTurnover: string;
  isTopLevel: boolean;
}
export interface HomeWalletBriefs {
  coupon: string;
  reward: string;
  promote: string;
}
export interface HomeAggregate {
  banners: HomeBanner[];
  levelProgress: HomeLevelProgress;
  walletBriefs: HomeWalletBriefs;
  recommendedGoods: MallGoods[];
  warehouseListings: WarehouseMarketListing[];
}
