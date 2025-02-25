import React, { useState } from 'react';
import { useAuth } from '../../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';

const Login: React.FC = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await login(email, password);
      setError('');
    } catch (err) {
      setError('Invalid credentials');
    }
  };

  return (
    <div className='max-w-md mx-auto mt-10 p-6 bg-white rounded-lg shadow'>
      <h2 className='text-2xl font-bold mb-4 text-[var(--foreground)]'>
        Login
      </h2>
      <form onSubmit={handleSubmit} className='space-y-4'>
        <div>
          <label className='block text-sm font-medium text-gray-700'>
            Email
          </label>
          <input
            type='email'
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className='mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50'
            required
          />
        </div>
        <div>
          <label className='block text-sm font-medium text-gray-700'>
            Password
          </label>
          <input
            type='password'
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className='mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50'
            required
          />
        </div>
        {error && <p className='text-red-500'>{error}</p>}
        <button
          type='submit'
          className='w-full bg-[var(--color-primary)] text-white p-2 rounded-md hover:bg-opacity-90 transition-colors'
        >
          Login
        </button>
      </form>
      <p className='mt-4 text-center'>
        Don’t have an account?{' '}
        <a
          href='/register'
          className='text-[var(--color-primary)] hover:text-opacity-80'
        >
          Register here
        </a>
      </p>
    </div>
  );
};

export default Login;
