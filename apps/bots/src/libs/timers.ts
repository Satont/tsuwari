import { Logger } from '@nestjs/common';
import { config } from '@tsuwari/config';
import { Timer } from '@tsuwari/prisma';
import { Queue } from '@tsuwari/shared';
import { HelixStreamData } from '@twurple/api/lib/index.js';

import { Bots, staticApi } from '../bots.js';
import { nestApp } from '../nest/index.js';
import { ParserService } from '../nest/parser/parser.service.js';
import { prisma } from './prisma.js';
import { redis, redisOm } from './redis.js';

const logger = new Logger('Timers');

export const timersQueue = new Queue<Timer>(async function (taskId: string) {
  const timer = await prisma.timer.findFirst({
    where: {
      id: taskId,
    },
    include: {
      channel: true,
    },
  });

  if (!timer || !timer.enabled) {
    return;
  }

  const rawStream = await redis.get(`streams:${timer.channelId}`);
  if (!rawStream) return;
  const stream = JSON.parse(rawStream) as HelixStreamData & { parsedMessages?: number };

  stream.parsedMessages = stream.parsedMessages ?? 0;

  if (
    timer.messageInterval > 0 &&
    timer.lastTriggerMessageNumber - stream.parsedMessages + timer.messageInterval > 0
  ) {
    return;
  }

  const responses = timer.responses as Array<string>;

  const bot = Bots.cache.get(timer.channel.botId);
  const user = await staticApi.users.getUserById(timer.channelId);

  const response = responses[timer.last];
  if (!response) return;

  if (!bot || !user) {
    return;
  }

  if (bot._authProvider) {
    const service = nestApp.get(ParserService);
    const parsedResponses = await service.parseResponse({
      channelId: timer.channelId,
      text: response,
    });

    if (parsedResponses) {
      for (const r of parsedResponses) {
        if (config.isProd) {
          bot.say(user.name, r);
        } else {
          logger.log(`${user.name} -> ${r}`);
        }
      }
    }
  }

  await prisma.timer.update({
    where: {
      id: timer.id,
    },
    data: {
      last: ++timer.last % (timer.responses as string[]).length,
      lastTriggerMessageNumber: stream.parsedMessages as number,
    },
  });
});

const getId = (t: Timer | string) => (typeof t === 'string' ? t : t.id);
export async function addTimerToQueue(timerOrId: Timer | string) {
  const id = getId(timerOrId);
  let timer: Timer | null;

  if (typeof id === 'string') {
    timer = await prisma.timer.findFirst({ where: { id: id as string } });
    if (!timer?.enabled) return;
  } else {
    timer = timerOrId as Timer;
  }

  removeTimerFromQueue(timerOrId);
  if (timer) {
    timersQueue.addTimerToQueue(timer.id, timer, {
      interval: timer.timeInterval * (config.isDev ? 1000 : 60000),
    });
  }
}

export function removeTimerFromQueue(timer: Timer | string) {
  const id = getId(timer);

  timersQueue.removeTimerFromQueue(id);
}
