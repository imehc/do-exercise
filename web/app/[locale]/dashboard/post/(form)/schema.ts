import { z } from 'zod';
import { Status } from '~/do-exercise-api';

export const postSchema = z.object({
  name: z
    .string({ required_error: 'nameFail.required' })
    .min(1, 'nameFail.required'),
  code: z
    .string({ required_error: 'codeFail.required' })
    .min(1, 'codeFail.required'),
  remark: z.string().optional(),
  status: z
    .union([z.literal(Status.disabled), z.literal(Status.enabled)])
    .optional(),
  sort: z.coerce.number().optional(),
});

export type PostFormValues = z.infer<typeof postSchema>;
