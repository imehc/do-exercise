import { z } from 'zod'
import { t } from '@lingui/core/macro'

export const getPasswordRule = () =>
  z
    .string({ required_error: t`密码不能为空` })
    .min(6, { message: t`密码至少为6个字符` })
    .max(16, { message: t`密码最多为16个字符` })
    .regex(/[a-zA-Z]/, { message: t`密码必须包含至少一个字母` })
    .regex(/[0-9]/, { message: t`密码必须包含至少一个数字` })
    // eslint-disable-next-line no-useless-escape
    .regex(/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/, {
      message: t`密码必须包含至少一个特殊字符`,
    })

export const getUsernameRule = () =>
  z
    .string({ required_error: t`请输入用户名` })
    .min(2, { message: t`用户名长度不能小于2个字符` })
    .max(10, { message: t`用户名长度不能大于10个字符` })
    .regex(/^[a-zA-Z0-9]+$/, { message: t`用户名不能包含特殊字符` })
    .regex(/^[a-zA-Z]/, { message: t`用户名必须以字母开头` })
    .regex(/[a-zA-Z]/, { message: t`用户名必须包含字母` })

export const getBaseSchema = () =>
  z.object({
    password: z.string().optional(),
    confirmPassword: z.string().optional(),
  })

export const getSchema = () =>
  getBaseSchema()
    .extend({
      username: getUsernameRule(),
      nickname: z
        .string()
        .max(10, { message: t`昵称长度不能超过10个字符` })
        .optional(),
      email: z
        .string()
        .email({ message: t`请输入正确的邮箱地址` })
        .optional()
        .or(z.literal('')),
      avatar: z.string().optional(),
      roleIds: z
        .array(z.number(), {
          required_error: t`请选择关联角色`,
        })
        .min(1, t`请选择关联角色`),

      isEdit: z.boolean().default(false).optional(),
    })
    .superRefine((data, ctx) => {
      if (data.isEdit) {
        return true
      }

      // 校验密码是否填写并符合规则
      if (!data.password || data.password.trim() === '') {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path: ['password'],
          message: t`请输入密码`,
        })
      } else {
        // 手动触发所有密码规则
        const passwordResult = getPasswordRule().safeParse(data.password)
        if (!passwordResult.success) {
          passwordResult.error.issues.forEach((issue) => {
            ctx.addIssue({
              ...issue,
              path: ['password'], // 强制绑定到 password 字段
            })
          })
        }
      }

      // 校验确认密码是否填写
      if (!data.confirmPassword || data.confirmPassword.trim() === '') {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path: ['confirmPassword'],
          message: t`请输入确认密码`,
        })
      }

      // 校验两次密码是否一致
      if (
        data.password &&
        data.confirmPassword &&
        data.password !== data.confirmPassword
      ) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path: ['confirmPassword'],
          message: t`两次输入的密码不一致`,
        })
      }
    })

export const resetPasswordSchema = z
  .object({
    password: getPasswordRule(),
    confirmPassword: getPasswordRule(),
  })
  .superRefine((data, ctx) => {
    if (
      data.password &&
      data.confirmPassword &&
      data.password !== data.confirmPassword
    ) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        path: ['confirmPassword'],
        message: t`两次输入的密码不一致`,
      })
    }
  })

export type ActionSysUserFormValues = z.infer<ReturnType<typeof getSchema>>
export type ActionResetPasswordSysUserFormValues = z.infer<
  typeof resetPasswordSchema
>
