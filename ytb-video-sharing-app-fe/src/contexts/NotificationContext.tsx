import { createContext, useEffect, useState } from 'react';
import { getWebSocket } from '../services/websocket';

interface Notification {
  title: string;
  shared_by: string;
  thumbnail: string;
}

interface EventMessage {
  type: string;
  payload: Notification;
}

interface NotificationContextProps {
  notification: Notification | null;
}

export const NotificationContext = createContext<
  NotificationContextProps | undefined
>(undefined);

export const NotificationProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [notification, setNotification] = useState<Notification | null>(null);

  useEffect(() => {
    const ws = getWebSocket();

    if (ws) {
      ws.onmessage = (event) => {
        try {
          const data: EventMessage = JSON.parse(event.data);
          console.log('WebSocket message received:', data);

          switch (data.type) {
            case 'new_video': {
              setNotification({
                title: data.payload.title,
                shared_by: data.payload.shared_by,
                thumbnail: data.payload.thumbnail,
              });

              setTimeout(() => {
                setNotification(null);
              }, 6000);
              break;
            }
          }
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      };
    }
  }, [notification, setNotification]);

  return (
    <NotificationContext.Provider value={{ notification }}>
      {children}
    </NotificationContext.Provider>
  );
};
