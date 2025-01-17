<script setup lang="ts">
	import { StyleProvider, Themes } from '@varlet/ui';
	interface BasicMenu {
		type: 'link';
		label: string;
		to?: () => void;
		icon?: string;
	}
	interface OptionMenu {
		type: 'option';
		label: string;
		icon: string;
		menus: BasicMenu[];
	}
	type Menu = BasicMenu | OptionMenu;

	const menus = ref<Menu[]>([
		{
			type: 'option',
			label: '主题',
			icon: 'palette-outline',
			menus: [
				{
					label: 'Md2 亮色',
					type: 'link',
					to: () => StyleProvider(null),
				},
				{
					label: 'Md2 暗色',
					type: 'link',
					to: () => StyleProvider(Themes.dark),
				},
				{
					label: 'Md3 亮色',
					type: 'link',
					to: () => StyleProvider(Themes.md3Light),
				},
				{
					label: 'Md3 暗色',
					type: 'link',
					to: () => StyleProvider(Themes.md3Dark),
				},
			],
		},
		{
			type: 'link',
			label: '登出',
			icon: 'bell-outline',
		},
	]);

	const router = useRouter();
</script>

<template>
	<div class="flex h-full items-center justify-between">
		<var-app-bar title-position="center" round class="px-6">
			<template #left>
				<!-- TODO: -->
			</template>
			<template #right>
				<var-space justify="flex-end">
					<template v-for="(item, idx) in menus" :key="idx">
						<var-menu-select
							v-if="item.type === 'option'"
							placement="bottom"
							trigger="hover"
						>
							<var-button round text color="transparent" text-color="#fff">
								<var-icon :name="item.icon" :size="24" />
							</var-button>
							<template #options>
								<var-menu-option
									v-for="(menu, idx2) in item.menus"
									:key="idx2"
									:label="menu.label"
									@click="menu.to"
								/>
							</template>
						</var-menu-select>
						<var-button
							v-if="item.type === 'link'"
							round
							text
							color="transparent"
							text-color="#fff"
							@click="item.to"
						>
							<var-icon :name="item.icon" :size="24" />
						</var-button>
					</template>
				</var-space>
			</template>
		</var-app-bar>
		<div>
			<!-- TODO: logo -->
		</div>
		<!-- <div>
			
		</div> -->
	</div>
</template>
