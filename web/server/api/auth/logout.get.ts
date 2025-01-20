import { AuthApi } from '~/do-exercise-api';

export default defineEventHandler(async (event) => {
	const authApi = await apiInstance(event, AuthApi);
	const res = await wrapperResponseHandler(authApi.logout());
	if (res.ok) {
		await clearUserSession(event);
		return;
	}

	throw createError({
		statusCode: res.code,
		message: res.message,
	});
});
