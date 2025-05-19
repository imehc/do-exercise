export const callIsInternalIpTypes = new Map<true | false, string>([
  [
    true,
    'bg-yellow-100/30 text-yellow-800 dark:text-yellow-200 border-yellow-200',
  ],
  [
    false,
    'bg-emerald-100/30 text-emerald-900 dark:text-emerald-200 border-emerald-200',
  ],
])

export const callIsMobileTypes = new Map<true | false, string>([
  [
    true,
    'bg-purple-100/30 text-purple-900 dark:text-purple-200 border-purple-200',
  ],
  [false, 'bg-sky-100/30 text-sky-900 dark:text-sky-200 border-sky-200'],
])

export const callIsBotTypes = new Map<true | false, string>([
  [true, 'bg-amber-100/30 text-amber-900 dark:text-amber-200 border-amber-200'],
  [
    false,
    'bg-emerald-100/30 text-emerald-900 dark:text-emerald-200 border-emerald-200',
  ],
])

export const callIsSuccessTypes = new Map<true | false, string>([
  [true, 'bg-rose-100/30 text-rose-900 dark:text-rose-200 border-rose-200'],
  [false, 'bg-cyan-100/30 text-cyan-900 dark:text-cyan-200 border-cyan-200'],
])

export const callCodeTypes = new Map<number, string>([
  [
    200,
    'bg-green-100/30 text-green-900 dark:text-green-200 border-green-200', // 成功
  ],
  [
    400,
    'bg-yellow-100/30 text-yellow-900 dark:text-yellow-200 border-yellow-200', // 客户端错误（语义错误、请求无效）
  ],
  [
    401,
    'bg-orange-100/30 text-orange-900 dark:text-orange-200 border-orange-200', // 未认证（需要登录）
  ],
  [
    403,
    'bg-pink-100/30 text-pink-900 dark:text-pink-200 border-pink-200', // 被禁止（权限问题）
  ],
  [
    404,
    'bg-slate-100/30 text-slate-900 dark:text-slate-200 border-slate-200', // 找不到资源
  ],
  [
    500,
    'bg-destructive/10 dark:bg-destructive/50 text-destructive dark:text-primary border-destructive/10', // 服务器错误
  ],
])
