import { io, Socket } from 'socket.io-client';
import { useAuth } from '../contexts/AuthContext';

let socket: Socket | null = null;

export const connectWebSocket = async (token: string): Promise<Socket> => {
  if (socket?.connected) return socket;

  socket = io(import.meta.env.VITE_WS_URL, {
    auth: { token },
    reconnection: true,
    reconnectionAttempts: 5,
  });

  return new Promise((resolve, reject) => {
    socket?.on('connect', () => {
      resolve(socket!);
    });
    socket?.on('connect_error', (err) => {
      reject(err);
    });
    socket?.on('newVideo', (data: { title: string; userName: string }) => {
      addNotification(data); // Assume NotificationContext integration
    });
  });
};

export const disconnectWebSocket = () => {
  if (socket) {
    socket.disconnect();
    socket = null;
  }
};

export const startHeartbeat = (ws: Socket) => {
  setInterval(() => {
    if (ws.connected) ws.emit('ping');
  }, 5000); // 5 seconds heartbeat
};
