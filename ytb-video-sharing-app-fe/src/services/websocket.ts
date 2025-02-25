let socket: WebSocket | null = null;

export const connectWebSocket = async (otp: string, connID: string): Promise<WebSocket> => {
  if (socket && socket.readyState === WebSocket.OPEN) return socket;

  if (socket) {
    socket.close();
    socket = null;
  }

  return new Promise((resolve, reject) => {
    const ws = new WebSocket(`${import.meta.env.VITE_WS_URL}?connID=${connID}&otp=${otp}`);
    socket = ws;

    ws.onopen = () => {
      console.log('WebSocket connected');
      resolve(ws);
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.conn_id) {
          localStorage.setItem('connId', data.conn_id);
        }
      } catch (error) {
        console.error('Failed to parse message:', error);
      }
    };

    ws.onclose = async (event) => {
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
