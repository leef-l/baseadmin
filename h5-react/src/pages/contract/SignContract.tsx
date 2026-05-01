import { Button, Toast } from 'antd-mobile';
import { useEffect, useRef, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import PageHeader from '@/components/layout/PageHeader';
import SignaturePad from '@/components/SignaturePad';
import { contractApi, ContractTemplate, ContractType } from '@/api/contract';

export default function SignContract() {
  const [search] = useSearchParams();
  const contractType = (search.get('type') as ContractType) || 'register';
  const [tpl, setTpl] = useState<ContractTemplate | null>(null);
  const [signature, setSignature] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const sigRef = useRef('');
  const nav = useNavigate();

  useEffect(() => {
    contractApi.template(contractType).then(setTpl);
  }, [contractType]);

  const submit = async () => {
    const sig = sigRef.current || signature;
    if (!sig) {
      Toast.show({ icon: 'fail', content: '请先手写签名' });
      return;
    }
    setSubmitting(true);
    try {
      const r = await contractApi.sign({
        contractType,
        templateId: tpl?.templateId,
        signatureImage: sig,
      });
      Toast.show({ icon: 'success', content: `签署成功 ${r.contractNo}` });
      setTimeout(() => nav('/me/contracts', { replace: true }), 600);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="app-page bg-[#f5f5f7] min-h-screen pb-32">
      <PageHeader title="签署合同" />
      <div className="bg-white p-4 m-3 rounded-xl">
        <div className="text-sm font-medium mb-2">{tpl?.templateName || '合同协议'}</div>
        <div
          className="text-sm leading-relaxed text-gray-700 max-h-[50vh] overflow-y-auto pr-2"
          style={{ wordBreak: 'break-all' }}
          dangerouslySetInnerHTML={{ __html: tpl?.content || '加载中…' }}
        />
      </div>

      <div className="bg-white p-4 m-3 rounded-xl">
        <div className="text-sm font-medium mb-2">乙方签字（手写）</div>
        <SignaturePad
          onChange={(d) => {
            sigRef.current = d;
            setSignature(d);
          }}
        />
        <div className="text-xs text-gray-400 mt-2">提示：请在框内手写您的姓名作为电子签名</div>
      </div>

      <div
        className="fixed bottom-0 left-0 right-0 bg-white p-3"
        style={{
          borderTop: '1px solid #f0f0f0',
          paddingBottom: 'calc(12px + env(safe-area-inset-bottom))',
        }}
      >
        <Button
          block
          color="primary"
          size="large"
          loading={submitting}
          disabled={!signature}
          onClick={submit}
          style={{ borderRadius: 999 }}
        >
          确认签署
        </Button>
      </div>
    </div>
  );
}
