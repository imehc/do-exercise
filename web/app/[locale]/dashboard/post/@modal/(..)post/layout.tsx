'use client';

import { PropsWithChildren } from 'react';
import { Drawer, DrawerContent, DrawerTitle } from '~/components/ui/drawer';
import { VisuallyHidden } from '@radix-ui/react-visually-hidden';
import { useRouter } from '~/i18n/routing';

export default function PostLayout({ children }: PropsWithChildren) {
  const router = useRouter();
  const onClose = () => {
    router.back();
  };

  return (
    <Drawer open onClose={onClose}>
      <VisuallyHidden>
        <DrawerTitle>Post</DrawerTitle>
      </VisuallyHidden>
      <DrawerContent>{children}</DrawerContent>
    </Drawer>
  );
}
