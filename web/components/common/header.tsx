import type { PropsWithChildren } from 'react';
import { ModeToggle } from '~/components/model-toggle';
import { I18nToggle } from '~/components/i18n-toggle';

export default function Header({ children }: PropsWithChildren) {
  return (
    <div className="flex flex-col h-screen">
      <div className="h-14 flex justify-between items-center px-6 shadow">
        <div className="flex justify-start items-center">{/* TODO: */}</div>
        <div className="flex justify-end items-center flex-row-reverse gap-x-3">
          <ModeToggle />
          <I18nToggle />
        </div>
      </div>
      <div className="flex-1">{children}</div>
    </div>
  );
}
