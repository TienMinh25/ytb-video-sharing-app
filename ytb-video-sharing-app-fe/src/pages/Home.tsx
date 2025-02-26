import React, { useCallback, useEffect, useState } from 'react';
import NotificationPopup from '../components/Notification/NotificationHandler';
import VideoCard from '../components/Video/VideoCard';
import api from '../services/api';
import { ApiResponse } from '../types/response';
import { Video } from '../types/video';

interface VideoResponse {
  id: number;
  title: string;
  description: string;
  thumbnail: string;
  upvote: number;
  downvote: number;
  video_url: string;
  shared_by: string;
}

const Home: React.FC = () => {
  const [initialLoading, setInitialLoading] = useState(false);
  const [fetchingMore, setFetchingMore] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const [videos, setVideos] = useState<Video[]>([]);
  const [page, setPage] = useState(1);
  const limit = 6;

  const mapVideoResponse = (video: VideoResponse): Video => ({
    id: video.id,
    title: video.title,
    description: video.description,
    youtubeUrl: video.video_url,
    sharedBy: video.shared_by,
    upvote: video.upvote,
    downvote: video.downvote,
  });

  const fetchVideos = async (pageNumber: number) => {
    try {
      const res = await api.get<ApiResponse<VideoResponse[]>>(
        `/videos?page=${pageNumber}&limit=${limit}`,
      );
      const { data, metadata } = res.data;

      if (pageNumber === 1) {
        setVideos(data.map(mapVideoResponse));
      } else {
        setVideos((prev) => [...prev, ...data.map(mapVideoResponse)]);
      }

      setHasMore(metadata.pagination.is_next);
      setPage(pageNumber + 1);
    } catch (error) {
      console.error('Failed to fetch videos', error);
      setHasMore(false);
    }
  };

  useEffect(() => {
    setInitialLoading(true);
    fetchVideos(1).finally(() => setInitialLoading(false));
  }, []);

  const fetchMoreVideos = useCallback(async () => {
    if (fetchingMore || !hasMore) return;

    setFetchingMore(true);
    await fetchVideos(page);
    setFetchingMore(false);
  }, [page, fetchingMore, hasMore]);

  useEffect(() => {
    const handleScroll = () => {
      const { scrollTop, clientHeight, scrollHeight } =
        document.documentElement;
      if (scrollTop + clientHeight >= scrollHeight - 100 && hasMore) {
        fetchMoreVideos();
      }
    };

    window.addEventListener('scroll', handleScroll);
    return () => {
      window.removeEventListener('scroll', handleScroll);
    };
  }, [fetchMoreVideos, hasMore]);

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
                <VideoCard key={video.id} video={video} />
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default Home;
