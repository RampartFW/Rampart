import React from 'react';

const Dashboard: React.FC = () => {
  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
        <div className="card p-6">
          <h3 className="text-sm font-medium text-gray-500 truncate">Active Rules</h3>
          <p className="mt-2 text-3xl font-bold text-gray-900 dark:text-white">128</p>
        </div>
        <div className="card p-6">
          <h3 className="text-sm font-medium text-gray-500 truncate">Cluster Nodes</h3>
          <p className="mt-2 text-3xl font-bold text-gray-900 dark:text-white">3</p>
        </div>
        <div className="card p-6">
          <h3 className="text-sm font-medium text-gray-500 truncate">Snapshots</h3>
          <p className="mt-2 text-3xl font-bold text-gray-900 dark:text-white">25</p>
        </div>
        <div className="card p-6">
          <h3 className="text-sm font-medium text-gray-500 truncate">Audit Events (24h)</h3>
          <p className="mt-2 text-3xl font-bold text-gray-900 dark:text-white">42</p>
        </div>
      </div>
      
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <div className="card p-6 min-h-[300px]">
          <h3 className="text-lg font-medium mb-4">Recent Activity</h3>
          <div className="space-y-4">
            {[1, 2, 3].map(i => (
              <div key={i} className="flex items-center gap-4 py-3 border-b border-gray-100 dark:border-gray-700 last:border-0">
                <div className="h-10 w-10 rounded-full bg-indigo-100 dark:bg-indigo-900/30 flex items-center justify-center text-indigo-600 dark:text-indigo-400">
                  <span className="text-xs font-bold">UA</span>
                </div>
                <div>
                  <p className="text-sm font-medium">Policy Applied</p>
                  <p className="text-xs text-gray-500">ersin updated production-web-tier • 2 hours ago</p>
                </div>
              </div>
            ))}
          </div>
        </div>
        
        <div className="card p-6 min-h-[300px]">
          <h3 className="text-lg font-medium mb-4">System Status</h3>
          <div className="space-y-4">
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-600 dark:text-gray-400">Backend</span>
              <span className="text-sm font-medium px-2 py-1 bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-400 rounded">nftables</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-600 dark:text-gray-400">Raft Role</span>
              <span className="text-sm font-medium px-2 py-1 bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-400 rounded">Leader</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-600 dark:text-gray-400">Uptime</span>
              <span className="text-sm font-medium">12d 4h 22m</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
