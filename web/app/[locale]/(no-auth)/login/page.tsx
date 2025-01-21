import { getTranslations } from 'next-intl/server';
import { Button } from '~/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '~/components/ui/card';
import { Input } from '~/components/ui/input';
import { Label } from '~/components/ui/label';
import Image from 'next/image';

export default async function LoginPage() {
  const t = await getTranslations('LoginPage');

  return (
    <div className="flex justify-center items-center h-full">
      <Card className="w-[350px]">
        <CardHeader>
          <CardTitle>{t('title')}</CardTitle>
          <CardDescription>{t('subTitle')}</CardDescription>
        </CardHeader>
        <CardContent>
          <form>
            <div className="grid w-full items-center gap-4">
              <div className="flex flex-col space-y-1.5">
                <Label htmlFor="username">{t('account')}</Label>
                <Input id="username" placeholder={t('accountPlaceholder')} />
              </div>
              <div className="flex flex-col space-y-1.5">
                <Label htmlFor="password">{t('password')}</Label>
                <Input
                  id="password"
                  placeholder={t('passwordPlaceholder')}
                  type="password"
                />
              </div>
              <div className="flex justify-between items-end gap-x-3">
                <div className="flex flex-col space-y-1.5">
                  <Label htmlFor="captcha">{t('captcha')}</Label>
                  <Input id="captcha" placeholder={t('captchaPlaceholder')} />
                </div>
                <div className="h-9 aspect-[3/1] relative">
                  <Image fill src="/next.svg" alt="captcha" />
                </div>
              </div>
            </div>
          </form>
        </CardContent>
        <CardFooter className="flex justify-end gap-x-3">
          <Button>{t('signOn')}</Button>
        </CardFooter>
      </Card>
    </div>
  );
}
