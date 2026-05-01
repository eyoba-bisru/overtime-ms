import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../api/client';
import type { OvertimeProgram } from '../types';
import Layout from '../components/Layout';
import { useToast } from '../context/ToastContext';

export default function CreateOvertimePage() {
  const navigate = useNavigate();
  const { showToast } = useToast();
  const [loading, setLoading] = useState(false);
  const [form, setForm] = useState({
    date: '',
    start_time: '',
    end_time: '',
    job_done: '',
    program: 'night' as OvertimeProgram,
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    setForm(prev => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      await api.post('/overtime', form);
      showToast({ type: 'success', title: 'Success', message: 'Overtime request submitted successfully' });
      navigate('/overtime/my');
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || 'Failed to create request';
      showToast({ type: 'error', title: 'Submission Error', message: msg });
    } finally {
      setLoading(false);
    }
  };

  return (
    <Layout>
      <div className="page-header">
        <h1 className="page-title">New Overtime Request</h1>
        <p className="page-subtitle">Submit a new overtime work request</p>
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
            <button type="submit" className="btn btn-primary" disabled={loading}>
              {loading ? <span className="spinner" /> : 'Submit Request'}
            </button>
            <button type="button" className="btn btn-ghost" onClick={() => navigate('/overtime/my')}>Cancel</button>
          </div>
        </form>
      </div>
    </Layout>
  );
}
