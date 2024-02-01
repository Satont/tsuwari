import type { Settings } from '@twir/api/messages/overlays_dudes/overlays_dudes';
import type { DudesJumpRequest } from '@twir/grpc/websockets/websockets';
import { useWebSocket } from '@vueuse/core';
import { defineStore, storeToRefs } from 'pinia';
import { ref, watch } from 'vue';

import { useDudesSettings } from './use-dudes-settings';
import { useDudes } from './use-dudes.js';

import type { TwirWebSocketEvent } from '@/api.js';
import { generateSocketUrlWithParams } from '@/helpers.js';
import type { ChannelData } from '@/types.js';

const nameBoxDefaults: Partial<Settings['nameBoxSettings']> = {
	strokeThickness: 0,
	dropShadow: false,
	dropShadowAlpha: 0,
	dropShadowBlur: 0,
	dropShadowDistance: 0,
	dropShadowAngle: 0,
};

const messageBoxDefaults: Partial<Settings['messageBoxSettings']> = {
	enabled: false,
	padding: 0,
	borderRadius: 0,
};

export const useDudesSocket = defineStore('dudes-socket', () => {
	const dudesStore = useDudes();
	const { dudes } = storeToRefs(dudesStore);

	const { updateSettings, updateChannelData, loadFont } = useDudesSettings();
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

	watch(data, async (d) => {
		const parsedData = JSON.parse(d) as TwirWebSocketEvent;
		if (parsedData.eventName === 'settings') {
			const data = parsedData.data as Required<Settings & ChannelData>;

			updateChannelData({
				channelId: data.channelId,
				channelName: data.channelName,
				channelDisplayName: data.channelDisplayName,
			});

			const fontFamily = await loadFont(
				data.nameBoxSettings.fontFamily,
				data.nameBoxSettings.fontWeight,
				data.nameBoxSettings.fontStyle,
			);

			updateSettings({
				dude: {
					...data.dudeSettings,
					sounds: {
						enabled: data.dudeSettings.soundsEnabled,
						volume: data.dudeSettings.soundsVolume,
					},
				},
				nameBox: {
					...nameBoxDefaults,
					...data.nameBoxSettings,
					fontFamily,
				},
				messageBox: {
					...messageBoxDefaults,
					...data.messageBoxSettings,
					fontFamily,
				},
			});
		}

		if (parsedData.eventName === 'jump') {
			const userData = parsedData.data as DudesJumpRequest;
			const dude = dudes.value?.getDude(userData.userDisplayName);
			if (dude) {
				dudesStore.jumpDude(userData);
			} else {
				dudesStore.createNewDude(userData.userDisplayName, userData.userColor);
			}
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
