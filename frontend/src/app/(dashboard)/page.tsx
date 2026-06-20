'use client';
import { useEffect, useState } from 'react';
import { api } from '@/lib/api';
import { useAuth } from '@/lib/auth';
import { DashboardSummary, ProfitReport, TopProduct, DashboardAlert, formatMoney, formatNumber } from '@/lib/types';

export default function DashboardPage() {
  const { isAdmin } = useAuth();
  const [summary, setSummary] = useState<DashboardSummary | null>(null);
  const [profit, setProfit] = useState<ProfitReport | null>(null);
  const [topProducts, setTopProducts] = useState<TopProduct[]>([]);
  const [alerts, setAlerts] = useState<DashboardAlert[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    Promise.all([
      api.get<DashboardSummary>('/api/dashboard/summary'),
      api.get<DashboardAlert[]>('/api/dashboard/alerts'),
      isAdmin ? api.get<ProfitReport>('/api/dashboard/profit') : Promise.resolve(null),
      isAdmin ? api.get<TopProduct[]>('/api/dashboard/top-products?limit=5') : Promise.resolve([]),
    ]).then(([s, a, p, t]) => {
      setSummary(s);
      setAlerts(a || []);
      setProfit(p);
      setTopProducts(t || []);
    }).catch(console.error)
      .finally(() => setLoading(false));
  }, [isAdmin]);

  if (loading) return <div className="loading-page"><div className="spinner" /></div>;

  return (
    <>
      <div className="topbar">
        <h1>📊 Dashboard</h1>
        {isAdmin && (
          <div className="topbar-actions">
            <button className="btn btn-primary" onClick={async () => {
              try {
                await api.post('/api/dashboard/trigger-report', {});
                alert("Hisobot Telegramga jo'natildi!");
              } catch (e: any) { alert(e.message); }
            }}>📲 Telegramga Hisobot Yuborish</button>
          </div>
        )}
      </div>
      <div className="page-content fade-in">
        {/* Stats */}
        <div className="stats-grid">
          <div className="stat-card blue">
            <div className="stat-icon">📥</div>
            <div className="stat-value">{formatMoney(summary?.total_purchases || 0)}</div>
            <div className="stat-label">Jami Kirim</div>
          </div>
          <div className="stat-card green">
            <div className="stat-icon">💰</div>
            <div className="stat-value">{formatMoney(summary?.total_sales || 0)}</div>
            <div className="stat-label">Jami Sotuv</div>
          </div>
          <div className="stat-card yellow">
            <div className="stat-icon">↩️</div>
            <div className="stat-value">{formatMoney(summary?.total_returns || 0)}</div>
            <div className="stat-label">Vozvratlar</div>
          </div>
          <div className="stat-card red">
            <div className="stat-icon">🗑️</div>
            <div className="stat-value">{formatMoney(summary?.total_write_offs || 0)}</div>
            <div className="stat-label">Yo&apos;qotishlar</div>
          </div>
        </div>

        <div className="stats-grid">
          <div className="stat-card cyan">
            <div className="stat-icon">📦</div>
            <div className="stat-value">{summary?.total_products || 0}</div>
            <div className="stat-label">Mahsulotlar soni</div>
          </div>
          <div className="stat-card purple">
            <div className="stat-icon">📁</div>
            <div className="stat-value">{summary?.total_categories || 0}</div>
            <div className="stat-label">Kategoriyalar</div>
          </div>
          <div className="stat-card blue">
            <div className="stat-icon">🚚</div>
            <div className="stat-value">{summary?.total_suppliers || 0}</div>
            <div className="stat-label">Ta&apos;minotchilar</div>
          </div>
          <div className="stat-card green">
            <div className="stat-icon">🏦</div>
            <div className="stat-value">{formatMoney(summary?.inventory_value || 0)}</div>
            <div className="stat-label">Ombor qiymati</div>
          </div>
        </div>

        {/* Alerts Section */}
        {alerts && alerts.length > 0 && (
          <div style={{ marginBottom: 32 }}>
            <h2 style={{ marginBottom: 16 }}>⚠️ Ogohlantirishlar</h2>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
              {alerts.map((a, i) => (
                <div key={i} style={{
                  padding: '16px 20px', borderRadius: 'var(--radius-md)', display: 'flex', justifyContent: 'space-between',
                  alignItems: 'center', background: a.type === 'low_stock' ? 'rgba(239,68,68,0.1)' : 'rgba(245,158,11,0.1)',
                  borderLeft: `4px solid ${a.type === 'low_stock' ? 'var(--accent-red)' : 'var(--accent-yellow)'}`
                }}>
                  <div>
                    <div style={{ fontWeight: 600, color: 'var(--text-main)', marginBottom: 4 }}>{a.product_name}</div>
                    <div style={{ fontSize: '0.85rem', color: 'var(--text-muted)' }}>{a.message}</div>
                  </div>
                  <div style={{ fontWeight: 700, color: a.type === 'low_stock' ? 'var(--accent-red)' : 'var(--accent-yellow)' }}>
                    {a.value}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Admin-only profit section */}
        {isAdmin && profit && (
          <div style={{ marginBottom: 32 }}>
            <div className={`profit-card ${profit.net_profit < 0 ? 'negative' : ''}`}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: 20 }}>
                <div>
                  <div style={{ color: 'var(--text-muted)', fontSize: '0.85rem', marginBottom: 8 }}>Sof Foyda</div>
                  <div className="profit-value">{formatMoney(profit.net_profit)}</div>
                </div>
                <div style={{ display: 'flex', gap: 32, flexWrap: 'wrap' }}>
                  <div>
                    <div style={{ color: 'var(--text-muted)', fontSize: '0.75rem' }}>Sotuv daromadi</div>
                    <div style={{ fontSize: '1.1rem', fontWeight: 700, color: 'var(--accent-green)' }}>{formatMoney(profit.total_sales_revenue)}</div>
                  </div>
                  <div>
                    <div style={{ color: 'var(--text-muted)', fontSize: '0.75rem' }}>Tan narx</div>
                    <div style={{ fontSize: '1.1rem', fontWeight: 700, color: 'var(--accent-yellow)' }}>{formatMoney(profit.total_cost_of_sold)}</div>
                  </div>
                  <div>
                    <div style={{ color: 'var(--text-muted)', fontSize: '0.75rem' }}>Yo&apos;qotishlar</div>
                    <div style={{ fontSize: '1.1rem', fontWeight: 700, color: 'var(--accent-red)' }}>{formatMoney(profit.total_write_off_loss)}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Top Products */}
        {isAdmin && topProducts && topProducts.length > 0 && (
          <div className="table-container">
            <div className="table-header">
              <h2>🏆 Eng ko&apos;p sotilgan mahsulotlar</h2>
            </div>
            <div className="table-scroll">
              <table>
                <thead>
                  <tr>
                    <th>#</th>
                    <th>Mahsulot</th>
                    <th>Sotilgan miqdor</th>
                    <th>Jami daromad</th>
                  </tr>
                </thead>
                <tbody>
                  {topProducts.map((p, i) => (
                    <tr key={p.product_id}>
                      <td>{i + 1}</td>
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
      </div>
    </>
  );
}
