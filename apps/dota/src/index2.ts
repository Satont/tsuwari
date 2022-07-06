import { dirname, resolve } from 'path';
import { fileURLToPath } from 'url';

import { config } from '@tsuwari/config';
import protobufjs from 'protobufjs';
import SteamUser from 'steam-user';
import SteamID from 'steamid';

const client = new SteamUser({
  autoRelogin: true,
});

const proto = await new protobufjs.Root().load(resolve(dirname(fileURLToPath(import.meta.url)), '..', 'protos', 'dota2', 'dota_gcmessages_client_watch.proto'), {
  keepCase: true,
});

client.logOn({
  accountName: config.STEAM_USERNAME,
  password: config.STEAM_PASSWORD,
});

client.on('loggedOn', (e, d) => {
  console.log(`Logged into Steam as ${e.client_supplied_steamid} ${e.vanity_url}`);
  client.setPersona(1);
  client.gamesPlayed([570], true);
});

type RP = SteamUser.RichPresence & {
  watching_server?: string,
  param0?: string, // lobby type
  param1?: string,
  param2?: string, // hero
  WatchableGameID?: string,
}
function converUsers(users: Record<string, { richPresence: RP }>) {
  return Object.entries(users).map(user => {
    const USERRP = user[1].richPresence;
    return {
      userId: new SteamID(user[0]).accountid,
      steamId: user[0],
      richPresence: {
        ...USERRP,
        watching_server: USERRP.watching_server ? new SteamID(USERRP.watching_server).getSteamID64() : USERRP.watching_server,
        createdAt: new Date(),
        lobbyId: USERRP.WatchableGameID,
      },
    };
  });
}

client.on('appLaunched', async (appid) => {
  console.log('Dota launched', appid);
  const accs = [911977148, 1102609846, 70388657, 86738694, 164217422].map(SteamID.fromIndividualAccountID).map(id => id.getSteamID64());
  const type = proto.root.lookupType('CMsgClientToGCFindTopSourceTVGames');

  client.requestRichPresence(570, accs, 'english', (error, data) => {
    const users = converUsers(data.users);

    const lobbyIds = new Set(users.filter(u => u.richPresence.lobbyId).map(u => u.richPresence.lobbyId));
    if (!lobbyIds.size) return;

    const newMsg = type.encode({
      search_key: '',
      league_id: 0,
      hero_id: 0,
      start_game: 0,
      game_list_index: 0,
      lobby_ids: [...lobbyIds.values()],
    });
    client.sendToGC(570, 8009, {}, Buffer.from(newMsg.finish()));
    setInterval(() => client.sendToGC(570, 8009, {}, Buffer.from(newMsg.finish())), 5000);
  });
});

client.on('disconnected', (e, w) => console.log(e, w));

client.on('receivedFromGC', async (app, msg, payload) => {
  console.log(msg);
  if (msg === 4009) {
    const proto = await new protobufjs.Root().load(resolve(dirname(fileURLToPath(import.meta.url)), '..', 'protos', 'dota2', 'gcsdk_gcmessages.proto'), {
      keepCase: true,
    });
    const type = proto.root.lookupType('CMsgConnectionStatus');
    console.log(type.decode(payload).toJSON());

    const helloType = proto.root.lookupType('CMsgClientHello');
    client.sendToGC(570, 4006, {}, Buffer.from(helloType.encode({}).finish()));
    return;
  }

  if (msg === 8010) {
    const type = proto.root.lookupType('CMsgGCToClientFindTopSourceTVGamesResponse');

    const data = type.decode(payload).toJSON() as {
      game_list?: Array<{
        lobby_type: number,
        game_mode: number,
        average_mmr: number,
        players: Array<{ account_id: number, hero_id: number }>,
        weekend_tourney_bracket_round?: string,
        weekend_tourney_skill_level?: string,
        match_id?: string
      }>
    };

    if (data.game_list) {
      const match = data.game_list.find(g => g.players.find(p => p.account_id == 164217422));
      console.log(match, { depth: null });
    }
  }
});

client.on('error', (e) => {
  console.log(e);
});