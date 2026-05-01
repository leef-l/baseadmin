import { Avatar, Card, InfiniteScroll, List, Tabs } from 'antd-mobile';
import { useEffect, useState } from 'react';
import PageHeader from '@/components/layout/PageHeader';
import { meApi } from '@/api/me';
import { MeProfile, TeamMember } from '@/api/types';

type Scope = 'direct' | 'all';

export default function Team() {
  const [profile, setProfile] = useState<MeProfile | null>(null);
  const [scope, setScope] = useState<Scope>('direct');
  const [list, setList] = useState<TeamMember[]>([]);
  const [pageNum, setPageNum] = useState(1);
  const [hasMore, setHasMore] = useState(true);

  useEffect(() => {
    meApi.profile().then(setProfile);
  }, []);

  useEffect(() => {
    setList([]);
    setPageNum(1);
    setHasMore(true);
  }, [scope]);

  const loadMore = async () => {
    const res = await meApi.team({ scope, pageNum, pageSize: 20 });
    setList((arr) => [...arr, ...res.list]);
    setPageNum((n) => n + 1);
    setHasMore(list.length + res.list.length < res.total);
  };

  return (
    <div className="app-page bg-[#f5f5f7] min-h-screen">
      <PageHeader title="我的团队" />

      <Card className="m-3 rounded-xl gradient-primary text-white">
        <div className="grid grid-cols-3 gap-2 text-center">
          <Stat label="团队总数" value={profile?.teamCount} />
          <Stat label="活跃用户" value={profile?.activeCount} />
          <Stat label="直推人数" value={profile?.directCount} />
        </div>
        <div className="text-center mt-4 text-xs opacity-80">
          团队总业绩：¥{((profile?.teamTurnover || 0) / 100).toFixed(2)}
        </div>
      </Card>

      <Tabs activeKey={scope} onChange={(k) => setScope(k as Scope)}>
        <Tabs.Tab title="直推" key="direct" />
        <Tabs.Tab title="全部团队" key="all" />
      </Tabs>

      <List>
        {list.map((m) => (
          <List.Item
            key={m.memberId}
            prefix={<Avatar src={m.avatar || ''} style={{ '--size': '40px' } as any} />}
            description={
              <div className="text-xs text-gray-500">
                {m.phone} · {m.levelName || '普通会员'} · 加入 {m.joinedAt}
              </div>
            }
            extra={
              <span
                className={`text-xs px-2 py-0.5 rounded ${
                  m.isQualified === 1 ? 'bg-orange-100 text-primary' : 'bg-gray-100 text-gray-500'
                }`}
              >
                {m.isQualified === 1 ? '有效' : '未激活'}
              </span>
            }
          >
            {m.nickname || '会员'}
          </List.Item>
        ))}
      </List>
      <InfiniteScroll loadMore={loadMore} hasMore={hasMore} />
    </div>
  );
}

function Stat({ label, value }: { label: string; value: number | undefined }) {
  return (
    <div>
      <div className="text-2xl font-bold">{value ?? 0}</div>
      <div className="text-xs opacity-80 mt-1">{label}</div>
    </div>
  );
}
