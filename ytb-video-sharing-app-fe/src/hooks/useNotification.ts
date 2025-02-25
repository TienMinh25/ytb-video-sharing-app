import { useContext } from 'react';
import { NotificationContext } from '../contexts/NotificationContext';

export const useNotification = () => {
  const context = useContext(NotificationContext);
  if (!context)
    throw new Error(
      'useNotifications must be used within a NotificationProvider',
    );
  return context;
};
