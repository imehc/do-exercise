'use client';

import { Button, ButtonProps } from '~/components/ui/button';
import { Icon } from '@iconify/react';

interface SubmieButtonProps extends ButtonProps {
  loading?: boolean;
}

export default function SubmitButton({
  loading = false,
  disabled,
  children,
  ...props
}: SubmieButtonProps) {
  return (
    <Button type="submit" {...props} disabled={disabled || loading}>
      {loading && <Icon icon="line-md:loading-loop" className="animate-spin" />}
      {children}
    </Button>
  );
}
