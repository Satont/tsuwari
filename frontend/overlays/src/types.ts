import type { Settings as BrbOverlaySettings } from '@twir/api/messages/overlays_be_right_back/overlays_be_right_back';
import type { Settings as KappagenOverlaySettings } from '@twir/api/messages/overlays_kappagen/overlays_kappagen';
import type { Emote, KappagenAnimations } from '@twirapp/kappagen';

// emotes start
export type SevenTvEmote = {
	id: string;
	name: string;
	flags: number;
	data: {
		host: {
			url: string;
			files: Array<{ name: string; format: string; height: number; width: number }>;
		};
	};
};

export type SevenTvChannelResponse = {
	user: {
		id: string;
	};
	emote_set: {
		id: string;
		emotes: Array<SevenTvEmote>;
	};
};

export type SevenTvGlobalResponse = {
	emotes: Array<SevenTvEmote>;
};

export type BttvEmote = {
	code: string;
	imageType: string;
	id: string;
	height?: number;
	width?: number;
	modifier?: boolean;
};

export type BttvChannelResponse = {
	channelEmotes: Array<BttvEmote>;
	sharedEmotes: Array<BttvEmote>;
};

export type BttvGlobalResponse = Array<BttvEmote>;

export type FfzEmote = {
	name: string;
	urls: Record<string, string>;
	height: number;
	width: number;
	modifier: boolean;
	modifier_flags?: number;
};

export type FfzChannelResponse = {
	sets: {
		[x: string]: {
			emoticons: FfzEmote[];
		};
	};
};

export type FfzGlobalResponse = {
	sets: {
		[x: string]: {
			emoticons: FfzEmote[];
		};
	};
};
// emotes end

// brb start
export type BrbSetSettingsFn = (settings: BrbOverlaySettings) => void;
export type BrbOnStartFn = (minutes: number, text: string) => void;
export type BrbOnStopFn = () => void;
// brb end

// kappagen start
export type KappagenSettings = KappagenOverlaySettings & { channelName: string; channelId: string };
export type KappagenSpawnAnimatedEmotesFn = (
	emotes: Emote[],
	animation: KappagenAnimations
) => void;
export type KappagenSpawnEmotesFn = (emotes: Emote[]) => void;
export type KappagenSetSettingsFn = (settings: KappagenSettings) => void;
// kappagen end

// dudes start
export type ChannelData = {
	channelDisplayName: string
	channelId: string
	channelName: string
}

export type UserData = {
	channelId: string
	userDisplayName: string
	userId: string
	userName: string
}
// dudes end
