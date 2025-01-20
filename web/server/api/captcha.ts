import { CaptchaApi } from '~/do-exercise-api';

export default defineEventHandler(async (event) => {
	const captchaApi = await apiInstance(event, CaptchaApi);
	const res = await wrapperResponseHandler(captchaApi.getCaptcha());
	if (res.ok) {
		return res.data;
	}
	throw createError({
		statusCode: res.code,
		message: res.message,
	});
});
