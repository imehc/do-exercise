'use client';

import { Input } from '~/components/ui/input';
import { parseAsInteger, parseAsString, useQueryStates } from 'nuqs';
import { useTransition } from 'react';
import { postListSchema } from './schema';

export function DataTableSearch() {
  const [, startTransition] = useTransition();
  const [{ name }, setQueryState] = useQueryStates(
    {
      pageSize: parseAsInteger.withDefault(10),
      page: parseAsInteger.withDefault(0),
      name: parseAsString,
    },
    {
      shallow: false,
      startTransition,
    }
  );

  return (
    <div className="flex items-center space-x-4">
      <Input
        placeholder="按名称搜索..."
        value={name ?? ''}
        onChange={(event) => {
          setQueryState((state) => {
            return {
              ...state,
              name: postListSchema.shape.name.parse(event.target.value),
              page: 0,
            };
          });
        }}
        className="w-[280px]"
      />
    </div>
  );
}
