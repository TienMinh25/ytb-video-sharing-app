export interface Video {
  id: number;
  title: string;
  description: string;
  youtubeUrl: string;
  sharedBy: string;
  upvote: number;
  downvote: number;
}

export interface VideoShareRequest {
  title: string;
  description: string;
  video_url: string;
  thumbnail: string;
  upvote: number;
  downvote: number;
}
