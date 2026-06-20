'use client';
import { useEffect, useState } from 'react';
import { api } from '@/lib/api';
import { InventoryItem, formatMoney, formatNumber } from '@/lib/types';

export default function InventoryPage() {
  const [items, setItems] = useState<InventoryItem[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get<InventoryItem[]>('/api/dashboard/inventory')
      .then(i => setItems(i || []))
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <div className="loading-page"><div className="spinner" /></div>;

  const totalValue = items.reduce((sum, i) => sum + i.stock_value, 0);
  const totalItems = items.reduce((sum, i) => sum + i.current_stock, 0);

  return (
    <>
      <div className="topbar">
        <h1>📦 Ombor Qoldig&apos;i</h1>
      </div>
      <div className="page-content fade-in">
        <div className="stats-grid" style={{ marginBottom: 24 }}>
          <div className="stat-card cyan">
            <div className="stat-icon">📦</div>
            <div className="stat-value">{formatNumber(totalItems)}</div>
            <div className="stat-label">Jami mahsulot miqdori</div>
          </div>
          <div className="stat-card green">
            <div className="stat-icon">🏦</div>
            <div className="stat-value">{formatMoney(totalValue)}</div>
            <div className="stat-label">Ombor umumiy qiymati</div>
          </div>
        </div>

        <div className="table-container">
          <div className="table-scroll">
            <table>
              <thead>
                <tr>
                  <th>Mahsulot</th><th>Kategoriya</th><th>O&apos;lchov</th>
                  <th>Tan narxi</th><th>Sotuv narxi</th><th>Qoldiq</th><th>Qiymat</th>
                </tr>
              </thead>
              <tbody>
                {(!items || items.length === 0) ? (
                  <tr><td colSpan={7}><div className="empty-state"><div className="empty-icon">📦</div><p>Ombor bo&apos;sh</p></div></td></tr>
                ) : items.map(i => (
                  <tr key={i.product_id}>
                    <td style={{ fontWeight: 600 }}>{i.product_name}</td>
                    <td><span className="badge badge-purchase">{i.category_name}</span></td>
                    <td>{i.unit}</td>
                    <td>{formatMoney(i.cost_price)}</td>
                    <td>{formatMoney(i.sale_price)}</td>
                    <td style={{
                      fontWeight: 700,
                      color: i.current_stock <= 0 ? 'var(--accent-red)' : i.current_stock < 10 ? 'var(--accent-yellow)' : 'var(--accent-green)'
                    }}>
                      {formatNumber(i.current_stock)} {i.unit}
                    </td>
                    <td style={{ fontWeight: 600 }}>{formatMoney(i.stock_value)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </>
  );
}
