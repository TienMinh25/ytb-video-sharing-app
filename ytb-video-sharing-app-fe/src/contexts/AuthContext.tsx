import React, {
  createContext,
  ReactNode,
  useEffect,
  useRef,
  useState,
} from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../services/api';
import {
  connectWebSocket,
  disconnectWebSocket,
  startHeartbeat,
} from '../services/websocket';
import {
  LoginRequest,
  RefreshTokenResponse,
  RegisterRequest,
  TokenResponse,
  User,
} from '../types/auth';
import { ApiResponse } from '../types/response';

interface AuthContextType {
  user: User | null;
  token: string | null;
  refreshToken: string | null;
  connId: string | null;
  login: (
    inputs: LoginRequest,
    setErr: React.Dispatch<React.SetStateAction<string>>,
    event: React.FormEvent<HTMLFormElement>,
  ) => Promise<void>;
  logout: () => void;
  refreshAccessToken: () => Promise<void>;
  register: (
    inputs: RegisterRequest,
    setErr: React.Dispatch<React.SetStateAction<string>>,
    event: React.FormEvent<HTMLFormElement>,
  ) => Promise<void>;
  loading: boolean;
}

export const AuthContext = createContext<AuthContextType | undefined>(
  undefined,
);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({
  children,
}) => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [user, setUser] = useState<User | null>(
    localStorage.getItem('user')
      ? JSON.parse(localStorage.getItem('user')!)
      : null,
  );
  const [token, setToken] = useState<string | null>(
    localStorage.getItem('accessToken'),
  );
  const [refreshToken, setRefreshToken] = useState<string | null>(
    localStorage.getItem('refreshToken'),
  );
  const [connId, setConnId] = useState<string | null>(
    localStorage.getItem('connId'),
  );
  const wsRef = useRef<WebSocket | null>(null);

  const register = async (
    inputs: RegisterRequest,
    setErr: React.Dispatch<React.SetStateAction<string>>,
    e: React.FormEvent<HTMLFormElement>,
  ) => {
    e.preventDefault();
    setLoading(true);
    try {
      if (
        !inputs.fullname.trim() ||
        !inputs.email.trim() ||
        !inputs.password.trim()
      ) {
        setErr('Please fill up all fields!');
        return;
      }

      const res = await api.post<ApiResponse<TokenResponse>>(
        '/accounts/register',
        inputs,
      );
      const { data } = res.data;

      const user: User = {
        id: data.id,
        email: data.email,
        fullname: data.fullname,
        avatarURL: data.avatar_url,
      };

      setToken(data.access_token);
      setRefreshToken(data.refresh_token);
      setUser(user);

      localStorage.setItem('user', JSON.stringify(user));
      localStorage.setItem('accessToken', data.access_token);
      localStorage.setItem('refreshToken', data.refresh_token);

      await setupWebSocket(data.access_token);
      navigate('/');
    } catch (err: any) {
      setErr(err.response?.data?.error || 'Registration failed');
    } finally {
      setLoading(false);
    }
  };

  const login = async (
    inputs: LoginRequest,
    setErr: React.Dispatch<React.SetStateAction<string>>,
    e: React.FormEvent<HTMLFormElement>,
  ) => {
    e.preventDefault();
    setLoading(true);
    try {
      const res = await api.post<ApiResponse<TokenResponse>>(
        '/accounts/login',
        inputs,
      );
      const { data } = res.data;

      const user: User = {
        id: data.id,
        email: data.email,
        fullname: data.fullname,
        avatarURL: data.avatar_url,
      };

      setToken(data.access_token);
      setRefreshToken(data.refresh_token);
      setUser(user);

      localStorage.setItem('accessToken', data.access_token);
      localStorage.setItem('refreshToken', data.refresh_token);
      localStorage.setItem('user', JSON.stringify(user));

      await setupWebSocket(data.access_token);
      navigate('/');
    } catch (err: any) {
      setErr(err.response?.data?.error || 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  const logout = async () => {
    try {
      if (token && refreshToken && connId) {
        await api.post(
          '/accounts/logout',
          {},
          {
            headers: {
              Authorization: `Bearer ${token}`,
              'X-Authorization': refreshToken,
              connID: connId,
            },
          },
        );
      }
    } catch (err) {
      console.error('Logout failed:', err);
    } finally {
      setUser(null);
      setToken(null);
      setRefreshToken(null);
      setConnId(null);

      localStorage.removeItem('accessToken');
      localStorage.removeItem('refreshToken');
      localStorage.removeItem('connId');
      localStorage.removeItem('user');

      disconnectWebSocket();
      navigate('/');
    }
  };

  const refreshAccessToken = async () => {
    if (!refreshToken) throw new Error('No refresh token available');
    try {
      const res = await api.post<ApiResponse<RefreshTokenResponse>>(
        '/accounts/refresh-token',
        {},
        { headers: { 'X-Authorization': refreshToken } },
      );
      const { access_token, refresh_token } = res.data.data;

      setToken(access_token);
      setRefreshToken(refresh_token);
      localStorage.setItem('accessToken', access_token);
      localStorage.setItem('refreshToken', refresh_token);

      await setupWebSocket(access_token); // Reconnect WebSocket với token mới
    } catch (err) {
      console.error('Token refresh failed:', err);
      logout();
      throw err;
    }
  };

  const setupWebSocket = async (accessToken: string) => {
    try {
      const ws = await connectWebSocket(
        accessToken,
        refreshAccessToken,
        (connId) => setConnId(connId), // Callback để set connId vào state
      );
      wsRef.current = ws;
      startHeartbeat(ws); // Bắt đầu heartbeat ngay
      console.log('WebSocket setup complete');
    } catch (err) {
      console.error('WebSocket setup failed:', err);
    }
  };

  useEffect(() => {
    if (!token) {
      disconnectWebSocket();
      setConnId(null);
      return;
    }

    console.log('Setting up WebSocket with token:', token);
    setupWebSocket(token);

    // Cleanup chỉ khi component unmount hoặc token thay đổi
    return () => {
      // Chỉ đóng khi thực sự cần, để debug thì comment dòng dưới
      // disconnectWebSocket();
      console.log('Cleanup WebSocket for token:', token);
    };
  }, [token]);

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        refreshToken,
        connId,
        login,
        logout,
        refreshAccessToken,
        register,
        loading,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};
