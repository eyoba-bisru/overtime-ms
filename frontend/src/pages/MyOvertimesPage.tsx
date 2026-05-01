import { useState, useEffect, useCallback } from 'react';
import { Link } from 'react-router-dom';
import api from '../api/client';
import type { Overtime, PaginationMeta } from '../types';
import Layout from '../components/Layout';
import ConfirmModal from '../components/ConfirmModal';

function SkeletonTable({ cols }: { cols: number }) {
  return (
    <>
      {Array.from({ length: 5 }).map((_, i) => (
        <div className="skeleton-row" key={i}>
          {Array.from({ length: cols }).map((_, j) => (
            <div key={j} className={`skeleton skeleton-cell ${j === 0 ? 'skeleton-cell-md' : ''}`} />
          ))}
        </div>
      ))}
    </>
  );
}

export default function MyOvertimesPage() {
  const [data, setData] = useState<Overtime[]>([]);
  const [meta, setMeta] = useState<PaginationMeta | null>(null);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [deleteTarget, setDeleteTarget] = useState<string | null>(null);
  const [deleteLoading, setDeleteLoading] = useState(false);

  const fetchData = useCallback(async () => {
    setLoading(true);
    try {
      const res = await api.get(`/overtime/my?page=${page}&page_size=15`);
      setData(res.data.data || []);
      setMeta(res.data.meta || null);
    } catch { /* handled by interceptor */ }
    setLoading(false);
  }, [page]);

  const handleDelete = async () => {
    if (!deleteTarget) return;
    setDeleteLoading(true);
    try {
      await api.delete(`/overtime/${deleteTarget}`);
      setDeleteTarget(null);
      fetchData();
    } catch { /* handled by interceptor */ }
    setDeleteLoading(false);
  };

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  // Compute stats
  const totalCount = meta?.total ?? data.length;
  const pendingCount = data.filter(o => o.status === 'pending').length;
  const approvedCount = data.filter(o => o.status === 'approved').length;
  const rejectedCount = data.filter(o => o.status === 'rejected').length;

  return (
    <Layout>
      <div className="page-header">
        <h1 className="page-title">My Overtime Requests</h1>
        <p className="page-subtitle">View and track your submitted overtime requests</p>
      </div>

      <div className="stats-grid">
        <div className="stat-card">
          <div className="stat-icon">📋</div>
          <div className="stat-value">{totalCount}</div>
          <div className="stat-label">Total Requests</div>
        </div>
        <div className="stat-card">
          <div className="stat-icon">⏳</div>
          <div className="stat-value">{pendingCount}</div>
          <div className="stat-label">Pending</div>
        </div>
        <div className="stat-card">
          <div className="stat-icon">✅</div>
          <div className="stat-value">{approvedCount}</div>
          <div className="stat-label">Approved</div>
        </div>
        <div className="stat-card">
          <div className="stat-icon">❌</div>
          <div className="stat-value">{rejectedCount}</div>
          <div className="stat-label">Rejected</div>
        </div>
      </div>

      <div className="card">
        {loading ? (
          <SkeletonTable cols={9} />
        ) : data.length === 0 ? (
          <div className="empty-state">
            <div className="empty-state-icon">📋</div>
            <div className="empty-state-text">No overtime requests yet</div>
          </div>
        ) : (
          <>
            <div className="table-container">
              <table>
                <thead>
                  <tr>
                    <th>Department</th>
                    <th>Date</th>
                    <th>Start</th>
                    <th>End</th>
                    <th>Duration</th>
                    <th>Program</th>
                    <th>Job Done</th>
                    <th>Status</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {data.map(ot => (
                    <tr key={ot.id}>
                      <td>{ot.department_name}</td>
                      <td>{ot.date}</td>
                      <td>{ot.start_time}</td>
                      <td>{ot.end_time}</td>
                      <td>{ot.duration.toFixed(1)}h</td>
                      <td className="text-capitalize">{ot.program}</td>
                      <td className="text-ellipsis">{ot.job_done}</td>
                      <td><span className={`badge badge-${ot.status}`}>{ot.status}</span></td>
                      <td>
                        {ot.status === 'pending' && (
                          <div className="btn-group">
                            <Link to={`/overtime/edit/${ot.id}`} className="btn btn-ghost btn-sm">
                              ✎ Edit
                            </Link>
                            <button onClick={() => setDeleteTarget(ot.id)} className="btn btn-ghost btn-sm" style={{ color: 'var(--danger)' }}>
                              🗑 Delete
                            </button>
                          </div>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
            {meta && meta.total_pages > 1 && (
              <div className="pagination">
                <div className="pagination-info">
                  Page {meta.page} of {meta.total_pages} · {meta.total} total
                </div>
                <div className="pagination-buttons">
                  <button className="btn btn-ghost btn-sm" disabled={page <= 1} onClick={() => setPage(p => p - 1)}>← Prev</button>
                  <button className="btn btn-ghost btn-sm" disabled={page >= meta.total_pages} onClick={() => setPage(p => p + 1)}>Next →</button>
                </div>
              </div>
            )}
          </>
        )}
      </div>

      {deleteTarget && (
        <ConfirmModal
          title="Delete Overtime Request"
          message="Are you sure you want to delete this request? This action cannot be undone."
          confirmLabel="Delete"
          variant="danger"
          loading={deleteLoading}
          onConfirm={handleDelete}
          onCancel={() => setDeleteTarget(null)}
        />
      )}
    </Layout>
  );
}
