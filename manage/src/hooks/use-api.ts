import { Configuration, ConfigurationParameters } from "~/do-exercise-api";

export function useApi<
T extends new (conf?: Configuration) => InstanceType<T>,
>(Api: T, conf?: ConfigurationParameters) {
    const accessToken = ''
    const _conf = new Configuration({
      basePath: '/api',
      accessToken: accessToken,
      headers: conf?.headers,
      ...conf,
    });
  
    const instance: InstanceType<T> = new Api(_conf);
  
    return instance;
}

