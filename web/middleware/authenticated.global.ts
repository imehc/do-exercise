export default defineNuxtRouteMiddleware(async (to) => {
	const { loggedIn, clear } = useUserSession();

	if (to.path.endsWith('/login')) {
		await clear();
		return;
	}

	const locale = useCookie('i18n_redirected');
	if (!loggedIn.value) {
		if (import.meta.server) {
			clear();
		}
		if (import.meta.client) {
			try {
				const path =
					!locale || locale.value === 'zh' ? '/login' : `${locale.value}/login`;
				return await navigateTo(path);
			} catch (error) {
				return await navigateTo('/login');
			}
		}
	}
});
