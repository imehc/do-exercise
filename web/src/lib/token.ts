import { addSeconds } from 'date-fns'
import type { LoginResult, Token } from '~/do-exercise-api'

// applyToken 将登录/切换租户返回的 LoginResult 归一化为 Token 并写入前端缓存结构：
// 换算绝对过期时间，并给可选字段补默认值（登录成功路径必然携带 token 字段）。
export function applyToken(token: LoginResult): Token {
  return {
    ...token,
    accessToken: token.accessToken ?? '',
    refreshToken: token.refreshToken ?? '',
    mustChangePassword: token.mustChangePassword ?? false,
    expireTime: addSeconds(new Date(), token.expireTime ?? 0).getTime(),
    refreshExpireTime: addSeconds(
      new Date(),
      token.refreshExpireTime ?? 0
    ).getTime(),
  } as Token
}