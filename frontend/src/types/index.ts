export type Role = 'admin' | 'applicant' | 'checker' | 'approver' | 'finance';
export type OvertimeStatus = 'pending' | 'checked' | 'approved' | 'rejected';
export type OvertimeProgram = 'night' | 'weekend' | 'holiday';

export interface User {
  id: string;
  email: string;
  name: string;
  role: Role;
  is_blocked: boolean;
  force_password_change: boolean;
  created_at: string;
  updated_at: string;
}

export interface Overtime {
  id: string;
  user_id: string;
  date: string;
  start_time: string;
  end_time: string;
  job_done: string;
  program: OvertimeProgram;
  status: OvertimeStatus;
  duration: number;
  created_at: string;
  updated_at: string;
}

export interface PaginationMeta {
  page: number;
  page_size: number;
  total: number;
  total_pages: number;
}

export interface APIResponse<T = unknown> {
  success: boolean;
  message?: string;
  data?: T;
  error?: string;
  meta?: PaginationMeta;
}
