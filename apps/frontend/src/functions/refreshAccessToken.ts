// eslint-disable-next-line import/no-cycle
import { api } from '@/plugins/api';

export const refreshAccessToken = async () => {
  const refreshToken = localStorage.getItem('refreshToken');

  if (!refreshToken) {
    throw new Error('Refresh token is empty.');
  }

  try {
    const request = await api.post<{
      accessToken: string;
      refreshToken: string;
    }>('/auth/token', { refreshToken });
    const data = request.data;

    console.log('data', data);

    localStorage.setItem('accessToken', data.accessToken);
    localStorage.setItem('refreshToken', data.refreshToken);
    // eslint-disable-next-line no-empty
  } catch (error: any) {}
};
