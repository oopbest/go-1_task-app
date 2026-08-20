// โครงสร้างข้อมูล Task ที่ตรงกับ Go Backend
export interface Task {
  id: number;
  title: string;
  description: string;
  completed: boolean;
  user_id: number;
  created_at: string;
}

// ข้อมูล Pagination Metadata
export interface TasksResponse {
  data: Task[];
  page: number;
  limit: number;
  total_items: number;
  total_pages: number;
}

// User Model (ตรงกับ Go Backend models.User)
export interface User {
  id: number;
  email: string;
  created_at: string;
}

// Auth Response จาก /auth/login หรือ /auth/register
export interface AuthResponse {
  token: string;
  user: User;
}

// Filter Query Parameters
export interface TaskFilters {
  page?: number;
  limit?: number;
  search?: string;
  completed?: boolean;
  sort?: string;
  order?: "asc" | "desc";
}
