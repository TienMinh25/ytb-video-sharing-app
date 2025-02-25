import React, { useState } from 'react';
import { useAuth } from '../../hooks/useAuth';
import api from '../../services/api';

const VideoShareForm: React.FC = () => {
  const [youtubeUrl, setYoutubeUrl] = useState('');
  const [error, setError] = useState('');
  const { token } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const videoId = extractVideoId(youtubeUrl);
      const youtubeData = await fetchYouTubeMetadata(videoId);
      await api.post(
        '/videos',
        {
          youtubeUrl,
          title: youtubeData.title,
          description: youtubeData.description,
        },
        {
          headers: { Authorization: `Bearer ${token}` },
        },
      );
      setYoutubeUrl('');
      setError('');
    } catch (err) {
      setError('Failed to share video');
    }
  };

  return (
    <div className='max-w-md mx-auto mt-10 p-6 bg-white rounded-lg shadow'>
      <h2 className='text-2xl font-bold mb-4 text-[var(--foreground)]'>
        Share a Funny Movie
      </h2>
      <form onSubmit={handleSubmit} className='space-y-4'>
        <div>
          <label className='block text-sm font-medium text-gray-700'>
            YouTube URL
          </label>
          <input
            type='url'
            value={youtubeUrl}
            onChange={(e) => setYoutubeUrl(e.target.value)}
            className='mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50'
            required
          />
        </div>
        {error && <p className='text-red-500'>{error}</p>}
        <button
          type='submit'
          className='w-full bg-[var(--color-primary)] text-white p-2 rounded-md hover:bg-opacity-90 transition-colors'
        >
          Share
        </button>
      </form>
    </div>
  );
};

// Helper functions
const extractVideoId = (url: string): string => {
  const match = url.match(/[?&]v=([^&]+)/);
  return match ? match[1] : '';
};

const fetchYouTubeMetadata = async (videoId: string) => {
  // Mock implementation (replace with actual YouTube API call)
  const response = await fetch(
    `https://www.googleapis.com/youtube/v3/videos?part=snippet&id=${videoId}&key=YOUR_YOUTUBE_API_KEY`,
  );
  const data = await response.json();
  return {
    title: data.items[0].snippet.title,
    description: data.items[0].snippet.description,
  };
};

export default VideoShareForm;
