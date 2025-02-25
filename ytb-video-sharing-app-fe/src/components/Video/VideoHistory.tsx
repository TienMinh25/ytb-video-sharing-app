import React, { useState, useEffect } from 'react';
import api from '../../services/api';
import { useAuth } from '../../contexts/AuthContext';

interface VideoHistory {
  id: number;
  title: string;
  youtubeUrl: string;
  watchedAt: string;
}

const VideoHistory: React.FC = () => {
  const [history, setHistory] = useState<VideoHistory[]>([]);
  const { token } = useAuth();

  useEffect(() => {
    fetchHistory();
  }, []);

  const fetchHistory = async () => {
    try {
      const response = await api.get('/videos/history', {
        headers: { Authorization: `Bearer ${token}` },
      });
      setHistory(response.data);
    } catch (err) {
      console.error('Failed to fetch video history', err);
    }
  };

  return (
    <div className='mt-8'>
      <h2 className='text-2xl font-bold mb-4 text-[var(--foreground)]'>
        Watched Videos
      </h2>
      {history.length === 0 ? (
        <p className='text-gray-500'>No watched videos yet.</p>
      ) : (
        <div className='space-y-4'>
          {history.map((video) => (
            <div key={video.id} className='bg-white rounded-lg shadow p-4'>
              <h3 className='text-xl font-semibold text-[var(--foreground)]'>
                {video.title}
              </h3>
              <p className='text-gray-600'>
                Watched on: {new Date(video.watchedAt).toLocaleString()}
              </p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default VideoHistory;
