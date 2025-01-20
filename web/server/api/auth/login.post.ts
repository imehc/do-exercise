import { addSeconds } from 'date-fns';
import { AuthApi } from '~/do-exercise-api';
import { loginSchema } from '~/schemas/auth/login';

export default defineEventHandler(async (event) => {
	const body = await readValidatedBody(event, loginSchema.parse);

	const authApi = await apiInstance(event, AuthApi);
	const res = await wrapperResponseHandler(
		authApi.login({ loginRequest: body })
	);
	if (res.ok) {
		const expiresIn = addSeconds(new Date(), res.data.expiresIn).getTime();
		// doc: https://nuxt.com/modules/auth-utils#session-management
		await setUserSession(event, {
			user: {
				// TODO: get user info from response
				id: 'dashboard',
				accessToken: res.data.accessToken,
				refreshToken: res.data.refreshToken,
				expiresIn: expiresIn,
			},
			secure: {
				accessToken: res.data.accessToken,
				refreshToken: res.data.refreshToken,
				expiresIn: expiresIn,
				tokenType: res.data.tokenType,
			},
		});
		return;
	}

	throw createError({
		statusCode: res.code,
		message: res.message,
	});
});
