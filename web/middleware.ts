import { NextRequest, NextResponse } from 'next/server';
import createIntlMiddleware from 'next-intl/middleware';
import { auth, signOut } from '~/helper/auth';
import { routing } from '~/i18n/routing';
import { locales } from '~/i18n/i18';
import { isAfter, isBefore } from 'date-fns';

const protectedPages = ['/dashboard/*'];
const authPages = ['/auth/signin'];

const intlMiddleware = createIntlMiddleware(routing);

const testPagesRegex = (pages: string[], pathname: string) => {
  const regex = `^(/(${locales.join('|')}))?(${pages
    .map((p) => p.replace('/*', '.*'))
    .join('|')})/?$`;
  return new RegExp(regex, 'i').test(pathname);
};

const handleAuth = async (
  req: NextRequest,
  isAuthPage: boolean,
  isProtectedPage: boolean
) => {
  'use server';
  const session = await auth();
  const isAuth = !!session?.user;
  // TODO: 尝试使用refresh_token刷新【当前会执行多次会导致二次调用刷新接口无效】
  // 相关issue: https://authjs.dev/guides/refresh-token-rotation
  const isVaild = session && isAfter(new Date(session.expiresIn), Date.now()); // 有效的
  const isError = session?.error === 'RefreshTokenError'; // 没有错误

  if ((!isAuth && isProtectedPage) || isError) {
    await signOut({ redirect: false });
    let from = req.nextUrl.pathname;
    if (req.nextUrl.search) {
      from += req.nextUrl.search;
    }
    return NextResponse.redirect(
      new URL(`/auth/signin?from=${encodeURIComponent(from)}`, req.url)
    );
  }

  if (isAuth && isAuthPage && isVaild && !isError) {
    const url = req.nextUrl.clone();
    const fromValue = url.searchParams.get('from');
    return NextResponse.redirect(
      new URL(fromValue ?? '/dashboard', req.nextUrl)
    );
  }

  return intlMiddleware(req);
};

export default async function middleware(req: NextRequest) {
  const isAuthPage = testPagesRegex(authPages, req.nextUrl.pathname);
  const isProtectedPage = testPagesRegex(protectedPages, req.nextUrl.pathname);

  return await handleAuth(req, isAuthPage, isProtectedPage);
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico).*)'],
};
