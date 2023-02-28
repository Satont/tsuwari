import { useCallback, useEffect, useRef, useState } from 'react';
import { useParams } from 'react-router-dom';
import { io, Socket } from 'socket.io-client';

declare global {
  interface Window {
    webkitAudioContext: typeof AudioContext
  }
}

export const TTS: React.FC = () => {
  const { apiKey } = useParams();
  const [tts, setTTS] = useState<Socket | null>(null);

  useEffect(() => {
    const tts = io(
      `${`${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}`}/tts`,
      {
        transports: ['websocket'],
        autoConnect: true,
        auth: (cb) => {
          cb({ apiKey });
        },
      },
    );

    setTTS(tts);

    return () => {
      tts.removeAllListeners();
      tts.disconnect();
    };
  }, []);

  const queueRef = useRef<Array<Record<string, string>>>([]);
  const currentAudioBuffer = useRef<AudioBufferSourceNode | null>(null);

  useEffect(() => {
    if (tts) {
      tts.on('connect', () => {
        console.log('connected');
      });

      tts.on('say', (data) => {
        queueRef.current.push(data);

        if (queueRef.current.length === 1) {
          processQueue();
        }
      });

      tts.on('skip', () => {
        currentAudioBuffer.current?.stop();
      });
    }

    return () => {
      tts?.removeAllListeners();
    };
  }, [tts]);

  const processQueue = useCallback(async () => {
    if (queueRef.current.length === 0) {
      return;
    }

    await say(queueRef.current[0]);
    queueRef.current = queueRef.current.slice(1);

    // Process the next item in the queue
    processQueue();
  }, []);

  const say = async (data: Record<string, string>) => {
    if (!apiKey) return;

    const query = new URLSearchParams(data);

    const audioContext = new (window.AudioContext || window.webkitAudioContext)();
    const gainNode = audioContext.createGain();

    const req = await fetch(`/api/v1/tts/say?${query}`, {
      headers: {
        'Api-Key': apiKey,
      },
    });
    const arrayBuffer = await req.arrayBuffer();

    const source = audioContext.createBufferSource();
    currentAudioBuffer.current = source;

    source.buffer = await audioContext.decodeAudioData(arrayBuffer);

    gainNode.gain.value = parseInt(data.volume) / 100;
    source.connect(gainNode);
    gainNode.connect(audioContext.destination);

    return new Promise((resolve) => {
      source.onended = () => {
        currentAudioBuffer.current = null;
        resolve(null);
      };

      source.start(0);
    });
  };

  return <></>;
};
