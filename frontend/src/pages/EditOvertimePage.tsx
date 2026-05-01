import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import api from '../api/client';
import type { Overtime, OvertimeProgram } from '../types';
import Layout from '../components/Layout';
import { useToast } from '../context/ToastContext';

export default function EditOvertimePage() {
  const navigate = useNavigate();
  const { id } = useParams();
  const { showToast } = useToast();
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [form, setForm] = useState({
    date: '',
    start_time: '',
    end_time: '',
    job_done: '',
    program: 'night' as OvertimeProgram,
  });

  useEffect(() => {
    const fetchData = async () => {
      try {
        const res = await api.get(`/overtime/${id}`);
        const ot: Overtime = res.data.data;
        
        if (ot.status !== 'pending') {
          showToast({ type: 'error', title: 'Action Denied', message: 'Only pending requests can be edited.' });
          navigate('/overtime/my');
          return;
        }

        setForm({
          date: ot.date,
          start_time: ot.start_time,
          end_time: ot.end_time,
          job_done: ot.job_done,
          program: ot.program,
        });
      } catch {
        navigate('/overtime/my');
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, [id, navigate, showToast]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    setForm(prev => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    try {
      await api.patch(`/overtime/${id}`, form);
      showToast({ type: 'success', title: 'Success', message: 'Overtime request updated successfully' });
      navigate('/overtime/my');
    } catch {
      // Interceptor handles toast
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) return <Layout><div className="loading-page"><span className="spinner" /></div></Layout>;

  return (
    <Layout>
      <div className="page-header">
        <h1 className="page-title">Edit Overtime Request</h1>
        <p className="page-subtitle">Update your pending overtime work request</p>
      </div>
      <div className="card max-w-form">
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label className="form-label">Date</label>
            <input type="date" name="date" className="form-input" value={form.date} onChange={handleChange} required />
          </div>
          <div className="form-row">
            <div className="form-group">
              <label className="form-label">Start Time</label>
              <input type="time" name="start_time" className="form-input" value={form.start_time} onChange={handleChange} required />
            </div>
            <div className="form-group">
              <label className="form-label">End Time</label>
              <input type="time" name="end_time" className="form-input" value={form.end_time} onChange={handleChange} required />
            </div>
          </div>
          <div className="form-group">
            <label className="form-label">Program</label>
            <select name="program" className="form-select" value={form.program} onChange={handleChange}>
              <option value="night">Night</option>
              <option value="weekend">Weekend</option>
              <option value="holiday">Holiday</option>
            </select>
          </div>
          <div className="form-group">
            <label className="form-label">Job Done</label>
            <textarea
              name="job_done"
              className="form-input"
              rows={3}
              placeholder="Describe the work performed..."
              value={form.job_done}
              onChange={handleChange}
              required
              style={{ resize: 'vertical' }}
            />
          </div>
          <div className="btn-group">
            <button type="submit" className="btn btn-primary" disabled={submitting}>
              {submitting ? <span className="spinner" /> : 'Save Changes'}
            </button>
            <button type="button" className="btn btn-ghost" onClick={() => navigate('/overtime/my')}>Cancel</button>
          </div>
        </form>
      </div>
    </Layout>
  );
}
