import { apiInstance } from '~/helper/api';
import { PostApi } from '~/do-exercise-api';
import { z } from 'zod';
import { PostForm } from '../../../../(form)/form';

const schema = z.object({
  id: z.coerce.number(),
});

export default async function EditPostPage({ params }: RouteParams) {
  const param = await params;
  const postApi = await apiInstance(PostApi);
  const { id } = schema.parse(param);
  const data = await postApi.getPost({ id });
  return <PostForm data={data} />;
}
