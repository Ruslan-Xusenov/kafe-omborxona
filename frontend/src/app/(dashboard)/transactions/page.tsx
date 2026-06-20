'use client';
import { useEffect, useState } from 'react';
import { api } from '@/lib/api';
import { Transaction, TransactionType, TRANSACTION_LABELS, formatMoney, formatNumber, formatDate } from '@/lib/types';
import { Confirm, useModal } from '@/components/ui/Modal';
import { useAuth } from '@/lib/auth';
import Link from 'next/link';

export default function TransactionsPage() {
  const { isAdmin } = useAuth();
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);
  const confirm = useModal();
  const [deleteId, setDeleteId] = useState(0);
  const [filter, setFilter] = useState({ type: '', date_from: '', date_to: '' });

  const load = () => {
    const params = new URLSearchParams();
    if (filter.type) params.set('type', filter.type);
    if (filter.date_from) params.set('date_from', filter.date_from);
    if (filter.date_to) params.set('date_to', filter.date_to);
    api.get<Transaction[]>(`/api/transactions?${params}`)
      .then(t => setTransactions(t || []))
      .catch(console.error)
      .finally(() => setLoading(false));
  };

  useEffect(() => { load(); }, []);

  const handleDelete = async () => {
    try { await api.delete(`/api/transactions/${deleteId}`); load(); }
    catch (err: unknown) { alert(err instanceof Error ? err.message : 'Xatolik'); }
  };

  if (loading) return <div className="loading-page"><div className="spinner" /></div>;

  return (
    <>
      <div className="topbar">
        <h1>🔄 Tranzaksiyalar</h1>
        <div className="topbar-actions" style={{ display: 'flex', gap: 12 }}>
          <button className="btn btn-ghost" onClick={() => {
            const params = new URLSearchParams();
            if (filter.type) params.set('type', filter.type);
            if (filter.date_from) params.set('date_from', filter.date_from);
            if (filter.date_to) params.set('date_to', filter.date_to);
            const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
            window.location.href = `${API_URL}/api/export/transactions?token=${localStorage.getItem('token')}&${params}`;
          }}>📥 Excelga Yuklash</button>
          <Link href="/transactions/new" className="btn btn-primary">+ Yangi tranzaksiya</Link>
        </div>
      </div>
      <div className="page-content fade-in">
        <div className="filter-bar">
          <div className="input-group">
            <label>Turi</label>
            <select value={filter.type} onChange={e => setFilter({...filter, type: e.target.value})}>
              <option value="">Barchasi</option>
              <option value="purchase">Kirim</option>
              <option value="return">Vozvrat</option>
              <option value="sale">Sotuv</option>
              <option value="write_off">Spisaniye</option>
            </select>
          </div>
          <div className="input-group">
            <label>Boshlanish</label>
            <input className="input" type="date" value={filter.date_from} onChange={e => setFilter({...filter, date_from: e.target.value})} />
          </div>
          <div className="input-group">
            <label>Tugash</label>
            <input className="input" type="date" value={filter.date_to} onChange={e => setFilter({...filter, date_to: e.target.value})} />
          </div>
          <button className="btn btn-primary btn-sm" onClick={load} style={{ marginTop: 'auto' }}>🔍 Qidirish</button>
        </div>

        <div className="table-container">
          <div className="table-scroll">
            <table>
              <thead>
                <tr>
                  <th>ID</th><th>Sana</th><th>Turi</th><th>Mahsulot</th><th>Ta&apos;minotchi</th>
                  <th>Muddati</th>
                  <th>Miqdor</th><th>Narx</th><th>Summa</th><th>Xodim</th>
                  {isAdmin && <th>Amallar</th>}
                </tr>
              </thead>
              <tbody>
                {(!transactions || transactions.length === 0) ? (
                  <tr><td colSpan={isAdmin ? 11 : 10}><div className="empty-state"><div className="empty-icon">🔄</div><p>Tranzaksiyalar topilmadi</p></div></td></tr>
                ) : transactions.map(t => (
                  <tr key={t.id}>
                    <td>{t.id}</td>
                    <td>{formatDate(t.created_at)}</td>
                    <td><span className={`badge badge-${t.type}`}>{TRANSACTION_LABELS[t.type as TransactionType]}</span></td>
                    <td style={{ fontWeight: 600 }}>{t.product_name}</td>
                    <td>{t.supplier_name || '—'}</td>
                    <td>{t.expiry_date ? formatDate(t.expiry_date).split(',')[0] : '—'}</td>
                    <td>{formatNumber(t.quantity)}</td>
                    <td>{formatMoney(t.unit_price)}</td>
                    <td style={{ fontWeight: 600 }}>{formatMoney(t.total_amount)}</td>
                    <td>{t.user_name}</td>
                    {isAdmin && (
                      <td>
                        <button className="btn btn-ghost btn-sm" onClick={() => { setDeleteId(t.id); confirm.show(); }}>🗑️</button>
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <Confirm open={confirm.open} onClose={confirm.hide} onConfirm={handleDelete}
        title="Tranzaksiyani o'chirish" message="Bu tranzaksiyani bekor qilishga ishonchingiz komilmi?" />
    </>
  );
}
