import { NavLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import type { Role } from '../types';

interface NavItem {
  path: string;
  label: string;
  icon: string;
  roles: Role[];
}

const navItems: NavItem[] = [
  { path: '/overtime/my', label: 'My Overtimes', icon: '📋', roles: ['applicant'] },
  { path: '/overtime/create', label: 'New Request', icon: '➕', roles: ['applicant'] },
  { path: '/overtime/pending', label: 'Pending Review', icon: '⏳', roles: ['checker'] },
  { path: '/overtime/checked', label: 'Checked Review', icon: '✅', roles: ['approver'] },
  { path: '/overtime/approved', label: 'Approved Records', icon: '💰', roles: ['finance'] },
  { path: '/admin/users', label: 'User Management', icon: '👥', roles: ['admin'] },
  { path: '/admin/departments', label: 'Departments', icon: '🏢', roles: ['admin'] },
  { path: '/admin/overtime', label: 'All Overtimes', icon: '📊', roles: ['admin'] },
];

export default function Layout({ children }: { children: React.ReactNode }) {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  if (!user) return null;

  const filteredItems = navItems.filter(item => item.roles.includes(user.role));

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div className="app-layout">
      <aside className="sidebar">
        <div className="sidebar-header">
          <div className="sidebar-logo">OvertimeMS</div>
          <div className="sidebar-subtitle">Management System</div>
        </div>
        <nav className="sidebar-nav">
          <div className="nav-section-title">Navigation</div>
          {filteredItems.map(item => (
            <NavLink
              key={item.path}
              to={item.path}
              className={({ isActive }) => `nav-link ${isActive ? 'active' : ''}`}
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
          <button className="nav-link" onClick={handleLogout} style={{ marginTop: 8 }}>
            <span className="nav-icon">🚪</span>
            Logout
          </button>
        </div>
      </aside>
      <main className="main-content">
        {children}
      </main>
    </div>
  );
}
