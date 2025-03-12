"use client"

import { ColumnDef } from "@tanstack/react-table"
import { Checkbox } from "~/components/ui/checkbox"
import { Post } from "~/do-exercise-api"
import { renderTableCell } from "~/helper/table-cell"

export const columns: ColumnDef<Post>[] = [
    {
        id: "select",
        header: ({ table }) => (
            <Checkbox
                checked={
                    table.getIsAllPageRowsSelected() ||
                    (table.getIsSomePageRowsSelected() && "indeterminate")
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
        accessorKey: "id",
        header: "ID"
    },
    {
        accessorKey: "code",
        header: "岗位编码",
    },
    {
        accessorKey: "name",
        header: "岗位名称",
    },
    {
        accessorKey: "status",
        header: "状态",
        cell: ({ row }) => renderTableCell({ row, key: 'status', type: 'status' }),
    },
    {
        accessorKey: "remark",
        header: "岗位描述",
        cell: ({ row }) => renderTableCell({ row, key: 'remark' })
    },
    {
        accessorKey: "createAt",
        header: "创建时间",
        cell: ({ row }) => renderTableCell({ row, key: 'createAt', type: 'time' })
    },
    {
        accessorKey: "actions",
        header: "操作",
        // cell: ({ row }) => (
        //     <Button variant="destructive" size="sm">删除</Button>
        // ),
    },
]
