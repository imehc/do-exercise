import type { Pinia } from 'pinia';

declare module '@nuxt/types' {
	interface Context {
		$pinia: Pinia;
	}

	interface NuxtApp {
		$pinia: Pinia;
	}
}
