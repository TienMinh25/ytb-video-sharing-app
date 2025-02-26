// auth.ts
import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000/api/v1';

export const refreshAccessToken = async () => {
  const refreshToken = localStorage.getItem('refreshToken');
  if (!refreshToken) throw new Error('No refresh token');

  try {
    const res = await axios.post(
      `${API_URL}/accounts/refresh-token`,
      {},
      {
        headers: { 'X-Authorization': refreshToken },
      },
    );
    const { access_token, refresh_token } = res.data.data;
    localStorage.setItem('accessToken', access_token);
    localStorage.setItem('refreshToken', refresh_token);

    window.dispatchEvent(new Event('tokenRefreshed'));
    return access_token;
  } catch (err) {
    console.error('Refresh token failed:', err);
    throw new Error('Refresh token failed');
  }
};
