import {
  ChannelModerationSetting,
  SettingsType,
} from '@tsuwari/typeorm/entities/ChannelModerationSetting';
import { ChannelModerationWarn } from '@tsuwari/typeorm/entities/ChannelModerationWarn';
import { ChannelPermit } from '@tsuwari/typeorm/entities/ChannelPermit';
import { TwitchPrivateMessage } from '@twurple/chat/lib/commands/TwitchPrivateMessage';
import tlds from 'tlds' assert { type: 'json' };

import { typeorm } from './typeorm.js';

const clipsRegexps = [/.*(clips.twitch.tv\/)(\w+)/, /.*(www.twitch.tv\/\w+\/clip\/)(\w+)/];
const urlRegexps = [
  new RegExp(
    `(www)? ??\\.? ?[a-zA-Z0-9]+([a-zA-Z0-9-]+) ??\\. ?(${tlds.join('|')})(?=\\P{L}|$)`,
    'iu',
  ),
  new RegExp(`[a-zA-Z0-9]+([a-zA-Z0-9-]+)?\\.(${tlds.join('|')})(?=\\P{L}|$)`, 'iu'),
];
const symbolsRegexp = /([^\s\u0500-\u052F\u0400-\u04FF\w]+)/;

const repository = typeorm.getRepository(ChannelModerationSetting);

export class ModerationParser {
  private async getModerationSettings(channelId: string) {
    const settings = await repository.findBy({
      channelId,
    });

    return settings;
  }

  private async returnByWarnedState(
    cacheKey: SettingsType,
    userId: string,
    settings: ChannelModerationSetting,
  ) {
    const repository = typeorm.getRepository(ChannelModerationWarn);

    const query: Partial<ChannelModerationWarn> = {
      channelId: settings.channelId,
      userId,
      reason: cacheKey,
    };
    const existedWarning = await repository.findOneBy(query);

    if (!existedWarning) {
      repository.save(query);
      return {
        message: settings.warningMessage,
        delete: true,
      };
    } else {
      repository.delete({ id: existedWarning.id });

      return {
        time: settings.banTime,
        message: settings.banMessage,
      };
    }
  }

  async parse(message: string, state: TwitchPrivateMessage) {
    if (!state?.channelId) return;
    const settings = await this.getModerationSettings(state.channelId);

    const results = await Promise.all(
      settings.map((item) => {
        const key = item.type;
        if (!item) return;

        if (state.userInfo.isMod || state.userInfo.isBroadcaster) return;
        if (!item || !item.enabled) return;
        if (!item.vips && state.userInfo.isVip) return;
        if (!item.subscribers && state.userInfo.isSubscriber) return;

        return this[`${key}Parser`](message, item, state);
      }),
    );

    return results.find((r) => typeof r !== 'undefined');
  }

  async linksParser(
    message: string,
    settings: ChannelModerationSetting,
    state: TwitchPrivateMessage,
  ) {
    const containLink = urlRegexps.some((r) => r.test(message));
    if (!containLink) return;

    if (!settings.checkClips && clipsRegexps.some((r) => r.test(message))) return;

    const repository = typeorm.getRepository(ChannelPermit);
    const permit = await repository.findOneBy({
      channelId: state.channelId!,
      userId: state.userInfo.userId!,
    });
    if (permit) {
      repository.delete({
        id: permit.id,
      });
      return;
    }

    return this.returnByWarnedState(SettingsType.links, state.userInfo.userId, settings);
  }

  async blacklistsParser(
    message: string,
    settings: ChannelModerationSetting,
    state: TwitchPrivateMessage,
  ) {
    if (!Array.isArray(settings.blackListSentences)) return;
    const blackListed = settings.blackListSentences.some((b) => message.includes(b as string));
    if (!blackListed) return;

    return this.returnByWarnedState(SettingsType.blacklists, state.userInfo.userId, settings);
  }

  async symbolsParser(
    message: string,
    settings: ChannelModerationSetting,
    state: TwitchPrivateMessage,
  ) {
    if (!settings.maxPercentage) return;

    const matched = message.match(symbolsRegexp);
    if (!matched) return;

    let symbolsCount = 0;

    for (const item of matched) {
      symbolsCount = symbolsCount + item.length;
    }

    const check = Math.ceil((symbolsCount * 100) / message.length) >= settings.maxPercentage;
    if (!check) return;

    return this.returnByWarnedState(SettingsType.symbols, state.userInfo.userId, settings);
  }

  async longMessageParser(
    message: string,
    settings: ChannelModerationSetting,
    state: TwitchPrivateMessage,
  ) {
    if (!settings.triggerLength) return;
    if (message.length <= settings.triggerLength) return;

    return this.returnByWarnedState(SettingsType.longMessage, state.userInfo.userId, settings);
  }

  async capsParser(
    message: string,
    settings: ChannelModerationSetting,
    state: TwitchPrivateMessage,
  ) {
    if (!settings.maxPercentage) return;

    let capsCount = 0;

    for (const emote of state.parseEmotes().filter((o) => o.type === 'emote')) {
      if ('name' in emote) {
        message = message.replace(emote['name'], '').trim();
      }
    }

    for (let i = 0; i < message.length; i++) {
      const char = message.charAt(i);
      if (char !== char.toLowerCase()) {
        capsCount += 1;
      }
    }

    const check = Math.ceil((capsCount * 100) / message.length) >= settings.maxPercentage;
    if (!check) return;

    return this.returnByWarnedState(SettingsType.caps, state.userInfo.userId, settings);
  }

  async emotesParser(
    _message: string,
    settings: ChannelModerationSetting,
    state: TwitchPrivateMessage,
  ) {
    if (!settings.triggerLength) return;

    const emotesLength = state.parseEmotes().filter((o) => o.type === 'emote').length;
    if (emotesLength < settings.triggerLength) return;

    return this.returnByWarnedState(SettingsType.emotes, state.userInfo.userId, settings);
  }
}
