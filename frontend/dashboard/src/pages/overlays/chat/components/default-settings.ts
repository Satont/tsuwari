import type { Settings } from '@twir/api/messages/overlays_chat/overlays_chat';

export type ChatSettingsWithOptionalId = Omit<Settings, 'id'> & { id?: string }

export const defaultChatSettings: ChatSettingsWithOptionalId = {
	fontFamily: 'inter',
	fontSize: 20,
	fontWeight: 400,
	fontStyle: 'normal',
	hideBots: false,
	hideCommands: false,
	messageHideTimeout: 0,
	messageShowDelay: 0,
	preset: 'clean',
	showBadges: true,
	showAnnounceBadge: true,
	textShadowColor: 'rgba(0,0,0,1)',
	textShadowSize: 0,
	chatBackgroundColor: 'rgba(0, 0, 0, 0)',
	direction: 'top',
	paddingContainer: 0,
};
