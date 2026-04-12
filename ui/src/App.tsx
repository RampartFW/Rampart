import React from 'react';
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

const App: React.FC = () => {
  return (
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
  );
};

export default App;
