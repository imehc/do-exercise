<script setup lang="ts">
	const a = 'ss';
</script>

<template>
	<div>
		<text class="bg-slate-600 text-red-500">123</text>
		<div class="bg-primary text-md text-on-primary">hello</div>
		<var-space :size="[10, 10]">
			<var-button>默认按钮</var-button>
			<var-button type="primary">主要按钮</var-button>
			<var-button type="info">信息按钮</var-button>
			<var-button type="success">成功按钮</var-button>
			<var-button type="warning">警告按钮</var-button>
			<var-button type="danger">危险按钮</var-button>
		</var-space>
	</div>
</template>
