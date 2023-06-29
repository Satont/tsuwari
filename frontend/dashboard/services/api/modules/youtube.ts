import { useMutation, useQuery } from '@tanstack/react-query';
import { V1 } from '@twir/types/api';

import { useContext } from 'react';

import { authFetcher } from '@/services/api';
import { SelectedDashboardContext } from '@/services/selectedDashboardProvider';

type Youtube = V1['CHANNELS']['MODULES']['YouTube']

export const useYoutubeModule = () => {
  const dashboard = useContext(SelectedDashboardContext);
  const getUrl = () => `/api/v1/channels/${dashboard.id}/modules/youtube-sr`;

  return {
    useSettings: () => useQuery<Youtube['GET']>({
      queryKey: [getUrl()],
      queryFn: () => authFetcher(getUrl()),
      retry: false,
    }),
    useSearch: () => useMutation({
      mutationFn: ({ query, type }: {query: string, type: 'channel' | 'video'}) => {
        return authFetcher(`${getUrl()}/search?type=${type}&query=${query}`) as Promise<Youtube['SEARCH']>;
      },
      mutationKey: [`${getUrl()}/search`],
    }),
    useUpdate: () => useMutation({
      mutationFn: (body: Youtube['POST']) => {
        return authFetcher(`${getUrl()}`, {
          method: 'POST',
          body: JSON.stringify(body),
          headers: {
            'Content-Type': 'application/json',
          },
        });
      },
      mutationKey: [getUrl()],
    }),
  };
};
