import { AuthApi } from '~/do-exercise-api'
import { useApi } from '~/hooks'

export function Dashboard() {
  const authApi = useApi(AuthApi)
  authApi.login({
    login: {
      username: '',
      password: '',
      captchaId: '',
      captcha: '',
      publicKey: ''
    }
  })
  authApi
  return (
    <div>
      <h1>PageDashboard</h1>
    </div>
  )
}

export default Dashboard
