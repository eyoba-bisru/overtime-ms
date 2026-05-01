import { useState, useEffect } from 'react';
import api from '../api/client';
import Layout from '../components/Layout';
import { useToast } from '../context/ToastContext';
import type { Department } from '../types';

export default function DepartmentManagementPage() {
  const [departments, setDepartments] = useState<Department[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreate, setShowCreate] = useState(false);
  const [createName, setCreateName] = useState('');
  const [createLoading, setCreateLoading] = useState(false);
  const [editDept, setEditDept] = useState<Department | null>(null);
  const [editLoading, setEditLoading] = useState(false);
  const { showToast } = useToast();

  const fetchDepartments = async () => {
    setLoading(true);
    try {
      const res = await api.get('/admin/departments');
      setDepartments(res.data.data || []);
    } catch {
      showToast({ type: 'error', title: 'Error', message: 'Failed to fetch departments' });
    }
    setLoading(false);
  };

  useEffect(() => {
    fetchDepartments();
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!createName.trim()) return;
    setCreateLoading(true);
    try {
      await api.post('/admin/departments', { name: createName });
      showToast({ type: 'success', title: 'Success', message: 'Department created successfully' });
      setCreateName('');
      setShowCreate(false);
      await fetchDepartments();
    } catch (err: any) {
      const msg = err.response?.data?.error || 'Failed to create department';
      showToast({ type: 'error', title: 'Error', message: msg });
    }
    setCreateLoading(false);
  };

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editDept || !editDept.name.trim()) return;
    setEditLoading(true);
    try {
      await api.patch(`/admin/departments/${editDept.id}`, { name: editDept.name });
      showToast({ type: 'success', title: 'Success', message: 'Department updated successfully' });
      setEditDept(null);
      await fetchDepartments();
    } catch (err: any) {
      const msg = err.response?.data?.error || 'Failed to update department';
      showToast({ type: 'error', title: 'Error', message: msg });
    }
    setEditLoading(false);
  };

  const handleDelete = async (id: string) => {
    if (!window.confirm('Are you sure you want to delete this department? This might affect users linked to it.')) return;
    try {
      await api.delete(`/admin/departments/${id}`);
      showToast({ type: 'warning', title: 'Deleted', message: 'Department removed successfully' });
      await fetchDepartments();
    } catch (err: any) {
      const msg = err.response?.data?.error || 'Failed to delete department';
      showToast({ type: 'error', title: 'Error', message: msg });
    }
  };

  return (
    <Layout>
      <div className="page-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div>
          <h1 className="page-title">Department Management</h1>
          <p className="page-subtitle">Manage organizational departments</p>
        </div>
        <button className="btn btn-primary" onClick={() => setShowCreate(true)}>➕ Add Department</button>
      </div>

      {showCreate && (
        <div className="modal-overlay" onClick={() => setShowCreate(false)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h3 className="modal-title">Add New Department</h3>
              <button className="modal-close" onClick={() => setShowCreate(false)}>✕</button>
            </div>
            <form onSubmit={handleCreate}>
              <div className="modal-body">
                <div className="form-group">
                  <label className="form-label">Department Name</label>
                  <input
                    className="form-input"
                    value={createName}
                    onChange={e => setCreateName(e.target.value)}
                    placeholder="e.g. Finance"
                    required
                    autoFocus
                  />
                </div>
              </div>
              <div className="modal-footer">
                <button type="button" className="btn btn-ghost" onClick={() => setShowCreate(false)}>Cancel</button>
                <button type="submit" className="btn btn-primary" disabled={createLoading}>
                  {createLoading ? <span className="spinner" /> : 'Create'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {editDept && (
        <div className="modal-overlay" onClick={() => setEditDept(null)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h3 className="modal-title">Edit Department</h3>
              <button className="modal-close" onClick={() => setEditDept(null)}>✕</button>
            </div>
            <form onSubmit={handleUpdate}>
              <div className="modal-body">
                <div className="form-group">
                  <label className="form-label">Department Name</label>
                  <input
                    className="form-input"
                    value={editDept.name}
                    onChange={e => setEditDept({ ...editDept, name: e.target.value })}
                    required
                    autoFocus
                  />
                </div>
              </div>
              <div className="modal-footer">
                <button type="button" className="btn btn-ghost" onClick={() => setEditDept(null)}>Cancel</button>
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
                  <th>Department Name</th>
                  <th>ID</th>
                  <th style={{ textAlign: 'right' }}>Actions</th>
                </tr>
              </thead>
              <tbody>
                {departments.map(d => (
                  <tr key={d.id}>
                    <td style={{ fontWeight: 600 }}>{d.name}</td>
                    <td style={{ fontFamily: 'monospace', fontSize: '0.8rem', color: 'var(--text-secondary)' }}>{d.id}</td>
                    <td style={{ textAlign: 'right' }}>
                      <div className="btn-group" style={{ justifyContent: 'flex-end' }}>
                        <button className="btn btn-sm btn-ghost" onClick={() => setEditDept(d)}>Edit</button>
                        <button className="btn btn-sm btn-danger" onClick={() => handleDelete(d.id)}>Delete</button>
                      </div>
                    </td>
                  </tr>
                ))}
                {departments.length === 0 && (
                  <tr>
                    <td colSpan={3} style={{ textAlign: 'center', padding: 40, color: 'var(--text-secondary)' }}>
                      No departments found.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </Layout>
  );
}
