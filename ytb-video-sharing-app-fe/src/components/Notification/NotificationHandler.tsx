import React, { useEffect, useState } from 'react';
import { getWebSocket } from '../../services/websocket';

interface Notification {
  title: string;
  shared_by: string;
  thumbnail: string;
}

const NotificationHandler: React.FC = () => {
  const [notification, setNotification] = useState<Notification | null>(null);

  useEffect(() => {
    const ws = getWebSocket();

    if (ws) {
      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          console.log('Received WebSocket message:', data);

          if (data.type === 'new_video') {
            const { title, shared_by, thumbnail } = data.payload;
            setNotification({ title, shared_by, thumbnail });

            // Tự động ẩn sau 5 giây
            const timer = setTimeout(() => {
              setNotification(null);
            }, 5000); // 5000ms = 5 giây

            // Cleanup timer nếu component unmount hoặc notification thay đổi trước khi hết 5s
            return () => clearTimeout(timer);
          }
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      };
    }

    // Cleanup khi component unmount
    return () => {
      if (ws) {
        ws.onmessage = null; // Xóa handler để tránh memory leak
      }
    };
  }, []); // Chạy một lần khi mount

  if (!notification) return null;

  return (
    <div
      style={{
        position: 'fixed',
        top: '20px',
        right: '20px',
        background: '#fff',
        padding: '20px',
        borderRadius: '10px',
        border: '1px solid #ddd',
        boxShadow: '0 4px 20px rgba(0,0,0,0.1)',
        zIndex: 1000,
        maxWidth: '300px',
        width: 'auto',
        fontFamily: 'Arial, sans-serif',
      }}
    >
      <h4
        style={{
          margin: '0 0 10px',
          fontSize: '18px',
          fontWeight: 'bold',
          color: '#333',
        }}
      >
        New Video Shared!
      </h4>
      <p style={{ margin: '0 0 10px', fontSize: '14px', color: '#555' }}>
        <strong>Title:</strong> {notification.title}
      </p>
      <p style={{ margin: '0 0 10px', fontSize: '14px', color: '#555' }}>
        <strong>Shared by:</strong> {notification.shared_by}
      </p>
      {notification.thumbnail && (
        <img
          src={notification.thumbnail}
          alt='thumbnail'
          style={{
            maxWidth: '100px',
            width: '100%',
            height: 'auto',
            borderRadius: '5px',
            marginBottom: '10px',
          }}
        />
      )}
      <button
        onClick={() => setNotification(null)}
        style={{
          padding: '8px 12px',
          backgroundColor: '#007bff',
          color: '#fff',
          border: 'none',
          borderRadius: '5px',
          cursor: 'pointer',
          fontSize: '14px',
          transition: 'background-color 0.3s',
        }}
      >
        Close
      </button>
    </div>
  );
};

export default NotificationHandler;
