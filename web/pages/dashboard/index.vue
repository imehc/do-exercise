<script setup lang="ts">
	definePageMeta({
		title: 'dashboard',
		layout: 'system',
	});
	const { user } = useUserSession();
</script>

<template>
	<div>
		<div>{{ JSON.stringify(user?.id) }}</div>
		<div>{{ JSON.stringify(user?.accessToken) }}</div>
		<div>{{ JSON.stringify(user?.refreshToken) }}</div>
		<div>{{ JSON.stringify(user?.expiresIn) }}</div>
	</div>
</template>
