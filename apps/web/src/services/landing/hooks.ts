import { useQuery } from '@tanstack/vue-query';

import { getStats } from './api.js';

export const useStats = () =>
  useQuery(['stats'], getStats, {
    retry: false,
    refetchOnReconnect: true,
    refetchOnWindowFocus: true,
    refetchInterval: 2500,
  });

const formatter = new Intl.NumberFormat(undefined, { notation: 'standard' });
export const useStatsFormatter = () => {
  return {
    format: formatter.format.bind(formatter),
  };
};
