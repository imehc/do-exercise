import { Suspense } from 'react'
import { Link, Outlet } from 'react-router'
import { Loading } from '~/components'
import { useRouterMenus } from '~/provider'

const LayoutPage: React.FC = () => {
  const menus = useRouterMenus()

  return (
    <div>
      <ul>
        {menus.map(menu => (
          <li key={menu.route}>
            <Link to={menu.route}>{menu.name}</Link>
          </li>
        ))}
      </ul>
      <Suspense fallback={<Loading />}>
        <Outlet />
      </Suspense>
    </div>
  )
}

export default LayoutPage
