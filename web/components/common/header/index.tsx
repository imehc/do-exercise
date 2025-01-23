import type { PropsWithChildren } from 'react';
import { ModeToggle } from '~/components/common/header/model-toggle';
import { I18nToggle } from '~/components/common/header/i18n-toggle';
import { SignOutButton } from './sign-out-button';

export default function Header({ children }: PropsWithChildren) {
  return (
    <div className="flex flex-col h-screen">
      <div className="h-14 flex justify-between items-center px-6 shadow">
        <div className="flex justify-start items-center">{/* TODO: */}</div>
        <div className="flex justify-end items-center flex-row-reverse gap-x-3">
          <SignOutButton />
          <ModeToggle />
          <I18nToggle />
        </div>
      </div>
      <div className="flex-1">{children}</div>
    </div>
  );
}
