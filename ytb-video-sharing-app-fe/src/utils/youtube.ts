export const extractVideoId = (url: string): string => {
  const match = url.match(/[?&]v=([^&]+)/);
  return match ? match[1] : '';
};

export const fetchYouTubeMetadata = async (videoId: string) => {
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
