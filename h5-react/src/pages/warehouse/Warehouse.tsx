import { Tabs } from 'antd-mobile';
import { useState } from 'react';
import MyWarehouse from './MyWarehouse';
import Market from './Market';
import MyTrades from './MyTrades';
import { formatCountdown, useConsignWindow } from '@/hooks/usePurchaseWindow';

type Key = 'my' | 'market' | 'trades';

export default function Warehouse() {
  const [active, setActive] = useState<Key>('market');
  const window = useConsignWindow();
  return (
    <div className="bg-[#f5f5f7] min-h-screen">
      <div className="bg-white sticky top-0 z-30">
        <div className="px-4 pt-3 text-base font-bold">仓库</div>
        {window.cfg && (
          <div
            className={`mx-3 mb-2 text-xs px-3 py-1.5 rounded-lg ${
              window.isInWindow ? 'bg-green-50 text-green-700' : 'bg-orange-50 text-orange-700'
            }`}
          >
            {window.isInWindow ? (
              <>寄售已开放（{window.cfg.consignStart} 起）</>
            ) : (
              <>
                寄售时间 {window.cfg.consignStart} 起
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
