import { Empty } from '@tsuwari/grpc/generated/websockets/google/protobuf/empty';
import {
  ObsAudioDecreaseVolumeMessage,
  ObsAudioIncreaseVolumeMessage,
  ObsAudioSetVolumeMessage,
  ObsSetSceneMessage,
  ObsToggleAudioMessage,
  ObsToggleSourceMessage,
} from '@tsuwari/grpc/generated/websockets/websockets';
import SocketIo from 'socket.io';

import { authMiddleware, io } from '../libs/io.js';

const sockets: Map<string, SocketIo.Socket> = new Map();
export const obsNameSpace = io.of('obs');

obsNameSpace.use(authMiddleware);
obsNameSpace.on('connection', async (socket) => {
  const channelId = socket.handshake.auth.channelId;
  sockets.set(channelId, socket);
});

export const onSetAudio = async (data: ObsAudioSetVolumeMessage): Promise<Empty> => {
  const socket = sockets.get(data.channelId);
  if (!socket) return {};

  socket.emit('setVolume', data);

  return {};
};

export const onAudioIncrease = async (data: ObsAudioIncreaseVolumeMessage): Promise<Empty> => {
  const socket = sockets.get(data.channelId);
  if (!socket) return {};

  socket.emit('increaseVolume', data);

  return {};
};

export const onAudioDecrease = async (data: ObsAudioDecreaseVolumeMessage): Promise<Empty> => {
  const socket = sockets.get(data.channelId);
  if (!socket) return {};

  socket.emit('decreaseVolume', data);

  return {};
};

export const onToggleAudioSource = async (data: ObsToggleAudioMessage): Promise<Empty> => {
  const socket = sockets.get(data.channelId);
  if (!socket) return {};

  socket.emit('toggleAudioSource', data);

  return {};
};

export const onSetScene = async (data: ObsSetSceneMessage): Promise<Empty> => {
  const socket = sockets.get(data.channelId);
  if (!socket) return {};

  socket.emit('setScene', data);

  return {};
};

export const onToggleSource = async (data: ObsToggleSourceMessage): Promise<Empty> => {
  const socket = sockets.get(data.channelId);
  if (!socket) return {};

  socket.emit('toggleSource', data);

  return {};
};

