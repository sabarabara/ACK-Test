'use client';

import { useEffect, useRef, useState } from 'react';

type FrameMessage = {
  frameId: number;
  data: string;
};

type AckMessage = {
  type: 'ack';
  frameId: number;
};

export default function Page() {
  const socketRef = useRef<WebSocket | null>(null);
  const [log, setLog] = useState<string[]>([]);
  const pendingFrames = useRef<Map<number, NodeJS.Timeout>>(new Map());
  const frameIdRef = useRef<number>(1);

  // ログ追加用ヘルパー
  const addLog = (message: string) => {
    setLog((prev) => [message, ...prev.slice(0, 50)]);
  };

  useEffect(() => {
    const socket = new WebSocket('ws://localhost:8080/ws');
    socketRef.current = socket;

    socket.onopen = () => {
      addLog('✅ WebSocket connected');
      startSendingFrames();
    };

    socket.onmessage = (event) => {
      try {
        const data: AckMessage = JSON.parse(event.data);
        if (data.type === 'ack') {
          addLog(`✅ ACK received for frame ${data.frameId}`);
          // 再送タイマーを止める
          const timer = pendingFrames.current.get(data.frameId);
          if (timer) clearTimeout(timer);
          pendingFrames.current.delete(data.frameId);
        }
      } catch (err) {
        addLog(`⚠️ Invalid ACK: ${event.data}`);
      }
    };

    socket.onerror = (err) => {
      addLog('❌ WebSocket error');
      console.error(err);
    };

    socket.onclose = () => {
      addLog('❌ WebSocket closed');
    };

    return () => {
      socket.close();
    };
  }, []);

  const startSendingFrames = () => {
    const interval = setInterval(() => {
      const frameId = frameIdRef.current++;
      const message: FrameMessage = {
        frameId,
        data: `frame-${frameId}`,
      };

      const socket = socketRef.current;
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(message));
        addLog(`📤 Sent frameId ${frameId}`);

        // 再送タイマー（3秒後に再送）
        const retryTimer = setTimeout(() => {
          addLog(`⏳ Retry frameId ${frameId}`);
          socket.send(JSON.stringify(message));
        }, 3000);

        pendingFrames.current.set(frameId, retryTimer);
      }
    }, 1000); // 1秒ごとに送信

    // 停止したいときは clearInterval(interval)
  };

  return (
    <main className="p-6 font-mono">
      <h1 className="text-xl font-bold mb-4">🧪 WebSocket Frame Sender</h1>
      <div className="bg-gray-100 p-4 rounded shadow h-96 overflow-auto">
        {log.map((line, i) => (
          <div key={i}>{line}</div>
        ))}
      </div>
    </main>
  );
}
