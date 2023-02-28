import { useQuery } from '@tanstack/react-query';
import { HelixUserData } from '@twurple/api';

import { authFetcher } from '@/services/api/fetchWrappers';


export const useTwitch = () => {
  return {
    useGetUsersByNames: (names: string) => useQuery<HelixUserData[]>({
      queryKey: [`/api/v1/twitch/users`, names],
      queryFn: () => {
        if (!names.length) return [];
        return authFetcher(`/api/v1/twitch/users?names=${names}`);
      },
    }),
  };
};

export const useTwitchUsersByNames = (names: string[]) =>
  useQuery<HelixUserData[]>({
    queryKey: [`/api/v1/twitch/users`, names.join(',')],
    queryFn: () => authFetcher(`/api/auth/profile?names=${names.join(',')}`),
  });

export const useTwitchUsersByIds = (ids: string[]) =>
  useQuery<HelixUserData[]>({
    queryKey: [`/api/v1/twitch/users`, ids.join(',')],
    queryFn: () => authFetcher(`/api/auth/profile?ids=${ids.join(',')}`),
  });

