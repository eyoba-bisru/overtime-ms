import { createContext, useContext, useState, useEffect, type ReactNode } from 'react';
import type { User, Role } from '../types';
import api from '../api/client';

interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<{ forcePasswordChange: boolean }>;
  logout: () => void;
  setUser: (user: User | null) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const stored = localStorage.getItem('oms_user');
    if (stored) {
      try {
        setUser(JSON.parse(stored));
      } catch {
        localStorage.removeItem('oms_user');
      }
    }
    setLoading(false);
  }, []);

  const login = async (email: string, password: string) => {
    const res = await api.post('/auth/login', { email, password });
    const forcePasswordChange = res.data?.data?.force_password_change ?? false;
    const userData = res.data?.data?.user;

    if (userData) {
      const u: User = {
        ...userData,
        is_blocked: false,
        created_at: '',
        updated_at: '',
      };
      setUser(u);
      localStorage.setItem('oms_user', JSON.stringify(u));
    }

    return { forcePasswordChange };
  };

  const logout = () => {
    document.cookie = 'token=; Max-Age=0; path=/';
    localStorage.removeItem('oms_user');
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, loading, login, logout, setUser }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}
