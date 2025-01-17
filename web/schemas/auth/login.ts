import { z } from 'zod';

export const loginSchema = z.object({
	username: z
		.string()
		.trim()
		.min(1, '用户名不能为空')
		.min(4, '用户名最少为4个字符')
		.max(8, '用户名最多为8个字符')
		.regex(/^[a-zA-Z0-9]+$/, '用户名只能包含字母和数字')
		.regex(/^[a-zA-Z]/, '用户名必须以字母开头'),
	password: z
		.string()
		.trim()
		.min(1, '密码不能为空')
		.min(6, '密码最少为6个字符')
		.max(16, '密码最多为16个字符'),
	captcha: z.string().trim().min(1, '验证码不能为空'),
	captchaId: z
		.string({ required_error: '验证码Id不能为空' })
		.min(1, '验证码Id不能为空'),
});
