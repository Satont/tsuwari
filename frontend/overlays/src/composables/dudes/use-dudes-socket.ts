import { useWebSocket } from '@vueuse/core';
import type { Settings } from '@twir/api/messages/overlays_dudes/overlays_dudes';
import { defineStore } from 'pinia';
import { ref, watch } from 'vue';

import { useDudesSettings } from './use-dudes-settings';

import type { TwirWebSocketEvent } from '@/api.js';
import { generateSocketUrlWithParams } from '@/helpers.js';

export const useDudesSocket = defineStore('dudes-socket', () => {
	const { updateSettings } = useDudesSettings();
	const overlayId = ref('');
	const dudesUrl = ref('');
	const { data, send, open, close, status } = useWebSocket(
		dudesUrl,
		{
			immediate: false,
			autoReconnect: {
				delay: 500,
			},
			onConnected() {
				send(JSON.stringify({ eventName: 'getSettings' }));
			},
		},
	);

	watch(data, (d) => {
		const parsedData = JSON.parse(d) as TwirWebSocketEvent<Required<Settings>>;
		if (parsedData.eventName === 'settings') {
			const data = parsedData.data;
			updateSettings({
				dude: {
					...data.dudeSettings,
					sounds: {
						enabled: data.dudeSettings.soundsEnabled,
						volume: data.dudeSettings.soundsVolume,
					},
				},
				nameBox: data.nameBoxSettings,
				messageBox: data.messageBoxSettings,
			});
		}
	});

	function destroy() {
		if (status.value === 'OPEN') {
			close();
		}
	}

	function connect(apiKey: string, id: string) {
		overlayId.value = id;
		dudesUrl.value = generateSocketUrlWithParams('/overlays/dudes', {
			apiKey,
			id,
		});
		open();
	}

	return {
		destroy,
		connect,
	};
});
