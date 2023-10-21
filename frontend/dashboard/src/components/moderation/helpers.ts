import {
	IconAbc,
	IconAsteriskSimple,
	IconLanguageOff,
	IconLinkOff,
	IconListLetters,
	IconMessageOff,
	IconMoodOff,
	type SVGProps,
} from '@tabler/icons-vue';
import type { ItemCreateMessage } from '@twir/grpc/generated/api/api/moderation';
import type { FunctionalComponent } from 'vue';

type Item = ItemCreateMessage & {
	icon: FunctionalComponent<SVGProps>
}

export const availableSettings: Item[] = [
	{
		icon: IconLinkOff,
		deniedChatLanguages: [],
		banMessage: 'No links allowed',
		banTime: 600,
		checkClips: false,
		denyList: [],
		enabled: true,
		excludedRoles: [],
		maxPercentage: 0,
		maxWarnings: 3,
		triggerLength: 0,
		type: 'links',
		warningMessage: 'No links allowed [warning]',
	},
	{
		icon: IconLanguageOff,
		deniedChatLanguages: [],
		banMessage: 'Language not allowed',
		banTime: 600,
		checkClips: false,
		denyList: [],
		enabled: true,
		excludedRoles: [],
		maxPercentage: 0,
		maxWarnings: 3,
		triggerLength: 0,
		type: 'language',
		warningMessage: 'Language not allowed [warning]',
	},
	{
		icon: IconListLetters,
		deniedChatLanguages: [],
		banMessage: 'Bad word',
		banTime: 600,
		checkClips: false,
		denyList: [],
		enabled: true,
		excludedRoles: [],
		maxPercentage: 0,
		maxWarnings: 3,
		triggerLength: 0,
		type: 'deny_list',
		warningMessage: 'Bad word [warning]',
	},
	{
		icon: IconMessageOff,
		deniedChatLanguages: [],
		banMessage: 'Too long message',
		banTime: 600,
		checkClips: false,
		denyList: [],
		enabled: true,
		excludedRoles: [],
		maxPercentage: 0,
		maxWarnings: 3,
		triggerLength: 0,
		type: 'long_message',
		warningMessage: 'Too long message [warning]',
	},
	{
		icon: IconAbc,
		deniedChatLanguages: [],
		banMessage: 'Too much caps',
		banTime: 600,
		checkClips: false,
		denyList: [],
		enabled: true,
		excludedRoles: [],
		maxPercentage: 0,
		maxWarnings: 3,
		triggerLength: 0,
		type: 'caps',
		warningMessage: 'Too much caps [warning]',
	},
	{
		icon: IconMoodOff,
		deniedChatLanguages: [],
		banMessage: 'Too much emotes',
		banTime: 600,
		checkClips: false,
		denyList: [],
		enabled: true,
		excludedRoles: [],
		maxPercentage: 0,
		maxWarnings: 3,
		triggerLength: 0,
		type: 'emotes',
		warningMessage: 'Too much emotes [warning]',
	},
	{
		icon: IconAsteriskSimple,
		deniedChatLanguages: [],
		banMessage: 'Too much symbols',
		banTime: 600,
		checkClips: false,
		denyList: [],
		enabled: true,
		excludedRoles: [],
		maxPercentage: 0,
		maxWarnings: 3,
		triggerLength: 0,
		type: 'symbols',
		warningMessage: 'Too much symbols [warning]',
	},
];

export const availableSettingsTypes = availableSettings.map(n => n.type);
