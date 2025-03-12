'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '~/components/ui/card';
import { Input } from '~/components/ui/input';
import Image from 'next/image';
import SubmitButton from '~/components/common/submit-button';
import { useLocale, useTranslations } from 'next-intl';
import { type SigninSchema, signinSchema } from './schema';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '~/components/ui/form';
import { signinAction } from './action';
import type { Captcha } from '~/do-exercise-api';
import type { ResponseData } from '~/helper/format-response';
import { toast } from 'sonner';
import { useDebounceEffect } from 'ahooks';
import clsx from 'clsx';
import { useEffect } from 'react';
import { useRouter } from '~/i18n/routing';

interface SigninFormProps {
  captchaRes: ResponseData<Captcha>;
}

export default function SigninForm({ captchaRes }: SigninFormProps) {
  const t = useTranslations('SigninPage');

  const form = useForm<SigninSchema>({
    resolver: zodResolver(signinSchema),
    defaultValues: {
      username: '',
      password: '',
      captcha: '',
      captchaId: '',
    },
  });

  useDebounceEffect(() => {
    if (captchaRes.ok) return;
    toast.error(captchaRes.message);
  }, [captchaRes]);

  const router = useRouter();
  const locale = useLocale();

  const isSubmitting = form.formState.isValid && form.formState.isSubmitting;
  const isPending = form.formState.isSubmitSuccessful || isSubmitting;

  const onSubmit = async (data: SigninSchema) => {
    const res = await signinAction(data);
    if (!res.ok) {
      if (res.message) {
        toast.error(res.message);
      } else if (res.i18n) {
        toast.error(t(res.message));
      }
      router.refresh();
      form.reset();
      throw new Error('signin fail');
    }
    if (res.ok) {
      if (res.href) {
        router.replace(res.href, { locale });
      }
      if (res.i18n) {
        toast.success(t(res.i18n));
      }
    }
  };

  useEffect(() => {
    if (!captchaRes.ok) return;
    form.setValue('captchaId', captchaRes.data.captchaId);
  }, [captchaRes, form]);

  return (
    <div className="flex justify-center items-center h-full">
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)}>
          <Card className="w-[350px]">
            <CardHeader>
              <CardTitle>{t('title')}</CardTitle>
              <CardDescription>{t('subTitle')}</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid w-full items-center gap-4">
                <div className="flex flex-col space-y-1.5">
                  <FormField
                    control={form.control}
                    name="username"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('account')}</FormLabel>
                        <FormControl>
                          <Input
                            disabled={isPending}
                            startIcon="material-symbols:person-4-outline"
                            placeholder={t('accountPlaceholder')}
                            {...field}
                          />
                        </FormControl>
                        <FormMessage namespace="SigninPage" />
                      </FormItem>
                    )}
                  />
                </div>
                <div className="flex flex-col space-y-1.5">
                  <FormField
                    control={form.control}
                    name="password"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('password')}</FormLabel>
                        <FormControl>
                          <Input
                            disabled={isPending}
                            startIcon="material-symbols:lock-outline"
                            type="password"
                            placeholder={t('passwordPlaceholder')}
                            {...field}
                          />
                        </FormControl>
                        <FormMessage namespace="SigninPage" />
                      </FormItem>
                    )}
                  />
                </div>
                <div className="flex justify-between items-end gap-x-3">
                  <div className="flex flex-col space-y-1.5">
                    <FormField
                      control={form.control}
                      name="captcha"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>{t('captcha')}</FormLabel>
                          <FormControl>
                            <Input
                              disabled={isPending}
                              startIcon="material-symbols:health-and-safety-outline"
                              placeholder={t('captchaPlaceholder')}
                              {...field}
                            />
                          </FormControl>
                          <FormMessage namespace="SigninPage" />
                        </FormItem>
                      )}
                    />
                  </div>
                  <div className="h-9 aspect-[3/1] relative rounded-md">
                    {captchaRes.ok ? (
                      <Image
                        className={clsx([
                          !isPending
                            ? 'cursor-pointer pointer-events-auto'
                            : 'pointer-events-none cursor-none',
                        ])}
                        fill
                        src={captchaRes.data.picPath}
                        alt="captcha"
                        onClick={() => router.refresh()}
                      />
                    ) : null}
                  </div>
                </div>
              </div>
            </CardContent>
            <CardFooter className="flex justify-end gap-x-3">
              <SubmitButton className="w-full" loading={isPending}>
                {t('signIn')}
              </SubmitButton>
            </CardFooter>
          </Card>
        </form>
      </Form>
    </div>
  );
}
