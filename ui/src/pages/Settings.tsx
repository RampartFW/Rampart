import React from 'react';
import { Shield, Key, Database } from 'lucide-react';

const Settings: React.FC = () => {
  return (
    <div className="space-y-6 max-w-4xl">
      <div>
        <h2 className="text-2xl font-bold">Settings</h2>
        <p className="text-sm text-gray-500">Configure Rampart system and user preferences</p>
      </div>

      <div className="space-y-6">
        <div className="card p-6">
          <div className="flex items-center gap-3 mb-6">
            <Shield className="h-5 w-5 text-indigo-500" />
            <h3 className="text-lg font-medium">Backend Configuration</h3>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-1">
              <label className="text-xs font-bold uppercase text-gray-500">Active Backend</label>
              <div className="text-sm font-medium">nftables</div>
            </div>
            <div className="space-y-1">
              <label className="text-xs font-bold uppercase text-gray-500">Table Name</label>
              <div className="text-sm font-medium font-mono">rampart</div>
            </div>
            <div className="space-y-1">
              <label className="text-xs font-bold uppercase text-gray-500">Dual Stack (IPv4/IPv6)</label>
              <div className="text-sm font-medium text-green-600">Enabled</div>
            </div>
            <div className="space-y-1">
              <label className="text-xs font-bold uppercase text-gray-500">Atomic Apply</label>
              <div className="text-sm font-medium text-green-600">Supported (Native)</div>
            </div>
          </div>
        </div>

        <div className="card p-6">
          <div className="flex items-center gap-3 mb-6">
            <Key className="h-5 w-5 text-indigo-500" />
            <h3 className="text-lg font-medium">API Keys</h3>
          </div>
          <div className="space-y-4">
            <div className="flex justify-between items-center p-3 bg-gray-50 dark:bg-gray-700/50 rounded-lg">
              <div>
                <div className="text-sm font-medium">admin-key-01</div>
                <div className="text-xs text-gray-500 font-mono">rmp_********************</div>
              </div>
              <span className="px-2 py-0.5 text-[10px] font-bold bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400 rounded-full uppercase">Active</span>
            </div>
            <button className="btn btn-secondary w-full">Create New API Key</button>
          </div>
        </div>

        <div className="card p-6">
          <div className="flex items-center gap-3 mb-6">
            <Database className="h-5 w-5 text-indigo-500" />
            <h3 className="text-lg font-medium">Snapshot Retention</h3>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-1">
              <label className="text-xs font-bold uppercase text-gray-500">Max Snapshot Count</label>
              <div className="text-sm font-medium">100</div>
            </div>
            <div className="space-y-1">
              <label className="text-xs font-bold uppercase text-gray-500">Max Age</label>
              <div className="text-sm font-medium">30 days (720h)</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Settings;
