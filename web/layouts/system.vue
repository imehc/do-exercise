<script setup lang="ts"></script>

<template>
	<div class="fixed inset-0 flex h-screen w-screen flex-col justify-start">
		<div class="h-14 w-full">
			<md-header />
		</div>
		<div class="flex w-full flex-1 justify-start">
			<div class="w-full sm:w-48 lg:w-60"></div>
			<div class="hidden sm:block sm:flex-1 lg:block lg:flex-1">
				<slot />
			</div>
		</div>
	</div>
</template>
