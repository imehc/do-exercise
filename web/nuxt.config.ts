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
	],
	app: {
		pageTransition: { name: 'page', mode: 'out-in' },
		layoutTransition: { name: 'layout', mode: 'out-in' },
	},
	routeRules: {
		'/': { redirect: 'dashboard' },
		'/**': { appMiddleware: 'authenticated' },
	},
	// https://nuxt.com/docs/getting-started/transitions#view-transitions-api-experimental
	experimental: {
		viewTransition: true,
	},
});
