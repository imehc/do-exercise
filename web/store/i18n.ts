export const useI18nStore = defineStore(
	'i18n',
	() => {
		const { setLocale } = useI18n();
		const messages = {
			zh: Locale.zhCN,
			en: Locale.enUS,
		};
		type Lang = keyof typeof messages;
		const locale = ref<Lang>('zh');

		const switchLocale = (lang: Lang) => {
			locale.value = lang;
			setLocale(lang);
		};

		return {
			messages,
			locale,
			setLocal: switchLocale,
		};
	},
	{ persist: true }
);
