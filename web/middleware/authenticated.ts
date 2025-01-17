export default defineNuxtRouteMiddleware((to) => {
	const { loggedIn, clear } = useUserSession();
	if (to.path === '/login') return;
	// redirect the user to the login screen if they're not authenticated
	if (!loggedIn.value) {
		clear();
		return navigateTo('/login', { replace: true });
	}
});
