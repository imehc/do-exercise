import { useState } from 'react'
import { useMutation, useQuery } from '@tanstack/react-query'
import { Link } from '@tanstack/react-router'
import { IconArrowsExchange } from '@tabler/icons-react'
import { Trans } from '@lingui/react/macro'
import { useAtomValue, useSetAtom } from 'jotai'
import { originTokenAtom } from '~/atoms'
import { AuthApi, TenantOption } from '~/do-exercise-api'
import { applyToken } from '~/lib/token'
import { encryptPassword } from '~/utils/encrypt'
import { ensureHttpPrefix } from '~/utils/url'
import { useApi } from '~/hooks/use-api'
import { usePublicKey } from '~/hooks/use-public-key'
import { useLogout, useUserProfile } from '~/hooks/use-user'
import { Avatar, AvatarFallback, AvatarImage } from '~/components/ui/avatar'
import { Button } from '~/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '~/components/ui/dialog'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '~/components/ui/dropdown-menu'
import { PasswordInput } from '~/components/password-input'

export function ProfileDropdown() {
  const { data: userProfile, isLoading: userProfileIsLoading } =
    useUserProfile()
  const { mutate: logout, isPending: logoutIsPending } = useLogout()
  const authApi = useApi(AuthApi)
  const token = useAtomValue(originTokenAtom)
  const setToken = useSetAtom(originTokenAtom)

  const [switchOpen, setSwitchOpen] = useState(false)
  const [switchPassword, setSwitchPassword] = useState('')
  const { data: publicKeyData } = usePublicKey()

  // 当前账号可切换的租户（用于顶栏展示当前租户与切换入口；不含平台保留租户）
  const { data: myTenants = [] } = useQuery({
    queryKey: ['myTenants'],
    queryFn: () => authApi.myTenants(),
    enabled: !!token.accessToken,
    retry: false,
  })

  const currentTenant = myTenants.find((t) => t.tenantId === token.tenantId)

  const { mutate: switchTenant, isPending: switchIsPending } = useMutation({
    mutationFn: (tenant: TenantOption) => {
      const publicKey = publicKeyData?.publicKey
      if (!publicKey) throw new Error('public key unavailable')
      const password = encryptPassword(switchPassword, publicKey)
      if (!password) throw new Error('password encryption failed')
      return authApi.switchTenant({
        switchTenantRequest: {
          tenantId: tenant.tenantId ?? '',
          password,
          publicKey,
        },
      })
    },
    onSuccess: (data) => {
      setSwitchOpen(false)
      setSwitchPassword('')
      setToken(applyToken(data))
      // 切换租户后整页刷新，让权限菜单与数据按新租户重建
      window.location.reload()
    },
  })

  const isPending = userProfileIsLoading || logoutIsPending

  return (
    <>
      <DropdownMenu modal={false}>
        <DropdownMenuTrigger asChild disabled={userProfileIsLoading}>
          <Button variant='ghost' className='relative h-8 w-8 rounded-full'>
            <Avatar className='h-8 w-8'>
              <AvatarImage
                src={ensureHttpPrefix(userProfile?.avatar)}
                alt={userProfile?.username}
              />
              <AvatarFallback>ST</AvatarFallback>
            </Avatar>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent className='w-56' align='end' forceMount>
          <DropdownMenuLabel className='font-normal'>
            <div className='flex flex-col space-y-1'>
              <p className='text-sm leading-none font-medium'>
                {userProfile?.username}
              </p>
              <p className='text-muted-foreground text-xs leading-none'>
                {currentTenant?.name ?? userProfile?.email}
              </p>
            </div>
          </DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuGroup>
            {myTenants.length > 1 && (
              <DropdownMenuItem
                className='cursor-pointer'
                onClick={() => setSwitchOpen(true)}
              >
                <IconArrowsExchange className='mr-2 h-4 w-4' />
                <Trans>切换租户</Trans>
              </DropdownMenuItem>
            )}
            <DropdownMenuItem asChild>
              <Link to='/settings'>
                <Trans>设置</Trans>
                {/* <DropdownMenuShortcut>⌘S</DropdownMenuShortcut> */}
              </Link>
            </DropdownMenuItem>
          </DropdownMenuGroup>
          <DropdownMenuSeparator />
          <DropdownMenuItem
            className='cursor-pointer'
            onClick={() => logout()}
            disabled={isPending}
          >
            <Trans>退出登录</Trans>
            {/* <DropdownMenuShortcut>⇧⌘Q</DropdownMenuShortcut> */}
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      <Dialog open={switchOpen} onOpenChange={setSwitchOpen}>
        <DialogContent className='sm:max-w-md'>
          <DialogHeader>
            <DialogTitle>
              <Trans>切换租户</Trans>
            </DialogTitle>
            <DialogDescription>
              <Trans>选择要进入的租户</Trans>
            </DialogDescription>
          </DialogHeader>
          <PasswordInput
            value={switchPassword}
            onChange={(event) => setSwitchPassword(event.target.value)}
            placeholder='请输入当前租户口令'
            disabled={switchIsPending}
          />
          <div className='grid gap-2'>
            {myTenants.map((tenant) => {
              const isCurrent = tenant.tenantId === token.tenantId
              return (
                <Button
                  key={tenant.tenantId}
                  variant={isCurrent ? 'secondary' : 'outline'}
                  className='flex h-auto items-center justify-between px-4 py-3'
                  disabled={isCurrent || switchIsPending || !switchPassword}
                  onClick={() => switchTenant(tenant)}
                >
                  <span>{tenant.name}</span>
                  <span className='text-muted-foreground text-xs'>
                    {tenant.code}
                  </span>
                </Button>
              )
            })}
          </div>
        </DialogContent>
      </Dialog>
    </>
  )
}
