import { config } from '@tsuwari/config';
import { ApiClient } from '@twurple/api';
import { ClientCredentialsAuthProvider } from '@twurple/auth';

export const twitchApi = new ApiClient({
  authProvider: new ClientCredentialsAuthProvider(
    config.TWITCH_CLIENTID,
    config.TWITCH_CLIENTSECRET,
  ),
});
