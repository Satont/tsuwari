import { seoLocales } from '@/data/seo.js';
import type { Locale } from '@/locales/index.js';
import { defaultLocale } from '@/locales/index.js';
import { htmlLayout } from '@/utils/htmlLayout.js';
import type { PageContext } from '@/utils/pageContext.js';

export async function render(pageContext: PageContext) {
  const locale: Locale = defaultLocale;
  pageContext.locale = locale;

  const seoInfo = seoLocales[locale];
  const documentHtml = htmlLayout(seoInfo, pageContext);

  return {
    documentHtml,
    pageContext: {
      enableEagerStreaming: true,
      locale,
    },
  };
}
