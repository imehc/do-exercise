import { auth } from '~/helper/auth';

export default async function DashboardPage() {
  const session = await auth();
  return <div>
    <div>DashboardPage</div>
    <div>{session?.accessToken}</div>
    <div>{session?.refreshToken}</div>
    <div>{session?.expiresIn}</div>
    <div>{session?.tokenType}</div>
  </div>;
}
