export interface User {
  id: number;
  email: string;
  fullname: string;
  avatarURL: string;
}

export interface TokenResponse {
  access_token: string;
  refresh_token: string;
  avatar_url: string;
  email: string;
  fullname: string;
  id: number;
  otp: string;
}

export interface RegisterRequest {
  fullname: string;
  email: string;
  password: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RefreshTokenResponse {
  access_token: string;
  refresh_token: string;
}
