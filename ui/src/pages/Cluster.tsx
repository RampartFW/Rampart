import React from 'react';
import { Server, Shield, Activity } from 'lucide-react';

const Cluster: React.FC = () => {
  const nodes = [
    { id: 'node-1', address: '10.0.1.1:7946', state: 'leader', backend: 'nftables', rules: 128, sync: 'Healthy', health: 'ok' },
    { id: 'node-2', address: '10.0.1.2:7946', state: 'follower', backend: 'nftables', rules: 128, sync: 'Healthy', health: 'ok' },
    { id: 'node-3', address: '10.0.1.3:7946', state: 'follower', backend: 'iptables', rules: 128, sync: 'Healthy', health: 'ok' },
  ];

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold">Cluster Status</h2>
          <p className="text-sm text-gray-500">Manage and monitor Rampart nodes in the Raft cluster</p>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="card p-6">
          <div className="flex items-center gap-3 mb-4">
            <Activity className="h-5 w-5 text-indigo-500" />
            <h3 className="text-sm font-bold uppercase tracking-wider text-gray-500">Raft State</h3>
          </div>
          <div className="space-y-2">
            <div className="flex justify-between text-sm">
              <span className="text-gray-500">Current Term</span>
              <span className="font-mono font-medium">5</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-gray-500">Commit Index</span>
              <span className="font-mono font-medium">1,248</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-gray-500">Applied Index</span>
              <span className="font-mono font-medium">1,248</span>
            </div>
          </div>
        </div>
        
        <div className="card p-6 md:col-span-2">
          <div className="flex items-center gap-3 mb-4">
            <Shield className="h-5 w-5 text-indigo-500" />
            <h3 className="text-sm font-bold uppercase tracking-wider text-gray-500">Consensus Health</h3>
          </div>
          <div className="flex items-center gap-8 h-full pb-6">
            <div className="flex flex-col items-center">
              <div className="text-3xl font-bold text-green-500">3/3</div>
              <div className="text-xs text-gray-500 uppercase font-medium">Nodes Online</div>
            </div>
            <div className="flex flex-col items-center">
              <div className="text-3xl font-bold text-blue-500">100%</div>
              <div className="text-xs text-gray-500 uppercase font-medium">Quorum Reached</div>
            </div>
            <div className="flex flex-col items-center">
              <div className="text-3xl font-bold text-indigo-500">12ms</div>
              <div className="text-xs text-gray-500 uppercase font-medium">Avg Latency</div>
            </div>
          </div>
        </div>
      </div>

      <div className="card overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
          <thead className="bg-gray-50 dark:bg-gray-800/50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Node ID</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">State</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Backend</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Rules</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Sync</th>
              <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Health</th>
            </tr>
          </thead>
          <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
            {nodes.map((node) => (
              <tr key={node.id}>
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="flex items-center gap-3">
                    <Server className="h-5 w-5 text-gray-400" />
                    <div>
                      <div className="text-sm font-medium text-gray-900 dark:text-white">{node.id}</div>
                      <div className="text-xs text-gray-500 font-mono">{node.address}</div>
                    </div>
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`px-2 py-1 text-xs font-bold rounded-full uppercase ${
                    node.state === 'leader' ? 'bg-indigo-100 text-indigo-700 dark:bg-indigo-900/30 dark:text-indigo-400' : 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-400'
                  }`}>
                    {node.state}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {node.backend}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {node.rules}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className="text-sm text-green-600 dark:text-green-400 font-medium">{node.sync}</span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-right">
                  <div className="inline-block h-3 w-3 rounded-full bg-green-500"></div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default Cluster;
