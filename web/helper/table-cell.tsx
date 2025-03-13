import { Row } from '@tanstack/react-table';
import { format } from 'date-fns';
import { StatusTag } from '~/components/other/status-tag';
import { Status } from '~/do-exercise-api';

interface Props<T> {
  row: Row<T>;
  key: keyof T;
  type?: 'status' | 'time';
  // TODO: 根据不同类型渲染不同的格式
}

export function renderTableCell<T extends object>({
  row,
  key,
  type,
}: Props<T>) {
  const value = row?.getValue?.(key as string);
  switch (type) {
    case 'status':
      return <StatusTag status={value as Status} />;
    case 'time':
      return format(value as Date, 'yyyy-MM-dd hh:mm:ss');
    default:
      return value || '-';
  }
}
