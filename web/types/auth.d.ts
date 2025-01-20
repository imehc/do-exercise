import type { TokenResponse } from '~/do-exercise-api';

declare module '#auth-utils' {
	interface User {
		id: string;
		accessToken: string;
		refreshToken: string;
		expiresIn: number;
	}

	interface UserSession {
		user: User;
	}

	// eslint-disable-next-line @typescript-eslint/no-empty-object-type
	interface SecureSessionData extends TokenResponse {}
}

export {};
