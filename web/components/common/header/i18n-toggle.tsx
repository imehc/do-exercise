'use client';

import { Icon } from '@iconify/react';
import { Button } from '~/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '~/components/ui/dropdown-menu';
import { usePathname, useRouter } from '~/i18n/routing';
import { languages } from '~/i18n/i18';

export function I18nToggle() {
  const pathname = usePathname();
  const router = useRouter();

  const onValueChange = (newLocale: string) => {
    router.replace(pathname, { locale: newLocale });
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="icon">
          <Icon
            className="absolute h-[1.2rem] w-[1.2rem]"
            icon="material-symbols:translate"
          />
          <span className="sr-only">Toggle language</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        {languages.map((item) => (
          <DropdownMenuItem key={item.code} onClick={() => onValueChange(item.lang)}>
            {item.label}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
