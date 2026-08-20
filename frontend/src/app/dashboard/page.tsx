'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/auth-context';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api';
import { Task, TaskFilters, TasksResponse } from '@/types';
import { toast } from 'sonner';
import {
  CheckCircle2,
  Circle,
  Plus,
  Trash2,
  Edit3,
  Search,
  LogOut,
  Zap,
  ShieldCheck,
  Server,
  Layers,
  Sparkles,
  ChevronLeft,
  ChevronRight,
  ArrowUpDown,
  X,
  Loader2,
  Calendar,
  User as UserIcon,
} from 'lucide-react';

export default function DashboardPage() {
  const { user, logout, isLoading: isAuthLoading } = useAuth();
  const router = useRouter();
  const queryClient = useQueryClient();

  // Filters State
  const [page, setPage] = useState(1);
  const [limit] = useState(6);
  const [search, setSearch] = useState('');
  const [debouncedSearch, setDebouncedSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState<'all' | 'completed' | 'pending'>('all');
  const [sortBy, setSortBy] = useState('created_at');
  const [order, setOrder] = useState<'asc' | 'desc'>('desc');

  // Modals State
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | null>(null);
  const [taskTitle, setTaskTitle] = useState('');
  const [taskDescription, setTaskDescription] = useState('');

  // Debounce Search
  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedSearch(search);
      setPage(1);
    }, 300);
    return () => clearTimeout(handler);
  }, [search]);

  // Protect Route
  useEffect(() => {
    if (!isAuthLoading && !user) {
      router.replace('/login');
    }
  }, [user, isAuthLoading, router]);

  // Construct Query Filters
  const queryFilters: TaskFilters = {
    page,
    limit,
    search: debouncedSearch || undefined,
    completed: statusFilter === 'all' ? undefined : statusFilter === 'completed',
    sort: sortBy,
    order,
  };

  // 1. Fetch Tasks with React Query
  const { data: tasksData, isLoading: isTasksLoading } = useQuery({
    queryKey: ['tasks', queryFilters],
    queryFn: () => api.tasks.getAll(queryFilters),
    enabled: !!user,
  });

  // 2. Optimistic Mutation: Toggle Complete
  const toggleMutation = useMutation({
    mutationFn: (task: Task) =>
      api.tasks.update(task.id, {
        title: task.title,
        description: task.description,
        completed: !task.completed,
      }),
    onMutate: async (toggledTask) => {
      await queryClient.cancelQueries({ queryKey: ['tasks'] });
      const previousData = queryClient.getQueryData<TasksResponse>(['tasks', queryFilters]);

      queryClient.setQueryData<TasksResponse>(['tasks', queryFilters], (old) => {
        if (!old) return old;
        return {
          ...old,
          data: old.data.map((t: Task) =>
            t.id === toggledTask.id ? { ...t, completed: !t.completed } : t
          ),
        };
      });

      return { previousData };
    },
    onError: (_err, _variables, context) => {
      if (context?.previousData) {
        queryClient.setQueryData(['tasks', queryFilters], context.previousData);
      }
      toast.error('ไม่สามารถเปลี่ยนสถานะงานได้');
    },
    onSuccess: (updated) => {
      toast.success(updated.completed ? '🎉 ทำงานเสร็จสมบูรณ์!' : 'เปิดงานอีกครั้ง');
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
    },
  });

  // 3. Mutation: Create Task
  const createMutation = useMutation({
    mutationFn: (body: { title: string; description?: string }) => api.tasks.create(body),
    onSuccess: () => {
      toast.success('สร้างงานใหม่สำเร็จ!', { description: 'งานถูกส่งเข้า Go Background Worker Pool เรียบร้อย' });
      setIsCreateOpen(false);
      setTaskTitle('');
      setTaskDescription('');
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
    },
    onError: (err: Error) => toast.error('สร้างงานไม่สำเร็จ', { description: err.message }),
  });

  // 4. Mutation: Edit Task
  const editMutation = useMutation({
    mutationFn: ({ id, body }: { id: number; body: { title: string; description?: string; completed: boolean } }) =>
      api.tasks.update(id, body),
    onSuccess: () => {
      toast.success('อัปเดตงานสำเร็จ!');
      setEditingTask(null);
      setTaskTitle('');
      setTaskDescription('');
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
    },
    onError: (err: Error) => toast.error('อัปเดตงานไม่สำเร็จ', { description: err.message }),
  });

  // 5. Mutation: Delete Task
  const deleteMutation = useMutation({
    mutationFn: (id: number) => api.tasks.delete(id),
    onSuccess: () => {
      toast.success('ลบงานเรียบร้อยแล้ว');
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
    },
    onError: (err: Error) => toast.error('ลบงานไม่สำเร็จ', { description: err.message }),
  });

  if (isAuthLoading || !user) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-slate-950">
        <Loader2 className="w-10 h-10 animate-spin text-blue-500" />
      </div>
    );
  }

  const tasks = tasksData?.data || [];
  const totalPages = tasksData?.total_pages || 1;
  const totalItems = tasksData?.total_items || 0;

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 pb-16">
      {/* Top Navbar */}
      <header className="sticky top-0 z-30 border-b border-slate-800/80 bg-slate-950/80 backdrop-blur-xl">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <div className="w-10 h-10 rounded-xl bg-gradient-to-tr from-blue-600 to-emerald-500 flex items-center justify-center shadow-lg shadow-blue-500/20">
              <CheckCircle2 className="w-6 h-6 text-white" />
            </div>
            <div>
              <span className="font-bold text-lg text-white tracking-tight">TaskFlow SaaS</span>
              <span className="hidden sm:inline-block ml-2 px-2 py-0.5 text-[10px] font-semibold tracking-wide bg-blue-500/10 text-blue-400 border border-blue-500/20 rounded-full">
                Go API + Redis + gRPC
              </span>
            </div>
          </div>

          <div className="flex items-center space-x-4">
            <div className="hidden md:flex items-center space-x-2 text-xs text-slate-400 bg-slate-900 border border-slate-800 px-3 py-1.5 rounded-lg">
              <UserIcon className="w-3.5 h-3.5 text-blue-400" />
              <span>{user.email}</span>
              <span className="text-slate-600">|</span>
              <span className="text-emerald-400">ID #{user.id}</span>
            </div>

            <button
              onClick={logout}
              className="flex items-center space-x-1.5 text-xs text-slate-400 hover:text-red-400 bg-slate-900/60 hover:bg-red-500/10 border border-slate-800 hover:border-red-500/20 px-3 py-2 rounded-xl transition-all cursor-pointer"
            >
              <LogOut className="w-3.5 h-3.5" />
              <span className="hidden sm:inline">ออกจากระบบ</span>
            </button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-8">
        {/* Architecture Status Banner */}
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-8">
          <div className="bg-slate-900/60 border border-slate-800/80 rounded-xl p-3.5 flex items-center space-x-3">
            <div className="w-8 h-8 rounded-lg bg-blue-500/10 text-blue-400 flex items-center justify-center">
              <Server className="w-4 h-4" />
            </div>
            <div>
              <p className="text-[11px] text-slate-400">Backend API</p>
              <p className="text-xs font-semibold text-white">Gin REST :8080</p>
            </div>
          </div>

          <div className="bg-slate-900/60 border border-slate-800/80 rounded-xl p-3.5 flex items-center space-x-3">
            <div className="w-8 h-8 rounded-lg bg-amber-500/10 text-amber-400 flex items-center justify-center">
              <Zap className="w-4 h-4" />
            </div>
            <div>
              <p className="text-[11px] text-slate-400">Redis Cache</p>
              <p className="text-xs font-semibold text-white">~500 µs Latency</p>
            </div>
          </div>

          <div className="bg-slate-900/60 border border-slate-800/80 rounded-xl p-3.5 flex items-center space-x-3">
            <div className="w-8 h-8 rounded-lg bg-emerald-500/10 text-emerald-400 flex items-center justify-center">
              <Layers className="w-4 h-4" />
            </div>
            <div>
              <p className="text-[11px] text-slate-400">gRPC Microservice</p>
              <p className="text-xs font-semibold text-white">Protobuf :50051</p>
            </div>
          </div>

          <div className="bg-slate-900/60 border border-slate-800/80 rounded-xl p-3.5 flex items-center space-x-3">
            <div className="w-8 h-8 rounded-lg bg-purple-500/10 text-purple-400 flex items-center justify-center">
              <ShieldCheck className="w-4 h-4" />
            </div>
            <div>
              <p className="text-[11px] text-slate-400">Total Tasks</p>
              <p className="text-xs font-semibold text-white">{totalItems} งานทั้งหมด</p>
            </div>
          </div>
        </div>

        {/* Action Header & Search Filter Bar */}
        <div className="bg-slate-900/80 border border-slate-800/80 rounded-2xl p-4 sm:p-5 shadow-xl mb-6">
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
            {/* Search Input */}
            <div className="relative flex-1">
              <Search className="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400" />
              <input
                type="text"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                placeholder="ค้นหางาน (Search task title & description)..."
                className="w-full pl-10 pr-4 py-2 bg-slate-950/80 border border-slate-800 rounded-xl text-sm text-slate-100 placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 transition-all"
              />
              {search && (
                <button
                  onClick={() => setSearch('')}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-white"
                >
                  <X className="w-3.5 h-3.5" />
                </button>
              )}
            </div>

            {/* Filter Controls & Add Task Button */}
            <div className="flex flex-wrap items-center gap-2.5">
              {/* Status Filter */}
              <div className="flex items-center bg-slate-950/80 border border-slate-800 rounded-xl p-1 text-xs">
                <button
                  onClick={() => { setStatusFilter('all'); setPage(1); }}
                  className={`px-3 py-1.5 rounded-lg transition-all ${
                    statusFilter === 'all' ? 'bg-blue-600 text-white font-medium shadow' : 'text-slate-400 hover:text-white'
                  }`}
                >
                  ทั้งหมด
                </button>
                <button
                  onClick={() => { setStatusFilter('pending'); setPage(1); }}
                  className={`px-3 py-1.5 rounded-lg transition-all ${
                    statusFilter === 'pending' ? 'bg-amber-600 text-white font-medium shadow' : 'text-slate-400 hover:text-white'
                  }`}
                >
                  รอดำเนินการ
                </button>
                <button
                  onClick={() => { setStatusFilter('completed'); setPage(1); }}
                  className={`px-3 py-1.5 rounded-lg transition-all ${
                    statusFilter === 'completed' ? 'bg-emerald-600 text-white font-medium shadow' : 'text-slate-400 hover:text-white'
                  }`}
                >
                  เสร็จแล้ว
                </button>
              </div>

              {/* Sort Field & Order Toggle */}
              <button
                onClick={() => {
                  setSortBy(sortBy === 'created_at' ? 'title' : 'created_at');
                  setOrder(order === 'asc' ? 'desc' : 'asc');
                }}
                className="flex items-center space-x-1.5 text-xs bg-slate-950/80 border border-slate-800 text-slate-300 hover:text-white px-3 py-2 rounded-xl transition-colors cursor-pointer"
                title="สลับลำดับ เรียงเก่า-ใหม่"
              >
                <ArrowUpDown className="w-3.5 h-3.5 text-blue-400" />
                <span>{order === 'desc' ? 'ใหม่ ➡️ เก่า' : 'เก่า ➡️ ใหม่'}</span>
              </button>

              {/* Create Task Button */}
              <button
                onClick={() => {
                  setTaskTitle('');
                  setTaskDescription('');
                  setIsCreateOpen(true);
                }}
                className="flex items-center space-x-1.5 text-xs font-semibold bg-gradient-to-r from-blue-600 to-emerald-600 hover:from-blue-500 hover:to-emerald-500 text-white px-4 py-2 rounded-xl shadow-lg shadow-blue-500/20 transition-all cursor-pointer"
              >
                <Plus className="w-4 h-4" />
                <span>เพิ่มงานใหม่</span>
              </button>
            </div>
          </div>
        </div>

        {/* Task Cards Grid */}
        {isTasksLoading ? (
          <div className="py-20 flex flex-col items-center justify-center space-y-3 text-slate-400">
            <Loader2 className="w-8 h-8 animate-spin text-blue-500" />
            <p className="text-sm">กำลังดึงข้อมูลจาก Redis Cache & PostgreSQL...</p>
          </div>
        ) : tasks.length === 0 ? (
          <div className="bg-slate-900/40 border border-slate-800/60 rounded-2xl py-16 px-4 text-center">
            <div className="w-12 h-12 rounded-2xl bg-slate-800/80 text-slate-400 flex items-center justify-center mx-auto mb-4">
              <Sparkles className="w-6 h-6" />
            </div>
            <h3 className="text-base font-semibold text-white">ไม่พบรายการงาน</h3>
            <p className="text-xs text-slate-400 mt-1 max-w-sm mx-auto">
              {search ? 'ลองค้นหาด้วยคำค้นอื่น หรือล้างตัวกรอง' : 'กดปุ่ม "เพิ่มงานใหม่" ด้านบนเพื่อเริ่มสร้างงานแรกของคุณ'}
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {tasks.map((task) => (
              <div
                key={task.id}
                className={`group relative bg-slate-900/80 border rounded-2xl p-5 shadow-lg transition-all duration-200 hover:-translate-y-0.5 ${
                  task.completed
                    ? 'border-emerald-500/20 bg-emerald-950/5'
                    : 'border-slate-800 hover:border-slate-700'
                }`}
              >
                {/* Header: Checkbox & Actions */}
                <div className="flex items-start justify-between gap-3">
                  <button
                    onClick={() => toggleMutation.mutate(task)}
                    className="mt-0.5 text-slate-400 hover:text-emerald-400 transition-colors cursor-pointer flex-shrink-0"
                  >
                    {task.completed ? (
                      <CheckCircle2 className="w-5 h-5 text-emerald-400 fill-emerald-400/10" />
                    ) : (
                      <Circle className="w-5 h-5 hover:text-blue-400" />
                    )}
                  </button>

                  <div className="flex-1 min-w-0">
                    <h4
                      className={`text-sm font-semibold truncate ${
                        task.completed ? 'line-through text-slate-400' : 'text-slate-100'
                      }`}
                    >
                      {task.title}
                    </h4>
                    {task.description && (
                      <p className="text-xs text-slate-400 mt-1.5 line-clamp-2 leading-relaxed">
                        {task.description}
                      </p>
                    )}
                  </div>

                  {/* Actions (Edit / Delete) */}
                  <div className="flex items-center space-x-1 opacity-80 group-hover:opacity-100 transition-opacity">
                    <button
                      onClick={() => {
                        setEditingTask(task);
                        setTaskTitle(task.title);
                        setTaskDescription(task.description || '');
                      }}
                      className="p-1.5 text-slate-400 hover:text-blue-400 hover:bg-blue-500/10 rounded-lg transition-colors cursor-pointer"
                      title="แก้ไขงาน"
                    >
                      <Edit3 className="w-3.5 h-3.5" />
                    </button>
                    <button
                      onClick={() => {
                        if (confirm(`ยืนยันการลบงาน: "${task.title}" ?`)) {
                          deleteMutation.mutate(task.id);
                        }
                      }}
                      className="p-1.5 text-slate-400 hover:text-red-400 hover:bg-red-500/10 rounded-lg transition-colors cursor-pointer"
                      title="ลบงาน"
                    >
                      <Trash2 className="w-3.5 h-3.5" />
                    </button>
                  </div>
                </div>

                {/* Footer: ID & Date Badge */}
                <div className="mt-4 pt-3.5 border-t border-slate-800/60 flex items-center justify-between text-[11px] text-slate-500">
                  <span className="font-mono text-slate-500">#Task {task.id}</span>
                  <div className="flex items-center space-x-1 text-slate-400">
                    <Calendar className="w-3 h-3 text-slate-500" />
                    <span>{new Date(task.created_at).toLocaleDateString('th-TH')}</span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Pagination Controls */}
        {totalPages > 1 && (
          <div className="mt-8 flex items-center justify-between bg-slate-900/60 border border-slate-800/80 rounded-2xl px-5 py-3 text-xs">
            <span className="text-slate-400">
              หน้า <strong className="text-white">{page}</strong> จาก <strong>{totalPages}</strong> (ทั้งหมด {totalItems} งาน)
            </span>

            <div className="flex items-center space-x-2">
              <button
                disabled={page <= 1}
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                className="flex items-center space-x-1 px-3 py-1.5 bg-slate-950 border border-slate-800 rounded-lg text-slate-300 hover:text-white disabled:opacity-30 disabled:cursor-not-allowed transition-all"
              >
                <ChevronLeft className="w-3.5 h-3.5" />
                <span>ก่อนหน้า</span>
              </button>

              <button
                disabled={page >= totalPages}
                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                className="flex items-center space-x-1 px-3 py-1.5 bg-slate-950 border border-slate-800 rounded-lg text-slate-300 hover:text-white disabled:opacity-30 disabled:cursor-not-allowed transition-all"
              >
                <span>ถัดไป</span>
                <ChevronRight className="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        )}
      </main>

      {/* Modal Dialog: Create / Edit Task */}
      {(isCreateOpen || editingTask) && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm animate-in fade-in duration-150">
          <div className="bg-slate-900 border border-slate-800 rounded-2xl w-full max-w-md p-6 shadow-2xl">
            <div className="flex items-center justify-between mb-5">
              <h3 className="text-base font-bold text-white flex items-center space-x-2">
                <CheckCircle2 className="w-5 h-5 text-blue-400" />
                <span>{editingTask ? 'แก้ไขงาน' : 'เพิ่มงานใหม่'}</span>
              </h3>
              <button
                onClick={() => {
                  setIsCreateOpen(false);
                  setEditingTask(null);
                }}
                className="text-slate-400 hover:text-white p-1 rounded-lg hover:bg-slate-800 transition-colors"
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            <form
              onSubmit={(e) => {
                e.preventDefault();
                if (!taskTitle.trim()) {
                  toast.error('กรุณากรอกชื่องาน');
                  return;
                }
                if (editingTask) {
                  editMutation.mutate({
                    id: editingTask.id,
                    body: { title: taskTitle, description: taskDescription, completed: editingTask.completed },
                  });
                } else {
                  createMutation.mutate({ title: taskTitle, description: taskDescription });
                }
              }}
              className="space-y-4"
            >
              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-2">
                  ชื่องาน (Title) *
                </label>
                <input
                  type="text"
                  value={taskTitle}
                  onChange={(e) => setTaskTitle(e.target.value)}
                  placeholder="เช่น พัฒนาฟีเจอร์ Next.js Dashboard"
                  className="w-full px-3.5 py-2.5 bg-slate-950 border border-slate-800 rounded-xl text-slate-100 placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 transition-all text-sm"
                  required
                  autoFocus
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-2">
                  รายละเอียด (Description)
                </label>
                <textarea
                  value={taskDescription}
                  onChange={(e) => setTaskDescription(e.target.value)}
                  placeholder="รายละเอียดเพิ่มเติมของงาน..."
                  rows={3}
                  className="w-full px-3.5 py-2.5 bg-slate-950 border border-slate-800 rounded-xl text-slate-100 placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 transition-all text-sm resize-none"
                />
              </div>

              <div className="flex items-center justify-end space-x-2.5 pt-3">
                <button
                  type="button"
                  onClick={() => {
                    setIsCreateOpen(false);
                    setEditingTask(null);
                  }}
                  className="px-4 py-2 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-xl text-xs font-medium transition-colors"
                >
                  ยกเลิก
                </button>
                <button
                  type="submit"
                  disabled={createMutation.isPending || editMutation.isPending}
                  className="px-4 py-2 bg-gradient-to-r from-blue-600 to-emerald-600 hover:from-blue-500 hover:to-emerald-500 text-white rounded-xl text-xs font-semibold shadow-lg shadow-blue-500/20 transition-all flex items-center space-x-1.5 disabled:opacity-50"
                >
                  {(createMutation.isPending || editMutation.isPending) && (
                    <Loader2 className="w-3.5 h-3.5 animate-spin" />
                  )}
                  <span>{editingTask ? 'บันทึกการแก้ไข' : 'บันทึกงาน'}</span>
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
