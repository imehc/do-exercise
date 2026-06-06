/**
 * 检查字符串是否以 http:// 或 https:// 开头，如果不是则添加指定的前缀
 * @param url 要检查的字符串
 * @returns 处理后的字符串
 */
export function ensureHttpPrefix(url?: string) {
  if (!url) return
  if (
    url.startsWith('http://') ||
    url.startsWith('https://') ||
    url.startsWith('data:') ||
    url.startsWith('blob:')
  ) {
    return url
  }
  if (url.startsWith('/oss/')) {
    return url
  }
  return `/oss${url.startsWith('/') ? url : `/${url}`}`
}
