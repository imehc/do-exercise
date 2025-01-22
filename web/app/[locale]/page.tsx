import { Button } from '~/components/ui/button';
import { Link } from '~/i18n/routing';

export default function HomePage() {
  return (
    <div>
      <Link href="/auth/signin">
        <Button>signin</Button>
      </Link>
    </div>
  );
}
