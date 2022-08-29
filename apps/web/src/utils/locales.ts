import type { Locale } from '@/types/locale';

export const locales: Locale[] = ['en', 'ru'];
export const defaultLocale = 'en';

interface LocaleTypes {
  landing: ReturnType<() => typeof import('@/locales/landing/en.json')>;
  app: ReturnType<() => typeof import('@/locales/app/en.json')>;
}

export type LocaleType = keyof LocaleTypes;

export async function loadLocaleMessages<L extends LocaleType>(
  localeType: L,
  locale: Locale,
): Promise<LocaleTypes[L]> {
  return (await import(`../locales/${localeType}/${locale}.json`)).default;
}
