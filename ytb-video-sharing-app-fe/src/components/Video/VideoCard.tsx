import React from 'react';

interface Video {
  id: number;
  title: string;
  youtubeUrl: string;
  sharedBy: string;
  upvote: number;
  downvote: number;
}

const VideoCard: React.FC<{ video: Video; onView: () => void }> = ({
  video,
  onView,
}) => {
  return (
    <div className='bg-white rounded-lg shadow p-4 mb-4'>
      <iframe
        src={`https://www.youtube.com/embed/${extractVideoId(
          video.youtubeUrl,
        )}`}
        title={video.title}
        className='w-full h-48 rounded-md'
        allowFullScreen
        onLoad={onView}
      />
      <h3 className='text-xl font-semibold mt-2 text-[var(--foreground)]'>
        {video.title}
      </h3>
      <p className='text-gray-600'>Shared by: {video.sharedBy}</p>
      <div className='flex space-x-4 mt-2'>
        <span className='text-green-500'>👍 {video.upvote}</span>
        <span className='text-red-500'>👎 {video.downvote}</span>
      </div>
    </div>
  );
};

const extractVideoId = (url: string): string => {
  const match = url.match(/[?&]v=([^&]+)/);
  return match ? match[1] : '';
};

export default VideoCard;
