import { PropsWithChildren } from 'react'
import { WithPermission } from '~/provider'
import { Main } from '../layout/main'

interface Props extends PropsWithChildren {
  title?: React.ReactNode
  subTitle?: React.ReactNode
  element?: React.ReactNode
  actionElemnt?: React.ReactNode
}

export function MainBody({
  title,
  subTitle,
  element,
  actionElemnt,
  children,
}: Props) {
  return (
    <>
      <Main>
        <div className='mb-2 flex flex-wrap items-center justify-between space-y-2'>
          <div>
            <h2 className='text-2xl font-bold tracking-tight'>{title}</h2>
            <p className='text-muted-foreground'>{subTitle}</p>
          </div>
          <WithPermission permission='create'>{element}</WithPermission>
        </div>
        <div className='-mx-4 flex-1 overflow-auto px-4 py-1 lg:flex-row lg:space-y-0 lg:space-x-12'>
          {children}
        </div>
      </Main>
      {actionElemnt}
    </>
  )
}
