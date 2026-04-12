import React from 'react';
import { Camera, RotateCcw, FileText, Trash2 } from 'lucide-react';

const Snapshots: React.FC = () => {
  const snapshots = [
    { id: '01JQXYZ789', timestamp: '2026-04-11 10:30:00', trigger: 'manual', rules: 12, description: 'Before maintenance window' },
    { id: '01JQXYZ456', timestamp: '2026-04-11 04:00:00', trigger: 'scheduled', rules: 12, description: 'Daily backup' },
    { id: '01JQXYZ123', timestamp: '2026-04-10 15:22:00', trigger: 'pre-apply', rules: 10, description: 'Pre-apply: policy update' },
    { id: '01JQXYZ000', timestamp: '2026-04-10 10:00:00', trigger: 'manual', rules: 10, description: 'Initial stable state' },
  ];

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold">Snapshots</h2>
          <p className="text-sm text-gray-500">Restore your firewall to a previous state at any time</p>
        </div>
        <button className="btn btn-primary">
          <Camera className="h-4 w-4 mr-2" />
          Create Snapshot
        </button>
      </div>

      <div className="card overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
          <thead className="bg-gray-50 dark:bg-gray-800/50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Created</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Trigger</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Rules</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Description</th>
              <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
            {snapshots.map((snap) => (
              <tr key={snap.id} className="hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors group">
                <td className="px-6 py-4 whitespace-nowrap text-sm">
                  <div className="font-medium text-gray-900 dark:text-white">{snap.timestamp}</div>
                  <div className="text-xs text-gray-400 font-mono">{snap.id}</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`px-2 py-0.5 text-[10px] font-bold rounded-full uppercase tracking-wider ${
                    snap.trigger === 'manual' ? 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400' :
                    snap.trigger === 'scheduled' ? 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400' :
                    'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-400'
                  }`}>
                    {snap.trigger}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {snap.rules} rules
                </td>
                <td className="px-6 py-4 text-sm text-gray-500">
                  {snap.description}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <div className="flex justify-end gap-3 opacity-0 group-hover:opacity-100 transition-opacity">
                    <button className="text-indigo-600 hover:text-indigo-900 dark:text-indigo-400 dark:hover:text-indigo-300 flex items-center gap-1">
                      <RotateCcw className="h-4 w-4" /> Rollback
                    </button>
                    <button className="text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-gray-300">
                      <FileText className="h-4 w-4" />
                    </button>
                    <button className="text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300">
                      <Trash2 className="h-4 w-4" />
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default Snapshots;
