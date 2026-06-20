'use client';
import { useEffect, useState } from 'react';
import { api } from '@/lib/api';
import { Supplier, Debt, Transaction, formatMoney, formatDate } from '@/lib/types';
import { Modal, Confirm, useModal } from '@/components/ui/Modal';
import { useAuth } from '@/lib/auth';

export default function DebtsPage() {
  const { isAdmin } = useAuth();
  const [debts, setDebts] = useState<Debt[]>([]);
  const [suppliers, setSuppliers] = useState<Supplier[]>([]);
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);
  const [deleteId, setDeleteId] = useState(0);

  const modal = useModal();
  const payModal = useModal();
  const confirm = useModal();

  const [form, setForm] = useState({ supplier_id: 0, transaction_id: 0, total_debt: '', due_date: '' });
  const [payForm, setPayForm] = useState({ id: 0, amount: '' });

  const load = () => {
    Promise.all([
      api.get<Debt[]>('/api/debts'),
      api.get<Supplier[]>('/api/suppliers'),
      api.get<Transaction[]>('/api/transactions?type=purchase')
    ]).then(([d, s, t]) => {
      setDebts(d || []);
      setSuppliers(s || []);
      setTransactions(t || []);
    }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(() => { load(); }, []);

  const openCreate = () => {
    setForm({ supplier_id: suppliers[0]?.id || 0, transaction_id: 0, total_debt: '', due_date: '' });
    modal.show();
  };

  const openPay = (d: Debt) => {
    setPayForm({ id: d.id, amount: String(d.total_debt - d.paid_amount) });
    payModal.show();
  };

  const handleSave = async () => {
    try {
      await api.post('/api/debts', {
        supplier_id: form.supplier_id,
        transaction_id: form.transaction_id || null,
        total_debt: Number(form.total_debt),
        due_date: form.due_date || null
      });
      modal.hide();
      load();
    } catch (err: unknown) { alert(err instanceof Error ? err.message : 'Xatolik'); }
  };

  const handlePay = async () => {
    try {
      await api.post(`/api/debts/${payForm.id}/pay`, { amount: Number(payForm.amount) });
      payModal.hide();
      load();
    } catch (err: unknown) { alert(err instanceof Error ? err.message : 'Xatolik'); }
  };

  const handleDelete = async () => {
    try {
      await api.delete(`/api/debts/${deleteId}`);
      load();
    } catch (err: unknown) { alert(err instanceof Error ? err.message : 'Xatolik'); }
  };

  const statusLabel = {
    unpaid: { text: "To'lanmagan", class: 'red' },
    partial: { text: "Qisman", class: 'yellow' },
    paid: { text: "To'landi", class: 'green' }
  };

  if (loading) return <div className="loading-page"><div className="spinner" /></div>;

  return (
    <>
      <div className="topbar">
        <h1>🤝 Qarzdorlik (Debts)</h1>
        <div className="topbar-actions">
          <button className="btn btn-primary" onClick={openCreate}>+ Qarz Qayd Etish</button>
        </div>
      </div>
      <div className="page-content fade-in">
        <div className="table-container">
          <div className="table-scroll">
            <table>
              <thead>
                <tr>
                  <th>ID</th><th>Sana</th><th>Ta'minotchi</th><th>Tranzaksiya</th><th>Jami Qarz</th><th>To'landi</th><th>Qoldiq</th><th>Holati</th><th>Muddat</th><th>Amallar</th>
                </tr>
              </thead>
              <tbody>
                {(!debts || debts.length === 0) ? (
                  <tr><td colSpan={10}><div className="empty-state"><div className="empty-icon">🤝</div><p>Qarzlar topilmadi</p></div></td></tr>
                ) : debts.map(d => {
                  const remains = d.total_debt - d.paid_amount;
                  return (
                  <tr key={d.id}>
                    <td>{d.id}</td>
                    <td>{formatDate(d.created_at)}</td>
                    <td style={{ fontWeight: 600 }}>{d.supplier_name}</td>
                    <td>{d.transaction_id ? `#${d.transaction_id}` : '—'}</td>
                    <td>{formatMoney(d.total_debt)}</td>
                    <td style={{ color: 'var(--accent-green)' }}>{formatMoney(d.paid_amount)}</td>
                    <td style={{ fontWeight: 700, color: remains > 0 ? 'var(--accent-red)' : 'var(--text-main)' }}>
                      {formatMoney(remains)}
                    </td>
                    <td><span className={`badge badge-${statusLabel[d.status].class}`}>{statusLabel[d.status].text}</span></td>
                    <td>{d.due_date ? formatDate(d.due_date).split(',')[0] : '—'}</td>
                    <td>
                      <div style={{ display: 'flex', gap: 8 }}>
                        {d.status !== 'paid' && <button className="btn btn-primary btn-sm" onClick={() => openPay(d)}>To'lash</button>}
                        {isAdmin && <button className="btn btn-ghost btn-sm" onClick={() => { setDeleteId(d.id); confirm.show(); }}>🗑️</button>}
                      </div>
                    </td>
                  </tr>
                )})}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <Modal open={modal.open} onClose={modal.hide} title="Yangi Qarz Qayd Etish">
        <div className="form-grid">
          <div className="input-group full-width">
            <label>Ta'minotchi</label>
            <select value={form.supplier_id} onChange={e => setForm({...form, supplier_id: Number(e.target.value)})}>
              <option value={0}>Tanlang</option>
              {suppliers.map(s => <option key={s.id} value={s.id}>{s.name}</option>)}
            </select>
          </div>
          <div className="input-group full-width">
            <label>Bog'liq Tranzaksiya (Ixtiyoriy)</label>
            <select value={form.transaction_id} onChange={e => setForm({...form, transaction_id: Number(e.target.value)})}>
              <option value={0}>Bog'lanmagan</option>
              {transactions.filter(t => t.supplier_id === form.supplier_id).map(t => (
                <option key={t.id} value={t.id}>#{t.id} - {t.product_name} ({formatMoney(t.total_amount)})</option>
              ))}
            </select>
          </div>
          <div className="input-group">
            <label>Qarz summasi</label>
            <input className="input" type="number" value={form.total_debt} onChange={e => setForm({...form, total_debt: e.target.value})} placeholder="0" />
          </div>
          <div className="input-group">
            <label>To'lash muddati (Ixtiyoriy)</label>
            <input className="input" type="date" value={form.due_date} onChange={e => setForm({...form, due_date: e.target.value})} />
          </div>
        </div>
        <div className="modal-actions" style={{ marginTop: 24 }}>
          <button className="btn btn-ghost" onClick={modal.hide}>Bekor</button>
          <button className="btn btn-primary" onClick={handleSave} disabled={!form.supplier_id || !form.total_debt}>Saqlash</button>
        </div>
      </Modal>

      <Modal open={payModal.open} onClose={payModal.hide} title="Qarzni To'lash">
        <div className="form-grid">
          <div className="input-group full-width">
            <label>To'lov summasi</label>
            <input className="input" type="number" value={payForm.amount} onChange={e => setPayForm({...payForm, amount: e.target.value})} />
          </div>
        </div>
        <div className="modal-actions" style={{ marginTop: 24 }}>
          <button className="btn btn-ghost" onClick={payModal.hide}>Bekor</button>
          <button className="btn btn-primary" onClick={handlePay} disabled={!payForm.amount}>To'lash</button>
        </div>
      </Modal>

      <Confirm open={confirm.open} onClose={confirm.hide} onConfirm={handleDelete}
        title="Qarzni o'chirish" message="Bu ma'lumotni o'chirishga ishonchingiz komilmi?" />
    </>
  );
}
