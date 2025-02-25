import React, { useContext } from 'react';
import { NotificationContext } from '../../contexts/NotificationContext';

const NotificationPopup: React.FC = () => {
  const { notifications, clearNotification } = useContext(NotificationContext);

  if (!notifications.length) return null;

  return (
    <div className='fixed top-4 right-4 z-50'>
      {notifications.map((notif, index) => (
        <div
          key={index}
          className='bg-green-500 text-white p-4 rounded-md shadow-lg mb-2 animate-slide-in'
        >
          <p>
            New video shared: "{notif.title}" by {notif.userName}
          </p>
          <button
            onClick={() => clearNotification(index)}
            className='mt-2 text-sm text-gray-200 hover:text-white'
          >
            Dismiss
          </button>
        </div>
      ))}
    </div>
  );
};

export default NotificationPopup;
