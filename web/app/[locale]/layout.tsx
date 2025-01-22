import type { PropsWithChildren } from 'react';
import { ThemeProvider } from '~/components/theme-provider';
import { routing } from '~/i18n/routing';
import { notFound } from 'next/navigation';
import { getMessages } from 'next-intl/server';
import { NextIntlClientProvider } from 'next-intl';
import type { LocaleType } from '~/i18n/i18';
import { Toaster } from '~/components/ui/sonner';
import '~/app/globals.css';

type Props = PropsWithChildren & { params: Promise<{ locale: LocaleType }> };

export default async function RootLayout({ children, params }: Props) {
  const { locale } = await params;
  // Ensure that the incoming `locale` is valid
  if (!routing.locales.includes(locale)) {
    notFound();
  }
  // Providing all messages to the client
  // side is the easiest way to get started
  const messages = await getMessages();

  return (
    <html lang={locale} suppressHydrationWarning>
      <head />
      <body>
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          <NextIntlClientProvider messages={messages}>
            {children}
          </NextIntlClientProvider>
          <Toaster position='top-center' richColors />
        </ThemeProvider>
      </body>
    </html>
  );
}
