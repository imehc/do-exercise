import { ResponseError } from '~/do-exercise-api';

export type ResponseData<T = unknown> =
  | {
      data: T;
      ok: true;
    }
  | {
      ok: false;
      code: number;
      message: string;
    };

// 处理 response
export async function handleResponse<T>(
  response: Promise<T>,
  refresh?: () => void
): Promise<ResponseData<T>> {
  try {
    const data = await response;
    refresh?.();
    return {
      data,
      ok: true,
    };
  } catch (error) {
    if (error instanceof ResponseError) {
      const text = await error.response.text();
      return {
        ok: false,
        code: error.response.status,
        message: text || error.message,
      };
    }
    return {
      ok: false,
      code: 500,
      message: 'Internal Server Error',
    };
  }
}
