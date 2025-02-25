import React from 'react';
import VideoShareForm from '../components/Video/VideoShareForm';

const Share: React.FC = () => {
  return (
    <div className='container mx-auto p-4'>
      <h2 className='text-3xl font-bold mb-6 text-[var(--foreground)]'>
        Share a Funny Movie
      </h2>
      <VideoShareForm />
    </div>
  );
};

export default Share;
