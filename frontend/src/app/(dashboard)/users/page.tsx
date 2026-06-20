'use client';
import { useEffect, useState } from 'react';
import { api } from '@/lib/api';
import { useAuth } from '@/lib/auth';
import { User, formatDate } from '@/lib/types';
import { Modal, Confirm, useModal } from '@/components/ui/Modal';
import { useRouter } from 'next/navigation';

export default function UsersPage() {
  const { isAdmin, loading: authLoading } = useAuth();
  const router = useRouter();
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const modal = useModal();
  const confirm = useModal();
  const [editing, setEditing] = useState<User | null>(null);
  const [deleteId, setDeleteId] = useState(0);
  const [form, setForm] = useState({ username: '', password: '', full_name: '', role: 'warehouse_manager' as string });

  useEffect(() => {
    if (!authLoading && !isAdmin) { router.push('/'); return; }
  }, [isAdmin, authLoading, router]);

  const load = () => {
    api.get<User[]>('/api/users').then(u => setUsers(u || [])).catch(console.error).finally(() => setLoading(false));
  };
  useEffect(() => { if (isAdmin) load(); }, [isAdmin]);

  const openCreate = () => {
    setEditing(null);
    setForm({ username: '', password: '', full_name: '', role: 'warehouse_manager' });
    modal.show();
  };
  const openEdit = (u: User) => {
    setEditing(u);
    setForm({ username: u.username, password: '', full_name: u.full_name, role: u.role });
    modal.show();
  };

  const handleSave = async () => {
    try {
      if (editing) {
        await api.put(`/api/users/${editing.id}`, form);
      } else {
        await api.post('/api/users', form);
      }
      modal.hide(); load();
    } catch (err: unknown) { alert(err instanceof Error ? err.message : 'Xatolik'); }
  };

  const handleDelete = async () => {
    try { await api.delete(`/api/users/${deleteId}`); load(); }
    catch (err: unknown) { alert(err instanceof Error ? err.message : 'Xatolik'); }
  };

  if (!isAdmin) return null;
  if (loading) return <div className="loading-page"><div className="spinner" /></div>;

  return (
    <>
      <div className="topbar">
        <h1>👥 Foydalanuvchilar</h1>
        <div className="topbar-actions">
          <button className="btn btn-primary" onClick={openCreate}>+ Yangi foydalanuvchi</button>
        </div>
      </div>
      <div className="page-content fade-in">
        <div className="table-container">
          <div className="table-scroll">
            <table>
              <thead><tr><th>ID</th><th>Login</th><th>To&apos;liq ism</th><th>Roli</th><th>Yaratilgan</th><th>Amallar</th></tr></thead>
              <tbody>
                {(!users || users.length === 0) ? (
                  <tr><td colSpan={6}><div className="empty-state"><div className="empty-icon">👥</div><p>Foydalanuvchilar topilmadi</p></div></td></tr>
                ) : users.map(u => (
                  <tr key={u.id}>
                    <td>{u.id}</td>
                    <td style={{ fontWeight: 600 }}>{u.username}</td>
                    <td>{u.full_name}</td>
                    <td><span className={`badge badge-${u.role}`}>{u.role === 'admin' ? 'Administrator' : 'Ombor mudiri'}</span></td>
                    <td>{formatDate(u.created_at)}</td>
                    <td>
                      <div style={{ display: 'flex', gap: 8 }}>
                        <button className="btn btn-ghost btn-sm" onClick={() => openEdit(u)}>✏️</button>
                        <button className="btn btn-ghost btn-sm" onClick={() => { setDeleteId(u.id); confirm.show(); }}>🗑️</button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <Modal open={modal.open} onClose={modal.hide} title={editing ? 'Foydalanuvchini tahrirlash' : 'Yangi foydalanuvchi'}>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
          <div className="input-group">
            <label>Login</label>
            <input className="input" value={form.username} onChange={e => setForm({...form, username: e.target.value})} />
          </div>
          <div className="input-group">
            <label>{editing ? 'Yangi parol (bo\'sh qoldirsa o\'zgarmaydi)' : 'Parol'}</label>
            <input className="input" type="password" value={form.password} onChange={e => setForm({...form, password: e.target.value})}
              required={!editing} />
          </div>
          <div className="input-group">
            <label>To&apos;liq ism</label>
            <input className="input" value={form.full_name} onChange={e => setForm({...form, full_name: e.target.value})} />
          </div>
          <div className="input-group">
            <label>Rol</label>
            <select value={form.role} onChange={e => setForm({...form, role: e.target.value})}>
              <option value="admin">Administrator</option>
              <option value="warehouse_manager">Ombor mudiri</option>
            </select>
          </div>
        </div>
        <div className="modal-actions">
          <button className="btn btn-ghost" onClick={modal.hide}>Bekor</button>
          <button className="btn btn-primary" onClick={handleSave}>Saqlash</button>
        </div>
      </Modal>

      <Confirm open={confirm.open} onClose={confirm.hide} onConfirm={handleDelete}
        title="O'chirishni tasdiqlash" message="Bu foydalanuvchini o'chirishga ishonchingiz komilmi?" />
    </>
  );
}
