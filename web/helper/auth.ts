import NextAuth from 'next-auth';
import Credentials from 'next-auth/providers/credentials';
import { AuthApi } from '~/do-exercise-api';
import { apiInstance } from './api';
import { handleResponse } from './format-response';
import { JWT } from 'next-auth/jwt';

export const { handlers, signIn, signOut, auth } = NextAuth((req) => {
  return {
    providers: [
      Credentials({
        credentials: {
          username: { label: 'Username', type: 'string' },
          password: { label: 'Password', type: 'password' },
          captcha: { label: 'Captcha', type: 'string' },
          captchaId: { label: 'CaptchaId', type: 'string' },
        },
        authorize: async ({ username, password, captcha, captchaId }) => {
          if (!username || !password || !captcha || !captchaId) {
            return null;
          }
          const authApi = await apiInstance(AuthApi);
          const res = await handleResponse(
            authApi.login({
              login: {
                username: username as string,
                password: password as string,
                captcha: captcha as string,
                captchaId: captchaId as string,
              },
            })
          );
          if (res.ok) {
            return res.data;
          }
          throw new Error(res.message);
        },
      }),
    ],
    callbacks: {
      jwt: async ({ token, user, account, trigger }) => {
        if (account && user) {
          return {
            ...token,
            accessToken: user.accessToken,
            refreshToken: user.refreshToken,
            expiresIn: getExp(user.expiresIn),
            tokenType: user.tokenType,
            expires: 0,
          };
        }

        if (!token.expiresIn) {
          return null;
        }
        if (Date.now() < token.expiresIn) {
          return token;
        }
        return await refreshAccessToken(token);
      },
      session: async ({ session, token }) => {
        session.accessToken = token.accessToken;
        session.refreshToken = token.refreshToken;
        session.expiresIn = token.expiresIn;
        session.tokenType = token.tokenType;
        session.error = token.error;
        // TODO: 获取当前一些用户信息

        return session;
      },
    },
    events: {
      signOut: async () => {
        // TODO: 登出
        const authApi = await apiInstance(AuthApi);
        await handleResponse(authApi.logout());
      },
    },
  };
});

const getExp = (second: number) => {
  return Date.now() + second * 1000;
};

// ISSUE: 多次执行且使用的是仍然是已经使用过的refresh_token
async function refreshAccessToken(token: JWT): Promise<JWT> {
  // 这里使用客户端
  const authApi = new AuthApi();
  const res = await handleResponse(
    authApi.refreshToken({ refreshToken: token.refreshToken })
  );
  console.log('Refreshing access token');
  if (!res.ok) {
    return {
      ...token,
      error: 'RefreshTokenError',
    };
  }
  return {
    ...token,
    accessToken: res.data.accessToken,
    refreshToken: res.data.refreshToken,
    expiresIn: getExp(res.data.expiresIn),
    tokenType: res.data.tokenType,
  };
}
