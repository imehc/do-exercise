import { CaptchaApi } from '~/do-exercise-api';

export default defineEventHandler(async (event) => {
	const captchaApi = apiInstance(CaptchaApi);
	const res = await handleResponse(captchaApi.getCaptcha());
	if (res.ok) {
		return res.data;
	}
	throw createError({
		statusCode: res.code,
		message: res.message,
	});
});
