'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { useAuth } from '@/context/auth-context';
import { api } from '@/lib/api';
import { toast } from 'sonner';
import { CheckCircle2, Lock, Mail, ArrowRight, Loader2, ShieldCheck, Zap } from 'lucide-react';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const { login } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email.trim() || !password.trim()) {
      toast.error('กรุณากรอก Email และ Password');
      return;
    }

    setIsLoading(true);
    try {
      const res = await api.auth.login({ email, password });
      toast.success('เข้าสู่ระบบสำเร็จ!', {
        description: `ยินดีต้อนรับคุณ ${res.user.email}`,
      });
      login(res.token, res.user);
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Email หรือ Password ไม่ถูกต้อง';
      toast.error('เข้าสู่ระบบไม่สำเร็จ', {
        description: errorMsg,
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-950 px-4 relative overflow-hidden">
      {/* Background Gradient Orbs */}
      <div className="absolute -top-40 -left-40 w-96 h-96 bg-blue-600/20 rounded-full blur-3xl pointer-events-none" />
      <div className="absolute -bottom-40 -right-40 w-96 h-96 bg-emerald-600/20 rounded-full blur-3xl pointer-events-none" />

      <div className="w-full max-w-md z-10">
        {/* Logo & Header */}
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-tr from-blue-600 to-emerald-500 shadow-xl shadow-blue-500/20 mb-4">
            <CheckCircle2 className="w-8 h-8 text-white" />
          </div>
          <h1 className="text-3xl font-bold tracking-tight text-white">TaskFlow SaaS</h1>
          <p className="text-sm text-slate-400 mt-2">
            Enterprise Task Management powered by <span className="text-blue-400 font-medium">Go 1.25</span> &{' '}
            <span className="text-emerald-400 font-medium">Next.js 15</span>
          </p>
        </div>

        {/* Card */}
        <div className="bg-slate-900/80 backdrop-blur-xl border border-slate-800/80 rounded-2xl p-8 shadow-2xl">
          <form onSubmit={handleSubmit} className="space-y-5">
            <div>
              <label className="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-2">
                Email Address
              </label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400">
                  <Mail className="w-4 h-4" />
                </div>
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="name@example.com"
                  className="w-full pl-10 pr-4 py-2.5 bg-slate-950/60 border border-slate-800 rounded-xl text-slate-100 placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 transition-all text-sm"
                  required
                />
              </div>
            </div>

            <div>
              <label className="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-2">
                Password
              </label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400">
                  <Lock className="w-4 h-4" />
                </div>
                <input
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="••••••••"
                  className="w-full pl-10 pr-4 py-2.5 bg-slate-950/60 border border-slate-800 rounded-xl text-slate-100 placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 transition-all text-sm"
                  required
                />
              </div>
            </div>

            <button
              type="submit"
              disabled={isLoading}
              className="w-full py-3 px-4 bg-gradient-to-r from-blue-600 to-emerald-600 hover:from-blue-500 hover:to-emerald-500 text-white font-medium rounded-xl shadow-lg shadow-blue-500/25 transition-all duration-200 flex items-center justify-center space-x-2 disabled:opacity-50 disabled:cursor-not-allowed group cursor-pointer"
            >
              {isLoading ? (
                <>
                  <Loader2 className="w-4 h-4 animate-spin" />
                  <span>กำลังตรวจสอบ...</span>
                </>
              ) : (
                <>
                  <span>เข้าสู่ระบบ (Sign In)</span>
                  <ArrowRight className="w-4 h-4 group-hover:translate-x-1 transition-transform" />
                </>
              )}
            </button>
          </form>

          {/* Quick Demo Credentials */}
          <div className="mt-6 pt-6 border-t border-slate-800/80">
            <p className="text-xs text-slate-400 mb-2 font-medium">💡 บัญชีทดสอบที่เคยสมัครไว้ใน Go API:</p>
            <button
              type="button"
              onClick={() => {
                setEmail('test@example.com');
                setPassword('password123');
              }}
              className="w-full text-left p-2.5 rounded-lg bg-slate-950/40 border border-slate-800 text-xs text-slate-300 hover:bg-slate-800/50 transition-colors flex items-center justify-between"
            >
              <span>Email: <strong className="text-blue-400">test@example.com</strong> / Pass: <strong className="text-slate-400">password123</strong></span>
              <span className="text-[10px] text-blue-400 underline">คลิกเพื่อกรอก</span>
            </button>
          </div>

          <div className="mt-6 text-center text-xs text-slate-400">
            ยังไม่มีบัญชีใช้งาน?{' '}
            <Link href="/register" className="text-blue-400 hover:text-blue-300 font-semibold underline">
              สมัครสมาชิกที่นี่ (Register)
            </Link>
          </div>
        </div>

        {/* Feature Badges */}
        <div className="flex items-center justify-center space-x-6 mt-8 text-xs text-slate-500">
          <div className="flex items-center space-x-1.5">
            <Zap className="w-3.5 h-3.5 text-amber-400" />
            <span>Redis Cache (~500µs)</span>
          </div>
          <div className="flex items-center space-x-1.5">
            <ShieldCheck className="w-3.5 h-3.5 text-emerald-400" />
            <span>JWT Multi-Tenant</span>
          </div>
        </div>
      </div>
    </div>
  );
}
