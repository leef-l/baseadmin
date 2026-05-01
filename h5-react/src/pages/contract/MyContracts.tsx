import { Button, Empty, InfiniteScroll, List, Tag } from 'antd-mobile';
import { useState } from 'react';
import PageHeader from '@/components/layout/PageHeader';
import { contractApi, ContractItem } from '@/api/contract';

export default function MyContracts() {
  const [list, setList] = useState<ContractItem[]>([]);
  const [pageNum, setPageNum] = useState(1);
  const [hasMore, setHasMore] = useState(true);

  const loadMore = async () => {
    const res = await contractApi.list({ pageNum, pageSize: 20 });
    setList((arr) => [...arr, ...res.list]);
    setPageNum((n) => n + 1);
    setHasMore(list.length + res.list.length < res.total);
  };

  const download = (c: ContractItem) => {
    window.open(contractApi.downloadURL(c.contractId), '_blank');
  };

  return (
    <div className="app-page bg-[#f5f5f7] min-h-screen">
      <PageHeader title="我的合同" />
      {!list.length && !hasMore && <Empty description="暂无合同记录" />}
      <List>
        {list.map((c) => (
          <List.Item
            key={c.contractId}
            description={
              <div className="text-xs text-gray-500">
                {c.contractNo} · {c.signedAt}
              </div>
            }
            extra={
              <div className="flex items-center gap-2">
                <Tag color={c.pdfStatus === 2 ? 'success' : c.pdfStatus === 3 ? 'danger' : 'warning'}>
                  {c.pdfStatusText}
                </Tag>
                <Button size="small" color="primary" fill="outline" onClick={() => download(c)}>
                  查看 / 下载
                </Button>
              </div>
            }
          >
            {c.contractTypeText}
          </List.Item>
        ))}
      </List>
      <InfiniteScroll loadMore={loadMore} hasMore={hasMore} />
    </div>
  );
}
