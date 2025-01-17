declare module '#auth-utils' {
	interface User {
		id: string;
	}

	interface UserSession {
		user: User;
	}

	interface SecureSessionData {
		accessToken: string;
		refreshToken: string;
		expiresIn: number;
	}
}

export {};
