// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
	compatibilityDate: '2024-11-01',
	devtools: { enabled: true },
	typescript: {
		typeCheck: true,
	},
	debug: import.meta.env.DEV,
	devServer: {
		port: 6021,
		host: '0.0.0.0',
	},
	modules: [
		'@varlet/nuxt',
		'@nuxt/eslint',
		'@nuxtjs/tailwindcss',
		'@formkit/auto-animate',
		'@vueuse/nuxt',
		'@nuxt/image',
		'nuxt-auth-utils',
		'@vueuse/nuxt',
		'@nuxtjs/i18n',
		'@pinia/nuxt',
		'pinia-plugin-persistedstate',
		'@nuxt/icon',
	],
	app: {
		pageTransition: { name: 'page', mode: 'out-in' },
		layoutTransition: { name: 'layout', mode: 'out-in' },
	},
	// https://nuxt.com/docs/getting-started/transitions#view-transitions-api-experimental
	experimental: {
		viewTransition: true,
	},
	i18n: {
		detectBrowserLanguage: {
			useCookie: true,
			cookieKey: 'i18n_redirected',
			redirectOn: 'root', // recommended
		},
		strategy: 'prefix_except_default',
		locales: [
			{ code: 'en', language: 'en-US' },
			{ code: 'zh', language: 'zh-CN' },
		],
		lazy: true,
		defaultLocale: 'zh',
		vueI18n: './i18n.config.ts',
	},
	icon: { componentName: 'NuxtIcon' }, // https://icon-sets.iconify.design/
});
