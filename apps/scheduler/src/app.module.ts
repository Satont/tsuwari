import { Module } from '@nestjs/common';
import { ScheduleModule } from '@nestjs/schedule';
import { PrismaModule } from '@tsuwari/prisma';
import { RedisService, RedisModule } from '@tsuwari/shared';

import { DefaultCommandsCreatorModule } from './default-commands-creator/default-commands-creator.module.js';
import { DotaModule } from './dota/dota.module.js';
import { MicroservicesModule } from './microservices/microservices.module.js';
import { StreamStatusModule } from './streamstatus/streamstatus.module.js';

@Module({
  imports: [
    PrismaModule,
    RedisModule,
    ScheduleModule.forRoot(),
    StreamStatusModule,
    MicroservicesModule,
    DotaModule,
    DefaultCommandsCreatorModule,
    // IncreaseWatchedModule,
  ],
  providers: [RedisService],
})
export class AppModule { }
