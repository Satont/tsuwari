import { resolve } from 'path';

import dotenv from 'dotenv';
import { bool, cleanEnv, str } from 'envalid';

try {
  dotenv.config({ path: resolve(process.cwd(), '../../.env') });
  // eslint-disable-next-line no-empty
} catch {}

export const config = cleanEnv(process.env, {
  DATABASE_URL: str({
    default: 'postgresql://tsuwari:tsuwari@postgres:5432/tsuwari?schema=public',
  }),
  NODE_ENV: str({ choices: ['development', 'production'], default: 'development' }),
  TWITCH_CLIENTID: str(),
  TWITCH_CLIENTSECRET: str(),
  TWITCH_CALLBACKURL: str({ default: 'http://localhost:3005/login' }),
  JWT_EXPIRES_IN: str({ default: '5m' }),
  JWT_ACCESS_SECRET: str({ default: 'CoolSecretForAccess' }),
  JWT_REFRESH_SECRET: str({ default: 'CoolSecretForRefresh' }),
  REDIS_URL: str({ default: 'redis://localhost:6379/0' }),
  SAY_IN_CHAT: bool({ default: true }),
  HOSTNAME: str({ default: '' }),
  STEAM_USERNAME: str({ default: '' }),
  STEAM_PASSWORD: str({ default: '' }),
  STEAM_API_KEY: str({ default: '' }),
  MINIO_USER: str({ devDefault: '' }),
  MINIO_PASSWORD: str({ devDefault: '' }),
  MINIO_URL: str({ devDefault: '' }),
});
