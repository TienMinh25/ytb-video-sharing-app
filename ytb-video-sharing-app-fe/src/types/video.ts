export interface Video {
  id: number;
  title: string;
  youtubeUrl: string;
  sharedBy: string;
  upvote: number;
  downvote: number;
}

export interface VideoHistory {
  id: number;
  title: string;
  youtubeUrl: string;
  watchedAt: string;
}
