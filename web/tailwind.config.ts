import { presetVarlet } from '@varlet/preset-tailwindcss';
import type { Config } from 'tailwindcss';

export default <Partial<Config>>{
	content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
	presets: [presetVarlet()],
};
