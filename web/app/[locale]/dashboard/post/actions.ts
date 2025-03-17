'use server';

import { PostApi } from '~/do-exercise-api';
import { apiInstance } from '~/helper/api';
import { revalidatePath } from 'next/cache';
import { type ResponseData } from '~/helper/format-response';

export async function deletePostAction(id: number): Promise<ResponseData> {
  const postApi = await apiInstance(PostApi);
  try {
    await postApi.deletePost({
      id,
    });
    revalidatePath('/dashboard/post');
    return {
      ok: true,
      data: null,
      i18n: 'deleteSuccess',
    };
  } catch (error) {
    console.error(error);
    return {
      ok: false,
      message: '',
      i18n: 'deleteFail',
      code: 500,
    };
  }
}
