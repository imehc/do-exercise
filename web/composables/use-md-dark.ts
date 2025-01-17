import { Themes } from '@varlet/ui';

// 主题切换
export default function useMdDark() {
	const isDark = useDark({
		onChanged: (dark) => StyleProvider(dark ? Themes.md3Dark : Themes.md3Light),
	});
	const [value, toggle] = useToggle(isDark.value);

	toggle();

	return {
		isDark: value,
		toggleDark: toggle,
	};
}
