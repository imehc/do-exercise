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
import { useTranslations } from 'next-intl';
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
import type { CaptchaResponse } from '~/do-exercise-api';
import type { ResponseData } from '~/helper/format-response';
import { toast } from 'sonner';
import { useDebounceEffect } from "ahooks";

interface SigninFormProps {
  captchaRes: ResponseData<CaptchaResponse>;
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

  useDebounceEffect(()=>{
    if (captchaRes.ok) return;
    toast.error(captchaRes.message);
  },[captchaRes])

  if (captchaRes.ok) {
    form.setValue('captchaId', captchaRes.data.captchaId);
  }

  return (
    <div className="flex justify-center items-center h-full">
      <Form {...form}>
        <form onSubmit={form.handleSubmit(signinAction)}>
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
                    {captchaRes.ok && (
                      <Image fill src={captchaRes.data.picPath} alt="captcha" />
                    )}
                  </div>
                </div>
              </div>
            </CardContent>
            <CardFooter className="flex justify-end gap-x-3">
              <SubmitButton>{t('signIn')}</SubmitButton>
            </CardFooter>
          </Card>
        </form>
      </Form>
    </div>
  );
}
