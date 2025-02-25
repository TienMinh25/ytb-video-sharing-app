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
    if (error.response?.status === 401) {
      try {
        console.warn('Access token expired, refreshing...');
        await refreshAccessToken();

        console.log(
          localStorage.getItem('refreshToken') == '' ? 'dung roi' : 'sai roi',
        );
        error.config.headers.Authorization = `Bearer ${localStorage.getItem(
          'accessToken',
        )}`;
        return api(error.config); // Gửi lại request sau khi refresh thành công
      } catch (refreshError) {
        console.error('Failed to refresh token, logging out...', refreshError);
        // logout();
      }
    }
    return Promise.reject(error);
  },
);

export default api;
