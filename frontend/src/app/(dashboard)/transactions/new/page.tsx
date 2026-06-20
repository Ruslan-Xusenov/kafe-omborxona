'use client';
import { useEffect, useState } from 'react';
import { api } from '@/lib/api';
import { Product, Supplier } from '@/lib/types';
import { useRouter } from 'next/navigation';

export default function NewTransactionPage() {
  const router = useRouter();
  const [products, setProducts] = useState<Product[]>([]);
  const [suppliers, setSuppliers] = useState<Supplier[]>([]);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [form, setForm] = useState({
    product_id: 0, supplier_id: null as number | null,
    type: 'purchase' as string, quantity: '', unit_price: '', note: '', expiry_date: '',
  });

  useEffect(() => {
    Promise.all([
      api.get<Product[]>('/api/products'),
      api.get<Supplier[]>('/api/suppliers'),
    ]).then(([p, s]) => {
      setProducts(p || []); setSuppliers(s || []);
      if (p && p.length > 0) setForm(f => ({ ...f, product_id: p[0].id }));
    }).catch(console.error).finally(() => setLoading(false));
  }, []);

  const selectedProduct = products.find(p => p.id === form.product_id);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    try {
      await api.post('/api/transactions', {
        product_id: form.product_id,
        supplier_id: form.supplier_id || null,
        type: form.type,
        quantity: Number(form.quantity),
        unit_price: Number(form.unit_price),
        note: form.note,
        expiry_date: form.expiry_date || null,
      });
      router.push('/transactions');
    } catch (err: unknown) {
      alert(err instanceof Error ? err.message : 'Xatolik');
    } finally { setSubmitting(false); }
  };

  // Auto-fill price based on type
  const handleTypeChange = (type: string) => {
    setForm(f => {
      const price = type === 'sale' ? (selectedProduct?.sale_price || 0).toString()
        : (selectedProduct?.cost_price || 0).toString();
      return { ...f, type, unit_price: price };
    });
  };

  const handleProductChange = (id: number) => {
    const p = products.find(pr => pr.id === id);
    setForm(f => ({
      ...f, product_id: id,
      unit_price: f.type === 'sale' ? (p?.sale_price || 0).toString() : (p?.cost_price || 0).toString(),
    }));
  };

  if (loading) return <div className="loading-page"><div className="spinner" /></div>;

  const total = Number(form.quantity) * Number(form.unit_price);

  return (
    <>
      <div className="topbar">
        <h1>➕ Yangi Tranzaksiya</h1>
      </div>
      <div className="page-content fade-in">
        <div className="card" style={{ maxWidth: 700 }}>
          <form onSubmit={handleSubmit}>
            <div className="form-grid">
              <div className="input-group full-width">
                <label>Tranzaksiya turi</label>
                <div style={{ display: 'flex', gap: 10 }}>
                  {[
                    { val: 'purchase', label: '📥 Kirim', cls: 'btn-primary' },
                    { val: 'return', label: '↩️ Vozvrat', cls: 'btn-ghost' },
                    { val: 'sale', label: '💰 Sotuv', cls: 'btn-success' },
                    { val: 'write_off', label: '🗑️ Spisaniye', cls: 'btn-danger' },
                  ].map(t => (
                    <button key={t.val} type="button"
                      className={`btn btn-sm ${form.type === t.val ? t.cls : 'btn-ghost'}`}
                      onClick={() => handleTypeChange(t.val)}>
                      {t.label}
                    </button>
                  ))}
                </div>
              </div>

              <div className="input-group">
                <label>Shtrix-kod skaneri (yoki mahsulotni tanlang)</label>
                <input className="input" type="text" autoFocus placeholder="Skanerlang..." 
                  onChange={e => {
                    const code = e.target.value;
                    if (code) {
                      const found = products.find(p => p.barcode === code);
                      if (found) handleProductChange(found.id);
                    }
                  }} 
                />
              </div>

              <div className="input-group">
                <label>Mahsulot</label>
                <select value={form.product_id} onChange={e => handleProductChange(Number(e.target.value))}>
                  {products.map(p => <option key={p.id} value={p.id}>{p.name}</option>)}
                </select>
              </div>

              {(form.type === 'purchase' || form.type === 'return') && (
                <div className="input-group">
                  <label>Ta&apos;minotchi</label>
                  <select value={form.supplier_id || ''} onChange={e => setForm({...form, supplier_id: e.target.value ? Number(e.target.value) : null})}>
                    <option value="">— Tanlanmagan —</option>
                    {suppliers.map(s => <option key={s.id} value={s.id}>{s.name}</option>)}
                  </select>
                </div>
              )}

              <div className="input-group">
                <label>Miqdor ({selectedProduct?.unit || 'dona'})</label>
                <input className="input" type="number" step="0.001" min="0.001"
                  value={form.quantity} onChange={e => setForm({...form, quantity: e.target.value})} required />
              </div>

              <div className="input-group">
                <label>Narx (UZS)</label>
                <input className="input" type="number" step="0.01" min="0"
                  value={form.unit_price} onChange={e => setForm({...form, unit_price: e.target.value})} required />
              </div>

              {(form.type === 'purchase') && (
                <div className="input-group">
                  <label>Yaroqlilik muddati (ixtiyoriy)</label>
                  <input className="input" type="date" value={form.expiry_date} onChange={e => setForm({...form, expiry_date: e.target.value})} />
                </div>
              )}

              <div className="input-group full-width">
                <label>Izoh</label>
                <textarea className="input" rows={2} value={form.note}
                  onChange={e => setForm({...form, note: e.target.value})} placeholder="Qo'shimcha izoh..." />
              </div>

              {/* Total */}
              <div className="full-width" style={{
                background: 'rgba(59,130,246,0.1)', borderRadius: 'var(--radius-md)',
                padding: '16px 20px', display: 'flex', justifyContent: 'space-between', alignItems: 'center'
              }}>
                <span style={{ color: 'var(--text-muted)', fontWeight: 600 }}>Umumiy summa:</span>
                <span style={{ fontSize: '1.5rem', fontWeight: 800, color: 'var(--accent-blue)' }}>
                  {new Intl.NumberFormat('uz-UZ').format(Math.round(total || 0))} UZS
                </span>
              </div>
            </div>

            <div className="modal-actions" style={{ marginTop: 24 }}>
              <button type="button" className="btn btn-ghost" onClick={() => router.back()}>Bekor</button>
              <button type="submit" className="btn btn-primary btn-lg" disabled={submitting}>
                {submitting ? 'Saqlanmoqda...' : '✅ Saqlash'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </>
  );
}
