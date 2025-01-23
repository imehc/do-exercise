'use server';

import {
  BASE_PATH,
  Configuration,
  type ConfigurationParameters,
} from '~/do-exercise-api';
import { auth } from './auth';

export async function apiInstance<
  T extends new (conf?: Configuration) => InstanceType<T>
>(Api: T, conf?: ConfigurationParameters) {
  const session = await auth();

  const _conf = new Configuration({
    basePath: process.env.API_SERVER || BASE_PATH,
    accessToken: session?.accessToken,
    headers: conf?.headers,
    ...conf,
  });

  const instance: InstanceType<T> = new Api(_conf);

  return instance;
}
