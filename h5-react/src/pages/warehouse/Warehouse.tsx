import { Tabs } from 'antd-mobile';
import { useState } from 'react';
import MyWarehouse from './MyWarehouse';
import Market from './Market';
import MyTrades from './MyTrades';

type Key = 'my' | 'market' | 'trades';

export default function Warehouse() {
  const [active, setActive] = useState<Key>('market');
  return (
    <div className="bg-[#f5f5f7] min-h-screen">
      <div className="bg-white sticky top-0 z-30">
        <div className="px-4 pt-3 text-base font-bold">仓库</div>
        <Tabs activeKey={active} onChange={(k) => setActive(k as Key)}>
          <Tabs.Tab title="市场" key="market" />
          <Tabs.Tab title="我的库存" key="my" />
          <Tabs.Tab title="我的交易" key="trades" />
        </Tabs>
      </div>
      {active === 'market' && <Market />}
      {active === 'my' && <MyWarehouse />}
      {active === 'trades' && <MyTrades />}
    </div>
  );
}
