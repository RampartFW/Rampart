import React, { useState, useEffect, createContext, useContext } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import Layout from './components/Layout';
import Dashboard from './pages/Dashboard';
import Policies from './pages/Policies';
import Rules from './pages/Rules';
import Simulator from './pages/Simulator';
import Snapshots from './pages/Snapshots';
import AuditLog from './pages/AuditLog';
import Cluster from './pages/Cluster';
import Settings from './pages/Settings';

// --- Theme Management ---
export const ThemeContext = createContext({ isDark: true, toggle: () => {} });

export const useTheme = () => useContext(ThemeContext);

const App: React.FC = () => {
  const [isDark, setIsDark] = useState(() => {
    const saved = localStorage.getItem('rampart-ui-theme');
    return saved ? saved === 'dark' : true; // Default to dark
  });

  useEffect(() => {
    const root = window.document.documentElement;
    if (isDark) {
      root.classList.add('dark');
      localStorage.setItem('rampart-ui-theme', 'dark');
    } else {
      root.classList.remove('dark');
      localStorage.setItem('rampart-ui-theme', 'light');
    }
  }, [isDark]);

  const toggleTheme = () => setIsDark(!isDark);

  return (
    <ThemeContext.Provider value={{ isDark, toggle: toggleTheme }}>
      <Layout>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/policies" element={<Policies />} />
          <Route path="/rules" element={<Rules />} />
          <Route path="/simulator" element={<Simulator />} />
          <Route path="/snapshots" element={<Snapshots />} />
          <Route path="/audit" element={<AuditLog />} />
          <Route path="/cluster" element={<Cluster />} />
          <Route path="/settings" element={<Settings />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </Layout>
    </ThemeContext.Provider>
  );
};

export default App;
