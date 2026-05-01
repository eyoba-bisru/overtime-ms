import axios from 'axios';
import { eventBus } from '../utils/eventBus';

const api = axios.create({
  baseURL: '/api/v1',
  withCredentials: true,
  headers: { 'Content-Type': 'application/json' },
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    const message = error.response?.data?.error || error.message || 'An unexpected error occurred';
    
    if (error.response?.status === 401 || error.response?.status === 403) {
      // Don't show toast for 401/403 if we are already on login page
      if (window.location.pathname !== '/login') {
        const title = error.response?.status === 403 ? 'Access Denied' : 'Session Expired';
        eventBus.emit('SHOW_TOAST', {
          type: 'error',
          title: title,
          message: message,
        });
        window.location.href = '/login';
      }
    } else {
      eventBus.emit('SHOW_TOAST', {
        type: 'error',
        title: 'API Error',
        message: message,
      });
    }
    return Promise.reject(error);
  }
);

export default api;
