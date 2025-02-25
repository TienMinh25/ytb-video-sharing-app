import React, { useState, useEffect } from 'react';
import VideoCard from '../components/Video/VideoCard';
import VideoHistory from '../components/Video/VideoHistory';
import NotificationPopup from '../components/Notification/NotificationPopup';
import api from '../services/api';
import Header from '../components/Auth/Header';

interface Video {
  id: number;
  title: string;
  youtubeUrl: string;
  sharedBy: string;
  upvote: number;
  downvote: number;
}

const Home: React.FC = () => {
  const [videos, setVideos] = useState<Video[]>([]);

  useEffect(() => {
    fetchVideos();
  }, []);

  const fetchVideos = async () => {
    try {
      const response = await api.get('/videos');
      setVideos(response.data);
    } catch (err) {
      console.error('Failed to fetch videos', err);
    }
  };

  const handleVideoView = (videoId: number) => {
    // Mock implementation for tracking watched videos
    api.post('/videos/history', { videoId });
  };

  return (
    <div className='container mx-auto p-4'>
      <NotificationPopup />
      <div className='grid gap-6'>
        <div>
          <h3 className='text-2xl font-semibold mb-4 text-[var(--foreground)]'>
            Shared Videos
          </h3>
          {videos.length === 0 ? (
            <p className='text-gray-500'>No videos shared yet.</p>
          ) : (
            <div className='grid gap-4'>
              {videos.map((video) => (
                <VideoCard
                  key={video.id}
                  video={video}
                  onView={() => handleVideoView(video.id)}
                />
              ))}
            </div>
          )}
        </div>
        <VideoHistory />
      </div>
    </div>
  );
};

export default Home;
