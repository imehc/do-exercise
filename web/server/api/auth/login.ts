import { AuthApi } from '~/do-exercise-api';
import { loginSchema } from '~/schemas/auth/login';
import { apiInstance } from '~/utils/api';
import { handleResponse } from '~/utils/format-response';

export default defineEventHandler(async (event) => {
	const body = await readValidatedBody(event, loginSchema.parse);

	const authApi = apiInstance(AuthApi);
	const res = await handleResponse(authApi.login({ loginRequest: body }));
	if (res.ok) {
		// doc: https://nuxt.com/modules/auth-utils#session-management
		await setUserSession(event, {
			user: {
				// TODO: get user info from response
				id: 'dashboard',
			},
			secure: {
				accessToken: res.data.accessToken,
				refreshToken: res.data.refreshToken,
				expiresIn: res.data.expiresIn,
			},
		});
		return;
	}

	throw createError({
		statusCode: res.code,
		message: res.message,
	});
});
