import { NextRequest, NextResponse } from "next/server";
import createIntlMiddleware from "next-intl/middleware";
import { auth } from "~/helper/auth";
import { routing } from '~/i18n/routing';
import { locales } from "~/i18n/i18";

const protectedPages = ["/dashboard/*"];
const authPages = ["/auth/signin"];

const intlMiddleware = createIntlMiddleware(routing);

const testPagesRegex = (pages: string[], pathname: string) => {
  const regex = `^(/(${locales.join("|")}))?(${pages
    .map((p) => p.replace("/*", ".*"))
    .join("|")})/?$`;
  return new RegExp(regex, "i").test(pathname);
};


const handleAuth = async (req: NextRequest, isAuthPage: boolean, isProtectedPage: boolean) => {
  const session = await auth();
  const isAuth = !!session?.user;

  if (!isAuth && isProtectedPage) {
    let from = req.nextUrl.pathname;
    if (req.nextUrl.search) {
      from += req.nextUrl.search;
    }
    return NextResponse.redirect(new URL(`/auth/signin?from=${encodeURIComponent(from)}`, req.url));
  }

  if (isAuth && isAuthPage) {
    const url = req.nextUrl.clone();
    const fromValue = url.searchParams.get("from");
    return NextResponse.redirect(new URL(fromValue ?? "/", req.nextUrl));
  }

  return intlMiddleware(req);
};

export default async function middleware(req: NextRequest) {
  const isAuthPage = testPagesRegex(authPages, req.nextUrl.pathname);
  const isProtectedPage = testPagesRegex(protectedPages, req.nextUrl.pathname);

  return await handleAuth(req, isAuthPage, isProtectedPage);
}

export const config = {
  matcher: ["/((?!api|_next/static|_next/image|favicon.ico).*)"],
};