import type {
	Settings,
} from '@twir/api/messages/overlays_kappagen/overlays_kappagen';
import { ref } from 'vue';

import { animations } from './kappagen_animations';

const settings = ref<Settings>({
	emotes: {
		time: 5,
		max: 0,
		queue: 0,
		bttvEnabled: true,
		emojiStyle: 0,
		ffzEnabled: true,
		sevenTvEnabled: true,
	},
	animations: animations,
	enableRave: false,
	animation: {
		fadeIn: true,
		fadeOut: true,
		zoomIn: true,
		zoomOut: true,
	},
	cube: {
		speed: 6,
	},
	size: {
		// from 7 to 20
		ratioNormal: 0.05,
		// from 14 to 40
		ratioSmall: 0.02,
		min: 1,
		max: 256,
	},
	events: [],
	enableSpawn: true,
	excludedEmotes: [],
});

export const useKappagenFormSettings = () => {
	return {
		settings,
	};
};
