import { DataTable } from '~/components/ui/data-table';
import { columns } from './columns';
import { apiInstance } from '~/helper/api';
import { PostApi } from '~/do-exercise-api';
import { DataTableSearch } from './search';
import { postListSchema } from './schema';

export default async function PostPage({ searchParams }: RouteSearchParams) {
  const search = await searchParams;
  const postApi = await apiInstance(PostApi);
  const { page, pageSize, name } = postListSchema.parse(search);
  const { meta, data } = await postApi.getPosts({ page, pageSize, name });

  return (
    <div className="container mx-auto py-10">
      <DataTable columns={columns} pagination={meta} data={data}>
        <DataTableSearch />
      </DataTable>
    </div>
  );
}
