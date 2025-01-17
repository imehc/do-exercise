<template>
	<div
		class="fixed left-0 top-0 grid h-screen w-screen p-3 sm:grid-cols-4 sm:p-0 lg:grid-cols-7"
	>
		<div
			class="grid grid-rows-4 sm:col-span-2 sm:col-start-2 lg:col-span-2 lg:col-start-5"
		>
			<div class="row-span-2 row-start-2 flex items-center justify-center">
				<div class="flex min-w-72 max-w-96 items-center justify-center">
					<slot />
				</div>
			</div>
		</div>
	</div>
</template>
