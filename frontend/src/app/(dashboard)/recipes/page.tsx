'use client';
import { useEffect, useState } from 'react';
import { api } from '@/lib/api';
import { Product, Recipe, RecipeIngredient, formatNumber } from '@/lib/types';
import { Modal, Confirm, useModal } from '@/components/ui/Modal';

export default function RecipesPage() {
  const [recipes, setRecipes] = useState<Recipe[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const modal = useModal();
  const confirm = useModal();
  const [editing, setEditing] = useState<Recipe | null>(null);
  const [deleteId, setDeleteId] = useState(0);

  const [form, setForm] = useState({
    product_id: 0,
    name: '',
    ingredients: [] as RecipeIngredient[]
  });

  const load = () => {
    Promise.all([
      api.get<Recipe[]>('/api/recipes'),
      api.get<Product[]>('/api/products')
    ]).then(([r, p]) => {
      setRecipes(r || []);
      setProducts(p || []);
    }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(() => { load(); }, []);

  const openCreate = () => {
    setEditing(null);
    setForm({ product_id: products[0]?.id || 0, name: '', ingredients: [] });
    modal.show();
  };

  const openEdit = (r: Recipe) => {
    setEditing(r);
    setForm({ product_id: r.product_id, name: r.name, ingredients: r.ingredients || [] });
    modal.show();
  };

  const handleSave = async () => {
    try {
      if (editing) {
        await api.put(`/api/recipes/${editing.id}`, form);
      } else {
        await api.post('/api/recipes', form);
      }
      modal.hide();
      load();
    } catch (err: unknown) { alert(err instanceof Error ? err.message : 'Xatolik'); }
  };

  const handleDelete = async () => {
    try {
      await api.delete(`/api/recipes/${deleteId}`);
      load();
    } catch (err: unknown) { alert(err instanceof Error ? err.message : 'Xatolik'); }
  };

  const addIngredient = () => {
    setForm({
      ...form,
      ingredients: [...form.ingredients, { id: 0, recipe_id: 0, ingredient_id: products[0]?.id || 0, quantity: 1 }]
    });
  };

  const updateIngredient = (index: number, field: string, value: any) => {
    const newIngs = [...form.ingredients];
    newIngs[index] = { ...newIngs[index], [field]: value };
    setForm({ ...form, ingredients: newIngs });
  };

  const removeIngredient = (index: number) => {
    const newIngs = [...form.ingredients];
    newIngs.splice(index, 1);
    setForm({ ...form, ingredients: newIngs });
  };

  if (loading) return <div className="loading-page"><div className="spinner" /></div>;

  return (
    <>
      <div className="topbar">
        <h1>🍳 Retseptlar va Kalkulyatsiya</h1>
        <div className="topbar-actions">
          <button className="btn btn-primary" onClick={openCreate}>+ Yangi Retsept</button>
        </div>
      </div>
      <div className="page-content fade-in">
        <div className="table-container">
          <div className="table-scroll">
            <table>
              <thead>
                <tr>
                  <th>ID</th><th>Mahsulot (Taom)</th><th>Retsept Nomi</th><th>Tarkibi</th><th>Amallar</th>
                </tr>
              </thead>
              <tbody>
                {(!recipes || recipes.length === 0) ? (
                  <tr><td colSpan={5}><div className="empty-state"><div className="empty-icon">🍳</div><p>Retseptlar topilmadi</p></div></td></tr>
                ) : recipes.map(r => (
                  <tr key={r.id}>
                    <td>{r.id}</td>
                    <td style={{ fontWeight: 600 }}>{r.product_name}</td>
                    <td>{r.name}</td>
                    <td>
                      <div style={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
                        {r.ingredients?.map(ing => (
                          <span key={ing.id} style={{ fontSize: '0.85rem' }}>
                            • {ing.ingredient_name}: <b>{formatNumber(ing.quantity)} {ing.unit}</b>
                          </span>
                        ))}
                      </div>
                    </td>
                    <td>
                      <div style={{ display: 'flex', gap: 8 }}>
                        <button className="btn btn-ghost btn-sm" onClick={() => openEdit(r)}>✏️</button>
                        <button className="btn btn-ghost btn-sm" onClick={() => { setDeleteId(r.id); confirm.show(); }}>🗑️</button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <Modal open={modal.open} onClose={modal.hide} title={editing ? 'Retseptni tahrirlash' : 'Yangi retsept'}>
        <div className="form-grid">
          <div className="input-group full-width">
            <label>Sotiladigan Tayyor Mahsulot (Taom)</label>
            <select value={form.product_id} onChange={e => setForm({...form, product_id: Number(e.target.value)})}>
              {products.map(p => <option key={p.id} value={p.id}>{p.name}</option>)}
            </select>
          </div>
          <div className="input-group full-width">
            <label>Retsept Nomi (Ixtiyoriy)</label>
            <input className="input" value={form.name} onChange={e => setForm({...form, name: e.target.value})} placeholder="Masalan: Standard Porsiya" />
          </div>

          <div className="full-width" style={{ marginTop: 16 }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 }}>
              <label style={{ margin: 0 }}>Masalliqlar (Tarkibi)</label>
              <button type="button" className="btn btn-ghost btn-sm" onClick={addIngredient}>+ Qo&apos;shish</button>
            </div>
            
            {form.ingredients.length === 0 ? (
              <div style={{ padding: 20, textAlign: 'center', background: 'var(--bg-secondary)', borderRadius: 'var(--radius-md)', color: 'var(--text-muted)' }}>
                Hali masalliq qo&apos;shilmagan
              </div>
            ) : (
              <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                {form.ingredients.map((ing, i) => (
                  <div key={i} style={{ display: 'flex', gap: 12, alignItems: 'center', background: 'var(--bg-secondary)', padding: '10px 14px', borderRadius: 'var(--radius-md)' }}>
                    <select style={{ flex: 1 }} value={ing.ingredient_id} onChange={e => updateIngredient(i, 'ingredient_id', Number(e.target.value))}>
                      {products.map(p => <option key={p.id} value={p.id}>{p.name} ({p.unit})</option>)}
                    </select>
                    <input className="input" style={{ width: 100 }} type="number" step="0.001" placeholder="Miqdor" value={ing.quantity} onChange={e => updateIngredient(i, 'quantity', Number(e.target.value))} />
                    <button type="button" className="btn btn-ghost btn-sm" style={{ color: 'var(--accent-red)', padding: '6px 10px' }} onClick={() => removeIngredient(i)}>✕</button>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
        <div className="modal-actions" style={{ marginTop: 24 }}>
          <button className="btn btn-ghost" onClick={modal.hide}>Bekor</button>
          <button className="btn btn-primary" onClick={handleSave} disabled={form.ingredients.length === 0}>Saqlash</button>
        </div>
      </Modal>

      <Confirm open={confirm.open} onClose={confirm.hide} onConfirm={handleDelete}
        title="O'chirishni tasdiqlash" message="Bu retseptni o'chirishga ishonchingiz komilmi?" />
    </>
  );
}
