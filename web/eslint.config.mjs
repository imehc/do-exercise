import prettierPlugin from 'eslint-plugin-prettier';
// @ts-check
import withNuxt from './.nuxt/eslint.config.mjs';

export default withNuxt(
	// Your custom configs here
	{
		plugins: {
			prettier: prettierPlugin,
		},
		rules: {
			'prettier/prettier': [
				'error',
				{
					endOfLine: 'auto',
				},
			],
		},
	},
	{
		files: ['**/*.ts', '**/*.vue', '**/*.js'],
		ignores: ['**/node_modules/**', '**/dist/**', '**/do-exercise-api/**'],
		rules: {
			'@typescript-eslint/no-explicit-any': 'warn',
			'@typescript-eslint/no-unused-vars': 'off',

			'vue/multi-word-component-names': 'off',
			'vue/html-self-closing': 'off',

			'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
			'no-debugger': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
		},
	}
);
