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
	</div>
</template>
