import { getTranslations } from 'next-intl/server';

export type TranslateI18n = UnwrapPromise<ReturnType<typeof getTranslations>>
