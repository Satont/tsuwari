import { Injectable } from '@nestjs/common';
import { Command, CommandPermission } from '@tsuwari/prisma';
import { RedisService } from '@tsuwari/shared';
import { ChatUser } from '@twurple/chat';

export type CommandConditional = Command & { responses: (string | undefined)[] | undefined };

@Injectable()
export class HelpersService {
  constructor(
    private readonly redis: RedisService,
  ) { }

  async getChannelCommandsNamesFromRedis(channelId: string) {
    const channelCommandsKeys = await this.redis.keys(`commands:${channelId}:*`);

    if (!channelCommandsKeys.length) return;

    const channelCommandsNames = channelCommandsKeys.map((c) => c.split(':')[2]) as string[];
    if (!channelCommandsNames || !channelCommandsNames.length) return;

    return channelCommandsNames;
  }

  async getChannelCommandsByNamesFromRedis(channelId: string, names: string[]) {
    const result: CommandConditional[] = [];

    for (const name of names) {
      const command: CommandConditional = await this.redis.hgetall(
        `commands:${channelId}:${name}`,
      ) as unknown as CommandConditional;

      if (!Object.keys(command).length) continue;
      if ((JSON.parse(command.aliases as string) as Array<string>).includes(name)) continue;
      result.push(command);
    }

    return result;
  }

  getUserPermissions(userInfo: ChatUser): Record<CommandPermission, boolean> {
    return {
      BROADCASTER: userInfo.isBroadcaster,
      MODERATOR: userInfo.isMod,
      VIP: userInfo.isVip,
      SUBSCRIBER: userInfo.isSubscriber || userInfo.isFounder,
      FOLLOWER: true,
      VIEWER: true,
    };
  }

  hasPermission(perms: Record<CommandPermission, boolean>, searchForPermission: CommandPermission) {
    if (!searchForPermission) return true;

    const userPerms = Object.entries(perms) as [CommandPermission, boolean][];
    const permissionIndex = userPerms.find((perm) => perm[0] === searchForPermission);
    const commandPermissionIndex = userPerms.indexOf(permissionIndex!);

    const hasPerm = userPerms.some((p, index) => p[1] && index <= commandPermissionIndex);
    return hasPerm;
  }

  buildCooldownKey(command: CommandConditional, senderId: string) {
    if (command.cooldownType === 'GLOBAL') {
      return `commands:cooldowns:${command.id}`;
    } else {
      return `commands:cooldowns:${command.id}:${senderId}`;
    }
  }

  async isOnCooldown(command: CommandConditional, senderId: string) {
    if (!command.cooldown) return false;
    const item = await this.redis.get(this.buildCooldownKey(command, senderId));
    return item !== null;
  }

  setCommandCooldown(command: CommandConditional, senderId: string) {
    if (command.cooldown && command.cooldown <= 0) return;

    if (command.cooldownType === 'GLOBAL') {
      this.redis.set(`commands:cooldowns:${command.id}`, '', 'EX', command.cooldown!);
    } else {
      this.redis.set(`commands:cooldowns:${command.id}:${senderId}`, '', 'EX', command.cooldown!);
    }
  }

  findCommandInArrayOfNames(message: string, commands: string[]) {
    message = message.substring(1).trim();

    const msgArray = message.toLowerCase().split(' ');
    let isFound = false;
    let commandName = '';

    for (let i = 0, len = msgArray.length; i < len; i++) {
      const query = msgArray.join(' ');
      const find = commands.find((c) => c === query);
      if (!find) {
        msgArray.pop();
        continue;
      }

      commandName = find;
      isFound = true;
    }

    return {
      isFound,
      commandName,
      params: message.replace(new RegExp(`^${commandName}`), '').trim(),
    };
  }
}
