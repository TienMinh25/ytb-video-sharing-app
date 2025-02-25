import React, { useState } from 'react';
import api from '../../services/api';

const VideoShareForm: React.FC = () => {
  const [youtubeUrl, setYoutubeUrl] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    try {
      const videoId = extractVideoId(youtubeUrl);
      if (!videoId) throw new Error('Invalid YouTube URL');

      const youtubeData = await fetchYouTubeMetadata(videoId);

      await api.post(`/videos?conn_id=${localStorage.getItem('connID')}`, {
        video_url: youtubeUrl,
        title: youtubeData.title,
        description: youtubeData.description,
        downvote: youtubeData.downvote,
        upvote: youtubeData.upvote,
        thumbnail: youtubeData.thumbnail,
      });

      setYoutubeUrl('');
    } catch (err: any) {
      setError(err.message || 'Failed to share video');
    }
  };

  return (
    <div className='max-w-md mx-auto mt-10 p-6 bg-white rounded-lg shadow'>
      <h2 className='text-2xl font-bold mb-4 text-[var(--foreground)]'>
        Share a Funny Movie
      </h2>
      <form onSubmit={handleSubmit} className='space-y-4'>
        <InputField
          label='YouTube URL'
          value={youtubeUrl}
          onChange={(e) => setYoutubeUrl(e.target.value)}
        />
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

const InputField: React.FC<{
  label: string;
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
}> = ({ label, value, onChange }) => (
  <div>
    <label className='block text-sm font-medium text-gray-700'>{label}</label>
    <input
      type='url'
      value={value}
      onChange={onChange}
      className='mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-300 
                 focus:ring focus:ring-indigo-200 focus:ring-opacity-50 p-3 text-lg'
      required
    />
  </div>
);

const extractVideoId = (url: string): string | null => {
  const match = url.match(/[?&]v=([^&]+)/) || url.match(/youtu\.be\/(.+)$/);
  return match ? match[1] : null;
};

const fetchYouTubeMetadata = async (
  videoId: string,
): Promise<{
  title: string;
  description: string;
  upvote: number;
  downvote: number;
  thumbnail: string;
}> => {
  try {
    const response = await fetch(
      `https://www.googleapis.com/youtube/v3/videos?part=snippet&id=${videoId}&key=${
        import.meta.env.VITE_API_KEY
      }`,
    );
    const data = await response.json();
    console.log(data);
    if (!data.items?.length) throw new Error('Video not found');
    return {
      title: data.items[0].snippet.title,
      description: data.items[0].snippet.description,
      upvote: 0,
      downvote: 0,
      thumbnail: data.items[0].snippet.thumbnails.default.url,
    };
  } catch (error) {
    console.log(error);
    throw new Error('Failed to fetch YouTube metadata');
  }
};

export default VideoShareForm;
