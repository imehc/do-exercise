import { apiInstance } from '~/helper/api';
import SigninForm from './form';
import { AuthApi } from '~/do-exercise-api';
import { handleResponse } from '~/helper/format-response';

export default async function SigninPage() {
  const authApi = await apiInstance(AuthApi);
  const captchaRes = await handleResponse(authApi.getCaptcha());

  return <SigninForm captchaRes={captchaRes} />;
}
