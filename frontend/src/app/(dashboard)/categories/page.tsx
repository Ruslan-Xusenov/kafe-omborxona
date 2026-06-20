'use client';
import { useEffect, useState } from 'react';
import { api } from '@/lib/api';
import { Category, formatDate } from '@/lib/types';
import { Modal, Confirm, useModal } from '@/components/ui/Modal';

export default function CategoriesPage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const modal = useModal();
  const confirm = useModal();
  const [editing, setEditing] = useState<Category | null>(null);
  const [deleteId, setDeleteId] = useState(0);
  const [name, setName] = useState('');

  const load = () => {
    api.get<Category[]>('/api/categories')
      .then(c => setCategories(c || []))
      .catch(console.error)
      .finally(() => setLoading(false));
  };

  useEffect(() => { load(); }, []);

  const openCreate = () => { setEditing(null); setName(''); modal.show(); };
  const openEdit = (c: Category) => { setEditing(c); setName(c.name); modal.show(); };

  const handleSave = async () => {
    try {
      if (editing) {
        await api.put(`/api/categories/${editing.id}`, { name });
      } else {
        await api.post('/api/categories', { name });
      }
      modal.hide(); load();
    } catch (err: unknown) { alert(err instanceof Error ? err.message : 'Xatolik'); }
  };

  const handleDelete = async () => {
    try { await api.delete(`/api/categories/${deleteId}`); load(); }
    catch (err: unknown) { alert(err instanceof Error ? err.message : 'Xatolik'); }
  };

  if (loading) return <div className="loading-page"><div className="spinner" /></div>;

  return (
    <>
      <div className="topbar">
        <h1>📁 Kategoriyalar</h1>
        <div className="topbar-actions">
          <button className="btn btn-primary" onClick={openCreate}>+ Yangi kategoriya</button>
        </div>
      </div>
      <div className="page-content fade-in">
        <div className="table-container">
          <div className="table-scroll">
            <table>
              <thead><tr><th>ID</th><th>Nomi</th><th>Yaratilgan</th><th>Amallar</th></tr></thead>
              <tbody>
                {(!categories || categories.length === 0) ? (
                  <tr><td colSpan={4}><div className="empty-state"><div className="empty-icon">📁</div><p>Kategoriyalar topilmadi</p></div></td></tr>
                ) : categories.map(c => (
                  <tr key={c.id}>
                    <td>{c.id}</td>
                    <td style={{ fontWeight: 600 }}>{c.name}</td>
                    <td>{formatDate(c.created_at)}</td>
                    <td>
                      <div style={{ display: 'flex', gap: 8 }}>
                        <button className="btn btn-ghost btn-sm" onClick={() => openEdit(c)}>✏️</button>
                        <button className="btn btn-ghost btn-sm" onClick={() => { setDeleteId(c.id); confirm.show(); }}>🗑️</button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <Modal open={modal.open} onClose={modal.hide} title={editing ? 'Kategoriyani tahrirlash' : 'Yangi kategoriya'}>
        <div className="input-group">
          <label>Kategoriya nomi</label>
          <input className="input" value={name} onChange={e => setName(e.target.value)} placeholder="Masalan: Ichimliklar" />
        </div>
        <div className="modal-actions">
          <button className="btn btn-ghost" onClick={modal.hide}>Bekor</button>
          <button className="btn btn-primary" onClick={handleSave}>Saqlash</button>
        </div>
      </Modal>

      <Confirm open={confirm.open} onClose={confirm.hide} onConfirm={handleDelete}
        title="O'chirishni tasdiqlash" message="Bu kategoriyani o'chirishga ishonchingiz komilmi?" />
    </>
  );
}
