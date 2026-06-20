'use client';
import { useEffect, useState } from 'react';
import { api } from '@/lib/api';
import { Supplier, formatDate } from '@/lib/types';
import { Modal, Confirm, useModal } from '@/components/ui/Modal';

export default function SuppliersPage() {
  const [suppliers, setSuppliers] = useState<Supplier[]>([]);
  const [loading, setLoading] = useState(true);
  const modal = useModal();
  const confirm = useModal();
  const [editing, setEditing] = useState<Supplier | null>(null);
  const [deleteId, setDeleteId] = useState(0);
  const [form, setForm] = useState({ name: '', phone: '', address: '' });

  const load = () => {
    api.get<Supplier[]>('/api/suppliers').then(s => setSuppliers(s || [])).catch(console.error).finally(() => setLoading(false));
  };
  useEffect(() => { load(); }, []);

  const openCreate = () => { setEditing(null); setForm({ name: '', phone: '', address: '' }); modal.show(); };
  const openEdit = (s: Supplier) => { setEditing(s); setForm({ name: s.name, phone: s.phone, address: s.address }); modal.show(); };

  const handleSave = async () => {
    try {
      if (editing) await api.put(`/api/suppliers/${editing.id}`, form);
      else await api.post('/api/suppliers', form);
      modal.hide(); load();
    } catch (err: unknown) { alert(err instanceof Error ? err.message : 'Xatolik'); }
  };

  const handleDelete = async () => {
    try { await api.delete(`/api/suppliers/${deleteId}`); load(); }
    catch (err: unknown) { alert(err instanceof Error ? err.message : 'Xatolik'); }
  };

  if (loading) return <div className="loading-page"><div className="spinner" /></div>;

  return (
    <>
      <div className="topbar">
        <h1>🚚 Ta&apos;minotchilar</h1>
        <div className="topbar-actions">
          <button className="btn btn-primary" onClick={openCreate}>+ Yangi ta&apos;minotchi</button>
        </div>
      </div>
      <div className="page-content fade-in">
        <div className="table-container">
          <div className="table-scroll">
            <table>
              <thead><tr><th>ID</th><th>Nomi</th><th>Telefon</th><th>Manzil</th><th>Yaratilgan</th><th>Amallar</th></tr></thead>
              <tbody>
                {(!suppliers || suppliers.length === 0) ? (
                  <tr><td colSpan={6}><div className="empty-state"><div className="empty-icon">🚚</div><p>Ta&apos;minotchilar topilmadi</p></div></td></tr>
                ) : suppliers.map(s => (
                  <tr key={s.id}>
                    <td>{s.id}</td>
                    <td style={{ fontWeight: 600 }}>{s.name}</td>
                    <td>{s.phone || '—'}</td>
                    <td>{s.address || '—'}</td>
                    <td>{formatDate(s.created_at)}</td>
                    <td>
                      <div style={{ display: 'flex', gap: 8 }}>
                        <button className="btn btn-ghost btn-sm" onClick={() => openEdit(s)}>✏️</button>
                        <button className="btn btn-ghost btn-sm" onClick={() => { setDeleteId(s.id); confirm.show(); }}>🗑️</button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <Modal open={modal.open} onClose={modal.hide} title={editing ? "Ta'minotchini tahrirlash" : "Yangi ta'minotchi"}>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
          <div className="input-group">
            <label>Nomi</label>
            <input className="input" value={form.name} onChange={e => setForm({...form, name: e.target.value})} />
          </div>
          <div className="input-group">
            <label>Telefon</label>
            <input className="input" value={form.phone} onChange={e => setForm({...form, phone: e.target.value})} placeholder="+998 90 123 45 67" />
          </div>
          <div className="input-group">
            <label>Manzil</label>
            <input className="input" value={form.address} onChange={e => setForm({...form, address: e.target.value})} />
          </div>
        </div>
        <div className="modal-actions">
          <button className="btn btn-ghost" onClick={modal.hide}>Bekor</button>
          <button className="btn btn-primary" onClick={handleSave}>Saqlash</button>
        </div>
      </Modal>

      <Confirm open={confirm.open} onClose={confirm.hide} onConfirm={handleDelete}
        title="O'chirishni tasdiqlash" message="Bu ta'minotchini o'chirishga ishonchingiz komilmi?" />
    </>
  );
}
