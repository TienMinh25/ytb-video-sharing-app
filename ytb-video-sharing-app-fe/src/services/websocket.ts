let socket: WebSocket | null = null;

export const connectWebSocket = (
  token: string,
  refreshAccessToken: () => Promise<void>,
  onConnIdReceived: (connId: string) => void, // Callback để báo connId
): Promise<WebSocket> => {
  if (socket && socket.readyState === WebSocket.OPEN) {
    console.log('Reusing existing WebSocket connection');
    return Promise.resolve(socket);
  }

  if (socket) {
    socket.close();
    socket = null;
  }

  return new Promise((resolve, reject) => {
    const ws = new WebSocket(import.meta.env.VITE_WS_URL, [token]);
    socket = ws;

    ws.onopen = () => {
      console.log('WebSocket connected');
      resolve(ws); // Resolve nhưng không đóng, giữ ws sống
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.conn_id) {
          localStorage.setItem('connId', data.conn_id);
          onConnIdReceived(data.conn_id); // Gọi callback để set state
          console.log('Connection ID received:', data.conn_id);
        }
      } catch (error) {
        console.error('Failed to parse message:', error);
      }
    };

    ws.onerror = (err) => {
      console.error('WebSocket error:', err);
      reject(new Error('WebSocket connection failed'));
    };

    ws.onclose = (event) => {
      console.log('WebSocket closed:', event.code, event.reason);
      socket = null;
      if (event.reason.includes('token is expired')) {
        console.warn('Token expired, refreshing...');
        refreshAccessToken()
          .then(() => {
            const newToken = localStorage.getItem('accessToken');
            if (newToken)
              connectWebSocket(newToken, refreshAccessToken, onConnIdReceived);
          })
          .catch((err) => console.error('Failed to refresh token:', err));
      }
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

export const startHeartbeat = (ws: WebSocket) => {
  const intervalId = setInterval(() => {
    if (ws.readyState === WebSocket.OPEN) {
      ws.send('ping');
      console.log('Sent ping to server');
    } else {
      clearInterval(intervalId);
      console.log('Stopped heartbeat due to closed connection');
    }
  }, 5000); // Đồng bộ với BE (5s < 5 phút timeout)

  ws.onclose = () => {
    clearInterval(intervalId);
    console.log('Heartbeat stopped on close');
  };
};
