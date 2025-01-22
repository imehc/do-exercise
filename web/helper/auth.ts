import NextAuth from "next-auth"
import Credentials from "next-auth/providers/credentials"
import { LoginRequest } from "../do-exercise-api"

export const { handlers, signIn, signOut, auth } = NextAuth({
  providers: [
    Credentials({
      credentials: {
        username: {},
        password: {},
        captcha: {},
        captcha_id: {}
      },
      authorize:async(credentials)=>{
        
      }
    })
  ],
})