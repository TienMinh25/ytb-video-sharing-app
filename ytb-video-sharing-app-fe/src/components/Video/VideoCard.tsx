import React from 'react';
import { Video } from '../../types/video';

const VideoCard: React.FC<{ video: Video }> = ({ video }) => {
  return (
    <div className='flex bg-white rounded-lg shadow p-4 mb-4'>
      {/* left side */}
      <div className='w-4/5'>
        <iframe
          src={`https://www.youtube.com/embed/${extractVideoId(
            video.youtubeUrl,
          )}`}
          title={video.title}
          className='w-full h-80 rounded-md'
          allowFullScreen
        />
      </div>

      {/* right side */}
      <div className='w-3/5 pl-4'>
        <h3 className='text-xl font-semibold text-red-600'>{video.title}</h3>
        <p className='text-gray-600'>Shared by: {video.sharedBy}</p>
        <div className='flex space-x-4 mt-2'>
          <span className='text-green-500'>👍 {video.upvote}</span>
          <span className='text-red-500'>👎 {video.downvote}</span>
        </div>
        <p className='text-gray-700 mt-2'>
          <strong>Description:</strong> {video.description}
        </p>
      </div>
    </div>
  );
};

const extractVideoId = (url: string): string => {
  const match = url.match(/[?&]v=([^&]+)/);
  return match ? match[1] : '';
};

export default VideoCard;
