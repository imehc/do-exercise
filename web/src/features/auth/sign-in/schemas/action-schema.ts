import { z } from 'zod'
import { t } from '@lingui/core/macro'
import { getOptionalTenantCodeRule } from '~/features/tenant/schemas/action-schema'
import {
  getPasswordRule,
  getUsernameRule,
} from '~/features/user/schemas/action-schema'

export const getSignInActionSchema = () =>
  z.object({
    username: getUsernameRule(),
    password: getPasswordRule(),
    captchaId: z.string({
      error: t`验证码ID不能为空`,
    }),
    tenantId: z.string().optional(),
    // 租户编码选填。填了就直接进该租户（服务端只验一次口令），留空则按用户名
    // 归属的启用租户解析，多个归属时弹窗选择。
    // 规则复用创建租户的口径（2~32 位字母或数字），避免放过不可能存在的编码。
    tenantCode: getOptionalTenantCodeRule(),
    captcha: z
      .string({
        error: t`请输入验证码`,
      })
      .min(1, { error: t`请输入验证码` }),
    publicKey: z.string({
      error: t`公钥不能为空`,
    }),
  })

export type SignInActionFormValues = z.infer<
  ReturnType<typeof getSignInActionSchema>
>
