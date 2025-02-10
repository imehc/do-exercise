'use server';

import { AuthError } from 'next-auth';
import { isRedirectError } from 'next/dist/client/components/redirect-error';
import { signIn } from '~/helper/auth';
import type { LoginRequest } from '~/do-exercise-api';
import type { ResponseData } from '~/helper/format-response';

export const signinAction = async (
  formData: LoginRequest
): Promise<ResponseData> => {
  try {
    await signIn('credentials', { ...formData, redirect: false });
    return {
      ok: true,
      data: null,
      i18n: 'signinSuccess',
      href: '/dashboard', // 手动重定向
    };
  } catch (error) {
    if (isRedirectError(error)) {
      // doc https://next-auth.js.org/getting-started/client#using-the-redirect-false-option
      // 由于禁用重定向就不需要这个判断了
      return {
        ok: true,
        data: null,
        i18n: 'signinSuccess',
        href: '/dashboard',
      };
    }
    if (error instanceof AuthError) {
      const message = error.cause?.err?.message;
      return {
        ok: false,
        message: message ?? '',
        i18n: !message ? 'signinFailed' : undefined,
        code: 400,
      };
    }
    return {
      ok: false,
      message: '',
      i18n: 'signinFailed',
      code: 500,
    };
  }
};
