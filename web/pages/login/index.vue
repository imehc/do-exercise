<script setup lang="ts">
	import { Snackbar } from '@varlet/ui';
	import type { Form } from '@varlet/ui/types/form';
	import { FetchError } from 'ofetch';
	import type { z } from 'zod';
	import type { CaptchaResponse } from '~/do-exercise-api';
	import { loginSchema } from '~/schemas/auth/login';

	definePageMeta({
		layout: 'no-auth',
	});
	const { t, locale } = useI18n();
	const localeRoute = useLocaleRoute();

	const router = useRouter();
	const { data, status, refresh } =
		await useFetch<CaptchaResponse>('/api/captcha');
	const { fetch: refreshSession } = useUserSession();

	const formData = ref<z.infer<typeof loginSchema>>({
		username: '',
		password: '',
		captcha: '',
		captchaId: '',
	});

	watch(
		() => data.value,
		(val) => {
			if (status.value !== 'success') return;
			formData.value = { ...formData.value, captchaId: val!.captchaId };
		},
		{ immediate: true }
	);

	const form = ref<Form>();
	const linkPath = computed(() => {
		const route = localeRoute('dashboard', locale.value);
		return route != null ? route.path : '/';
	});
	const disabled = ref(false);

	const handleSubmit = async (valid: boolean) => {
		if (!valid) return;
		disabled.value = true;
		$fetch('/api/auth/login', {
			method: 'POST',
			body: formData.value,
		})
			.then(async () => {
				// Refresh the session on client-side and redirect to the home page
				Snackbar.success(t('loginSuccess'));
				await refreshSession();
				await router.replace(linkPath.value);
			})
			.catch((error) => {
				let msg = t('loginFailed');
				if (error instanceof FetchError) {
					if (error.response?._data.message) {
						msg = error.response?._data.message;
					}
				}
				refresh();
				Snackbar.error(msg);
			})
			.finally(() => {
				disabled.value = false;
			});
	};
</script>

<template>
	<div class="w-full">
		<var-card>
			<template #default>
				<var-form
					ref="form"
					:disabled="disabled"
					scroll-to-error="start"
					@submit="handleSubmit"
				>
					<var-space direction="column" :size="[12, 0]">
						<var-input
							v-model.trim="formData.username"
							variant="outlined"
							:placeholder="$t('loginFormLabel.username')"
							:rules="loginSchema.shape.username"
							:maxlength="8"
						>
							<template #prepend-icon>
								<var-icon
									class="prepend-icon mr-2"
									name="account-circle-outline"
								/>
							</template>
						</var-input>
						<var-input
							v-model.trim="formData.password"
							variant="outlined"
							:placeholder="$t('loginFormLabel.password')"
							type="password"
							:rules="loginSchema.shape.password"
							:maxlength="16"
						>
							<template #prepend-icon>
								<var-icon class="prepend-icon mr-2" name="lock-outline" />
							</template>
							<!-- <template #append-icon>
						<var-icon class="ml-2 append-icon" name="view" />
					</template> -->
						</var-input>
						<div class="flex items-start justify-start gap-x-3">
							<var-input
								v-model.trim="formData.captcha"
								variant="outlined"
								:placeholder="$t('loginFormLabel.captcha')"
								:rules="loginSchema.shape.captcha"
								:maxlength="6"
							/>
							<nuxt-img
								:loading="status === 'pending' ? 'lazy' : 'eager'"
								:src="data?.picPath"
								alt=""
								class="aspect-[3/1] h-[calc(var(--field-decorator-outlined-normal-margin-top)+var(--field-decorator-outlined-normal-margin-bottom)+24px)] cursor-pointer rounded-sm"
								:class="
									disabled ? 'pointer-events-none' : 'pointer-events-auto'
								"
								@click="refresh"
							/>
						</div>
						<var-button
							block
							:loading="disabled"
							type="success"
							native-type="submit"
						>
							{{ $t('login') }}
						</var-button>
					</var-space>
				</var-form>
			</template>
		</var-card>
	</div>
</template>
