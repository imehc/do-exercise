'use client';

import { ColumnDef } from '@tanstack/react-table';
import { MoreHorizontal } from 'lucide-react';
import { CommonAlertDialogContent } from '~/components/other';
import { AlertDialog, AlertDialogTrigger } from '~/components/ui/alert-dialog';
import { Button } from '~/components/ui/button';
import { Checkbox } from '~/components/ui/checkbox';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '~/components/ui/dropdown-menu';
import { Post } from '~/do-exercise-api';
import { renderTableCell } from '~/helper/table-cell';
import { Link } from '~/i18n/routing';
import { deletePostAction } from './actions';
import { useTranslations } from 'next-intl';
import { handleClientResponse } from '~/helper/client-response';

// TODO: 多语言
export const columns: ColumnDef<Post>[] = [
  {
    id: 'select',
    header: ({ table }) => (
      <Checkbox
        checked={
          table.getIsAllPageRowsSelected() ||
          (table.getIsSomePageRowsSelected() && 'indeterminate')
        }
        onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
        aria-label="Select all"
      />
    ),
    cell: ({ row }) => (
      <Checkbox
        checked={row.getIsSelected()}
        onCheckedChange={(value) => row.toggleSelected(!!value)}
        aria-label="Select row"
      />
    ),
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: 'id',
    header: 'ID',
  },
  {
    accessorKey: 'code',
    header: '岗位编码',
  },
  {
    accessorKey: 'name',
    header: '岗位名称',
  },
  {
    accessorKey: 'status',
    header: '状态',
    cell: ({ row }) => renderTableCell({ row, key: 'status', type: 'status' }),
  },
  {
    accessorKey: 'remark',
    header: '岗位描述',
    cell: ({ row }) => renderTableCell({ row, key: 'remark' }),
  },
  {
    accessorKey: 'createAt',
    header: '创建时间',
    cell: ({ row }) => renderTableCell({ row, key: 'createAt', type: 'time' }),
  },
  {
    accessorKey: 'actions',
    header: '操作',
    cell: ({ row }) => {
      const id: number = row.getValue('id');
      const name: string = row.getValue('name');

      return (
        <AlertDialog>
          <DropdownMenu modal={false}>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <span className="sr-only">Open menu</span>
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <Link href={`post/${id}`}>
                <DropdownMenuItem className="cursor-pointer">
                  编辑
                </DropdownMenuItem>
              </Link>
              <AlertDialogTrigger asChild>
                <DropdownMenuItem className="cursor-pointer !text-red-500/80 hover:!text-red-500">
                  删除
                </DropdownMenuItem>
              </AlertDialogTrigger>
            </DropdownMenuContent>
          </DropdownMenu>
          <DeletePostContent name={name} id={id} />
        </AlertDialog>
      );
    },
  },
];
interface Props {
  name: string;
  id: number;
}
function DeletePostContent({ name, id }: Props) {
  const t = useTranslations('PostForm');
  const onCancel = () => {
    console.log('cancel');
  };
  const onContinue = async () => {
    const res = await deletePostAction(id);
    handleClientResponse({ res, t });
  };

  return (
    <CommonAlertDialogContent
      title="确定删除吗?"
      subTitle={
        <>
          <span>删除</span>
          <span className="font-bold italic">{name}</span>
          <span>后将无法恢复</span>
        </>
      }
      onCancel={onCancel}
      onContinue={onContinue}
    />
  );
}
