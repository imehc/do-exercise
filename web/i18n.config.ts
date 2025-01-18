import en from './lang/en.json';
import zh from './lang/zh.json';

type MessageSchema = typeof en;

export default defineI18nConfig(() => ({
	legacy: false, // 是否兼容之前
	locale: 'zh',
	fallbackLocale: 'zh',
	messages: {
		en,
		zh,
	},
}));
