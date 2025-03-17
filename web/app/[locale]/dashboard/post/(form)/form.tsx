'use client';

import { useForm } from 'react-hook-form';
import { zodResolver as resolver } from '@hookform/resolvers/zod';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '~/components/ui/form';
import { Input } from '~/components/ui/input';
import { PostFormValues, postSchema } from './schema';
import { Button } from '~/components/ui/button';
import { useTranslations } from 'next-intl';
import { Switch } from '~/components/ui/switch';
import { Post, Status } from '~/do-exercise-api';
import { Link, useRouter } from '~/i18n/routing';
import { editPostAction } from './[id]/actions';
import type { ResponseData } from '~/helper/format-response';
import { createPostAction } from './new/actions';
import { handleClientResponse } from '~/helper/client-response';

interface Props {
  data?: Post;
}

export function PostForm({ data: values }: Props) {
  const t = useTranslations('PostForm');
  const form = useForm<PostFormValues>({
    resolver: resolver(postSchema),
    defaultValues: values,
  });

  const router = useRouter();

  const isSubmitting = form.formState.isValid && form.formState.isSubmitting;
  const isPending = form.formState.isSubmitSuccessful || isSubmitting;

  const onSubmit = async (data: PostFormValues) => {
    let res: ResponseData;
    if (values?.id) {
      res = await editPostAction(values.id, data);
    } else {
      res = await createPostAction(data);
    }
    handleClientResponse({
      res,
      t,
      onFail: () => {
        form.reset();
      },
      onSuccess: () => {
        form.reset();
        form.clearErrors();
        router.back();
      },
    });
  };

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="max-w-2xl mx-auto p-6 space-y-6 rounded-lg"
      >
        <h2 className="text-2xl font-semibold text-center mb-6">
          {values?.id ? t('updateTitle') : t('createTitle')}
        </h2>
        <div className="grid gap-6">
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel className="text-sm font-medium">
                  {t('name')}
                </FormLabel>
                <FormControl>
                  <Input
                    disabled={isPending}
                    className="w-full"
                    placeholder={t('namePlaceholder')}
                    {...field}
                  />
                </FormControl>
                <FormMessage namespace="PostForm" />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="code"
            render={({ field }) => (
              <FormItem>
                <FormLabel className="text-sm font-medium">
                  {t('code')}
                </FormLabel>
                <FormControl>
                  <Input
                    disabled={isPending}
                    className="w-full"
                    placeholder={t('codePlaceholder')}
                    {...field}
                  />
                </FormControl>
                <FormMessage namespace="PostForm" />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="remark"
            render={({ field }) => (
              <FormItem>
                <FormLabel className="text-sm font-medium">
                  {t('remark')}
                </FormLabel>
                <FormControl>
                  <Input
                    disabled={isPending}
                    className="w-full"
                    placeholder={t('remarkPlaceholder')}
                    {...field}
                  />
                </FormControl>
                <FormMessage namespace="PostForm" />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="status"
            render={({ field }) => (
              <FormItem className="flex flex-row items-center justify-between py-1 px-3 rounded-lg border bg-card/50 shadow-sm">
                <div>
                  <FormLabel className="text-sm font-medium">
                    {t('status')}
                  </FormLabel>
                </div>
                <FormControl>
                  <Switch
                    className="!mt-0"
                    disabled={isPending}
                    checked={field.value === Status.enabled}
                    onCheckedChange={(val) => {
                      field.onChange(val ? Status.enabled : Status.disabled);
                    }}
                  />
                </FormControl>
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="sort"
            render={({ field }) => (
              <FormItem>
                <FormLabel className="text-sm font-medium">
                  {t('sort')}
                </FormLabel>
                <FormControl>
                  <Input
                    disabled={isPending}
                    className="w-full"
                    placeholder={t('sortPlaceholder')}
                    {...field}
                  />
                </FormControl>
                <FormMessage namespace="PostForm" />
              </FormItem>
            )}
          />
        </div>
        <div className="flex justify-end gap-x-3 pt-4">
          <Link href="../post">
            <Button
              type="button"
              variant="outline"
              className="min-w-[120px]"
              onClick={() => router.back()}
            >
              {t('cancel')}
            </Button>
          </Link>
          <Button type="submit" className="min-w-[120px]" disabled={isPending}>
            {t('submit')}
          </Button>
        </div>
      </form>
    </Form>
  );
}
