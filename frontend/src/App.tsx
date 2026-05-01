import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './context/AuthContext';
import { ToastProvider } from './context/ToastContext';
import { ThemeProvider } from './context/ThemeContext';
import ErrorBoundary from './components/ErrorBoundary';
import LoginPage from './pages/LoginPage';
import ChangePasswordPage from './pages/ChangePasswordPage';
import MyOvertimesPage from './pages/MyOvertimesPage';
import CreateOvertimePage from './pages/CreateOvertimePage';
import ReviewPage from './pages/ReviewPage';
import UserManagementPage from './pages/UserManagementPage';
import DepartmentManagementPage from './pages/DepartmentManagementPage';
import EditOvertimePage from './pages/EditOvertimePage';
import type { Role } from './types';

function ProtectedRoute({ children, roles }: { children: React.ReactNode; roles?: Role[] }) {
  const { user, loading } = useAuth();

  if (loading) return <div className="loading-page"><span className="spinner" /></div>;
  if (!user) return <Navigate to="/login" replace />;
  if (user.force_password_change) return <Navigate to="/change-password" replace />;
  if (roles && !roles.includes(user.role)) return <Navigate to="/" replace />;

  return <>{children}</>;
}

function HomeRedirect() {
  const { user } = useAuth();
  if (!user) return <Navigate to="/login" replace />;
  switch (user.role) {
    case 'applicant': return <Navigate to="/overtime/my" replace />;
    case 'checker': return <Navigate to="/overtime/pending" replace />;
    case 'approver': return <Navigate to="/overtime/checked" replace />;
    case 'finance': return <Navigate to="/overtime/approved" replace />;
    case 'admin': return <Navigate to="/admin/users" replace />;
    default: return <Navigate to="/login" replace />;
  }
}

export default function App() {
  return (
    <ErrorBoundary>
      <ThemeProvider>
        <BrowserRouter>
          <ToastProvider>
            <AuthProvider>
              <Routes>
                <Route path="/login" element={<LoginPage />} />
                <Route path="/change-password" element={<ChangePasswordPage />} />

                <Route path="/" element={<ProtectedRoute><HomeRedirect /></ProtectedRoute>} />

                {/* Applicant */}
                <Route path="/overtime/my" element={
                  <ProtectedRoute roles={['applicant', 'admin']}>
                    <MyOvertimesPage />
                  </ProtectedRoute>
                } />
                <Route path="/overtime/create" element={
                  <ProtectedRoute roles={['applicant', 'admin']}>
                    <CreateOvertimePage />
                  </ProtectedRoute>
                } />
                <Route path="/overtime/edit/:id" element={
                  <ProtectedRoute roles={['applicant', 'admin']}>
                    <EditOvertimePage />
                  </ProtectedRoute>
                } />

                {/* Checker */}
                <Route path="/overtime/pending" element={
                  <ProtectedRoute roles={['checker', 'admin']}>
                    <ReviewPage
                      title="Pending Overtime Requests"
                      subtitle="Review and check or reject pending overtime submissions"
                      endpoint="/overtime/pending"
                      emptyIcon="⏳"
                      actions={[
                        { label: '✓ Check', action: 'check', className: 'btn-success', endpoint: (id) => `/overtime/${id}/check` },
                        { label: '✕ Reject', action: 'reject', className: 'btn-danger', endpoint: (id) => `/overtime/${id}/reject` },
                      ]}
                    />
                  </ProtectedRoute>
                } />

                {/* Approver */}
                <Route path="/overtime/checked" element={
                  <ProtectedRoute roles={['approver', 'admin']}>
                    <ReviewPage
                      title="Checked Overtime Requests"
                      subtitle="Review and approve or reject checked overtime submissions"
                      endpoint="/overtime/checked"
                      emptyIcon="✅"
                      actions={[
                        { label: '✓ Approve', action: 'approve', className: 'btn-success', endpoint: (id) => `/overtime/${id}/approve` },
                        { label: '✕ Reject', action: 'reject', className: 'btn-danger', endpoint: (id) => `/overtime/${id}/reject` },
                      ]}
                    />
                  </ProtectedRoute>
                } />

                {/* Finance */}
                <Route path="/overtime/approved" element={
                  <ProtectedRoute roles={['finance', 'admin']}>
                    <ReviewPage
                      title="Approved Overtime Records"
                      subtitle="View approved overtime records for payment processing"
                      endpoint="/overtime/approved"
                      emptyIcon="💰"
                      actions={[]}
                    />
                  </ProtectedRoute>
                } />

                {/* Admin */}
                <Route path="/admin/users" element={
                  <ProtectedRoute roles={['admin']}>
                    <UserManagementPage />
                  </ProtectedRoute>
                } />
                <Route path="/admin/departments" element={
                  <ProtectedRoute roles={['admin']}>
                    <DepartmentManagementPage />
                  </ProtectedRoute>
                } />
                <Route path="/admin/overtime" element={
                  <ProtectedRoute roles={['admin']}>
                    <ReviewPage
                      title="All Overtime Requests"
                      subtitle="View all overtime records system-wide"
                      endpoint="/admin/overtime"
                      emptyIcon="📊"
                      actions={[]}
                    />
                  </ProtectedRoute>
                } />

                <Route path="*" element={<Navigate to="/" replace />} />
              </Routes>
            </AuthProvider>
          </ToastProvider>
        </BrowserRouter>
      </ThemeProvider>
    </ErrorBoundary>
  );
}
