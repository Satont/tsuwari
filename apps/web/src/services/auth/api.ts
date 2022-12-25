import type { AuthUser } from '@tsuwari/shared';

import { authFetch } from './authFetch.js';
import { LOGIN_ROUTE_STATE, ORIGIN_STATE } from './locationHelpers.js';

export const getProfile = async (): Promise<AuthUser> => {
  const res = await authFetch('/api/auth/profile');
  if (!res.ok) {
    const error = await res.text();
    throw new Error(error);
  }
  return res.json();
};

export const API_LOGIN_ROUTE = `/api/auth?state=${ORIGIN_STATE}`;

/**
 * @returns Response object with new access token
 */
export const postRefreshToken = async () => {
  return await fetch('/api/auth/token', { method: 'post' });
};

export const authorizeByTwitchCode = async (code: string): Promise<{ accessToken: string }> => {
  const searchParams = new URLSearchParams({
    code,
    state: LOGIN_ROUTE_STATE,
  });

  const res = await fetch('/api/auth/token?' + searchParams);

  return res.json();
};

export const logout = async () => {
  const res = await authFetch('/api/auth/logout', { method: 'POST' });
  return res.ok;
};
