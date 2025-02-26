import React, {
  createContext,
  ReactNode,
  useEffect,
  useRef,
  useState,
} from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../services/api';
import { connectWebSocket, disconnectWebSocket } from '../services/websocket';
import { v4 as uuidv4 } from 'uuid';
import {
  LoginRequest,
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
  register: (
    inputs: RegisterRequest,
    setErr: React.Dispatch<React.SetStateAction<string>>,
    event: React.FormEvent<HTMLFormElement>,
  ) => Promise<void>;
  loading: boolean;
  logout: () => Promise<void>;
  checkTokenAndConnect: () => Promise<void>;
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
      const connID = uuidv4();

      localStorage.setItem('user', JSON.stringify(user));
      localStorage.setItem('accessToken', data.access_token);
      localStorage.setItem('refreshToken', data.refresh_token);
      localStorage.setItem('connID', connID);

      await setupWebSocket(data.otp, connID);
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

      const connID = uuidv4();
      localStorage.setItem('accessToken', data.access_token);
      localStorage.setItem('refreshToken', data.refresh_token);
      localStorage.setItem('user', JSON.stringify(user));
      localStorage.setItem('connID', connID);

      await setupWebSocket(data.otp, connID);
      navigate('/');
    } catch (err: any) {
      setErr(err.response?.data?.error || 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  const logout = async () => {
    const token = localStorage.getItem('accessToken');
    const refreshToken = localStorage.getItem('refreshToken');

    try {
      if (token && refreshToken && user?.id) {
        await api.post(
          `/accounts/logout/${user.id}`,
          {},
          {
            headers: {
              Authorization: `Bearer ${token}`,
              'X-Authorization': `Bearer ${refreshToken}`,
            },
          },
        );
      }
    } catch (err) {
      console.error('Logout failed:', err);
    } finally {
      localStorage.clear();
      setUser(null);
      setToken(null);
      setRefreshToken(null);
      setConnId(null);
      disconnectWebSocket();
      navigate('/');
    }
  };

  const setupWebSocket = async (otp: string, connID: string) => {
    try {
      const ws = await connectWebSocket(otp, connID);
      wsRef.current = ws;
      console.log('WebSocket setup complete');
    } catch (err) {
      console.error('WebSocket setup failed:', err);
    }
  };

  const checkTokenAndConnect = async () => {
    const savedToken = localStorage.getItem('accessToken');
    if (!savedToken) {
      console.log('No token found, skipping check');
      return;
    }

    try {
      const res = await api.get<ApiResponse<{ otp: string }>>(
        '/accounts/check-token',
        {
          headers: { Authorization: `Bearer ${savedToken}` },
        },
      );
      const otp = res.data.data.otp;
      const savedConnId = localStorage.getItem('connId') || uuidv4();

      if (!localStorage.getItem('connId')) {
        localStorage.setItem('connId', savedConnId);
      }
      setConnId(savedConnId);

      await setupWebSocket(otp, savedConnId);
      console.log('Token valid, WebSocket connected with OTP:', otp);
    } catch (err) {
      console.warn('Check token failed:', err);
    }
  };

  useEffect(() => {
    if (localStorage.getItem('accessToken')) {
      setToken(localStorage.getItem('accessToken'));
      checkTokenAndConnect();
    }

    const handleTokenRefreshed = () => {
      setToken(localStorage.getItem('accessToken'));
      checkTokenAndConnect();
    };

    window.addEventListener('tokenRefreshed', handleTokenRefreshed);
    return () => {
      window.removeEventListener('tokenRefreshed', handleTokenRefreshed);
    };
  }, []);

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        refreshToken,
        connId,
        login,
        register,
        loading,
        logout: logout,
        checkTokenAndConnect,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};
