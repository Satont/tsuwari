import { ChannelEvent, EventType } from '@tsuwari/typeorm/entities/ChannelEvent';
import { ChannelDonationEvent } from '@tsuwari/typeorm/entities/channelEvents/Donation';
import { ChannelIntegration } from '@tsuwari/typeorm/entities/ChannelIntegration';

import { donatelloStore, removeIntegration, typeorm } from '../index.js';
import { eventsGrpcClient } from '../libs/eventsGrpc.js';

export class Donatello {
  #interval: NodeJS.Timer;

  constructor(
    private readonly apiKey: string,
    private readonly twitchUserId: string,
  ) {}

  async init() {
    console.info(`Donatello: start polling ${this.twitchUserId}`);
    this.#interval = setInterval(() => {
      this.#checkNewDonates();
    }, 10 * 1000);
  }

  async #checkNewDonates() {
    const req = await fetch('https://donatello.to/api/v1/donates', {
      method: 'GET',
      headers: {
        'X-Token': this.apiKey,
      },
    });

    if (!req.ok) {
      console.error(`Donatello(${this.twitchUserId}): ${await req.text()}`);
      return;
    }

    const response = await req.json() as Response;
    const repository = typeorm.getRepository(ChannelDonationEvent);

    if (!response.content.length) return;

    const latestDonation = await repository.findOneBy({
      donateId: response.content.at(0)!.pubId,
    });

    // no new donations
    if (latestDonation) {
      return;
    }

    for (const donation of response.content) {
      const event = await typeorm.getRepository(ChannelEvent).create({
        channelId: this.twitchUserId,
        type: EventType.DONATION,
      });
      await typeorm.getRepository(ChannelEvent).save(event);

      await repository.save({
        amount: Number(donation.amount),
        eventId: event.id,
        currency: donation.currency,
        toUserId: this.twitchUserId,
        message: donation.message,
        username: donation.clientName,
        donateId: donation.pubId,
      });

      eventsGrpcClient.donate({
        amount: donation.amount,
        message: donation.message,
        currency: donation.currency,
        baseInfo: { channelId: this.twitchUserId },
        userName: donation.clientName,
      });
    }
  }

  async destroy() {
    clearInterval(this.#interval);
  }
}

export interface Response {
  content: Donation[]
  page: number
  size: number
  num: number
  first: boolean
  last: boolean
  total: number
}

export interface Donation {
  pubId: string
  clientName: string
  message: string
  amount: string
  currency: string
  goal: string
  isPublished: boolean
  createdAt: string
}

export async function addDonatelloIntegration(integration: ChannelIntegration) {
  if (!integration.integration || !integration.apiKey) {
    return;
  }

  if (donatelloStore.get(integration.channelId)) {
    await removeIntegration(integration);
  }
  const instance = new Donatello(
    integration.apiKey,
    integration.channelId,
  );
  await instance.init();

  return instance;
}
