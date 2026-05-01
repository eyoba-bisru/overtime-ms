import { useState } from 'react';
import { NavLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useTheme } from '../context/ThemeContext';
import type { Role } from '../types';

interface NavItem {
  path: string;
  label: string;
  icon: string;
  roles: Role[];
}

const navItems: NavItem[] = [
  { path: '/overtime/my', label: 'My Overtimes', icon: '📋', roles: ['applicant'] },
  { path: '/overtime/create', label: 'New Request', icon: '✨', roles: ['applicant'] },
  { path: '/overtime/pending', label: 'Pending Review', icon: '⏳', roles: ['checker'] },
  { path: '/overtime/checked', label: 'Checked Review', icon: '✅', roles: ['approver'] },
  { path: '/overtime/approved', label: 'Approved Records', icon: '💰', roles: ['finance'] },
  { path: '/admin/users', label: 'User Management', icon: '👥', roles: ['admin'] },
  { path: '/admin/departments', label: 'Departments', icon: '🏢', roles: ['admin'] },
  { path: '/admin/overtime', label: 'All Overtimes', icon: '📊', roles: ['admin'] },
];

export default function Layout({ children }: { children: React.ReactNode }) {
  const { user, logout } = useAuth();
  const { theme, toggleTheme } = useTheme();
  const navigate = useNavigate();
  const [sidebarOpen, setSidebarOpen] = useState(false);

  if (!user) return null;

  const filteredItems = navItems.filter(item => item.roles.includes(user.role));

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const closeSidebar = () => setSidebarOpen(false);

  return (
    <div className="app-layout">
      <button
        className="sidebar-toggle"
        onClick={() => setSidebarOpen(prev => !prev)}
        aria-label="Toggle menu"
        id="sidebar-toggle"
      >
        {sidebarOpen ? '✕' : '☰'}
      </button>

      {sidebarOpen && (
        <div className="sidebar-overlay visible" onClick={closeSidebar} />
      )}

      <aside className={`sidebar ${sidebarOpen ? 'open' : ''}`}>
        <div className="sidebar-header" style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <div>
            <div className="sidebar-logo">OvertimeMS</div>
            <div className="sidebar-subtitle">Management System</div>
          </div>
          <button
            className="theme-toggle"
            onClick={toggleTheme}
            aria-label="Toggle theme"
            id="theme-toggle"
            title={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
          >
            {theme === 'dark' ? '☀️' : '🌙'}
          </button>
        </div>
        <nav className="sidebar-nav">
          <div className="nav-section-title">Navigation</div>
          {filteredItems.map(item => (
            <NavLink
              key={item.path}
              to={item.path}
              className={({ isActive }) => `nav-link ${isActive ? 'active' : ''}`}
              onClick={closeSidebar}
              id={`nav-${item.path.replace(/\//g, '-').substring(1)}`}
            >
              <span className="nav-icon">{item.icon}</span>
              {item.label}
            </NavLink>
          ))}
        </nav>
        <div className="sidebar-footer">
          <div className="sidebar-user">
            <div className="sidebar-avatar">
              {user.name.charAt(0).toUpperCase()}
            </div>
            <div className="sidebar-user-info">
              <div className="sidebar-user-name">{user.name}</div>
              <div className="sidebar-user-role">{user.role}</div>
            </div>
          </div>
          <button className="nav-link" onClick={handleLogout} style={{ marginTop: 8 }} id="logout-btn">
            <span className="nav-icon">🚪</span>
            Logout
          </button>
        </div>
      </aside>
      <main className="main-content page-animate">
        {children}
      </main>
    </div>
  );
}
