import { ClientProxy as CP } from '@nestjs/microservices';
import { Command, CommandPermission, Response } from '@tsuwari/prisma';
import { rawDataSymbol } from '@twurple/common';
import { EventSubChannelUpdateEvent } from '@twurple/eventsub';
import { Observable } from 'rxjs';

export interface ClientProxyCommands {
  'streamstatuses.process': {
    input: string[],
    result: boolean
  },
  'bots.getDefaultCommands': {
    input: any,
    result: Array<{ name: string, description?: string, visible: boolean, permission: CommandPermission }>
  },
  'bots.getVariables': {
    input: any,
    result: Array<{
      name: string,
      example?: string,
      description?: string
    }>
  },
  'parseChatMessage': {
    input: string,
    result: string[]
  },
  'parseResponse': {
    input: {
      userId?: string,
      channelId: string,
      userName?: string,
      userDisplayName?: string,
      text: string
    };
    result: string[];
  },
  'setCommandCache': {
    input: Command & { responses?: Response[] },
    result: any,
  }
}

export interface ClientProxyEvents {
  'streams.online': {
    input: { streamId: string, channelId: string },
    result: any
  },
  'streams.offline': {
    input: { channelId: string },
    result: any
  },
  'stream.update': {
    input: EventSubChannelUpdateEvent[typeof rawDataSymbol],
    result: any,
  }
  'bots.joinOrLeave': {
    input: {
      action: 'join' | 'part',
      username: string,
      botId: string,
    },
    result: any
  },
  'bots.addTimerToQueue': {
    input: string,
    result: any
  },
  'bots.removeTimerFromQueue': ClientProxyEvents['bots.addTimerToQueue'],
  'dota.cacheAccountsMatches': {
    input: string[],
    result: any,
  }
}

export type ClientProxyResult<K extends keyof ClientProxyCommands> = Observable<ClientProxyCommands[K]['result']>
export type ClientProxyCommandsKey = keyof ClientProxyCommands
export type ClientProxyEventsKey = keyof ClientProxyEvents


export abstract class ClientProxy extends CP {
  abstract send<TEvent extends keyof ClientProxyCommands>(pattern: TEvent, data: ClientProxyCommands[TEvent]['input']): Observable<ClientProxyCommands[TEvent]['result']>;
  abstract emit<TEvent extends keyof ClientProxyEvents>(pattern: TEvent, data: ClientProxyEvents[TEvent]['input']): Observable<ClientProxyEvents[TEvent]['result']>;
}
