import { StyleProvider, Themes } from '@varlet/ui';
import { useStorage } from '@vueuse/core';
import type { MdThemeType } from 'global';

// 主题切换
export default function useMd() {
	const state = useStorage<MdThemeType>('md-theme', 'md3');
	/** 设置主题 */
	const setTheme = (theme?: MdThemeType) => {
		switch (theme) {
			case 'md2':
				{
					state.value = 'md2';
					StyleProvider(null);
				}
				break;
			case 'md2-dark':
				{
					state.value = 'md2-dark';
					StyleProvider(Themes.dark);
				}
				break;
			case 'md3-dark':
				{
					state.value = 'md3-dark';
					StyleProvider(Themes.md3Dark);
				}
				break;
			case 'md3':
			default:
				{
					state.value = 'md3';
					StyleProvider(Themes.md3Light);
				}
				break;
		}
	};

	const init = () => setTheme(state.value);

	return {
		init,
		theme: state.value,
		setTheme,
	};
}
