import { useState, useEffect } from 'react';
import api from '../api/client';
import type { User, Role } from '../types';
import Layout from '../components/Layout';
import { useToast } from '../context/ToastContext';

const ROLES: Role[] = ['applicant', 'checker', 'approver', 'finance', 'admin'];

export default function UserManagementPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreate, setShowCreate] = useState(false);
  const [createForm, setCreateForm] = useState({ email: '', name: '', password: '', role: 'applicant' as Role });
  const [createError, setCreateError] = useState('');
  const [createLoading, setCreateLoading] = useState(false);
  const [editUser, setEditUser] = useState<User | null>(null);
  const [editLoading, setEditLoading] = useState(false);
  const [tempPassword, setTempPassword] = useState<{ userId: string; password: string } | null>(null);
  const { showToast } = useToast();

  const fetchUsers = async () => {
    setLoading(true);
    try {
      const res = await api.get('/admin/users');
      setUsers(res.data.data || []);
    } catch { /* handled */ }
    setLoading(false);
  };

  useEffect(() => { fetchUsers(); }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setCreateError('');
    setCreateLoading(true);
    try {
      await api.post('/admin/users', createForm);
      showToast({ type: 'success', title: 'User Created', message: 'The new user has been created successfully.' });
      setShowCreate(false);
      setCreateForm({ email: '', name: '', password: '', role: 'applicant' });
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
        role: editUser.role
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

  const handleBlock = async (id: string, isBlocked: boolean) => {
    try {
      await api.patch(`/admin/users/${id}/block`, { is_blocked: isBlocked });
      showToast({ type: 'info', title: 'Status Changed', message: isBlocked ? 'User has been blocked.' : 'User has been unblocked.' });
      await fetchUsers();
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || 'Failed to change status';
      showToast({ type: 'error', title: 'Error', message: msg });
    }
  };

  const handleResetPassword = async (id: string) => {
    try {
      const res = await api.patch(`/admin/users/${id}/reset-password`);
      showToast({ type: 'success', title: 'Password Reset', message: 'Temporary password generated.' });
      setTempPassword({ userId: id, password: res.data.data?.temporary_password || '' });
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || 'Failed to reset password';
      showToast({ type: 'error', title: 'Error', message: msg });
    }
  };

  const handleDelete = async (id: string) => {
    if (!window.confirm('Are you sure you want to delete this user?')) return;
    try {
      await api.delete(`/admin/users/${id}`);
      showToast({ type: 'warning', title: 'User Deleted', message: 'User account has been removed.' });
      await fetchUsers();
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || 'Failed to delete user';
      showToast({ type: 'error', title: 'Deletion Failed', message: msg });
    }
  };

  return (
    <Layout>
      <div className="page-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
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

      <div className="card">
        {loading ? (
          <div style={{ textAlign: 'center', padding: 40 }}><span className="spinner" /></div>
        ) : (
          <div className="table-container">
            <table>
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Email</th>
                  <th>Role</th>
                  <th>Status</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {users.map(u => (
                  <tr key={u.id}>
                    <td style={{ fontWeight: 600 }}>{u.name}</td>
                    <td>{u.email}</td>
                    <td style={{ textTransform: 'capitalize' }}>{u.role}</td>
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
                          onClick={() => handleBlock(u.id, !u.is_blocked)}
                        >
                          {u.is_blocked ? 'Unblock' : 'Block'}
                        </button>
                        <button className="btn btn-sm btn-ghost" onClick={() => handleResetPassword(u.id)}>Reset PW</button>
                        <button className="btn btn-sm btn-danger" onClick={() => handleDelete(u.id)}>Delete</button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </Layout>
  );
}
