<script setup lang="ts">
	import type { MdThemeType } from 'global';
	import { FetchError } from 'ofetch';
	import { useI18nStore } from '~/store/i18n';

	const { locale, t } = useI18n();
	const localeRoute = useLocaleRoute();
	const { setLocal } = useI18nStore();
	const { setTheme } = useMd();
	const { clear, fetch: refreshSession } = useUserSession();

	const themeOptions = [
		{ key: 'md2', label: t('theme.md2') },
		{ key: 'md2-dark', label: t('theme.md2Dark') },
		{ key: 'md3', label: t('theme.md3') },
		{ key: 'md3-dark', label: t('theme.md3Dark') },
	] as const;

	const languageOptions = [
		{ key: 'zh', label: '中文' },
		{ key: 'en', label: 'English' },
	] as const;

	const currentTheme = ref<MdThemeType>();
	const currentLang = ref<typeof locale.value>();

	onNuxtReady(() => {
		currentTheme.value = useMd().theme;
		currentLang.value = locale.value;
	});

	const linkPath = computed(() => {
		const route = localeRoute('login', locale.value);
		return route != null ? route.path : '/';
	});
	const disabled = ref(false);
	const handleLogout = () => {
		disabled.value = true;
		$fetch('/api/auth/logout')
			.then(() => {
				Snackbar.success(t('logoutSuccess'));
				clear();
				refreshSession();
				navigateTo(linkPath.value, { replace: true });
			})
			.catch((error) => {
				let msg = t('logoutFailed');
				if (error instanceof FetchError) {
					if (error.response?._data.message) {
						msg = error.response?._data.message;
					}
				}
				Snackbar.error(msg);
			})
			.finally(() => {
				disabled.value = false;
			});
	};
</script>

<template>
	<div class="flex h-full items-center justify-between">
		<var-app-bar title-position="center" round class="px-1 sm:px-4 lg:px-6">
			<template #left>
				<nuxt-link-locale to="/dashboard">
					{{ $t('headerLeftMenu.home') }}
				</nuxt-link-locale>
			</template>
			<template #right>
				<var-space justify="flex-end">
					<var-menu-select
						v-model="currentTheme"
						placement="bottom"
						trigger="hover"
					>
						<var-button round text color="transparent" text-color="#fff">
							<var-icon name="palette-outline" :size="24" />
						</var-button>
						<template #options>
							<var-menu-option
								v-for="item in themeOptions"
								:key="item.key"
								:label="item.label"
								:value="item.key"
								@click="() => setTheme(item.key)"
							/>
						</template>
					</var-menu-select>
					<var-menu-select
						v-model="currentLang"
						placement="bottom"
						trigger="hover"
					>
						<var-button round text color="transparent" text-color="#fff">
							<var-icon name="translate" :size="24" />
						</var-button>
						<template #options>
							<var-menu-option
								v-for="item in languageOptions"
								:key="item.key"
								:label="item.label"
								:value="item.key"
								@click="() => setLocal(item.key)"
							/>
						</template>
					</var-menu-select>
					<!-- <nuxt-link-locale to="/dashboard/role"> -->
					<var-button
						round
						text
						color="transparent"
						text-color="#fff"
						@click="handleLogout"
					>
						<nuxt-icon name="material-symbols:logout" :size="24" />
					</var-button>
					<!-- </nuxt-link-locale> -->
				</var-space>
			</template>
		</var-app-bar>
	</div>
</template>
