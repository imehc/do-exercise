"use server";

import { AuthError } from "next-auth";
import { isRedirectError } from "next/dist/client/components/redirect-error";
import { signIn } from "~/helper/auth";
import type { LoginRequest } from "~/do-exercise-api";
import type { ResponseData } from "~/helper/format-response";

export const signinAction = async (formData: LoginRequest): Promise<ResponseData> => {
  try {
    await signIn("credentials", formData)
    return {
      ok: true,
      data: null,
      i18n: "signinSuccess",
    }
  } catch (error) {
    if (isRedirectError(error)) {
      return {
        ok: true,
        data: null,
        i18n: "signinSuccess",
        href: '/dashboard'
      }

    }
    if (error instanceof AuthError) {
      const message = error.cause?.err?.message
      return {
        ok: false,
        message: message ?? '',
        i18n: !message ? 'signinFailed' : undefined,
        code: 400
      }
    }
    return {
      ok: false,
      message: "",
      i18n: "signinFailed",
      code: 500,
    }
  }
}

