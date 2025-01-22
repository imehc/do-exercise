import NextAuth from "next-auth"
import Credentials from "next-auth/providers/credentials"

export const { handlers, signIn, signOut, auth } = NextAuth({
  providers: [
    Credentials({
      credentials: {
        username: {},
        password: {},
        captcha: {},
        captcha_id: {}
      },
      authorize: async (credentials) => {
        console.log(credentials, 'credentials')

        return {}
      }
    })
  ],
})