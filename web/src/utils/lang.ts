import { detect, fromUrl, fromNavigator } from '@lingui/detect-locale'

export function getLang(): string {
  return detect(fromUrl('lang'), fromNavigator()) || 'zh-CN'
}
