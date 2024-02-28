import type { Settings } from '@twir/api/messages/overlays_dudes/overlays_dudes';
import type { DudesJumpRequest, DudesUserPunishedRequest } from '@twir/grpc/websockets/websockets';
import {
	DudesSprite,
	type DudesGrowRequest,
	type DudesUserSettings,
} from '@twir/types/overlays';
import { useWebSocket } from '@vueuse/core';
import { defineStore, storeToRefs } from 'pinia';
import { onMounted, ref, watch } from 'vue';

import { getSprite } from './dudes-config.js';
import { useDudesSettings } from './use-dudes-settings';
import { useDudes } from './use-dudes.js';

import type { TwirWebSocketEvent } from '@/api.js';
import { generateSocketUrlWithParams, normalizeDisplayName } from '@/helpers.js';
import type { ChannelData } from '@/types.js';

declare global {
	interface GlobalEventHandlersEventMap {
		'get-user-settings': CustomEvent<string>;
	}
}

export const useDudesSocket = defineStore('dudes-socket', () => {
	const dudesStore = useDudes();
	const { dudes } = storeToRefs(dudesStore);

	const dudesSettingsStore = useDudesSettings();
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

	watch(data, async (recieviedData) => {
		if (!dudes.value) return;

		const parsedData = JSON.parse(recieviedData) as TwirWebSocketEvent;
		if (parsedData.eventName === 'settings') {
			const data = parsedData.data as Required<Settings & ChannelData>;

			dudesSettingsStore.updateChannelData({
				channelId: data.channelId,
				channelName: data.channelName,
				channelDisplayName: data.channelDisplayName,
			});

			updateSettingFromSocket(data);
		}

		if (parsedData.eventName === 'userSettings') {
			const data = parsedData.data as DudesUserSettings;
			const dudeName = normalizeDisplayName(data.userDisplayName, data.userName);
			dudesSettingsStore.dudesUserSettings.set(dudeName, data);

			const dude = dudesStore.createDude(dudeName, data.userId, data.dudeColor)?.dude;
			if (!dude) return;

			if (data.dudeColor) {
				dude.bodyTint(data.dudeColor);
			}

			if (data.dudeSprite) {
				dude.spriteName = getSprite(data.dudeSprite);
				dude.playAnimation('Run', true);
			}
		}

		if (parsedData.eventName === 'jump') {
			const data = parsedData.data as DudesJumpRequest;
			const dudeName = normalizeDisplayName(data.userDisplayName, data.userName);
			dudesStore.createDude(dudeName, data.userId, data.userColor)?.dude.jump();
		}

		if (parsedData.eventName === 'grow') {
			const data = parsedData.data as DudesGrowRequest;
			const dudeName = normalizeDisplayName(data.userDisplayName, data.userName);
			dudesStore.createDude(dudeName, data.userId, data.userColor)?.dude.grow();
		}

		if (parsedData.eventName === 'punished') {
			const data = parsedData.data as DudesUserPunishedRequest;
			const dudeName = normalizeDisplayName(data.userDisplayName, data.userName);
			dudes.value.removeDude(dudeName);
			dudesSettingsStore.dudesUserSettings.delete(data.userId);
		}
	});

	async function updateSettingFromSocket(data: Required<Settings>) {
		const fontFamily = await dudesSettingsStore.loadFont(
			data.nameBoxSettings.fontFamily,
			data.nameBoxSettings.fontWeight,
			data.nameBoxSettings.fontStyle,
		);

		dudesSettingsStore.updateSettings({
			ignore: data.ignoreSettings,
			overlay: {
				defaultSprite: data.dudeSettings.defaultSprite as keyof typeof DudesSprite,
				maxOnScreen: data.dudeSettings.maxOnScreen,
			},
			dudes: {
				dude: {
					...data.dudeSettings,
					sounds: {
						enabled: data.dudeSettings.soundsEnabled,
						volume: data.dudeSettings.soundsVolume,
					},
				},
				name: {
					...data.nameBoxSettings,
					fontFamily,
				},
				message: {
					...data.messageBoxSettings,
					fontFamily,
				},
				spitter: data.spitterEmoteSettings,
			},
		});
	}

	function destroy(): void {
		if (status.value === 'OPEN') {
			close();
		}
	}

	function connect(apiKey: string, id: string): void {
		overlayId.value = id;
		dudesUrl.value = generateSocketUrlWithParams('/overlays/dudes', {
			apiKey,
			id,
		});
		open();
	}

	onMounted(() => {
		document.addEventListener('get-user-settings', (event) => {
			if (status.value !== 'OPEN') return;
			send(JSON.stringify({ eventName: 'getUserSettings', data: event.detail }));
		});
	});

	return {
		destroy,
		connect,
		updateSettingFromSocket,
	};
});
