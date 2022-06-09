import { ModerationSettings, SettingsType } from '@tsuwari/prisma';
import { TwitchPrivateMessage } from '@twurple/chat/lib/commands/TwitchPrivateMessage';
import tlds from 'tlds' assert { type: 'json' };

import { prisma } from './prisma.js';
import { redis } from './redis.js';

const clipsRegexps = [/.*(clips.twitch.tv\/)(\w+)/g, /.*(www.twitch.tv\/\w+\/clip\/)(\w+)/g];
const urlRegexps = [
  new RegExp(`(www)? ??\\.? ?[a-zA-Z0-9]+([a-zA-Z0-9-]+) ??\\. ?(${tlds.join('|')})(?=\\P{L}|$)`, 'iu'),
  new RegExp(`[a-zA-Z0-9]+([a-zA-Z0-9-]+)?\\.(${tlds.join('|')})(?=\\P{L}|$)`, 'iu'),
];
const symbolsRegexp = /([^\s\u0500-\u052F\u0400-\u04FF\w]+)/g;

// @TODO: update redis cache on changes from panel
export class ModerationParser {
  async getModerationSettings(channelId: string) {
    const result = {} as Record<SettingsType, ModerationSettings>;
    const settingsKeys = Object.values(SettingsType);

    for (const key of settingsKeys) {
      const redisKey = `settings:moderation:${channelId}:${key}`;
      const cachedSettings = await redis.get(redisKey);

      if (cachedSettings) {
        result[key] = JSON.parse(cachedSettings) as ModerationSettings;
      } else {
        const entity = await prisma.moderationSettings.findFirst({ where: { channelId: channelId, type: key } });
        if (entity) {
          redis.set(redisKey, JSON.stringify(entity), 'EX', 5 * 60 * 60);
          result[key] = entity;
        }
      }
    }

    return result;
  }

  async returnByWarnedState(cacheKey: SettingsType, userId: string, settings: ModerationSettings) {
    const redisKey = `moderation:warnings:${cacheKey}:${userId}`;
    const isWarned = await redis.get(redisKey);

    if (isWarned === null) {
      redis.set(redisKey, '', 'EX', 60 * 60);
      return {
        message: settings.warningMessage,
        delete: true,
      };
    } else {
      redis.del(redisKey);
      return {
        time: settings.banTime,
        message: settings.banMessage,
      };
    }
  }

  async parse(message: string, state: TwitchPrivateMessage) {
    if (!state?.channelId) return;
    const settings = await this.getModerationSettings(state.channelId);

    const results = await Promise.all(Object.keys(settings).map((k) => {
      const key = k as SettingsType;
      const parserSettings = settings[key];

      if (state.userInfo.isMod || state.userInfo.isBroadcaster) return;
      if (!parserSettings || !parserSettings.enabled) return;
      if (!parserSettings.vips && state.userInfo.isVip) return;
      if (!parserSettings.subscribers && state.userInfo.isSubscriber) return;

      return this[`${key}Parser`](message, parserSettings, state);
    }));

    console.log(results);

    return results.find(r => typeof r !== 'undefined');
  }

  async linksParser(message: string, settings: ModerationSettings, state: TwitchPrivateMessage) {
    const containLink = urlRegexps.some(r => r.test(message));
    if (!containLink) return;

    if (!settings.checkClips && clipsRegexps.some(r => r.test(message))) return;

    const permit = await prisma.permit.findFirst({ where: { channelId: state.channelId!, userId: state.userInfo.userId } });
    if (permit) {
      await prisma.permit.delete({ where: { id: permit.id } });
      return;
    }

    return this.returnByWarnedState('links', state.userInfo.userId, settings);
  }

  async blacklistsParser(message: string, settings: ModerationSettings, state: TwitchPrivateMessage) {
    if (!Array.isArray(settings.blackListSentences)) return;
    const blackListed = settings.blackListSentences.some(b => message.includes(b as string));
    if (!blackListed) return;

    return this.returnByWarnedState('blacklists', state.userInfo.userId, settings);
  }

  async symbolsParser(message: string, settings: ModerationSettings, state: TwitchPrivateMessage) {
    if (!settings.maxPercentage) return;

    const matched = message.match(symbolsRegexp);
    if (!matched) return;

    let symbolsCount = 0;

    for (const item of matched) {
      symbolsCount = symbolsCount + item.length;
    }

    const check = Math.ceil(symbolsCount * 100 / message.length) >= settings.maxPercentage;
    if (!check) return;

    return this.returnByWarnedState('symbols', state.userInfo.userId, settings);
  }

  async longMessageParser(message: string, settings: ModerationSettings, state: TwitchPrivateMessage) {
    if (!settings.triggerLength) return;
    if (message.length <= settings.triggerLength) return;

    return this.returnByWarnedState('longMessage', state.userInfo.userId, settings);
  }

  async capsParser(message: string, settings: ModerationSettings, state: TwitchPrivateMessage) {
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

    const check = Math.ceil(capsCount * 100 / message.length) >= settings.maxPercentage;
    if (!check) return;

    return this.returnByWarnedState('caps', state.userInfo.userId, settings);
  }

  async emotesParser(_message: string, settings: ModerationSettings, state: TwitchPrivateMessage) {
    if (!settings.triggerLength) return;

    const emotesLength = state.parseEmotes().filter((o) => o.type === 'emote').length;
    if (emotesLength < settings.triggerLength) return;

    return this.returnByWarnedState('emotes', state.userInfo.userId, settings);
  }
}