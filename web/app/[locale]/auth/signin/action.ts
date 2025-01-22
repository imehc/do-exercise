"use server";

import { isRedirectError } from "next/dist/client/components/redirect-error";
import type { LoginRequest } from "~/do-exercise-api";
import { signIn } from "~/helper/auth";
import type { ResponseData } from "~/helper/format-response";

export const signinAction = async (formData: LoginRequest): Promise<ResponseData> => {
  try {
    await signIn("credentials", formData)
    return {
      ok: true,
      data: {
        message: 'signinSuccess',
      }
    }
  } catch (error) {
    if (isRedirectError(error)) {
      return {
        ok: true,
        data: {
          message: 'signinSuccess',
          path: '/dashboard'
        }

      }
    }
    return {
      ok: false,
      message: "signinFailed",
      code: 500,
    }
  }
}

