import { Button } from '~/components/ui/button';
import { Link } from '~/i18n/routing';

export default function HomePage() {
  return (
    <div>
      <Link href="/login">
        <Button>Login</Button>
      </Link>
    </div>
  );
}
