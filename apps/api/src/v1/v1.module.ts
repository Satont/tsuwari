import { Module } from '@nestjs/common';

import { CommandsModule } from './commands/commands.module.js';
import { GreetingsModule } from './greetings/greetings.module.js';
import { LastfmModule } from './integrations/lastfm/lastfm.module.js';
import { SpotifyModule } from './integrations/spotify/spotify.module.js';
import { KeywordsModule } from './keywords/keywords.module.js';
import { TimersModule } from './timers/timers.module.js';
import { VariablesModule } from './variables/variables.module.js';

@Module({
  imports: [CommandsModule, GreetingsModule, TimersModule, SpotifyModule, LastfmModule, KeywordsModule, VariablesModule],
})
export class V1Module { }
