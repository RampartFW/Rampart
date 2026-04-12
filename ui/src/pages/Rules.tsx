import React from 'react';
import { Search, Plus, BarChart2 } from 'lucide-react';

const Rules: React.FC = () => {
  const rules = [
    { id: '1', name: 'allow-ssh-bastion', priority: 10, protocol: 'tcp', dport: '22', source: '10.0.1.0/24', action: 'accept', pkts: '1.2k', bytes: '840KB' },
    { id: '2', name: 'allow-http', priority: 500, protocol: 'tcp', dport: '80', source: '0.0.0.0/0', action: 'accept', pkts: '124k', bytes: '12MB' },
    { id: '3', name: 'allow-https', priority: 500, protocol: 'tcp', dport: '443', source: '0.0.0.0/0', action: 'accept', pkts: '450k', bytes: '84MB' },
    { id: '4', name: 'deny-all-ssh', priority: 10, protocol: 'tcp', dport: '22', source: '0.0.0.0/0', action: 'drop', pkts: '4.5k', bytes: '2.4MB' },
    { id: '5', name: 'allow-dns-udp', priority: 300, protocol: 'udp', dport: '53', source: '0.0.0.0/0', action: 'accept', pkts: '12k', bytes: '1.1MB' },
  ];

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold">Active Firewall Rules</h2>
          <p className="text-sm text-gray-500">Real-time view of rules currently applied to backends</p>
        </div>
        <button className="btn btn-primary">
          <Plus className="h-4 w-4 mr-2" />
          Add Quick Rule
        </button>
      </div>

      <div className="card overflow-hidden">
        <div className="p-4 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50 flex flex-col sm:flex-row gap-4">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
            <input 
              type="text" 
              placeholder="Search rules by name, IP or port..." 
              className="pl-10 input"
            />
          </div>
          <div className="flex gap-2">
            <select className="input w-auto">
              <option>All Actions</option>
              <option>Accept</option>
              <option>Drop</option>
              <option>Reject</option>
            </select>
            <select className="input w-auto">
              <option>All Protocols</option>
              <option>TCP</option>
              <option>UDP</option>
              <option>ICMP</option>
            </select>
          </div>
        </div>
        
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead className="bg-gray-50 dark:bg-gray-800/50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Priority</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Match Condition</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Action</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Stats (Pkts/Bytes)</th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
              {rules.map((rule) => (
                <tr key={rule.id} className="hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-500">
                    {rule.priority.toString().padStart(3, '0')}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm font-medium text-gray-900 dark:text-white">{rule.name}</div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex flex-wrap gap-1">
                      <span className="px-2 py-0.5 rounded bg-gray-100 dark:bg-gray-700 text-[10px] font-bold uppercase tracking-tight">{rule.protocol}</span>
                      <span className="px-2 py-0.5 rounded bg-gray-100 dark:bg-gray-700 text-[10px] font-mono tracking-tight">:{rule.dport}</span>
                      <span className="px-2 py-0.5 rounded bg-gray-100 dark:bg-gray-700 text-[10px] font-mono tracking-tight">← {rule.source}</span>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className={`px-2 py-1 text-xs font-bold rounded-full ${
                      rule.action === 'accept' 
                        ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400' 
                        : 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400'
                    }`}>
                      {rule.action.toUpperCase()}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    <div className="flex items-center gap-2">
                      <BarChart2 className="h-3 w-3" />
                      {rule.pkts} / {rule.bytes}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <button className="text-indigo-600 hover:text-indigo-900 dark:text-indigo-400 dark:hover:text-indigo-300">View</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default Rules;
