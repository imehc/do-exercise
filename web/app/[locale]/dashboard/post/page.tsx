import { DataTable } from '~/components/ui/data-table';
import { columns } from './columns';
import { apiInstance } from '~/helper/api';
import { PostApi } from '~/do-exercise-api';
import { DataTableSearch } from './search';
import { postListSchema } from './schema';
import { Button } from '~/components/ui/button';
import { Link } from '~/i18n/routing';
import { Icon } from '@iconify/react';

export default async function PostPage({ searchParams }: RouteSearchParams) {
  const search = await searchParams;
  const postApi = await apiInstance(PostApi);
  const { page, pageSize, name } = postListSchema.parse(search);
  const { meta, data } = await postApi.getPosts({ page, pageSize, name });

  return (
    <div className="container mx-auto py-8">
      <div className="space-y-4">
        <div className="flex justify-between items-center">
          <h2 className="text-2xl font-semibold tracking-tight">岗位管理</h2>
          <Link href="post/new">
            <Button size="sm" className="h-9">
              <Icon icon="material-symbols:add" className="mr-2 h-4 w-4" />
              新建岗位
            </Button>
          </Link>
        </div>
        <div className="flex items-center justify-between">
          <DataTableSearch />
        </div>
        <DataTable 
          columns={columns} 
          pagination={meta} 
          data={data} 
        />
      </div>
    </div>
  );
}
