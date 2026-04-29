import { useState, useEffect, useCallback } from 'react';
import api from '../api/client';
import type { Overtime, PaginationMeta } from '../types';
import Layout from '../components/Layout';

export default function MyOvertimesPage() {
  const [data, setData] = useState<Overtime[]>([]);
  const [meta, setMeta] = useState<PaginationMeta | null>(null);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);

  const fetchData = useCallback(async () => {
    setLoading(true);
    try {
      const res = await api.get(`/overtime/my?page=${page}&page_size=15`);
      setData(res.data.data || []);
      setMeta(res.data.meta || null);
    } catch { /* handled by interceptor */ }
    setLoading(false);
  }, [page]);

  useEffect(() => { fetchData(); }, [fetchData]);

  return (
    <Layout>
      <div className="page-header">
        <h1 className="page-title">My Overtime Requests</h1>
        <p className="page-subtitle">View and track your submitted overtime requests</p>
      </div>
      <div className="card">
        {loading ? (
          <div style={{ textAlign: 'center', padding: 40 }}><span className="spinner" /></div>
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
                    <th>Date</th>
                    <th>Start</th>
                    <th>End</th>
                    <th>Duration</th>
                    <th>Program</th>
                    <th>Job Done</th>
                    <th>Status</th>
                  </tr>
                </thead>
                <tbody>
                  {data.map(ot => (
                    <tr key={ot.id}>
                      <td>{ot.date}</td>
                      <td>{ot.start_time}</td>
                      <td>{ot.end_time}</td>
                      <td>{ot.duration.toFixed(1)}h</td>
                      <td style={{ textTransform: 'capitalize' }}>{ot.program}</td>
                      <td style={{ maxWidth: 200, overflow: 'hidden', textOverflow: 'ellipsis' }}>{ot.job_done}</td>
                      <td><span className={`badge badge-${ot.status}`}>{ot.status}</span></td>
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
    </Layout>
  );
}
