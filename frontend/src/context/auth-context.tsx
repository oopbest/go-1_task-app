'use client';

import React, { createContext, useContext, useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { User } from '@/types';
import { getToken, removeToken, setToken } from '@/lib/api';

interface AuthContextType {
  user: User | null;
  token: string | null;
  login: (token: string, user: User) => void;
  logout: () => void;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setTokenState] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    // โหลดข้อมูล User และ Token เมื่อ Component Mount ฝั่ง Client
    const initAuth = () => {
      const savedToken = getToken();
      const savedUser = localStorage.getItem('task_user');

      if (savedToken && savedUser) {
        try {
          setUser(JSON.parse(savedUser));
          setTokenState(savedToken);
        } catch {
          removeToken();
          localStorage.removeItem('task_user');
        }
      }
      setIsLoading(false);
    };

    initAuth();
  }, []);

  const login = (newToken: string, newUser: User) => {
    setToken(newToken);
    localStorage.setItem('task_user', JSON.stringify(newUser));
    setTokenState(newToken);
    setUser(newUser);
    router.push('/dashboard');
  };

  const logout = () => {
    removeToken();
    localStorage.removeItem('task_user');
    setTokenState(null);
    setUser(null);
    router.push('/login');
  };

  return (
    <AuthContext.Provider value={{ user, token, login, logout, isLoading }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
