'use server';

import { PostApi } from '~/do-exercise-api';
import { apiInstance } from '~/helper/api';
import { PostFormValues } from '../schema';
import { revalidatePath } from 'next/cache';
import { type ResponseData } from '~/helper/format-response';

export async function createPostAction(
  values: PostFormValues
): Promise<ResponseData> {
  const postApi = await apiInstance(PostApi);
  try {
    await postApi.createPost({
      postCreate: values,
    });
    revalidatePath('/dashboard/post');
    return {
      ok: true,
      data: null,
      i18n: 'createSuccess',
    };
  } catch (error) {
    console.error(error);
    return {
      ok: false,
      message: '',
      i18n: 'createFail',
      code: 500,
    };
  }
}
