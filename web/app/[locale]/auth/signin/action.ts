"use server";

import type { LoginRequest } from "~/do-exercise-api";

export const signinAction = async (signin: LoginRequest) => {
  console.log(signin, 'signn')
}

