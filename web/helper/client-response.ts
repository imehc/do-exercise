import { useTranslations } from 'next-intl';
import { ResponseData } from './format-response';
import { toast } from 'sonner';

interface Props {
  res: ResponseData;
  t: ReturnType<typeof useTranslations>;
  onFail?(): void;
  onSuccess?(): void;
}

export function handleClientResponse({ res, t, onFail, onSuccess }: Props) {
  if (!res.ok) {
    if (res.message) {
      toast.error(res.message);
    } else if (res.i18n) {
      toast.error(t(res.message));
    }
    onFail?.();
  }
  if (res.ok) {
    if (res.i18n) {
      toast.success(t(res.i18n));
    }
    onSuccess?.();
  }
}
