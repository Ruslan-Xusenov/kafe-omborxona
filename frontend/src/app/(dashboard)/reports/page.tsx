'use client';
import { useEffect, useState } from 'react';
import { api } from '@/lib/api';
import { useAuth } from '@/lib/auth';
import { ProfitReport, TopProduct, formatMoney, formatNumber } from '@/lib/types';
import { useRouter } from 'next/navigation';

export default function ReportsPage() {
  const { isAdmin, loading: authLoading } = useAuth();
  const router = useRouter();
  const [profit, setProfit] = useState<ProfitReport | null>(null);
  const [topProducts, setTopProducts] = useState<TopProduct[]>([]);
  const [loading, setLoading] = useState(true);
  const [dateFrom, setDateFrom] = useState('');
  const [dateTo, setDateTo] = useState('');

  useEffect(() => {
    if (!authLoading && !isAdmin) { router.push('/'); return; }
  }, [isAdmin, authLoading, router]);

  const load = () => {
    setLoading(true);
    const params = new URLSearchParams();
    if (dateFrom) params.set('date_from', dateFrom);
    if (dateTo) params.set('date_to', dateTo);

    Promise.all([
      api.get<ProfitReport>(`/api/dashboard/profit?${params}`),
      api.get<TopProduct[]>(`/api/dashboard/top-products?limit=10&${params}`),
    ]).then(([p, t]) => { setProfit(p); setTopProducts(t || []); })
      .catch(console.error).finally(() => setLoading(false));
  };

  useEffect(() => { if (isAdmin) load(); }, [isAdmin]);

  if (!isAdmin) return null;

  return (
    <>
      <div className="topbar">
        <h1>💰 Moliyaviy Hisobot</h1>
      </div>
      <div className="page-content fade-in">
        <div className="filter-bar">
          <div className="input-group">
            <label>Boshlanish sanasi</label>
            <input className="input" type="date" value={dateFrom} onChange={e => setDateFrom(e.target.value)} />
          </div>
          <div className="input-group">
            <label>Tugash sanasi</label>
            <input className="input" type="date" value={dateTo} onChange={e => setDateTo(e.target.value)} />
          </div>
          <button className="btn btn-primary btn-sm" onClick={load} style={{ marginTop: 'auto' }}>📊 Hisobni ko&apos;rish</button>
        </div>

        {loading ? <div className="loading-page"><div className="spinner" /></div> : (
          <>
            {profit && (
              <div className={`profit-card ${profit.net_profit < 0 ? 'negative' : ''}`} style={{ marginBottom: 32 }}>
                <div style={{ textAlign: 'center', marginBottom: 24 }}>
                  <div style={{ color: 'var(--text-muted)', fontSize: '0.9rem', marginBottom: 8 }}>SOF FOYDA</div>
                  <div className="profit-value" style={{ fontSize: '3rem' }}>{formatMoney(profit.net_profit)}</div>
                </div>
                <div style={{ display: 'flex', justifyContent: 'center', gap: 48, flexWrap: 'wrap', marginTop: 24 }}>
                  <div style={{ textAlign: 'center' }}>
                    <div style={{ fontSize: '1.5rem', fontWeight: 700, color: 'var(--accent-green)' }}>{formatMoney(profit.total_sales_revenue)}</div>
                    <div style={{ color: 'var(--text-muted)', fontSize: '0.8rem', marginTop: 4 }}>Sotuv daromadi</div>
                  </div>
                  <div style={{ textAlign: 'center' }}>
                    <div style={{ fontSize: '1.5rem', fontWeight: 700, color: 'var(--accent-yellow)' }}>{formatMoney(profit.total_cost_of_sold)}</div>
                    <div style={{ color: 'var(--text-muted)', fontSize: '0.8rem', marginTop: 4 }}>Sotilgan tan narxi</div>
                  </div>
                  <div style={{ textAlign: 'center' }}>
                    <div style={{ fontSize: '1.5rem', fontWeight: 700, color: 'var(--accent-red)' }}>{formatMoney(profit.total_write_off_loss)}</div>
                    <div style={{ color: 'var(--text-muted)', fontSize: '0.8rem', marginTop: 4 }}>Yo&apos;qotishlar</div>
                  </div>
                </div>
                <div style={{ textAlign: 'center', marginTop: 24, padding: '16px', background: 'rgba(0,0,0,0.2)', borderRadius: 'var(--radius-md)', fontSize: '0.85rem', color: 'var(--text-secondary)' }}>
                  <strong>Formula:</strong> Foyda = Sotuv daromadi − (Sotilgan mahsulotlar tan narxi + Yo&apos;qotishlar)
                </div>
              </div>
            )}

            {topProducts && topProducts.length > 0 && (
              <div className="table-container">
                <div className="table-header">
                  <h2>🏆 Eng ko&apos;p sotilgan mahsulotlar (Top 10)</h2>
                </div>
                <div className="table-scroll">
                  <table>
                    <thead><tr><th>#</th><th>Mahsulot</th><th>Sotilgan miqdor</th><th>Jami daromad</th></tr></thead>
                    <tbody>
                      {topProducts.map((p, i) => (
                        <tr key={p.product_id}>
                          <td style={{ fontWeight: 700, color: i < 3 ? 'var(--accent-yellow)' : 'var(--text-muted)' }}>
                            {i < 3 ? ['🥇','🥈','🥉'][i] : i + 1}
                          </td>
                          <td style={{ fontWeight: 600 }}>{p.product_name}</td>
                          <td>{formatNumber(p.total_sold)}</td>
                          <td style={{ color: 'var(--accent-green)', fontWeight: 600 }}>{formatMoney(p.total_revenue)}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            )}
          </>
        )}
      </div>
    </>
  );
}
