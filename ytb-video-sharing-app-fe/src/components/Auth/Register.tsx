import React, { useState } from 'react';
import { useAuth } from '../../hooks/useAuth';
import { RegisterRequest } from '../../types/auth';

const Register: React.FC = () => {
  const [inputs, setInputs] = useState<RegisterRequest>({
    email: '',
    password: '',
    fullname: '',
  });
  const [error, setError] = useState('');
  const { register } = useAuth();

  return (
    <div className='max-w-md mx-auto mt-12 p-8 bg-white rounded-2xl shadow-lg'>
      <h2 className='text-3xl font-bold mb-6 text-center text-gray-800'>
        Register
      </h2>
      <form
        onSubmit={(event) => register(inputs, setError, event)}
        className='space-y-6'
      >
        <div>
          <label className='block text-sm font-medium text-gray-700'>
            Name
          </label>
          <input
            type='text'
            value={inputs.fullname}
            onChange={(e) => {
              setInputs((prev) => ({ ...prev, fullname: e.target.value }));
            }}
            className='mt-2 block w-full px-4 py-3 text-lg rounded-lg border border-gray-300 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-300'
            required
          />
        </div>
        <div>
          <label className='block text-sm font-medium text-gray-700'>
            Email
          </label>
          <input
            type='email'
            value={inputs.email}
            onChange={(e) => {
              setInputs((prev) => ({ ...prev, email: e.target.value }));
            }}
            className='mt-2 block w-full px-4 py-3 text-lg rounded-lg border border-gray-300 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-300'
            required
          />
        </div>
        <div>
          <label className='block text-sm font-medium text-gray-700'>
            Password
          </label>
          <input
            type='password'
            value={inputs.password}
            onChange={(e) => {
              setInputs((prev) => ({ ...prev, password: e.target.value }));
            }}
            className='mt-2 block w-full px-4 py-3 text-lg rounded-lg border border-gray-300 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-300'
            required
          />
        </div>
        {error && <p className='text-red-500 text-center'>{error}</p>}
        <button
          type='submit'
          className='w-full bg-blue-600 text-white p-3 rounded-lg text-lg font-semibold hover:bg-blue-700 transition-colors'
        >
          Register
        </button>
      </form>
      <p className='mt-6 text-center text-gray-600'>
        Already have an account?{' '}
        <a href='/login' className='text-blue-600 hover:underline'>
          Login here
        </a>
      </p>
    </div>
  );
};

export default Register;
