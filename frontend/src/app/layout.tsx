import type { Metadata } from 'next';
import './globals.css';
import { AuthProvider } from '@/lib/auth';

export const metadata: Metadata = {
  title: 'Kafe Omborxona - Boshqarish Tizimi',
  description: 'Kafe omboridagi mahsulotlar harakatini nazorat qilish va moliyaviy hisobotlarni avtomatlashtirish tizimi',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="uz">
      <body>
        <AuthProvider>{children}</AuthProvider>
      </body>
    </html>
  );
}
