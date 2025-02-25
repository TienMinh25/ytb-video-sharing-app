import React, { useContext } from 'react';
import { Link } from 'react-router-dom';
import { AuthContext } from '../contexts/AuthContext';

const Layout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { user, logout } = useContext(AuthContext);

  return (
    <div className='min-h-screen bg-[var(--background)]'>
      <header className='bg-white shadow-md p-4'>
        <div className='container flex justify-between items-center'>
          <h1 className='text-2xl font-bold text-[var(--foreground)]'>
            Funny Movies
          </h1>
          <nav className='space-x-4'>
            <Link to='/' className='text-blue-500 hover:text-blue-700'>
              Home
            </Link>
            {user && (
              <Link to='/share' className='text-blue-500 hover:text-blue-700'>
                Share
              </Link>
            )}
            {user ? (
              <button
                onClick={logout}
                className='text-red-500 hover:text-red-700 cursor-pointer'
              >
                Logout
              </button>
            ) : (
              <Link to='/login' className='text-blue-500 hover:text-blue-700'>
                Login/Register
              </Link>
            )}
          </nav>
        </div>
      </header>
      <main className='container'>{children}</main>
    </div>
  );
};

export default Layout;
