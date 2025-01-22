import { apiInstance } from '~/helper/api';
import SigninForm from './form';
import { CaptchaApi } from '~/do-exercise-api';
import { handleResponse } from '~/helper/format-response';

export default async function SigninPage() {
  const captchaApi = apiInstance(CaptchaApi);
  const captchaRes = await handleResponse(captchaApi.getCaptcha());

  return <SigninForm captchaRes={captchaRes} />;
}
