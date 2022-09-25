import { ChannelEvent, EventType } from '@tsuwari/typeorm/entities/ChannelEvent';
import { ChannelDonationEvent } from '@tsuwari/typeorm/entities/channelEvents/Donation';
import { ChannelIntegration } from '@tsuwari/typeorm/entities/ChannelIntegration';
import * as IO from 'socket.io-client';

import { typeorm } from '../index.js';
import { sendMessage } from '../libs/sender.js';

type Socket = typeof IO.Socket;

export class StreamLabs {
  #conn: Socket;

  constructor(token: string, private readonly twitchUserId: string) {
    this.#conn = IO.connect(`https://sockets.streamlabs.com?token=${token}`, {
      transports: ['websocket'],
    });

    this.#conn.on('event', (eventData: Event) => {
      if (eventData.type === 'donation') {
        eventData.message.forEach((m) => this.#handler(m));
      }
    });
  }

  async #handler(data: Message) {
    const event = await typeorm.getRepository(ChannelEvent).save({
      channelId: this.twitchUserId,
      type: EventType.DONATION,
    });
    await typeorm.getRepository(ChannelDonationEvent).save({
      event,
      amount: data.amount,
      currency: data.currency,
      toUserId: this.twitchUserId,
      message: data.message,
      username: data.from,
    });
    sendMessage({
      channelId: this.twitchUserId,
      message: `${data.from}: ${data.amount}${data.currency} ${data.message}`,
      channelName: '',
    });
  }
}

export type Event = {
  type: 'donation';
  message: Message[];
  for: string;
  event_id: string;
};

export type Message = {
  name: string;
  isTest: boolean;
  formatted_amount: string;
  amount: number;
  message: string;
  currency: string;
  to: { name: string };
  from: string;
  from_user_id: number;
  _id: string;
  priority: number;
};

export async function addStreamlabsIntegration(integration: ChannelIntegration) {
  if (
    !integration.accessToken ||
    !integration.refreshToken ||
    !integration.integration ||
    !integration.integration.clientId ||
    !integration.integration.clientSecret ||
    !integration.integration.redirectUrl
  ) {
    return;
  }

  const refresh = await fetch('https://www.twitchalerts.com/api/v1.0/token', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
    body: new URLSearchParams({
      grant_type: 'refresh_token',
      refresh_token: integration.refreshToken,
      redirect_url: integration.integration.redirectUrl,
      client_id: integration.integration.clientId,
      client_secret: integration.integration.clientSecret,
    }).toString(),
  });

  if (!refresh.ok) {
    console.error(await refresh.text());
    return;
  }

  const refreshResponse = await refresh.json();

  await typeorm
    .getRepository(ChannelIntegration)
    .update(
      { id: integration.id },
      { accessToken: refreshResponse.access_token, refreshToken: refreshResponse.refresh_token },
    );

  const socketRequest = await fetch(
    `https://streamlabs.com/api/v1.0/socket/token?access_token=${refreshResponse.access_token}`,
  );

  if (!socketRequest.ok) {
    console.error(await socketRequest.text());
    return;
  }

  const { socket_token } = await socketRequest.json();

  const instance = new StreamLabs(socket_token, integration.channelId);

  return instance;
}
