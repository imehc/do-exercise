import type { PropsWithChildren } from 'react';
import { SessionProvider } from 'next-auth/react';
import { ThemeProvider } from '~/components/theme-provider';
import { routing } from '~/i18n/routing';
import { notFound } from 'next/navigation';
import { getMessages } from 'next-intl/server';
import { NextIntlClientProvider } from 'next-intl';
import type { LocaleType } from '~/i18n/i18';
import { Toaster } from '~/components/ui/sonner';
import '~/app/globals.css';
import { auth } from '~/helper/auth';
import { ReactScan } from '~/components/scan';
import { NuqsAdapter } from 'nuqs/adapters/next/app';

type Props = PropsWithChildren & RouteParams<{ locale: LocaleType }>;

export default async function RootLayout({ children, params }: Props) {
  const { locale } = await params;
  // Ensure that the incoming `locale` is valid
  if (!routing.locales.includes(locale)) {
    notFound();
  }
  // Providing all messages to the client
  // side is the easiest way to get started
  const messages = await getMessages();

  const session = await auth();

  return (
    <html lang={locale} suppressHydrationWarning>
      <head />
      <ReactScan />
      <body>
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          <NextIntlClientProvider messages={messages}>
            <SessionProvider session={session}>
              <NuqsAdapter>{children}</NuqsAdapter>
            </SessionProvider>
          </NextIntlClientProvider>
          <Toaster position="top-center" richColors />
        </ThemeProvider>
      </body>
    </html>
  );
}
