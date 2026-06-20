'use client';
import { useState } from 'react';
import { useAuth } from '@/lib/auth';
import { useRouter } from 'next/navigation';

export default function LoginPage() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await login(username, password);
      router.push('/');
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Xatolik yuz berdi');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-page">
      <div className="login-card slide-up">
        <div className="login-header">
          <div className="logo-big">☕</div>
          <h1>Kafe Omborxona</h1>
          <p>Tizimga kirish</p>
        </div>

        {error && <div className="login-error">{error}</div>}

        <form className="login-form" onSubmit={handleSubmit}>
          <div className="input-group">
            <label htmlFor="username">Foydalanuvchi nomi</label>
            <input id="username" className="input" type="text"
              placeholder="admin" value={username}
              onChange={e => setUsername(e.target.value)} required />
          </div>
          <div className="input-group">
            <label htmlFor="password">Parol</label>
            <input id="password" className="input" type="password"
              placeholder="••••••••" value={password}
              onChange={e => setPassword(e.target.value)} required />
          </div>
          <button className="btn btn-primary" type="submit" disabled={loading}>
            {loading ? 'Kirish...' : 'Tizimga kirish'}
          </button>
        </form>
      </div>
    </div>
  );
}
