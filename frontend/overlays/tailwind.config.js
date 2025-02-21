/** @type {import('tailwindcss').Config} */
export default {
	content: [
		'./index.html',
		'./src/**/*.{js,ts,jsx,tsx,vue}',
		'../../libs/*/src/**/*.{vue,js,ts}',
	],
	theme: {
		extend: {},
	},
	plugins: [],
}
