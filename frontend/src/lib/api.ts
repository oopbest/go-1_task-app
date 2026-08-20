import Cookies from "js-cookie";
import { AuthResponse, Task, TaskFilters, TasksResponse } from "@/types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

// ฟังก์ชันดึง Token จาก Cookie หรือ LocalStorage
export function getToken(): string | undefined {
  if (typeof window === "undefined") return undefined;
  return (
    Cookies.get("task_token") || localStorage.getItem("task_token") || undefined
  );
}

// ฟังก์ชันเก็บ Token
export function setToken(token: string) {
  Cookies.set("task_token", token, { expires: 7 }); // อยู่ได้ 7 วัน
  localStorage.setItem("task_token", token);
}

// ฟังก์ชันลบ Token (Logout)
export function removeToken() {
  Cookies.remove("task_token");
  localStorage.removeItem("task_token");
}

// Core Fetcher รองรับ Auth Bearer Token อัตโนมัติ
async function apiFetch<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const token = getToken();
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string>),
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers,
  });

  const data = await response.json();

  if (!response.ok) {
    if (response.status === 401 && typeof window !== "undefined") {
      removeToken();
    }
    throw new Error(data.error || "Something went wrong");
  }

  return data as T;
}

// ==========================================
// 🚀 API Endpoints Functions
// ==========================================

export const api = {
  // 1. Auth Endpoints
  auth: {
    register: (body: { email: string; password: string }) =>
      apiFetch<AuthResponse>("/auth/register", {
        method: "POST",
        body: JSON.stringify(body),
      }),
    login: (body: { email: string; password: string }) =>
      apiFetch<AuthResponse>("/auth/login", {
        method: "POST",
        body: JSON.stringify(body),
      }),
  },

  // 2. Tasks Endpoints (CRUD + Pagination + Search)
  tasks: {
    getAll: (filters: TaskFilters = {}) => {
      const params = new URLSearchParams();
      if (filters.page) params.append("page", filters.page.toString());
      if (filters.limit) params.append("limit", filters.limit.toString());
      if (filters.search) params.append("search", filters.search);
      if (filters.completed !== undefined)
        params.append("completed", filters.completed.toString());
      if (filters.sort) params.append("sort", filters.sort);
      if (filters.order) params.append("order", filters.order);

      const queryString = params.toString();
      return apiFetch<TasksResponse>(
        `/tasks${queryString ? `?${queryString}` : ""}`
      );
    },

    getByID: (id: number) => apiFetch<Task>(`/tasks/${id}`),

    create: (body: {
      title: string;
      description?: string;
      completed?: boolean;
    }) =>
      apiFetch<Task>("/tasks", {
        method: "POST",
        body: JSON.stringify(body),
      }),

    update: (
      id: number,
      body: { title: string; description?: string; completed: boolean }
    ) =>
      apiFetch<Task>(`/tasks/${id}`, {
        method: "PUT",
        body: JSON.stringify(body),
      }),

    delete: (id: number) =>
      apiFetch<{ message: string }>(`/tasks/${id}`, {
        method: "DELETE",
      }),
  },

  // 3. Health Check
  health: () =>
    apiFetch<{ status: string; database: string; cache: string }>("/health"),
};
