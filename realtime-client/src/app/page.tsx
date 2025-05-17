// Next.js クライアント側: 複数 WebSocket クライアントを立ち上げ、フレームを送信しACKを待つ
'use client';

import { useEffect } from 'react';

const NUM_CLIENTS = 5; // クライアント数
const SEND_INTERVAL = 1000; // 送信間隔
const RETRY_TIMEOUT = 3000; // ACKが来ないときの再送時間

export default function Page() {
  useEffect(() => {
    //@ts-expect-error
    const clients = [];

    for (let i = 0; i < NUM_CLIENTS; i++) {
      let frameCounter = 0;
      const pendingFrames = new Map();
      const socket = new WebSocket('ws://localhost:8080/ws');

      socket.onopen = () => {
        console.log(`🟢 Client ${i} connected`);

        const sendFrame = () => {
          const frameId = frameCounter++;
          const message = JSON.stringify({
            clientId: i,
            frameId,
            data: `Frame from client ${i}`
          });
          console.log(`📤 Client ${i} sending frameId: ${frameId}`);
          socket.send(message);

          // 送信時間記録してACK待ちに登録
          const timeout = setTimeout(() => {
            console.warn(`⏳ Client ${i} frameId ${frameId} retrying...`);
            socket.send(message);
          }, RETRY_TIMEOUT);

          pendingFrames.set(frameId, timeout);
        };

        setTimeout(() => {
          setInterval(sendFrame, SEND_INTERVAL + Math.random() * 500);
        }, i * 500);
      };

      socket.onmessage = (e) => {
        const msg = JSON.parse(e.data);
        if (msg.ack) {
          const ackId = msg.frameId;
          if (pendingFrames.has(ackId)) {
            clearTimeout(pendingFrames.get(ackId));
            pendingFrames.delete(ackId);
            console.log(`✅ Client ${i} ACK received for frameId: ${ackId}`);
          }
        } else {
          console.log(`📩 Client ${i} received unknown message:`, msg);
        }
      };

      socket.onerror = (e) => {
        console.error(`❌ Client ${i} error`, e);
      };

      socket.onclose = () => {
        console.log(`🔴 Client ${i} disconnected`);
        pendingFrames.forEach((timeout) => clearTimeout(timeout));
      };

      clients.push(socket);
    }

    return () => {
      //@ts-expect-error
      clients.forEach((s) => s.close());
    };
  }, []);

  return (
    <main className="p-4 text-sm">
      <h1 className="text-lg font-bold">🔁 WebSocket Multi Client Test + Retry</h1>
      <p>クライアントがフレームを送信し、ACKが来ない場合は再送します。</p>
      <p>詳細はブラウザのコンソールをご確認ください。</p>
    </main>
  );
}