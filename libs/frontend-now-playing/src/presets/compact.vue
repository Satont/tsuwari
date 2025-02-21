<script setup lang="ts">
import type { Settings, Track } from '../types'

const props = defineProps<{
	track?: Track | null
	settings: Settings
}>()
</script>

<template>
	<div class="flex gap-2 w-full">
		<div
			class="w-[30%] relative flex aspect-square items-center justify-center overflow-hidden transition duration-500"
			:style="{
				backgroundImage: `url(${track?.imageUrl})`,
				backgroundSize: 'cover',
				backgroundPosition: 'center center',
				aspectRatio: '1 / 1',
				borderRadius: '0.5rem',
			}"
		>
		</div>

		<div class="flex flex-col z-40 w-full min-w-0 flex-col gap-2">
			<div class="flex flex-col relative h-full w-full overflow-x-clip rounded-lg px-3 py-1 transition-[background-color,color] duration-300 text-white" style="background-color: rgb(94, 55, 76); color: rgb(137, 211, 240);">
				<span class="text-lg font-bold drop-shadow-md">{{ props.track?.title }}</span>
				<span class="relative z-50 mb-1 overflow-clip overflow-ellipsis whitespace-nowrap text-sm">{{ props.track?.artist }}</span>
			</div>

			<div class="flex h-full w-full items-center justify-between gap-1 rounded-lg px-3 text-lg transition-all duration-300 ">
				duration
			</div>
		</div>
	</div>
</template>
