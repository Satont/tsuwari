import { resolve } from 'node:path';
import 'reflect-metadata';

import dotenv from 'dotenv';
import { DataSource } from 'typeorm';

dotenv.config({ path: resolve(process.cwd(), '../../.env') });

import { Bot } from './entities/Bot.js';
import { Channel } from './entities/Channel.js';
import { ChannelChatMessage } from './entities/ChannelChatMessage.js';
import { ChannelCommand } from './entities/ChannelCommand.js';
import { ChannelCustomvar } from './entities/ChannelCustomvar.js';
import { ChannelDotaAccount } from './entities/ChannelDotaAccount.js';
import { ChannelEvent } from './entities/ChannelEvent.js';
import { ChannelDonationEvent } from './entities/channelEvents/Donation.js';
import { ChannelFollowEvent } from './entities/channelEvents/Follow.js';
import { ChannelGreeting } from './entities/ChannelGreeting.js';
import { ChannelIntegration } from './entities/ChannelIntegration.js';
import { ChannelKeyword } from './entities/ChannelKeyword.js';
import { ChannelModerationSetting } from './entities/ChannelModerationSetting.js';
import { ChannelModerationWarn } from './entities/ChannelModerationWarn.js';
import { ChannelPermit } from './entities/ChannelPermit.js';
import { ChannelStream } from './entities/ChannelStream.js';
import { ChannelTimer } from './entities/ChannelTimer.js';
import { ChannelTimerResponse } from './entities/ChannelTimerResponse.js';
import { CommandResponse } from './entities/CommandResponse.js';
import { CommandUsage } from './entities/CommandUsage.js';
import { DashboardAccess } from './entities/DashboardAccess.js';
import { DotaGameMode } from './entities/DotaGameMode.js';
import { DotaHero } from './entities/DotaHero.js';
import { DotaMatch } from './entities/DotaMatch.js';
import { DotaMatchCard } from './entities/DotaMatchCard.js';
import { DotaMatchResult } from './entities/DotaMatchResult.js';
import { Integration } from './entities/Integration.js';
import { Notification } from './entities/Notification.js';
import { NotificationMessage } from './entities/NotificationMessage.js';
import { Token } from './entities/Token.js';
import { User } from './entities/User.js';
import { UserFile } from './entities/UserFile.js';
import { UserOnline } from './entities/UserOnline.js';
import { UserStats } from './entities/UserStats.js';
import { UserViewedNotification } from './entities/UserViewedNotification.js';

export * from 'typeorm';

export const AppDataSource = new DataSource({
  type: 'postgres',
  url: process.env.DATABASE_URL,
  entities: [
    Bot,
    Channel,
    ChannelCommand,
    ChannelCustomvar,
    ChannelDotaAccount,
    ChannelGreeting,
    ChannelIntegration,
    ChannelKeyword,
    ChannelModerationSetting,
    ChannelModerationWarn,
    ChannelPermit,
    ChannelTimer,
    ChannelTimerResponse,
    CommandResponse,
    CommandUsage,
    DashboardAccess,
    DotaGameMode,
    DotaHero,
    DotaMatch,
    DotaMatchCard,
    DotaMatchResult,
    Integration,
    Notification,
    NotificationMessage,
    Token,
    User,
    UserFile,
    UserStats,
    UserViewedNotification,
    UserOnline,
    ChannelEvent,
    ChannelFollowEvent,
    ChannelDonationEvent,
    ChannelStream,
    ChannelChatMessage,
  ],
  subscribers: [],
  migrations: ['src/migrations/*.ts'],
  migrationsTableName: 'typeorm_migrations',
});
