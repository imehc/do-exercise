import { Token } from "~/do-exercise-api";

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