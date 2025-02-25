import React, {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
} from 'react';
import api from '../services/api';
import { useNavigate } from 'react-router-dom';
import {
  connectWebSocket,
  disconnectWebSocket,
  startHeartbeat,
} from '../services/websocket';

interface AuthContextType {
  user: any | null;
  token: string | null;
  refreshToken: string | null;
  connId: string | null;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  refreshAccessToken: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({
  children,
}) => {
  const [user, setUser] = useState<any | null>(null);
  const [token, setToken] = useState<string | null>(
    localStorage.getItem('accessToken'),
  );
  const [refreshToken, setRefreshToken] = useState<string | null>(
    localStorage.getItem('refreshToken'),
  );
  const [connId, setConnId] = useState<string | null>(
    localStorage.getItem('connId'),
  );
  const navigate = useNavigate();

  const login = async (email: string, password: string) => {
    try {
      const response = await api.post('/accounts/login', { email, password });
      const {
        accessToken,
        refreshToken: newRefreshToken,
        user: userData,
      } = response.data;
      setToken(accessToken);
      setRefreshToken(newRefreshToken);
      setUser(userData);
      localStorage.setItem('accessToken', accessToken);
      localStorage.setItem('refreshToken', newRefreshToken);
      await setupWebSocket(accessToken);
      navigate('/');
    } catch (err) {
      console.error('Login failed', err);
      throw err;
    }
  };

  const logout = async () => {
    if (!token || !refreshToken || !connId) return;
    try {
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
    } catch (err) {
      console.error('Logout failed', err);
    } finally {
      setUser(null);
      setToken(null);
      setRefreshToken(null);
      setConnId(null);
      localStorage.removeItem('accessToken');
      localStorage.removeItem('refreshToken');
      localStorage.removeItem('connId');
      disconnectWebSocket();
      navigate('/login');
    }
  };

  const refreshAccessToken = async () => {
    if (!refreshToken) throw new Error('No refresh token available');
    try {
      const response = await api.post(
        '/accounts/refresh-token',
        {},
        {
          headers: { 'X-Authorization': refreshToken },
        },
      );
      const { accessToken } = response.data;
      setToken(accessToken);
      localStorage.setItem('accessToken', accessToken);
      return accessToken;
    } catch (err) {
      logout();
      throw err;
    }
  };

  const setupWebSocket = async (accessToken: string) => {
    try {
      const ws = await connectWebSocket(accessToken);
      ws.on('connect', () => {
        const connId = ws.id; // Assume socket.io-client provides an ID
        setConnId(connId);
        localStorage.setItem('connId', connId);
        startHeartbeat(ws);
      });
    } catch (err) {
      console.error('WebSocket setup failed', err);
    }
  };

  const value = {
    user,
    token,
    refreshToken,
    connId,
    login,
    logout,
    refreshAccessToken,
  };

  useEffect(() => {
    if (token) {
      const checkToken = async () => {
        try {
          await api.get('/accounts/me', {
            headers: { Authorization: `Bearer ${token}` },
          });
        } catch (err) {
          await refreshAccessToken();
        }
      };
      checkToken();
      setupWebSocket(token);
    }
  }, [token]);

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) throw new Error('useAuth must be used within an AuthProvider');
  return context;
};
