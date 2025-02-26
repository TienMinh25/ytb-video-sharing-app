import axios from 'axios';
import { refreshAccessToken } from './auth';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000/api/v1';

const api = axios.create({
  baseURL: API_URL,
  headers: { 'Content-Type': 'application/json' },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('accessToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // Nếu request là refresh-token và bị 401, không xử lý lại để tránh vòng lặp
    if (
      error.response?.status === 401 &&
      originalRequest.url.includes('/accounts/refresh-token') &&
      !originalRequest._retry
    ) {
      originalRequest._retry = true;
      try {
        console.warn('Access token expired, refreshing...');
        const newAccessToken = await refreshAccessToken();

        originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
        return api(originalRequest);
      } catch (refreshError) {
        console.error('Failed to refresh token, logging out...', refreshError);
        localStorage.clear();
        window.location.href = '/';
        return Promise.reject(refreshError);
      }
    }
    return Promise.reject(error);
  },
);

export default api;
