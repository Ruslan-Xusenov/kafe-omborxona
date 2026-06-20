'use client';
import { useAuth } from '@/lib/auth';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';
import Sidebar from '@/components/layout/Sidebar';

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && !user) {
      router.push('/login');
    }
  }, [user, loading, router]);

  if (loading) {
    return (
      <div className="loading-page">
        <div className="spinner" />
        <p style={{ color: 'var(--text-muted)' }}>Yuklanmoqda...</p>
      </div>
    );
  }

  if (!user) return null;

  return (
    <>
      <Sidebar />
      <main className="main-content">
        {children}
      </main>
    </>
  );
}
