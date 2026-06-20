'use client';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useAuth } from '@/lib/auth';

const navItems = [
  { section: 'Asosiy', items: [
    { href: '/', icon: '📊', label: 'Dashboard' },
    { href: '/inventory', icon: '📦', label: 'Ombor Qoldig\'i' },
  ]},
  { section: 'Boshqarish', items: [
    { href: '/products', icon: '🏷️', label: 'Mahsulotlar' },
    { href: '/recipes', icon: '🍳', label: 'Retseptlar' },
    { href: '/categories', icon: '📁', label: 'Kategoriyalar' },
    { href: '/suppliers', icon: '🚚', label: 'Ta\'minotchilar' },
  ]},
  { section: 'Operatsiyalar', items: [
    { href: '/transactions', icon: '🔄', label: 'Tranzaksiyalar' },
    { href: '/transactions/new', icon: '➕', label: 'Yangi Tranzaksiya' },
    { href: '/debts', icon: '🤝', label: 'Qarzlar' },
  ]},
];

const adminItems = { section: 'Admin', items: [
  { href: '/reports', icon: '💰', label: 'Moliyaviy Hisobot' },
  { href: '/users', icon: '👥', label: 'Foydalanuvchilar' },
]};

const systemItems = { section: 'Tizim', items: [
  { href: '/guide', icon: '📖', label: "Qo'llanma" },
]};

export default function Sidebar() {
  const pathname = usePathname();
  const { user, logout, isAdmin } = useAuth();

  const allNav = isAdmin ? [...navItems, adminItems, systemItems] : [...navItems, systemItems];

  return (
    <aside className="sidebar">
      <div className="sidebar-logo">
        <div className="logo-icon">☕</div>
        <div>
          <h1>Kafe Ombor</h1>
          <span>Boshqarish tizimi</span>
        </div>
      </div>

      <nav className="sidebar-nav">
        {allNav.map(section => (
          <div className="nav-section" key={section.section}>
            <div className="nav-section-title">{section.section}</div>
            {section.items.map(item => (
              <Link key={item.href} href={item.href}
                className={`nav-link ${pathname === item.href ? 'active' : ''}`}>
                <span className="nav-icon">{item.icon}</span>
                {item.label}
              </Link>
            ))}
          </div>
        ))}
      </nav>

      <div className="sidebar-footer">
        <div className="user-info">
          <div className="user-avatar">
            {user?.full_name?.charAt(0) || 'U'}
          </div>
          <div style={{ flex: 1 }}>
            <div className="user-name">{user?.full_name}</div>
            <div className="user-role">{user?.role === 'admin' ? 'Administrator' : 'Ombor mudiri'}</div>
          </div>
          <button className="btn btn-ghost btn-sm" onClick={logout} title="Chiqish">🚪</button>
        </div>
      </div>
    </aside>
  );
}
