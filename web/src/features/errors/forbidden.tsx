import { useNavigate, useRouter } from '@tanstack/react-router'
import { Trans } from '@lingui/react/macro'
import { Button } from '~/components/ui/button'

export default function ForbiddenError() {
  const navigate = useNavigate()
  const { history } = useRouter()
  return (
    <div className='h-svh'>
      <div className='m-auto flex h-full w-full flex-col items-center justify-center gap-2'>
        <h1 className='text-[7rem] leading-tight font-bold'>403</h1>
        <span className='font-medium'>Access Forbidden</span>
        <p className='text-muted-foreground text-center'>
          You don't have necessary permission <br />
          to view this resource.
        </p>
        <div className='mt-6 flex gap-4'>
          <Button variant='outline' onClick={() => history.go(-1)}>
            <Trans>返回</Trans>
          </Button>
          <Button onClick={() => navigate({ to: '/' })}>
            <Trans>返回首页</Trans>
          </Button>
        </div>
      </div>
    </div>
  )
}
