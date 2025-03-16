import { Post, Status } from '~/do-exercise-api';
import { PostForm } from '../form';

export default async function CreatePostPage() {
  const data = {
    name: '',
    code: '',
    remark: '',
    status: Status.disabled,
    sort: 0,
  } as Post;
  return <PostForm data={data} />;
}
