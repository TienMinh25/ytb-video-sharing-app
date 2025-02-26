let socket: WebSocket | null = null;

const API_WS_URL = import.meta.env.VITE_API_URL || 'ws://localhost:3001/ws';

export const connectWebSocket = async (
  otp: string,
  connID: string,
): Promise<WebSocket> => {
  if (socket && socket.readyState === WebSocket.OPEN) return socket;

  if (socket) {
    socket.close();
    socket = null;
  }

  return new Promise((resolve, reject) => {
    const ws = new WebSocket(`${API_WS_URL}?connID=${connID}&otp=${otp}`);
    socket = ws;

    ws.onopen = () => {
      console.log('WebSocket connected');
      resolve(ws);
    };

    ws.onclose = (event) => {
      console.log('WebSocket closed:', event.code, event.reason);
      socket = null;
    };

    ws.onerror = (err) => {
      console.error('WebSocket error:', err);
      reject(new Error('WebSocket connection failed'));
    };
  });
};

export const disconnectWebSocket = () => {
  if (socket) {
    socket.close();
    socket = null;
    console.log('WebSocket manually disconnected');
  }
};

export const getWebSocket = (): WebSocket | null => socket;
