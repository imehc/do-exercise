import { LocaleType } from '~/i18n/i18';
import { redirect } from '~/i18n/routing';

interface HomePageProps {
  params: Promise<{ locale: LocaleType }>;
}

export default async function HomePage({ params }: HomePageProps) {
  const { locale } = await params;

  return redirect({ href: '/dashboard', locale });
}
