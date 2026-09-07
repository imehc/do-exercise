import { z } from 'zod'
import { t } from '@lingui/core/macro'

export const getPasswordRule = () =>
  z
    .string({
      error: t`密码不能为空`,
    })
    .min(6, { error: t`密码至少为6个字符` })
    .max(16, { error: t`密码最多为16个字符` })
    .regex(/[a-zA-Z]/, { error: t`密码必须包含至少一个字母` })
    .regex(/[0-9]/, { error: t`密码必须包含至少一个数字` })
    // eslint-disable-next-line no-useless-escape
    .regex(/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/, {
      error: t`密码必须包含至少一个特殊字符`,
    })

export const getUsernameRule = () =>
  z
    .string({
      error: t`请输入用户名`,
    })
    .min(2, { error: t`用户名长度不能小于2个字符` })
    .max(10, { error: t`用户名长度不能大于10个字符` })
    .regex(/^[a-zA-Z0-9]+$/, { error: t`用户名不能包含特殊字符` })
    .regex(/^[a-zA-Z]/, { error: t`用户名必须以字母开头` })
    .regex(/[a-zA-Z]/, { error: t`用户名必须包含字母` })

export const getBaseSchema = () =>
  z.object({
    password: z.string().optional(),
    confirmPassword: z.string().optional(),
  })

export const getSchema = (checkUsername?: (username: string) => Promise<boolean>) =>
  getBaseSchema()
    .extend({
      username: getUsernameRule(),
      nickname: z
        .string()
        .max(10, { error: t`昵称长度不能超过10个字符` })
        .optional(),
      email: z
        .email({
          error: (issue) =>
            issue.input === undefined ? t`请输入您的邮箱` : t`邮箱无效`,
        })
        .optional()
        .or(z.literal('')),
      avatar: z.string().optional(),
      roleIds: z.array(z.number(), {
        error: t`请选择关联角色`,
      }),

      isEdit: z.boolean().default(false).optional(),
    })
    .check(async ({ issues, value }) => {
      if (value.isEdit) {
        return
      }

      // 创建用户允许不选择角色（角色由后续分配）

      // 校验用户名是否已存在（仅创建模式）
      if (value.username && checkUsername) {
        try {
          const exists = await checkUsername(value.username)
          if (exists) {
            issues.push({
              code: 'custom',
              path: ['username'],
              message: t`用户名已存在`,
              input: value.username,
            })
          }
        } catch {
          // 查重接口异常时放行，交由后端兜底校验
        }
      }

      // 校验密码是否填写并符合规则
      if (!value.password || value.password.trim() === '') {
        issues.push({
          code: 'custom',
          path: ['password'],
          message: t`请输入密码`,
          input: value,
        })
      } else {
        // 手动触发所有密码规则
        const passwordResult = getPasswordRule().safeParse(value.password)
        if (!passwordResult.success) {
          passwordResult.error.issues.forEach((issue) => {
            issues.push({
              code: 'custom',
              path: ['password'], // 强制绑定到 password 字段
              message: issue.message,
              input: value.password,
            })
          })
        }
      }

      // 校验确认密码是否填写
      if (!value.confirmPassword || value.confirmPassword.trim() === '') {
        issues.push({
          code: 'custom',
          path: ['confirmPassword'],
          message: t`请输入确认密码`,
          input: value,
        })
      }

      // 校验两次密码是否一致
      if (
        value.password &&
        value.confirmPassword &&
        value.password !== value.confirmPassword
      ) {
        issues.push({
          code: 'custom',
          path: ['confirmPassword'],
          message: t`两次输入的密码不一致`,
          input: value,
        })
      }
    })

export const resetPasswordSchema = z
  .object({
    password: getPasswordRule(),
    confirmPassword: getPasswordRule(),
  })
  .check(({ issues, value }) => {
    if (
      value.password &&
      value.confirmPassword &&
      value.password !== value.confirmPassword
    ) {
      issues.push({
        code: 'custom',
        path: ['confirmPassword'],
        message: t`两次输入的密码不一致`,
        input: value,
      })
    }
  })

export type ActionSysUserFormValues = z.infer<ReturnType<typeof getSchema>>
export type ActionResetPasswordSysUserFormValues = z.infer<
  typeof resetPasswordSchema
>
