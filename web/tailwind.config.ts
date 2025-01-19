import { presetVarlet } from '@varlet/preset-tailwindcss';
import type { Config } from 'tailwindcss';

export default <Partial<Config>>{
	content: [
		'./pages/**/*.{vue,js,ts,jsx,tsx}',
		'./layouts/**/*.{vue,js,ts,jsx,tsx}',
		'./components/**/*.{vue,js,ts,jsx,tsx}',
	],
	presets: [presetVarlet()],
};
