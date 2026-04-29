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
    
    if (error.response?.status === 401) {
      // Don't show toast for 401 if we are already on login page
      if (window.location.pathname !== '/login') {
        eventBus.emit('SHOW_TOAST', {
          type: 'error',
          title: 'Session Expired',
          message: 'Please login again.',
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
