import {
	BASE_PATH,
	Configuration,
	type ConfigurationParameters,
} from '~/do-exercise-api';

export function apiInstance<
	T extends new (conf?: Configuration) => InstanceType<T>,
>(Api: T, conf?: ConfigurationParameters) {
	// https://nuxt.com/docs/guide/directory-structure/server#server-utilities
	// https://github.com/atinux/nuxt-auth-utils/issues/91
	// const session = getUserSession();
	const accessToken = ''; // session.value.secure?.accessToken;

	const _conf = new Configuration({
		basePath: process.env.API_SERVER || BASE_PATH,
		accessToken,
		headers: conf?.headers,
		...conf,
	});

	const instance: InstanceType<T> = new Api(_conf);

	return instance;
}
