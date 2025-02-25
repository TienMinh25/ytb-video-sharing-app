export interface User {
  id: number;
  email: string;
  name: string;
}

export interface TokenResponse {
  accessToken: string;
  refreshToken: string;
  user: User;
}
