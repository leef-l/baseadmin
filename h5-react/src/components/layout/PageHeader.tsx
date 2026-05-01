import { LeftOutline } from 'antd-mobile-icons';
import { ReactNode } from 'react';
import { useNavigate } from 'react-router-dom';

interface Props {
  title: string;
  right?: ReactNode;
  onBack?: () => void;
  transparent?: boolean;
}

export default function PageHeader({ title, right, onBack, transparent }: Props) {
  const nav = useNavigate();
  return (
    <div
      className={`sticky top-0 z-50 flex items-center h-12 px-3 ${
        transparent ? 'bg-transparent' : 'bg-white border-b border-[#f0f0f0]'
      }`}
    >
      <div
        className="w-8 h-8 flex items-center justify-center"
        onClick={() => (onBack ? onBack() : nav(-1))}
      >
        <LeftOutline fontSize={20} />
      </div>
      <div className="flex-1 text-center text-base font-medium truncate px-2">{title}</div>
      <div className="w-8 h-8 flex items-center justify-center">{right}</div>
    </div>
  );
}
