import { TokenResponse } from "~/do-exercise-api";

type Token = TokenResponse

declare module "next-auth/jwt" {
  interface JWT extends Token {
    error?: "RefreshTokenError"
  }
}

declare module "next-auth" {
  interface Session extends Token {
    error?: "RefreshTokenError"
  }

  interface User extends Token { }
}