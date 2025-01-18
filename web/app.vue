<script setup lang="ts">
	import { Locale } from '@varlet/ui';
	import { useI18nStore } from './store/i18n';

	const store = useI18nStore();
	const { messages, locale } = storeToRefs(store);

	onNuxtReady(() => {
		useMd().init();
		Locale.add('en-US', Locale.enUS);
	});
</script>

<template>
	<var-locale-provider :locale="locale" :messages="messages">
		<nuxt-layout>
			<nuxt-page />
		</nuxt-layout>
		<nuxt-route-announcer />
	</var-locale-provider>
</template>

<style>
	.page-enter-active,
	.page-leave-active {
		transition: all 0.4s;
	}

	.page-enter-from,
	.page-leave-to {
		opacity: 0;
		filter: blur(1rem);
	}

	.layout-enter-active,
	.layout-leave-active {
		transition: all 0.4s;
	}

	.layout-enter-from,
	.layout-leave-to {
		filter: grayscale(1);
	}
	body {
		transition:
			background-color 0.25s,
			color 0.25s;
		color: var(--color-text);
		background-color: var(--color-body);
		color-scheme: var(--color-scheme);
	}
</style>
