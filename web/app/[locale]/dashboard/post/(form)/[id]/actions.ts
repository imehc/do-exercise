'use server';

import { PostApi } from '~/do-exercise-api';
import { apiInstance } from '~/helper/api';
import { PostFormValues } from '../schema';
import { revalidatePath } from 'next/cache';
import { type ResponseData } from '~/helper/format-response';

export async function editPostAction(
  id: number,
  values: PostFormValues
): Promise<ResponseData> {
  const postApi = await apiInstance(PostApi);
  try {
    await postApi.updatePost({
      id,
      postUpdate: values,
    });
    revalidatePath('/dashboard/post');
    return {
      ok: true,
      data: null,
      i18n: 'updateSuccess',
    };
  } catch (error) {
    console.error(error);
    return {
      ok: false,
      message: '',
      i18n: 'updateFail',
      code: 500,
    };
  }
}
