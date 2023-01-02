import { Empty } from '@tsuwari/grpc/generated/websocket/google/protobuf/empty';
import {
  YoutubeAddSongToQueueRequest,
  YoutubeRemoveSongFromQueueRequest,
} from '@tsuwari/grpc/generated/websocket/websocket';
import { IsNull } from '@tsuwari/typeorm';
import { RequestedSong } from '@tsuwari/typeorm/entities/RequestedSong';
import SocketIo from 'socket.io';

import { authMiddleware, io } from '../libs/io.js';
import { redis } from '../libs/redis.js';
import { typeorm } from '../libs/typeorm.js';

const sockets: Map<string, SocketIo.Socket> = new Map();
const repository = typeorm.getRepository(RequestedSong);

export const youtubeNamespace = io.of('youtube');

youtubeNamespace.use(authMiddleware);
youtubeNamespace.on('connection', async (socket) => {
  const channelId = socket.handshake.auth.channelId;
  sockets.set(channelId, socket);

  socket.on('currentQueue', async (cb) => {
    const songs = await repository.find({
      where: {
        channelId,
        deletedAt: IsNull(),
      },
      order: {
        queuePosition: 'asc',
      },
    });

    cb(songs);
  });

  socket.on('skip', async (id) => {
    const entity = await repository.findOneBy({ id });
    if (entity) {
      await repository.softDelete({ id });
    }
    redis.del(`songrequests:youtube:${channelId}:currentPlaying`);
  });

  socket.on('disconnect', () => {
    sockets.delete(channelId);
  });

  socket.on('play', async (data) => {
    await redis.set(
      `songrequests:youtube:${channelId}:currentPlaying`,
      data.id,
      'PX',
      data.duration,
    );
  });

  socket.on('pause', () => {
    redis.del(`songrequests:youtube:${channelId}:currentPlaying`);
  });

  socket.on('newOrder', async (videos: RequestedSong[]) => {
    const currentVideosCount = await repository.count({
      where: { channelId },
    });

    const filteredVideos = videos.filter((v) => v.queuePosition > 1);

    if (filteredVideos.some((v) => v.queuePosition > currentVideosCount + 2)) {
      return;
    }

    for (const video of filteredVideos) {
      await repository.update({ id: video.id }, { queuePosition: video.queuePosition });
    }
  });
});

export const onAddRequest = async (data: YoutubeAddSongToQueueRequest): Promise<Empty> => {
  const socket = sockets.get(data.channelId);
  if (!socket) return {};

  const entity = await repository.findOneBy({
    id: data.entityId,
  });
  socket.emit('newTrack', entity);

  return {};
};

export const onRemoveRequest = async (data: YoutubeRemoveSongFromQueueRequest): Promise<Empty> => {
  const socket = sockets.get(data.channelId);
  if (!socket) return {};

  const entity = await repository.findOne({
    where: { id: data.entityId },
    withDeleted: true,
  });
  socket.emit('removeTrack', entity);

  return {};
};
