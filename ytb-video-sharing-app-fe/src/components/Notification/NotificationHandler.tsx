import React from 'react';
import { useNotification } from '../../hooks/useNotification';

const NotificationHandler: React.FC = () => {
  const { notification } = useNotification();

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
    </div>
  );
};

export default NotificationHandler;
