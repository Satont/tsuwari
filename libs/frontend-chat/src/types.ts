import type { Settings as ChatSettings } from '@twir/grpc/generated/api/api/modules_chat_overlay';

export type MessageChunk = {
	type: 'text' | 'emote' | '3rd_party_emote';
	value: string;
}

export type BadgeVersion = {
	id: string,
	image_url_1x: string,
	image_url_2x: string,
	image_url_4x: string,
}

export type ChatBadge = {
	set_id: string,
	versions: Array<BadgeVersion>
}

export type Message = {
	internalId: string,
	id?: string,
	type: 'message' | 'system',
	chunks: MessageChunk[],
	sender?: string,
	senderColor?: string,
	senderDisplayName?: string
	badges?: Record<string, string>,
	isItalic: boolean;
	createdAt: Date;
	announceColor?: string;
	isAnnounce: boolean;
};

export type Settings = {
	channelId: string
	channelName: string
	channelDisplayName: string
	globalBadges: Map<string, ChatBadge>
	channelBadges: Map<string, BadgeVersion>
} & ChatSettings;
