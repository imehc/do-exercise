import { LocaleType } from '~/i18n/i18';
import { redirect } from '~/i18n/routing';

export default async function HomePage({ params }: RouteParams<{ locale: LocaleType }>) {
  const { locale } = await params;

  return redirect({ href: '/dashboard', locale });
}
