import { IconAlertTriangle } from '@tabler/icons-react'
import { Trans } from '@lingui/react/macro'
import { cn } from '~/lib/utils'
import { Badge } from '~/components/ui/badge'
import { Button } from '~/components/ui/button'
import { ScrollArea } from '~/components/ui/scroll-area'
import { Separator } from '~/components/ui/separator'
import { callMethodTypes } from '~/features/api/data/data'
import type { SysMenuApiBrief } from '~/do-exercise-api'
import type {
  MenuNodeSummary,
  MenuSelectionDiff,
} from '../utils/permission-diff'

interface Props {
  /** 勾选范围内的按钮权限 */
  permissions: MenuNodeSummary[]
  /** 勾选范围实际放开的接口（已去重） */
  apis: SysMenuApiBrief[]
  /** 漏勾的祖先菜单 */
  missingAncestors: MenuNodeSummary[]
  /** 勾了页面却没勾任何操作按钮的页面 */
  pagesWithoutActions: MenuNodeSummary[]
  /** 与保存前相比的增删；新建角色时传 null 不展示 */
  diff: MenuSelectionDiff | null
  /** 当前持有该角色的用户数；新建角色时传 null（还没有持有者） */
  userCount: number | null
  onFixAncestors: () => void
}

function NodeBadges({ nodes }: { nodes: MenuNodeSummary[] }) {
  return (
    <div className='flex flex-wrap gap-1'>
      {nodes.map((node) => (
        <Badge key={node.id} variant='outline' className='font-normal'>
          {node.permission ? `${node.label} · ${node.permission}` : node.label}
        </Badge>
      ))}
    </div>
  )
}

/**
 * 授权变更的影响面预览。
 *
 * 勾菜单这个动作离「这个角色能做什么」隔了两层：菜单 -> 权限标识 -> 接口。
 * 授权人只看菜单树无法判断自己放开了什么，这里把三层一次摊开，并把最常见的
 * 「漏勾父级导致页面进不去」做成可一键修复的提示。
 */
export function RolePermissionSummary({
  permissions,
  apis,
  missingAncestors,
  pagesWithoutActions,
  diff,
  userCount,
  onFixAncestors,
}: Props) {
  const hasDiff = !!diff && (diff.added.length > 0 || diff.removed.length > 0)
  const addedCount = diff?.added.length ?? 0
  const removedCount = diff?.removed.length ?? 0
  const permissionCount = permissions.length
  const apiCount = apis.length

  return (
    <div className='col-span-8 col-start-3 space-y-3 rounded-md border p-3'>
      {missingAncestors.length > 0 && (
        <div className='space-y-2 rounded-md border border-destructive/40 p-2'>
          <div className='flex items-center gap-2 text-destructive text-sm'>
            <IconAlertTriangle className='size-4' />
            <span>
              <Trans>
                以下父级菜单未勾选，子权限虽然放开，对应页面仍然打不开
              </Trans>
            </span>
          </div>
          <NodeBadges nodes={missingAncestors} />
          <Button
            type='button'
            size='sm'
            variant='outline'
            onClick={onFixAncestors}
          >
            <Trans>一并勾选父级</Trans>
          </Button>
        </div>
      )}

      {pagesWithoutActions.length > 0 && (
        <div className='space-y-2 rounded-md border border-amber-500/40 p-2'>
          <div className='flex items-center gap-2 text-amber-600 text-sm dark:text-amber-500'>
            <IconAlertTriangle className='size-4' />
            <span>
              <Trans>
                以下页面没有勾选任何操作权限，该角色打开后只能看、不能操作
              </Trans>
            </span>
          </div>
          <NodeBadges nodes={pagesWithoutActions} />
        </div>
      )}

      <div className='space-y-1'>
        <div className='font-medium text-sm'>
          <Trans>权限标识</Trans>
          <span className='text-muted-foreground ml-1 text-xs'>
            ({permissionCount})
          </span>
        </div>
        {permissions.length ? (
          <ScrollArea className='max-h-24'>
            <NodeBadges nodes={permissions} />
          </ScrollArea>
        ) : (
          <p className='text-muted-foreground text-xs'>
            <Trans>尚未勾选任何按钮权限，该角色只能浏览页面</Trans>
          </p>
        )}
      </div>

      <Separator />

      <div className='space-y-1'>
        <div className='font-medium text-sm'>
          <Trans>放开的接口</Trans>
          <span className='text-muted-foreground ml-1 text-xs'>
            ({apiCount})
          </span>
        </div>
        {apis.length ? (
          <ScrollArea className='max-h-32'>
            <div className='flex flex-col gap-1'>
              {apis.map((api) => (
                <Badge
                  key={`${api.method} ${api.path}`}
                  variant='outline'
                  className={cn('w-fit font-normal', callMethodTypes.get(api.method))}
                >
                  <span>{api.method}</span>
                  <span className='mx-1'>|</span>
                  <span>{api.path}</span>
                  {api.description && (
                    <span className='ml-1 opacity-70'>{api.description}</span>
                  )}
                </Badge>
              ))}
            </div>
          </ScrollArea>
        ) : (
          <p className='text-muted-foreground text-xs'>
            <Trans>当前勾选未绑定任何接口</Trans>
          </p>
        )}
      </div>

      {hasDiff && (
        <>
          <Separator />
          <div className='space-y-2'>
            <div className='font-medium text-sm'>
              <Trans>本次变更</Trans>
            </div>
            {/* 角色是多人共享对象：改一次授权就改了所有持有者的能力边界，
                所以把「影响多少人」和「改了哪几项」放在一起看 */}
            {userCount !== null && (
              <div
                className={cn(
                  'text-xs',
                  userCount > 0 ? 'text-destructive' : 'text-muted-foreground'
                )}
              >
                {userCount > 0 ? (
                  <Trans>保存后将立即影响 {userCount} 名在用用户</Trans>
                ) : (
                  <Trans>该角色当前没有用户，保存后仅对后续分配生效</Trans>
                )}
              </div>
            )}
            {addedCount > 0 && (
              <div className='space-y-1'>
                <div className='text-muted-foreground text-xs'>
                  <Trans>新增 {addedCount} 项</Trans>
                </div>
                <NodeBadges nodes={diff.added} />
              </div>
            )}
            {removedCount > 0 && (
              <div className='space-y-1'>
                <div className='text-destructive text-xs'>
                  <Trans>收回 {removedCount} 项</Trans>
                </div>
                <NodeBadges nodes={diff.removed} />
              </div>
            )}
          </div>
        </>
      )}
    </div>
  )
}
