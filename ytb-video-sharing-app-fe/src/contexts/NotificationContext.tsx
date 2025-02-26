import { createContext, useState } from 'react';

interface Notification {
  title: string;
  shared_by: string;
  thumbnail: string;
}

interface NotificationContextProps {
  notification: Notification | null;
  showNotification: (data: Notification) => void;
}

export const NotificationContext = createContext<
  NotificationContextProps | undefined
>(undefined);

export const NotificationProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [notification, setNotification] = useState<Notification | null>(null);

  const showNotification = (data: Notification) => {
    setNotification(data);
    setTimeout(() => setNotification(null), 5000); // Ẩn sau 5s
  };

  return (
    <NotificationContext.Provider value={{ notification, showNotification }}>
      {children}
    </NotificationContext.Provider>
  );
};
