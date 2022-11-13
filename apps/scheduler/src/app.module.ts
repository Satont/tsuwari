import { Module } from '@nestjs/common';
import { ScheduleModule } from '@nestjs/schedule';
import { RedisORMModule } from '@tsuwari/redis';
import { RedisModule, RedisService, TwitchApiService } from '@tsuwari/shared';

import { DefaultCommandsCreatorModule } from './default-commands-creator/default-commands-creator.module.js';
import { DeleteOldMessagesModule } from './delete-old-messages/deleteoldmessages.module.js';
import { DotaModule } from './dota/dota.module.js';
import { MicroservicesModule } from './microservices/microservices.module.js';
import { OnlineUsersModule } from './online-users/online-users.module.js';
import { StreamStatusModule } from './streamstatus/streamstatus.module.js';

@Module({
  imports: [
    RedisORMModule,
    RedisModule,
    ScheduleModule.forRoot(),
    StreamStatusModule,
    MicroservicesModule,
    DotaModule,
    DefaultCommandsCreatorModule,
    OnlineUsersModule,
    DeleteOldMessagesModule,
  ],
  providers: [TwitchApiService, RedisService],
})
export class AppModule {}
