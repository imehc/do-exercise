import { z } from 'zod'
import { t } from '@lingui/core/macro'
import {
  getPasswordRule,
  getUsernameRule,
} from '~/features/user/schemas/action-schema'

/**
 * 租户编码的约束单点：与服务端 `CreateSysTenantReq.Code` 的
 * `binding:"required,alphanum,min=2,max=32"`（`server/model/system/request/sys-tenant.go:9`）
 * 以及 `sys_tenant.code` 的 `varchar(32)` 保持一致。
 * 编码只在创建时可设（`UpdateSysTenantReq` 没有 Code），所以这里是唯一的写入口径。
 */
const TENANT_CODE_MIN = 2
const TENANT_CODE_MAX = 32
const TENANT_CODE_PATTERN = /^[a-zA-Z0-9]+$/

export const getTenantCodeRule = () =>
  z
    .string({
      error: t`请输入租户编码`,
    })
    .trim()
    .min(TENANT_CODE_MIN, { error: t`租户编码长度不能小于2个字符` })
    .max(TENANT_CODE_MAX, { error: t`租户编码长度不能大于32个字符` })
    .regex(TENANT_CODE_PATTERN, { error: t`租户编码只能包含字母和数字` })

/**
 * 登录页的租户编码是选填的：留空表示按用户名自动解析租户，
 * 非空时必须与创建时同一规则——否则前端会放过一个系统里根本不可能存在的编码，
 * 用户拿到的是「租户不存在」而不是「编码只能是字母和数字」。
 */
export const getOptionalTenantCodeRule = () =>
  z
    .string()
    .trim()
    .optional()
    .refine(
      (value) =>
        !value ||
        (value.length >= TENANT_CODE_MIN &&
          value.length <= TENANT_CODE_MAX &&
          TENANT_CODE_PATTERN.test(value)),
      { error: t`租户编码为 2~32 位字母或数字，留空则自动识别租户` }
    )

export const getSchema = () =>
  z
    .object({
      name: z
        .string({
          error: t`请输入租户名称`,
        })
        .trim()
        .min(2, { error: t`租户名称长度不能小于2个字符` })
        .max(20, { error: t`租户名称长度不能大于20个字符` }),
      code: getTenantCodeRule(),
      // 管理员账号模式：new=新建用户，existing=选择现有用户
      adminMode: z.string().default('new').optional(),
      // 选择现有用户时的用户 ID（admin_mode=existing）
      adminUserId: z.string().optional(),
      // 管理员用户名/密码仅在创建租户 new 模式下需要，编辑模式下不渲染，
      // 故字段级设可选，必填与格式校验放在 .check() 中仅对创建模式生效。
      adminUsername: z.string().optional(),
      adminPassword: z.string().optional(),
      confirmPassword: z.string().optional(),
      status: z.boolean().optional(),
      remark: z
        .string()
        .max(200, { error: t`备注长度不能超过200个字符` })
        .optional(),
      isEdit: z.boolean().default(false).optional(),
    })
    .check(({ issues, value }) => {
      if (value.isEdit) {
        return
      }

      // 选择现有用户模式：仅需选择用户
      if (value.adminMode === 'existing') {
        if (!value.adminUserId || value.adminUserId.trim() === '') {
          issues.push({
            code: 'custom',
            path: ['adminUserId'],
            message: t`请选择现有用户`,
            input: value,
          })
        }
        return
      }

      // 新建用户模式：校验管理员用户名是否填写并符合规则
      if (!value.adminUsername || value.adminUsername.trim() === '') {
        issues.push({
          code: 'custom',
          path: ['adminUsername'],
          message: t`请输入管理员用户名`,
          input: value,
        })
      } else {
        const usernameResult = getUsernameRule().safeParse(value.adminUsername)
        if (!usernameResult.success) {
          usernameResult.error.issues.forEach((issue) => {
            issues.push({
              code: 'custom',
              path: ['adminUsername'],
              message: issue.message,
              input: value.adminUsername,
            })
          })
        }
      }

      if (!value.adminPassword || value.adminPassword.trim() === '') {
        issues.push({
          code: 'custom',
          path: ['adminPassword'],
          message: t`请输入密码`,
          input: value,
        })
      } else {
        const passwordResult = getPasswordRule().safeParse(value.adminPassword)
        if (!passwordResult.success) {
          passwordResult.error.issues.forEach((issue) => {
            issues.push({
              code: 'custom',
              path: ['adminPassword'],
              message: issue.message,
              input: value.adminPassword,
            })
          })
        }
      }

      if (!value.confirmPassword || value.confirmPassword.trim() === '') {
        issues.push({
          code: 'custom',
          path: ['confirmPassword'],
          message: t`请输入确认密码`,
          input: value,
        })
      }

      if (
        value.adminPassword &&
        value.confirmPassword &&
        value.adminPassword !== value.confirmPassword
      ) {
        issues.push({
          code: 'custom',
          path: ['confirmPassword'],
          message: t`两次输入的密码不一致`,
          input: value,
        })
      }
    })

export type ActionTenantFormValues = z.infer<ReturnType<typeof getSchema>>
