import { JSEncrypt } from 'jsencrypt'

/**
 * 使用 RSA 公钥加密密码
 * @param password 明文密码
 * @param base64PublicKey Base64 编码的 PEM 公钥（Go 返回）
 * @returns 加密后的密文，如果失败则返回 null
 */
export function encryptPassword(
  password: string,
  base64PublicKey: string
): string | null {
  // 注意：Go 返回的可能是 PEM 格式的 base64，需要先解码
  try {
    // 浏览器环境用 atob，Node 用 Buffer
    const publicKeyPem = atob(base64PublicKey) // 解码 base64 为 PEM 格式字符串
    const encrypt = new JSEncrypt()
    encrypt.setPublicKey(publicKeyPem)

    const encrypted = encrypt.encrypt(password)
    return encrypted || null
  } catch (error) {
    console.error('Error during encryption:', error)
    return null
  }
}
