import path from 'node:path'

import { webUpdateNotice } from '@plugin-web-update-notification/vite'
import vue from '@vitejs/plugin-vue'
import autoprefixer from 'autoprefixer'
import tailwind from 'tailwindcss'
import { defineConfig } from 'vite'
import { watch } from 'vite-plugin-watch'

// https://vitejs.dev/config/
export default defineConfig({
	css: {
		postcss: {
			plugins: [tailwind(), autoprefixer()],
		},
	},
	plugins: [
		watch({
			onInit: true,
			pattern: 'src/**/*.ts',
			command: 'graphql-codegen',
		}),
		vue({
			script: {
				defineModel: true,
			},
		}),
		webUpdateNotice({
			hiddenDefaultNotification: true,
			checkInterval: 1 * 60 * 1000,
		}),
	],
	resolve: {
		alias: {
			'vue': 'vue/dist/vue.esm-bundler.js',
			'@': path.resolve(__dirname, './src'),
		},
	},
	base: '/overlays',
	server: {
		host: true,
		port: 3008,
	},
	clearScreen: false,
})
