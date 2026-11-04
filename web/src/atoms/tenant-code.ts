import { atomWithStorage } from 'jotai/utils'

/**
 * 登录页上次使用的租户编码。
 *
 * 只是一个输入框的记忆，不参与任何鉴权：服务端仍会校验编码对应的租户是否存在、
 * 是否启用，以及该租户下的用户名口令。填了它登录就只在该租户内验一次口令，
 * 避免同名账号跨租户逐行 bcrypt 比对（P2-6）。
 */
export const lastTenantCodeAtom = atomWithStorage('lastTenantCode', '')
