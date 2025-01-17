<script setup lang="ts">
	onNuxtReady(() => {
		useMdDark();
	});
</script>

<template>
	<div>
		<nuxt-layout>
			<nuxt-page />
		</nuxt-layout>
		<nuxt-route-announcer />
	</div>
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
