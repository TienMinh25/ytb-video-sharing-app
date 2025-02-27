import React from 'react';
import VideoShareForm from '../components/Video/VideoShareForm';
import { useNotification } from '../hooks/useNotification';
import NotificationHandler from '../components/Notification/NotificationHandler';

const Share: React.FC = () => {
  const { notification } = useNotification();

  return (
    <div className='container mx-auto p-4'>
      {notification !== null ? <NotificationHandler /> : <></>}
      <h2 className='text-3xl font-bold mb-6 text-[var(--foreground)]'>
        Share a Funny Movie
      </h2>
      <VideoShareForm />
    </div>
  );
};

export default Share;
