'use client';
import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { api } from './api';
import { User, LoginResponse } from './types';

interface AuthContextType {
  user: User | null;
  token: string | null;
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
  loading: boolean;
  isAdmin: boolean;
}

const AuthContext = createContext<AuthContextType>({
  user: null, token: null, login: async () => {}, logout: () => {},
  loading: true, isAdmin: false,
});

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const savedToken = localStorage.getItem('token');
    if (savedToken) {
      setToken(savedToken);
      api.get<User>('/api/auth/me')
        .then(u => setUser(u))
        .catch(() => { localStorage.removeItem('token'); })
        .finally(() => setLoading(false));
    } else {
      setLoading(false);
    }
  }, []);

  const login = async (username: string, password: string) => {
    const res = await api.post<LoginResponse>('/api/auth/login', { username, password });
    localStorage.setItem('token', res.token);
    setToken(res.token);
    setUser(res.user);
  };

  const logout = () => {
    localStorage.removeItem('token');
    setToken(null);
    setUser(null);
    window.location.href = '/login';
  };

  return (
    <AuthContext.Provider value={{ user, token, login, logout, loading, isAdmin: user?.role === 'admin' }}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => useContext(AuthContext);
