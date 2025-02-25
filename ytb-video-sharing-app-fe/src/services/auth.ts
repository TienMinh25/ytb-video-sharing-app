import api from './api';

export const refreshAccessToken = async () => {
  const refreshToken = localStorage.getItem('refreshToken');
  if (refreshToken == '') throw new Error('No refresh token');

  console.log(refreshToken);
  try {
    const res = await api.post(
      '/accounts/refresh-token',
      {},
      {
        headers: { 'X-Authorization': refreshToken },
      },
    );
    console.log('RUN REFRESH');
    const { access_token, refresh_token } = res.data.data;
    localStorage.setItem('accessToken', access_token);
    localStorage.setItem('refreshToken', refresh_token);

    return access_token;
  } catch {
    throw new Error('Refresh token failed');
  }
};
