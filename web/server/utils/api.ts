import type { EventHandlerRequest, H3Event } from 'h3';
import type { ConfigurationParameters, Middleware } from '~/do-exercise-api';
import { AuthApi, BASE_PATH, Configuration } from '~/do-exercise-api';

export async function apiInstance<
	T extends new (conf?: Configuration) => InstanceType<T>,
>(event: H3Event<EventHandlerRequest>, Api: T, conf?: ConfigurationParameters) {
	// https://nuxt.com/docs/guide/directory-structure/server#server-utilities

	// https://github.com/atinux/nuxt-auth-utils/issues/91
	const { secure } = await getUserSession(event);

	const _conf = new Configuration({
		basePath: process.env.API_SERVER || BASE_PATH,
		accessToken: secure?.accessToken,
		headers: conf?.headers,
		middleware: [await authMiddleware(event)],
		...conf,
	});

	const instance: InstanceType<T> = new Api(_conf);

	return instance;
}

const authMiddleware = async (
	event: H3Event<EventHandlerRequest>
): Promise<Middleware> => ({
	pre: async (context) => {},
	post: async (context) => {
		if (context.response.ok) return;
		switch (context.response.status) {
			case 401: {
				const { user, secure } = await getUserSession(event);
				if (!secure?.refreshToken) {
					return;
				}
				try {
					const authApi = new AuthApi();
					const auth = await authApi.refreshToken({
						refreshToken: secure.refreshToken,
					});
					setUserSession(event, { user, secure: auth });
					const header = {
						...context.init.headers,
						Authorization: `Bearer ${auth.accessToken}`,
					};
					context.init = {
						...context.init,
						headers: header,
					};
					return await context.fetch(context.url, context.init);
				} catch (error) {
					console.log(error);
				}
			}
		}
	},
	onError: async (context) => {},
});
