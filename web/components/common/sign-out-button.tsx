import { Icon } from '@iconify/react/dist/iconify.js';
import { Button } from '~/components/ui/button';
import { auth, signOut } from '~/helper/auth';

export async function SignOutButton() {
  const session = await auth();

  if (!session?.user) {
    return null;
  }

  return (
    <form
      action={async () => {
        'use server';
        await signOut();
      }}
    >
      <Button variant="outline" size="icon" type="submit">
        <Icon
          className="absolute h-[1.2rem] w-[1.2rem]"
          icon="material-symbols:logout"
        />
        <span className="sr-only">Sign Out</span>
      </Button>
    </form>
  );
}
