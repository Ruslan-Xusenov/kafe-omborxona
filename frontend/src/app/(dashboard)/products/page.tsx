'use client';
import { useEffect, useState } from 'react';
import { api } from '@/lib/api';
import { Product, Category, formatMoney, formatNumber } from '@/lib/types';
import { Modal, Confirm, useModal } from '@/components/ui/Modal';

export default function ProductsPage() {
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const modal = useModal();
  const confirm = useModal();
  const [editing, setEditing] = useState<Product | null>(null);
  const [deleteId, setDeleteId] = useState(0);
  const [form, setForm] = useState({ name: '', unit: 'dona', category_id: 0, cost_price: 0, sale_price: 0, min_stock: 0, barcode: '', initial_quantity: 0 });

  const load = () => {
    Promise.all([
      api.get<Product[]>('/api/products'),
      api.get<Category[]>('/api/categories'),
    ]).then(([p, c]) => { setProducts(p || []); setCategories(c || []); })
      .catch(console.error).finally(() => setLoading(false));
  };

  useEffect(() => { load(); }, []);

  const openCreate = () => {
    setEditing(null);
    setForm({ name: '', unit: 'dona', category_id: categories[0]?.id || 0, cost_price: 0, sale_price: 0, min_stock: 0, barcode: '', initial_quantity: 0 });
    modal.show();
  };

  const openEdit = (p: Product) => {
    setEditing(p);
    setForm({ name: p.name, unit: p.unit, category_id: p.category_id, cost_price: p.cost_price, sale_price: p.sale_price, min_stock: p.min_stock || 0, barcode: p.barcode || '', initial_quantity: 0 });
    modal.show();
  };

  const handleSave = async () => {
    try {
      if (editing) {
        await api.put(`/api/products/${editing.id}`, form);
      } else {
        const newProduct = await api.post<Product>('/api/products', form);
        if (form.initial_quantity > 0 && newProduct?.id) {
          await api.post('/api/transactions', {
            product_id: newProduct.id,
            type: 'purchase',
            quantity: form.initial_quantity,
            unit_price: form.cost_price,
            note: "Dastlabki miqdor"
          });
        }
      }
      modal.hide();
      load();
    } catch (err: unknown) { alert(err instanceof Error ? err.message : 'Xatolik'); }
  };

  const handleDelete = async () => {
    try {
      await api.delete(`/api/products/${deleteId}`);
      load();
    } catch (err: unknown) { alert(err instanceof Error ? err.message : 'Xatolik'); }
  };

  if (loading) return <div className="loading-page"><div className="spinner" /></div>;

  return (
    <>
      <div className="topbar">
        <h1>🏷️ Mahsulotlar</h1>
        <div className="topbar-actions">
          <button className="btn btn-primary" onClick={openCreate}>+ Yangi mahsulot</button>
        </div>
      </div>
      <div className="page-content fade-in">
        <div className="table-container">
          <div className="table-scroll">
            <table>
              <thead>
                <tr>
                  <th>ID</th><th>Nomi</th><th>Kategoriya</th><th>O&apos;lchov</th>
                  <th>Tan narxi</th><th>Sotuv narxi</th><th>Shtrix-kod</th><th>Min. Qoldiq</th><th>Qoldiq</th><th>Qiymat</th><th>Amallar</th>
                </tr>
              </thead>
              <tbody>
                {(!products || products.length === 0) ? (
                  <tr><td colSpan={11}><div className="empty-state"><div className="empty-icon">📦</div><p>Mahsulotlar topilmadi</p></div></td></tr>
                ) : products.map(p => (
                  <tr key={p.id}>
                    <td>{p.id}</td>
                    <td style={{ fontWeight: 600 }}>{p.name}</td>
                    <td><span className="badge badge-purchase">{p.category_name}</span></td>
                    <td>{p.unit}</td>
                    <td>{formatMoney(p.cost_price)}</td>
                    <td>{formatMoney(p.sale_price)}</td>
                    <td style={{ fontFamily: 'monospace', color: 'var(--text-muted)' }}>{p.barcode || '—'}</td>
                    <td>{formatNumber(p.min_stock || 0)}</td>
                    <td style={{ fontWeight: 700, color: (p.current_stock || 0) <= (p.min_stock || 0) ? 'var(--accent-red)' : 'var(--accent-green)' }}>
                      {formatNumber(p.current_stock || 0)}
                    </td>
                    <td>{formatMoney(p.stock_value || 0)}</td>
                    <td>
                      <div style={{ display: 'flex', gap: 8 }}>
                        <button className="btn btn-ghost btn-sm" onClick={() => openEdit(p)}>✏️</button>
                        <button className="btn btn-ghost btn-sm" onClick={() => { setDeleteId(p.id); confirm.show(); }}>🗑️</button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <Modal open={modal.open} onClose={modal.hide} title={editing ? 'Mahsulotni tahrirlash' : 'Yangi mahsulot'}>
        <div className="form-grid">
          <div className="input-group full-width">
            <label>Nomi</label>
            <input className="input" value={form.name} onChange={e => setForm({...form, name: e.target.value})} />
          </div>
          <div className="input-group">
            <label>O&apos;lchov birligi</label>
            <select value={form.unit} onChange={e => setForm({...form, unit: e.target.value})}>
              <option value="dona">Dona</option><option value="kg">Kg</option>
              <option value="litr">Litr</option><option value="paket">Paket</option>
            </select>
          </div>
          <div className="input-group">
            <label>Kategoriya</label>
            <select value={form.category_id} onChange={e => setForm({...form, category_id: Number(e.target.value)})}>
              {categories.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
            </select>
          </div>
          <div className="input-group">
            <label>Tan narxi (UZS)</label>
            <input className="input" type="number" value={form.cost_price} onChange={e => setForm({...form, cost_price: Number(e.target.value)})} />
          </div>
          <div className="input-group">
            <label>Sotuv narxi (UZS)</label>
            <input className="input" type="number" value={form.sale_price} onChange={e => setForm({...form, sale_price: Number(e.target.value)})} />
          </div>
          <div className="input-group">
            <label>Shtrix-kod</label>
            <input className="input" type="text" value={form.barcode} onChange={e => setForm({...form, barcode: e.target.value})} placeholder="Skanerlang yoki kiriting" />
          </div>
          <div className="input-group">
            <label>Minimal qoldiq</label>
            <input className="input" type="number" step="0.001" value={form.min_stock} onChange={e => setForm({...form, min_stock: Number(e.target.value)})} placeholder="Tugayotganini bilish uchun" />
          </div>
          {!editing && (
            <div className="input-group">
              <label>Dastlabki miqdor (Omborga kirim)</label>
              <input className="input" type="number" step="0.001" value={form.initial_quantity} onChange={e => setForm({...form, initial_quantity: Number(e.target.value)})} placeholder="0" />
            </div>
          )}
        </div>
        <div className="modal-actions">
          <button className="btn btn-ghost" onClick={modal.hide}>Bekor</button>
          <button className="btn btn-primary" onClick={handleSave}>Saqlash</button>
        </div>
      </Modal>

      <Confirm open={confirm.open} onClose={confirm.hide} onConfirm={handleDelete}
        title="O'chirishni tasdiqlash" message="Bu mahsulotni o'chirishga ishonchingiz komilmi?" />
    </>
  );
}
