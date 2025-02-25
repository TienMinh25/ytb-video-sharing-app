import React, { createContext, ReactNode, useState } from 'react';

interface Notification {
  title: string;
  userName: string;
}

interface NotificationContextType {
  notifications: Notification[];
  addNotification: (notif: Notification) => void;
  clearNotification: (index: number) => void;
}

export const NotificationContext = createContext<
  NotificationContextType | undefined
>(undefined);

export const NotificationProvider: React.FC<{ children: ReactNode }> = ({
  children,
}) => {
  const [notifications, setNotifications] = useState<Notification[]>([]);

  const addNotification = (notif: Notification) => {
    setNotifications((prev) => [...prev, notif]);
  };

  const clearNotification = (index: number) => {
    setNotifications((prev) => prev.filter((_, i) => i !== index));
  };

  const value = { notifications, addNotification, clearNotification };

  return (
    <NotificationContext.Provider value={value}>
      {children}
    </NotificationContext.Provider>
  );
};
