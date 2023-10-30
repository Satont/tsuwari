/* eslint-disable @typescript-eslint/ban-ts-comment */
import { Settings, Message } from '@twir/frontend-chat';
import { Client } from 'tmi.js';
import { Ref, computed, onUnmounted, ref, unref, watch } from 'vue';

import { useThirdPartyEmotes } from '../components/chat_tmi_emotes.js';
import { makeMessageChunks } from '../components/chat_tmi_helpers.js';

type MakeOptional<Type, Key extends keyof Type> = Omit<Type, Key> &
  Partial<Pick<Type, Key>>;

type AddMessageOpts = Omit<
	MakeOptional<Message, 'isItalic' | 'isAnnounce'>,
	'createdAt' | 'internalId'
> & { messageHideTimeout?: number, messageShowDelay?: number; }

const knownBots = new Set([
	'moobot',
	'fossabot',
	'wizebot',
	'twirapp',
	'nightbot',
	'streamlabs',
	'streamelements',
]);

export const useTmiChat = (settings: Ref<Settings>) => {
	let client: Client | null = null;
	const messages = ref<Message[]>([]);
	const maxMessages = ref(100);

	const channelName = computed(() => settings.value.channelName);
	const channelId = computed(() => settings.value.channelId);

	useThirdPartyEmotes(channelName, channelId);

	onUnmounted(async () => {
		destroy();
	});

	function addMessage(opts: AddMessageOpts) {
		if (opts.sender && settings.value.hideBots && knownBots.has(opts.sender)) {
			return;
		}

		if (settings.value.hideCommands && opts.chunks.at(0)?.value.startsWith('!')) {
			return;
		}

		const internalId = crypto.randomUUID();

		const showDelay = opts.messageShowDelay ?? settings.value.messageShowDelay;

		if (messages.value.length >= maxMessages.value) {
			messages.value = messages.value.slice(1);
		}


		// TODO: store only 100 messages in ref for perfomance
		setTimeout(() => {
			messages.value.push({
				...opts,
				isItalic: opts.isItalic ?? false,
				createdAt: new Date(),
				internalId,
				isAnnounce: opts.isAnnounce ?? false,
			});
		}, showDelay * 1000);

		const hideTimeout = opts.messageHideTimeout ?? settings.value.messageHideTimeout;
		if (hideTimeout) {
			setTimeout(() => {
				removeMessageByInternalId(internalId);
			}, hideTimeout * 1000);
		}
	}

	function removeMessageByInternalId(id: string) {
		messages.value = messages.value.filter((m) => m.internalId !== id);
	}

	function removeMessageById(id: string) {
		messages.value = messages.value.filter((m) => m.id !== id);
	}

	function removeMessageBySenderName(name: string) {
		messages.value = messages.value.filter((m) => m.sender !== name);
	}

	async function destroy() {
		if (client) {
			await client.disconnect();
			client.removeAllListeners();
			client = null;
		}
	}

	async function create(channel: string) {
		await destroy();

		client = new Client({
			connection: {
				secure: true,
				reconnect: true,
			},
			channels: [],
		});

		client.on('message', (_, tags, message) => {
			addMessage({
				id: tags.id,
				type: 'message',
				chunks: makeMessageChunks(message, tags.emotes),
				sender: tags.username!,
				senderColor: tags.color,
				senderDisplayName: tags['display-name'],
				badges: tags.badges as Record<string, string> | undefined,
				isItalic: tags['message-type'] === 'action',
			});
		});

		// @ts-ignore
		client.on('usernotice', (msgId, channel, tags, msg) => {
			if(msgId === 'announcement') {
				addMessage({
					id: msgId,
					type: 'message',
					// @ts-ignore
					chunks: makeMessageChunks(msg, tags.emotes),
					// @ts-ignore
					sender: tags.login,
					// @ts-ignore
					senderColor: tags.color,
					senderDisplayName: tags['display-name'],
					// @ts-ignore
					badges: tags.badges as Record<string, string> | undefined,
					isItalic: tags['message-type'] === 'action',
					isAnnounce: true,
				});
			}
		});

		client.on('messagedeleted', (_channel, _username, _msgText, userState) => {
			const msgId = userState['target-msg-id'];
			if (msgId) {
				removeMessageById(msgId);
			}
		});

		client.on('timeout', (_channel, username) => {
			removeMessageBySenderName(username);
		});

		client.on('ban', (_channel, username) => {
			removeMessageBySenderName(username);
		});

		client.on('connecting', () => {
			addMessage({
				type: 'system',
				chunks: [{
					type: 'text',
					value: 'Connecting to servers...',
				}],
				messageHideTimeout: 5,
			});
		});

		client.on('connected', async () => {
			addMessage({
				type: 'system',
				chunks: [{ type: 'text', value: 'Connected' }],
				messageHideTimeout: 6,
			});

			await client!.join(channel);
			addMessage({
				type: 'system',
				chunks: [{ type: 'text', value: `Joined channel ${channel}` }],
				messageHideTimeout: 7,
			});
		});

		await client.connect();
	}

	watch(() => settings.value.channelName, (v) => {
		const name = unref(v);
		if (!name) return;

		create(name);
	});

	return {
		messages,
	};
};
