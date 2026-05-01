import { useState, useEffect, useMemo } from 'react';
import api from '../api/client';
import Layout from '../components/Layout';
import { useToast } from '../context/ToastContext';
import ConfirmModal from '../components/ConfirmModal';
import type { User, Role, Department } from '../types';

const ROLES: Role[] = ['applicant', 'checker', 'approver', 'finance', 'admin'];

export default function UserManagementPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [departments, setDepartments] = useState<Department[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState('');
  const [showCreate, setShowCreate] = useState(false);
  const [createForm, setCreateForm] = useState({ email: '', name: '', password: '', role: 'applicant' as Role, department_id: '' });
  const [createLoading, setCreateLoading] = useState(false);
  const [editUser, setEditUser] = useState<User | null>(null);
  const [editLoading, setEditLoading] = useState(false);
  const [tempPassword, setTempPassword] = useState<{ userId: string; password: string } | null>(null);
  const { showToast } = useToast();

  // Confirm modal state
  const [confirm, setConfirm] = useState<{
    title: string; message: string; variant: 'danger' | 'warning'; confirmLabel: string;
    action: () => Promise<void>;
  } | null>(null);
  const [confirmLoading, setConfirmLoading] = useState(false);

  const fetchUsers = async () => {
    setLoading(true);
    try {
      const res = await api.get('/admin/users');
      setUsers(res.data.data || []);
    } catch { /* handled */ }
    setLoading(false);
  };

  const fetchDepartments = async () => {
    try {
      const res = await api.get('/admin/departments');
      const depts = res.data.data || [];
      setDepartments(depts);
      if (depts.length > 0) {
        setCreateForm(p => ({ ...p, department_id: depts[0].id }));
      }
    } catch { /* handled */ }
  };

  useEffect(() => {
    fetchUsers();
    fetchDepartments();
  }, []);

  const filteredUsers = useMemo(() => {
    if (!search.trim()) return users;
    const q = search.toLowerCase();
    return users.filter(u =>
      u.name.toLowerCase().includes(q) ||
      u.email.toLowerCase().includes(q) ||
      u.role.toLowerCase().includes(q)
    );
  }, [users, search]);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setCreateLoading(true);
    try {
      await api.post('/admin/users', createForm);
      showToast({ type: 'success', title: 'User Created', message: 'The new user has been created successfully.' });
      setShowCreate(false);
      setCreateForm({ email: '', name: '', password: '', role: 'applicant', department_id: departments[0]?.id || '' });
      await fetchUsers();
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || 'Failed to create user';
      showToast({ type: 'error', title: 'Creation Failed', message: msg });
    }
    setCreateLoading(false);
  };

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editUser) return;
    setEditLoading(true);
    try {
      await api.patch(`/admin/users/${editUser.id}`, {
        email: editUser.email,
        name: editUser.name,
        role: editUser.role,
        department_id: editUser.department_id
      });
      showToast({ type: 'success', title: 'User Updated', message: 'User profile has been updated successfully.' });
      setEditUser(null);
      await fetchUsers();
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || 'Failed to update user';
      showToast({ type: 'error', title: 'Update Failed', message: msg });
    }
    setEditLoading(false);
  };

  const handleConfirmAction = async () => {
    if (!confirm) return;
    setConfirmLoading(true);
    try {
      await confirm.action();
    } finally {
      setConfirmLoading(false);
      setConfirm(null);
    }
  };

  const requestBlock = (id: string, isBlocked: boolean) => {
    const action = isBlocked ? 'block' : 'unblock';
    setConfirm({
      title: `${isBlocked ? 'Block' : 'Unblock'} User`,
      message: `Are you sure you want to ${action} this user?`,
      variant: 'warning',
      confirmLabel: isBlocked ? 'Block' : 'Unblock',
      action: async () => {
        await api.patch(`/admin/users/${id}/block`, { is_blocked: isBlocked });
        showToast({ type: 'info', title: 'Status Changed', message: isBlocked ? 'User has been blocked.' : 'User has been unblocked.' });
        await fetchUsers();
      },
    });
  };

  const requestResetPassword = (id: string) => {
    setConfirm({
      title: 'Reset Password',
      message: 'Are you sure you want to reset this user\'s password? They will be given a temporary password.',
      variant: 'warning',
      confirmLabel: 'Reset Password',
      action: async () => {
        const res = await api.patch(`/admin/users/${id}/reset-password`);
        showToast({ type: 'success', title: 'Password Reset', message: 'Temporary password generated.' });
        setTempPassword({ userId: id, password: res.data.data?.temporary_password || '' });
      },
    });
  };

  const requestDelete = (id: string) => {
    setConfirm({
      title: 'Delete User',
      message: 'Are you sure you want to permanently delete this user? This action cannot be undone.',
      variant: 'danger',
      confirmLabel: 'Delete',
      action: async () => {
        await api.delete(`/admin/users/${id}`);
        showToast({ type: 'warning', title: 'User Deleted', message: 'User account has been removed.' });
        await fetchUsers();
      },
    });
  };

  return (
    <Layout>
      <div className="page-header flex-between">
        <div>
          <h1 className="page-title">User Management</h1>
          <p className="page-subtitle">Manage system users and their roles</p>
        </div>
        <button className="btn btn-primary" onClick={() => setShowCreate(true)}>➕ Create User</button>
      </div>

      {/* Temp password display */}
      {tempPassword && (
        <div className="alert alert-success" style={{ marginBottom: 16 }}>
          🔑 Temporary password for user: <strong style={{ fontFamily: 'monospace', marginLeft: 8 }}>{tempPassword.password}</strong>
          <button className="btn btn-ghost btn-sm" style={{ marginLeft: 'auto' }} onClick={() => setTempPassword(null)}>✕</button>
        </div>
      )}

      {/* Create user modal */}
      {showCreate && (
        <div className="modal-overlay" onClick={() => setShowCreate(false)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h3 className="modal-title">Create New User</h3>
              <button className="modal-close" onClick={() => setShowCreate(false)}>✕</button>
            </div>
            <form onSubmit={handleCreate}>
              <div className="modal-body">
                <div className="form-group">
                  <label className="form-label">Name</label>
                  <input className="form-input" value={createForm.name} onChange={e => setCreateForm(p => ({ ...p, name: e.target.value }))} required />
                </div>
                <div className="form-group">
                  <label className="form-label">Email</label>
                  <input type="email" className="form-input" value={createForm.email} onChange={e => setCreateForm(p => ({ ...p, email: e.target.value }))} required />
                </div>
                <div className="form-group">
                  <label className="form-label">Password</label>
                  <input type="password" className="form-input" value={createForm.password} onChange={e => setCreateForm(p => ({ ...p, password: e.target.value }))} required />
                </div>
                <div className="form-group">
                  <label className="form-label">Role</label>
                  <select className="form-select" value={createForm.role} onChange={e => setCreateForm(p => ({ ...p, role: e.target.value as Role }))}>
                    {ROLES.map(r => <option key={r} value={r}>{r}</option>)}
                  </select>
                </div>
                <div className="form-group">
                  <label className="form-label">Department</label>
                  <select className="form-select" value={createForm.department_id} onChange={e => setCreateForm(p => ({ ...p, department_id: e.target.value }))}>
                    {departments.map(d => <option key={d.id} value={d.id}>{d.name}</option>)}
                  </select>
                </div>
              </div>
              <div className="modal-footer">
                <button type="button" className="btn btn-ghost" onClick={() => setShowCreate(false)}>Cancel</button>
                <button type="submit" className="btn btn-primary" disabled={createLoading}>
                  {createLoading ? <span className="spinner" /> : 'Create User'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Edit user modal */}
      {editUser && (
        <div className="modal-overlay" onClick={() => setEditUser(null)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h3 className="modal-title">Edit User Profile</h3>
              <button className="modal-close" onClick={() => setEditUser(null)}>✕</button>
            </div>
            <form onSubmit={handleUpdate}>
              <div className="modal-body">
                <div className="form-group">
                  <label className="form-label">Name</label>
                  <input className="form-input" value={editUser.name} onChange={e => setEditUser(p => p ? ({ ...p, name: e.target.value }) : null)} required />
                </div>
                <div className="form-group">
                  <label className="form-label">Email</label>
                  <input type="email" className="form-input" value={editUser.email} onChange={e => setEditUser(p => p ? ({ ...p, email: e.target.value }) : null)} required />
                </div>
                <div className="form-group">
                  <label className="form-label">Role</label>
                  <select className="form-select" value={editUser.role} onChange={e => setEditUser(p => p ? ({ ...p, role: e.target.value as Role }) : null)}>
                    {ROLES.map(r => <option key={r} value={r}>{r}</option>)}
                  </select>
                </div>
                <div className="form-group">
                  <label className="form-label">Department</label>
                  <select className="form-select" value={editUser.department_id} onChange={e => setEditUser(p => p ? ({ ...p, department_id: e.target.value }) : null)}>
                    {departments.map(d => <option key={d.id} value={d.id}>{d.name}</option>)}
                  </select>
                </div>
              </div>
              <div className="modal-footer">
                <button type="button" className="btn btn-ghost" onClick={() => setEditUser(null)}>Cancel</button>
                <button type="submit" className="btn btn-primary" disabled={editLoading}>
                  {editLoading ? <span className="spinner" /> : 'Save Changes'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Search bar */}
      <div style={{ marginBottom: 16 }}>
        <div className="search-input-wrap">
          <span className="search-input-icon">🔍</span>
          <input
            className="form-input"
            placeholder="Search users by name, email, or role..."
            value={search}
            onChange={e => setSearch(e.target.value)}
          />
        </div>
      </div>

      <div className="card">
        {loading ? (
          <div className="p-loading"><span className="spinner" /></div>
        ) : (
          <div className="table-container">
            <table>
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Email</th>
                  <th>Department</th>
                  <th>Role</th>
                  <th>Status</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {filteredUsers.map(u => (
                  <tr key={u.id}>
                    <td className="font-semibold">{u.name}</td>
                    <td>{u.email}</td>
                    <td>{u.department?.name || 'No Department'}</td>
                    <td className="text-capitalize">{u.role}</td>
                    <td>
                      <span className={`badge ${u.is_blocked ? 'badge-rejected' : 'badge-approved'}`}>
                        {u.is_blocked ? 'Blocked' : 'Active'}
                      </span>
                    </td>
                    <td>
                      <div className="btn-group">
                        <button className="btn btn-sm btn-ghost" onClick={() => setEditUser(u)}>Edit</button>
                        <button
                          className={`btn btn-sm ${u.is_blocked ? 'btn-success' : 'btn-ghost'}`}
                          onClick={() => requestBlock(u.id, !u.is_blocked)}
                        >
                          {u.is_blocked ? 'Unblock' : 'Block'}
                        </button>
                        <button className="btn btn-sm btn-ghost" onClick={() => requestResetPassword(u.id)}>Reset PW</button>
                        <button className="btn btn-sm btn-danger" onClick={() => requestDelete(u.id)}>Delete</button>
                      </div>
                    </td>
                  </tr>
                ))}
                {filteredUsers.length === 0 && (
                  <tr>
                    <td colSpan={6} className="p-loading" style={{ color: 'var(--text-muted)' }}>
                      {search ? 'No users match your search.' : 'No users found.'}
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Confirm Modal */}
      {confirm && (
        <ConfirmModal
          title={confirm.title}
          message={confirm.message}
          confirmLabel={confirm.confirmLabel}
          variant={confirm.variant}
          loading={confirmLoading}
          onConfirm={handleConfirmAction}
          onCancel={() => setConfirm(null)}
        />
      )}
    </Layout>
  );
}
